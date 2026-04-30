"use client";

import { useEffect, useState } from "react";
import { FormMessage } from "@/components/forms";
import { api } from "@/lib/api";
import { formatDate } from "@/lib/format";
import type { Notification } from "@/types/api";

export default function NotificationsPage() {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [message, setMessage] = useState("");
  const [tone, setTone] = useState<"error" | "success">("success");

  async function loadNotifications() {
    try {
      const response = await api.listNotifications();
      setNotifications(response.notifications ?? []);
    } catch (error) {
      setTone("error");
      setMessage(
        error instanceof Error ? error.message : "Could not load notifications",
      );
    }
  }

  useEffect(() => {
    void loadNotifications();
  }, []);

  async function markRead(id: number) {
    try {
      await api.markNotificationRead(id);
      setTone("success");
      setMessage("Notification marked read.");
      await loadNotifications();
    } catch (error) {
      setTone("error");
      setMessage(error instanceof Error ? error.message : "Action failed");
    }
  }

  async function markAllRead() {
    try {
      await api.markAllNotificationsRead();
      setTone("success");
      setMessage("All notifications marked read.");
      await loadNotifications();
    } catch (error) {
      setTone("error");
      setMessage(error instanceof Error ? error.message : "Action failed");
    }
  }

  return (
    <section>
      <div className="mb-5 flex flex-wrap items-center justify-between gap-3">
        <h1 className="text-3xl font-bold text-slate-950">Notifications</h1>
        <div className="flex gap-2">
          <button
            className="rounded-xl border border-slate-300 px-4 py-2 text-sm font-semibold text-slate-700"
            onClick={loadNotifications}
            type="button"
          >
            Refresh
          </button>
          <button
            className="rounded-xl bg-sky-600 px-4 py-2 text-sm font-semibold text-white"
            onClick={markAllRead}
            type="button"
          >
            Mark all read
          </button>
        </div>
      </div>
      <FormMessage message={message} tone={tone} />
      <div className="mt-6 grid gap-4">
        {notifications.map((notification) => (
          <article
            className={`rounded-2xl border p-5 shadow-sm ${
              notification.read_at
                ? "border-slate-200 bg-white"
                : "border-sky-200 bg-sky-50"
            }`}
            key={notification.id}
          >
            <div className="flex flex-wrap items-start justify-between gap-4">
              <div>
                <p className="text-sm font-semibold uppercase tracking-wide text-slate-500">
                  {notification.type}
                </p>
                <h2 className="mt-2 text-xl font-bold text-slate-950">
                  {notification.title}
                </h2>
                <p className="mt-2 text-slate-700">{notification.body}</p>
                <p className="mt-3 text-sm text-slate-500">
                  {formatDate(notification.created_at)}
                </p>
              </div>
              {!notification.read_at ? (
                <button
                  className="rounded-xl border border-slate-300 px-4 py-2 text-sm font-semibold text-slate-700"
                  onClick={() => markRead(notification.id)}
                  type="button"
                >
                  Mark read
                </button>
              ) : null}
            </div>
          </article>
        ))}
        {notifications.length === 0 ? (
          <p className="rounded-2xl border border-dashed border-slate-300 bg-white p-8 text-center text-slate-600">
            No notifications yet.
          </p>
        ) : null}
      </div>
    </section>
  );
}
