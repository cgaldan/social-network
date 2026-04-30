"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import DashboardNav from "../../components/DashboardNav";
import { getCurrentUser, getUnreadNotificationCount } from "../../lib/api";

export default function DashboardLayout({ children }) {
  const router = useRouter();
  const [ready, setReady] = useState(false);
  const [unreadCount, setUnreadCount] = useState(0);

  useEffect(() => {
    const token = localStorage.getItem("authToken");
    if (!token) {
      router.replace("/login");
      return;
    }

    let cancelled = false;

    (async () => {
      try {
        const me = await getCurrentUser(token);
        if (me.user) {
          localStorage.setItem("currentUser", JSON.stringify(me.user));
        }
      } catch {
        localStorage.removeItem("authToken");
        localStorage.removeItem("currentUser");
        if (!cancelled) router.replace("/login");
        return;
      }
      try {
        const unread = await getUnreadNotificationCount(token);
        if (!cancelled) setUnreadCount(unread.unread_count ?? 0);
      } catch {
        if (!cancelled) setUnreadCount(0);
      }
      if (!cancelled) setReady(true);
    })();

    return () => {
      cancelled = true;
    };
  }, [router]);

  if (!ready) {
    return (
      <div className="app-shell">
        <p className="helper-text app-loading">Loading…</p>
      </div>
    );
  }

  return (
    <div className="app-shell">
      <DashboardNav unreadCount={unreadCount} />
      <main className="app-main">{children}</main>
    </div>
  );
}
