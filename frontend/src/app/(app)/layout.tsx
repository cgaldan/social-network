"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/components/AuthProvider";
import { WebSocketProvider } from "@/components/WebSocketProvider";
import { NotificationCountProvider } from "@/components/NotificationCountProvider";
import { MessageCountProvider } from "@/components/MessageCountProvider";
import { AppShell } from "@/components/AppShell";

export default function AppLayout({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && !user) router.replace("/login");
  }, [user, loading, router]);

  if (loading || !user) {
    return (
      <div className="flex min-h-screen items-center justify-center text-slate-500">
        Loading…
      </div>
    );
  }

  return (
    <WebSocketProvider>
      <NotificationCountProvider>
        <MessageCountProvider>
          <AppShell>{children}</AppShell>
        </MessageCountProvider>
      </NotificationCountProvider>
    </WebSocketProvider>
  );
}
