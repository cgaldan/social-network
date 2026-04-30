"use client";

import { useEffect, useState } from "react";
import {
  getStoredToken,
  listNotifications,
  markAllNotificationsRead,
  markNotificationRead,
} from "../../../lib/api";

export default function NotificationsPage() {
  const token = getStoredToken();
  const [items, setItems] = useState([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  const load = async () => {
    setLoading(true);
    setError("");
    try {
      const res = await listNotifications(token, { limit: 50, offset: 0 });
      setItems(res.notifications ?? []);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const markOne = async (id) => {
    setError("");
    try {
      await markNotificationRead(token, id);
      await load();
    } catch (e) {
      setError(e.message);
    }
  };

  const markAll = async () => {
    setError("");
    try {
      await markAllNotificationsRead(token);
      await load();
    } catch (e) {
      setError(e.message);
    }
  };

  return (
    <div className="page-stack">
      <section className="surface-card">
        <div className="toolbar">
          <h1>Notifications</h1>
          <div className="button-row">
            <button type="button" className="button-text" onClick={load}>
              Refresh
            </button>
            <button type="button" className="button-secondary" onClick={markAll}>
              Mark all read
            </button>
          </div>
        </div>
        {error ? <p className="error-message">{error}</p> : null}
        {loading ? (
          <p className="helper-text">Loading…</p>
        ) : items.length === 0 ? (
          <p className="helper-text">No notifications.</p>
        ) : (
          <ul className="notif-list">
            {items.map((n) => (
              <li
                key={n.id}
                className={
                  n.read_at ? "notif-item notif-read" : "notif-item notif-unread"
                }
              >
                <div className="notif-head">
                  <strong>{n.title}</strong>
                  <span className="notif-time">
                    {new Date(n.created_at).toLocaleString()}
                  </span>
                </div>
                <p className="notif-body">{n.body}</p>
                <p className="notif-meta">
                  {n.type}
                  {n.entity_type ? ` · ${n.entity_type} #${n.entity_id}` : ""}
                </p>
                {!n.read_at ? (
                  <button
                    type="button"
                    className="button-text"
                    onClick={() => markOne(n.id)}
                  >
                    Mark read
                  </button>
                ) : null}
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}
