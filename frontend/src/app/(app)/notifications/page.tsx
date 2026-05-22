"use client";

import Link from "next/link";
import { useCallback, useEffect, useState } from "react";
import { useWebSocket } from "@/components/WebSocketProvider";
import { useNotificationCount } from "@/components/NotificationCountProvider";
import { InfiniteScrollSentinel } from "@/components/InfiniteScrollSentinel";
import { usePaginatedList } from "@/hooks/usePaginatedList";
import { api, ApiError } from "@/lib/api";
import { formatRelative } from "@/lib/format";
import type { Notification } from "@/types/api";

const PAGE_SIZE = 20;

export default function NotificationsPage() {
  const { on } = useWebSocket();
  const { decrement, reset: resetUnreadCount } = useNotificationCount();
  const [error, setError] = useState<string | null>(null);

  const fetcher = useCallback(
    async ({ limit, offset }: { limit: number; offset: number }) => {
      const res = await api.notifications.list({ limit, offset });
      return { items: res.notifications ?? [], hasMore: res.has_more };
    },
    [],
  );
  const { items, hasMore, loading, error: fetchError, loadMore, setItems } =
    usePaginatedList<Notification>(fetcher, PAGE_SIZE, (n) => n.id);

  useEffect(() => {
    return on((msg) => {
      if (msg.type === "notification.created") {
        const payload = msg.payload as { notification: Notification };
        setItems((prev) => [payload.notification, ...prev]);
      }
    });
  }, [on, setItems]);

  const markRead = async (id: number) => {
    const target = items.find((n) => n.id === id);
    if (!target || target.read_at) return;
    try {
      await api.notifications.markRead(id);
      setItems((prev) =>
        prev.map((n) => (n.id === id ? { ...n, read_at: new Date().toISOString() } : n)),
      );
      decrement();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to mark as read");
    }
  };

  const markAll = async () => {
    try {
      await api.notifications.markAllRead();
      const now = new Date().toISOString();
      setItems((prev) => prev.map((n) => ({ ...n, read_at: n.read_at ?? now })));
      resetUnreadCount();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to mark all as read");
    }
  };

  const actOnGroup = async (
    notification: Notification,
    kind: "accept" | "decline",
  ) => {
    if (notification.entity_id == null) return;
    try {
      if (notification.entity_type === "group_invitation") {
        if (kind === "accept") await api.groups.acceptInvite(notification.entity_id);
        else await api.groups.declineInvite(notification.entity_id);
      } else if (notification.entity_type === "group_join_request") {
        if (kind === "accept") await api.groups.acceptJoin(notification.entity_id);
        else await api.groups.declineJoin(notification.entity_id);
      } else {
        return;
      }
      try {
        await api.notifications.markRead(notification.id);
      } catch {
        // best-effort; if it was already read this errors silently
      }
      setItems((prev) =>
        prev.map((n) =>
          n.id === notification.id ? { ...n, read_at: new Date().toISOString() } : n,
        ),
      );
      if (!notification.read_at) decrement();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : `Failed to ${kind}`);
    }
  };

  const unread = items.filter((n) => !n.read_at).length;

  const renderBody = (n: Notification) => {
    if (n.entity_type === "follow_request") {
      const actor = parseActorName(n.metadata);
      if (actor) return `${actor} requested to follow you.`;
    }
    return n.body;
  };

  return (
    <div className="space-y-6">
      <header className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold text-slate-900">
          Notifications
          {unread > 0 ? (
            <span className="ml-2 rounded-full bg-indigo-600 px-2 py-0.5 align-middle text-xs font-medium text-white">
              {unread}
            </span>
          ) : null}
        </h1>
        {unread > 0 ? (
          <button
            onClick={markAll}
            className="rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 hover:bg-slate-50"
          >
            Mark all as read
          </button>
        ) : null}
      </header>

      {(error || fetchError) ? (
        <p className="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{error ?? fetchError}</p>
      ) : null}

      {loading && items.length === 0 ? (
        <p className="text-sm text-slate-500">Loading…</p>
      ) : items.length === 0 ? (
        <p className="rounded-xl border border-dashed border-slate-300 bg-white p-6 text-center text-sm text-slate-500">
          No notifications yet.
        </p>
      ) : (
        <ul className="space-y-2">
          {items.map((n) => {
            const actionable =
              !n.read_at &&
              (n.entity_type === "group_invitation" || n.entity_type === "group_join_request");
            const linkable = !actionable && !!n.action_url;
            const cardClass = `flex items-start justify-between gap-3 rounded-xl border p-4 shadow-sm transition ${
              n.read_at
                ? "border-slate-200 bg-white"
                : "border-indigo-200 bg-indigo-50/50"
            } ${linkable ? "hover:bg-slate-50" : ""}`;
            const inner = (
              <>
                <div className="min-w-0 flex-1">
                  <p className="text-sm font-medium text-slate-900">{n.title}</p>
                  <p className="mt-1 text-sm text-slate-600">{renderBody(n)}</p>
                  <p className="mt-1 text-xs text-slate-400">{formatRelative(n.created_at)}</p>
                </div>
                <div className="flex shrink-0 flex-col items-end gap-2">
                  {actionable ? (
                    <div className="flex gap-2">
                      <button
                        onClick={() => actOnGroup(n, "accept")}
                        className="rounded-lg bg-indigo-600 px-3 py-1 text-xs font-medium text-white shadow-sm hover:bg-indigo-700"
                      >
                        Accept
                      </button>
                      <button
                        onClick={() => actOnGroup(n, "decline")}
                        className="rounded-lg border border-slate-300 bg-white px-3 py-1 text-xs font-medium text-slate-700 hover:bg-slate-50"
                      >
                        Decline
                      </button>
                    </div>
                  ) : null}
                  {!n.read_at && !actionable && !linkable ? (
                    <button
                      onClick={() => markRead(n.id)}
                      className="rounded-lg border border-slate-300 bg-white px-2.5 py-1 text-xs font-medium text-slate-700 hover:bg-slate-50"
                    >
                      Mark read
                    </button>
                  ) : null}
                </div>
              </>
            );
            return linkable ? (
              <li key={n.id}>
                <Link
                  href={n.action_url!}
                  onClick={() => !n.read_at && markRead(n.id)}
                  className={cardClass}
                >
                  {inner}
                </Link>
              </li>
            ) : (
              <li key={n.id} className={cardClass}>
                {inner}
              </li>
            );
          })}
        </ul>
      )}

      {items.length > 0 && hasMore ? (
        <InfiniteScrollSentinel onIntersect={loadMore} enabled={!loading} />
      ) : null}
      {loading && items.length > 0 ? (
        <p className="text-center text-xs text-slate-400">Loading more…</p>
      ) : null}
    </div>
  );
}

function parseActorName(metadata: string | null | undefined): string | null {
  if (!metadata) return null;
  try {
    const parsed = JSON.parse(metadata) as { actor_name?: string };
    return parsed.actor_name?.trim() || null;
  } catch {
    return null;
  }
}
