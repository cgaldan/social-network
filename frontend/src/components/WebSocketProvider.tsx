"use client";

import { createContext, useCallback, useContext, useEffect, useMemo, useRef, useState } from "react";
import { useAuth } from "@/components/AuthProvider";
import { connectWs, type WsClient, type WsListener } from "@/lib/ws";
import type { OnlineUser, WsEnvelope } from "@/types/api";

interface WebSocketContextValue {
  online: OnlineUser[];
  send: (msg: WsEnvelope) => void;
  on: (listener: WsListener) => () => void;
  connected: boolean;
}

const WebSocketContext = createContext<WebSocketContextValue | null>(null);

export function WebSocketProvider({ children }: { children: React.ReactNode }) {
  const { token, user } = useAuth();
  const clientRef = useRef<WsClient | null>(null);
  const listenersRef = useRef(new Set<WsListener>());
  const [online, setOnline] = useState<OnlineUser[]>([]);
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    if (!token || !user) {
      clientRef.current?.close();
      clientRef.current = null;
      setConnected(false);
      setOnline([]);
      return;
    }

    const client = connectWs(token);
    clientRef.current = client;
    setConnected(true);

    const unsub = client.on((msg) => {
      if (msg.type === "online_users") {
        const payload = msg.payload as { users?: OnlineUser[] };
        setOnline(payload.users ?? []);
      } else if (msg.type === "user_status") {
        const payload = msg.payload as OnlineUser;
        setOnline((prev) => {
          const without = prev.filter((u) => u.user_id !== payload.user_id);
          return payload.is_online ? [...without, payload] : without;
        });
      }
      listenersRef.current.forEach((l) => l(msg));
    });

    return () => {
      unsub();
      client.close();
      clientRef.current = null;
      setConnected(false);
    };
  }, [token, user]);

  const send = useCallback((msg: WsEnvelope) => {
    clientRef.current?.send(msg);
  }, []);

  const on = useCallback((listener: WsListener) => {
    listenersRef.current.add(listener);
    return () => {
      listenersRef.current.delete(listener);
    };
  }, []);

  const value = useMemo<WebSocketContextValue>(
    () => ({ online, send, on, connected }),
    [online, send, on, connected],
  );

  return <WebSocketContext.Provider value={value}>{children}</WebSocketContext.Provider>;
}

export function useWebSocket(): WebSocketContextValue {
  const ctx = useContext(WebSocketContext);
  if (!ctx) throw new Error("useWebSocket must be used inside WebSocketProvider");
  return ctx;
}
