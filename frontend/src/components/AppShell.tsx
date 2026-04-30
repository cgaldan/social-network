"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useEffect, useState } from "react";
import { useAuth } from "@/components/AuthProvider";
import { api } from "@/lib/api";
import { connectWebSocket } from "@/lib/ws";

const navItems = [
  { href: "/feed", label: "Feed" },
  { href: "/follow", label: "Follow" },
  { href: "/groups", label: "Groups" },
  { href: "/messages", label: "Messages" },
  { href: "/notifications", label: "Notifications" },
  { href: "/profile", label: "Profile" },
];

export function AppShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const { logout, user } = useAuth();
  const [unreadCount, setUnreadCount] = useState(0);
  const [wsStatus, setWsStatus] = useState("offline");

  useEffect(() => {
    api
      .unreadNotifications()
      .then((response) => setUnreadCount(response.unread_count))
      .catch(() => setUnreadCount(0));

    const socket = connectWebSocket({
      onOpen: () => setWsStatus("live"),
      onClose: () => setWsStatus("offline"),
      onError: () => setWsStatus("error"),
      onMessage: () => setUnreadCount((count) => count + 1),
    });

    return () => socket?.close();
  }, []);

  return (
    <div className="min-h-screen">
      <header className="sticky top-0 z-20 border-b border-slate-200 bg-white/90 backdrop-blur">
        <div className="mx-auto flex max-w-7xl flex-col gap-4 px-4 py-4 md:flex-row md:items-center md:justify-between">
          <Link className="text-xl font-bold text-slate-950" href="/feed">
            Social Network
          </Link>
          <nav className="flex flex-wrap gap-2">
            {navItems.map((item) => {
              const active = pathname.startsWith(item.href);
              return (
                <Link
                  className={`rounded-full px-4 py-2 text-sm font-semibold transition ${
                    active
                      ? "bg-sky-600 text-white"
                      : "text-slate-600 hover:bg-slate-100 hover:text-slate-950"
                  }`}
                  href={item.href}
                  key={item.href}
                >
                  {item.label}
                  {item.href === "/notifications" && unreadCount > 0
                    ? ` (${unreadCount})`
                    : ""}
                </Link>
              );
            })}
          </nav>
          <div className="flex items-center gap-3">
            <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold text-slate-600">
              WS: {wsStatus}
            </span>
            <span className="hidden text-sm text-slate-600 md:inline">
              {user?.nickname || user?.email}
            </span>
            <button
              className="rounded-full border border-slate-300 px-4 py-2 text-sm font-semibold text-slate-700 transition hover:border-red-300 hover:text-red-600"
              onClick={logout}
              type="button"
            >
              Logout
            </button>
          </div>
        </div>
      </header>
      <main className="mx-auto max-w-7xl px-4 py-8">{children}</main>
    </div>
  );
}
