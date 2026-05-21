"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useEffect } from "react";
import { useAuth } from "@/components/AuthProvider";
import { useWebSocket } from "@/components/WebSocketProvider";
import { useNotificationCount } from "@/components/NotificationCountProvider";
import { useMessageCount } from "@/components/MessageCountProvider";
import { Avatar } from "@/components/Avatar";
import { ThemeToggle } from "@/components/ThemeToggle";
import type { Notification } from "@/types/api";

const NAV = [
  { href: "/feed", label: "Feed", icon: "📰" },
  { href: "/groups", label: "Groups", icon: "👥" },
  { href: "/messages", label: "Messages", icon: "💬" },
  { href: "/notifications", label: "Notifications", icon: "🔔" },
  { href: "/follow", label: "Follow", icon: "➕" },
  { href: "/profile", label: "Profile", icon: "👤" },
];

export function AppShell({ children }: { children: React.ReactNode }) {
  const { user, logout } = useAuth();
  const { on, online, connected } = useWebSocket();
  const { count: unread, increment } = useNotificationCount();
  const { count: unreadMessages } = useMessageCount();
  const pathname = usePathname();
  const router = useRouter();

  useEffect(() => {
    return on((msg) => {
      if (msg.type === "notification.created") {
        const payload = msg.payload as { notification: Notification };
        if (!payload.notification.read_at) {
          increment();
        }
      }
    });
  }, [on, increment]);

  const onLogout = async () => {
    await logout();
    router.replace("/login");
  };

  return (
    <div className="flex min-h-screen bg-slate-50">
      <aside className="hidden w-64 flex-shrink-0 border-r border-slate-200 bg-white p-4 md:flex md:flex-col">
        <Link href="/feed" className="mb-6 flex items-center gap-2 px-2">
          <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-indigo-600 text-white">
            ◈
          </span>
          <span className="text-lg font-semibold text-slate-900">Social</span>
        </Link>

        <nav className="flex-1 space-y-1">
          {NAV.map((item) => {
            const active = pathname === item.href || pathname.startsWith(item.href + "/");
            let badge = 0;
            if (item.href === "/notifications") badge = unread;
            else if (item.href === "/messages") badge = unreadMessages;
            return (
              <Link
                key={item.href}
                href={item.href}
                className={`flex items-center justify-between rounded-lg px-3 py-2 text-sm font-medium transition ${
                  active
                    ? "bg-indigo-50 text-indigo-700"
                    : "text-slate-700 hover:bg-slate-100"
                }`}
              >
                <span className="flex items-center gap-3">
                  <span aria-hidden>{item.icon}</span>
                  <span>{item.label}</span>
                </span>
                {badge > 0 ? (
                  <span className="rounded-full bg-indigo-600 px-2 py-0.5 text-xs font-medium text-white">
                    {badge}
                  </span>
                ) : null}
              </Link>
            );
          })}
        </nav>

        {user ? (
          <div className="mt-4 rounded-lg border border-slate-200 p-3">
            <Link href="/profile" className="flex items-center gap-3">
              <Avatar
                src={user.avatar_path}
                firstName={user.first_name}
                lastName={user.last_name}
                nickname={user.nickname}
                size={36}
              />
              <div className="min-w-0">
                <div className="truncate text-sm font-medium text-slate-900">
                  {user.first_name} {user.last_name}
                </div>
                <div className="truncate text-xs text-slate-500">@{user.nickname}</div>
              </div>
            </Link>
            <button
              onClick={onLogout}
              className="mt-3 w-full rounded-md border border-slate-200 px-3 py-1.5 text-xs font-medium text-slate-600 hover:bg-slate-50"
            >
              Sign out
            </button>
          </div>
        ) : null}

        <div className="mt-3 flex items-center justify-between gap-2 px-2 text-xs text-slate-500">
          <div className="flex items-center gap-2">
            <span
              className={`h-2 w-2 rounded-full ${connected ? "bg-emerald-500" : "bg-slate-300"}`}
              aria-hidden
            />
            {connected ? `${online.length} online` : "Offline"}
          </div>
          <ThemeToggle />
        </div>
      </aside>

      <div className="flex flex-1 flex-col">
        <header className="flex items-center justify-between border-b border-slate-200 bg-white px-4 py-3 md:hidden">
          <Link href="/feed" className="text-lg font-semibold text-slate-900">
            Social
          </Link>
          <div className="flex items-center gap-3">
            <ThemeToggle />
            {user ? (
              <Link href="/profile">
                <Avatar
                  src={user.avatar_path}
                  firstName={user.first_name}
                  lastName={user.last_name}
                  nickname={user.nickname}
                  size={32}
                />
              </Link>
            ) : null}
          </div>
        </header>

        <main className="flex-1">
          <div className="mx-auto max-w-3xl px-4 py-6">{children}</div>
        </main>

        <nav className="sticky bottom-0 grid grid-cols-6 border-t border-slate-200 bg-white md:hidden">
          {NAV.map((item) => {
            const active = pathname === item.href || pathname.startsWith(item.href + "/");
            return (
              <Link
                key={item.href}
                href={item.href}
                className={`flex flex-col items-center py-2 text-xs ${
                  active ? "text-indigo-600" : "text-slate-600"
                }`}
              >
                <span aria-hidden>{item.icon}</span>
                <span className="mt-0.5">{item.label}</span>
              </Link>
            );
          })}
        </nav>
      </div>
    </div>
  );
}
