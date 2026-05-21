"use client";

import { createContext, useCallback, useContext, useEffect, useState } from "react";
import { useAuth } from "@/components/AuthProvider";
import { api } from "@/lib/api";

interface ContextValue {
  count: number;
  increment: (by?: number) => void;
  decrement: (by?: number) => void;
  reset: () => void;
  refetch: () => Promise<void>;
}

const NotificationCountContext = createContext<ContextValue | null>(null);

export function NotificationCountProvider({ children }: { children: React.ReactNode }) {
  const { user } = useAuth();
  const [count, setCount] = useState(0);

  const refetch = useCallback(async () => {
    if (!user) {
      setCount(0);
      return;
    }
    try {
      const res = await api.notifications.unreadCount();
      setCount(res.unread_count);
    } catch {
      // ignore — leave stale count rather than zero it on transient errors
    }
  }, [user]);

  useEffect(() => {
    void refetch();
  }, [refetch]);

  const increment = useCallback((by = 1) => {
    setCount((c) => c + by);
  }, []);
  const decrement = useCallback((by = 1) => {
    setCount((c) => Math.max(0, c - by));
  }, []);
  const reset = useCallback(() => setCount(0), []);

  return (
    <NotificationCountContext.Provider value={{ count, increment, decrement, reset, refetch }}>
      {children}
    </NotificationCountContext.Provider>
  );
}

export function useNotificationCount(): ContextValue {
  const ctx = useContext(NotificationCountContext);
  if (!ctx)
    throw new Error("useNotificationCount must be used inside NotificationCountProvider");
  return ctx;
}
