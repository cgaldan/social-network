"use client";

import Image from "next/image";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
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

const SIDEBAR_STORAGE_KEY = "sidebar-collapsed";

export function AppShell({ children }: { children: React.ReactNode }) {
  const { user, logout } = useAuth();
  const { on, online, connected } = useWebSocket();
  const { count: unread, increment } = useNotificationCount();
  const { count: unreadMessages } = useMessageCount();
  const pathname = usePathname();
  const router = useRouter();

  // Collapse state. SSR can't read localStorage, so we default to expanded
  // and hydrate from storage on mount. There's a brief flash for users who
  // prefer collapsed; fine for a sidebar.
  const [collapsed, setCollapsed] = useState(false);
  useEffect(() => {
    if (typeof window === "undefined") return;
    const stored = window.localStorage.getItem(SIDEBAR_STORAGE_KEY);
    if (stored !== null) setCollapsed(stored === "1");
  }, []);
  useEffect(() => {
    if (typeof window === "undefined") return;
    window.localStorage.setItem(SIDEBAR_STORAGE_KEY, collapsed ? "1" : "0");
  }, [collapsed]);

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
      <aside
        className={`sticky top-0 hidden h-screen flex-shrink-0 overflow-y-auto border-r border-slate-200 bg-white transition-[width,padding] duration-200 ease-in-out md:flex md:flex-col ${
          collapsed ? "w-16 p-2" : "w-64 p-4"
        }`}
      >
        <div
          className={`mb-6 flex ${
            collapsed ? "flex-col items-center gap-2" : "items-center justify-between px-2"
          }`}
        >
          <Link href="/feed" className="flex items-center gap-2">
            <Image
              src="/cyclops.png"
              alt="Social"
              width={32}
              height={32}
              className="h-8 w-8"
              priority
            />
            {!collapsed ? (
              <span className="text-lg font-semibold text-slate-900">Femus</span>
            ) : null}
          </Link>
          <button
            type="button"
            onClick={() => setCollapsed((c) => !c)}
            aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
            title={collapsed ? "Expand sidebar" : "Collapse sidebar"}
            className="rounded p-1 text-slate-400 transition hover:bg-slate-100 hover:text-slate-700"
          >
            {collapsed ? "▶" : "◀"}
          </button>
        </div>

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
                title={collapsed ? item.label : undefined}
                className={`relative flex items-center rounded-lg text-sm font-medium transition ${
                  collapsed ? "justify-center px-2 py-2" : "justify-between px-3 py-2"
                } ${
                  active
                    ? "bg-indigo-50 text-indigo-700"
                    : "text-slate-700 hover:bg-slate-100"
                }`}
              >
                <span className={`flex items-center ${collapsed ? "" : "gap-3"}`}>
                  <span aria-hidden className="text-base">
                    {item.icon}
                  </span>
                  {!collapsed ? <span>{item.label}</span> : null}
                </span>
                {badge > 0 ? (
                  collapsed ? (
                    <span
                      aria-label={`${badge} unread`}
                      className="absolute right-1 top-1 h-2 w-2 rounded-full bg-indigo-600 ring-2 ring-white"
                    />
                  ) : (
                    <span className="rounded-full bg-indigo-600 px-2 py-0.5 text-xs font-medium text-white">
                      {badge}
                    </span>
                  )
                ) : null}
              </Link>
            );
          })}
        </nav>

        {user ? (
          collapsed ? (
            <div className="mt-4 flex flex-col items-center gap-2">
              <Link href="/profile" title={`@${user.nickname}`}>
                <Avatar
                  src={user.avatar_path}
                  firstName={user.first_name}
                  lastName={user.last_name}
                  nickname={user.nickname}
                  size={36}
                />
              </Link>
              <button
                onClick={onLogout}
                aria-label="Sign out"
                title="Sign out"
                className="rounded-md p-1.5 text-slate-500 transition hover:bg-slate-100 hover:text-red-600"
              >
                ⏻
              </button>
            </div>
          ) : (
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
          )
        ) : null}

        <div
          className={`mt-3 flex items-center text-xs text-slate-500 ${
            collapsed ? "flex-col gap-2" : "justify-between gap-2 px-2"
          }`}
        >
          <div
            className="flex items-center gap-2"
            title={connected ? `${online.length} online` : "Offline"}
          >
            <span
              className={`h-2 w-2 rounded-full ${connected ? "bg-emerald-500" : "bg-slate-300"}`}
              aria-hidden
            />
            {!collapsed ? (connected ? `${online.length} online` : "Offline") : null}
          </div>
          <ThemeToggle />
        </div>
      </aside>

      <div className="flex flex-1 flex-col">
        <header className="flex items-center justify-between border-b border-slate-200 bg-white px-4 py-3 md:hidden">
          <Link href="/feed" className="flex items-center gap-2 text-lg font-semibold text-slate-900">
            <Image
              src="/cyclops.png"
              alt=""
              width={28}
              height={28}
              className="h-7 w-7"
            />
            Femus
          </Link>
          <div className="flex items-center gap-3">
            <ThemeToggle />
            {user ? (
              <>
                <Link href="/profile">
                  <Avatar
                    src={user.avatar_path}
                    firstName={user.first_name}
                    lastName={user.last_name}
                    nickname={user.nickname}
                    size={32}
                  />
                </Link>
                <button
                  onClick={onLogout}
                  aria-label="Sign out"
                  title="Sign out"
                  className="rounded-md p-1.5 text-slate-500 transition hover:bg-slate-100 hover:text-red-600"
                >
                  ➜]
                </button>
              </>
            ) : null}
          </div>
        </header>

        <main className="flex-1">
          <div className="mx-auto max-w-3xl px-4 py-6">{children}</div>
        </main>

        <nav className="sticky bottom-0 grid grid-cols-6 border-t border-slate-200 bg-white md:hidden">
          {NAV.map((item) => {
            const active = pathname === item.href || pathname.startsWith(item.href + "/");
            let badge = 0;
            if (item.href === "/notifications") badge = unread;
            else if (item.href === "/messages") badge = unreadMessages;
            return (
              <Link
                key={item.href}
                href={item.href}
                className={`flex flex-col items-center py-2 text-xs ${
                  active ? "text-indigo-600" : "text-slate-600"
                }`}
              >
                <span aria-hidden className="relative inline-flex">
                  {item.icon}
                  {badge > 0 ? (
                    <span
                      aria-label={`${badge} unread`}
                      className="absolute -right-2 -top-1 inline-flex min-w-[16px] items-center justify-center rounded-full bg-red-600 px-1 text-[10px] font-semibold leading-none text-white ring-2 ring-white"
                    >
                      {badge > 99 ? "99+" : badge}
                    </span>
                  ) : null}
                </span>
                <span className="mt-0.5">{item.label}</span>
              </Link>
            );
          })}
        </nav>
      </div>
    </div>
  );
}
