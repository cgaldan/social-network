"use client";

import { createContext, useCallback, useContext, useEffect, useState } from "react";
import { useAuth } from "@/components/AuthProvider";
import { useWebSocket } from "@/components/WebSocketProvider";
import { api } from "@/lib/api";
import type { Message, MessageCreatedPayload } from "@/types/api";

interface ContextValue {
  count: number;
  increment: (by?: number) => void;
  decrement: (by?: number) => void;
  reset: () => void;
  refetch: () => Promise<void>;
}

const MessageCountContext = createContext<ContextValue | null>(null);

export function MessageCountProvider({ children }: { children: React.ReactNode }) {
  const { user } = useAuth();
  const { on } = useWebSocket();
  const [count, setCount] = useState(0);

  const refetch = useCallback(async () => {
    if (!user) {
      setCount(0);
      return;
    }
    try {
      const res = await api.conversations.list({ limit: 100 });
      const total = (res.conversations ?? []).reduce(
        (sum, c) => sum + (c.unread_count ?? 0),
        0,
      );
      setCount(total);
    } catch {
      // ignore — keep stale count
    }
  }, [user]);

  useEffect(() => {
    void refetch();
  }, [refetch]);

  useEffect(() => {
    return on((msg) => {
      if (msg.type !== "message.created") return;
      const payload = msg.payload as MessageCreatedPayload;
      const incoming: Message = payload.message;
      if (incoming.sender_id !== user?.id) {
        setCount((c) => c + 1);
      }
    });
  }, [on, user?.id]);

  const increment = useCallback((by = 1) => setCount((c) => c + by), []);
  const decrement = useCallback((by = 1) => setCount((c) => Math.max(0, c - by)), []);
  const reset = useCallback(() => setCount(0), []);

  return (
    <MessageCountContext.Provider value={{ count, increment, decrement, reset, refetch }}>
      {children}
    </MessageCountContext.Provider>
  );
}

export function useMessageCount(): ContextValue {
  const ctx = useContext(MessageCountContext);
  if (!ctx) throw new Error("useMessageCount must be used inside MessageCountProvider");
  return ctx;
}
