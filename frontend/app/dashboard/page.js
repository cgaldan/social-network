"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { getCurrentUser, logout } from "../../lib/api";

export default function DashboardPage() {
  const router = useRouter();
  const [currentUser, setCurrentUser] = useState(null);
  const [loggingOut, setLoggingOut] = useState(false);

  useEffect(() => {
    const raw = localStorage.getItem("currentUser");
    if (raw) {
      try {
        setCurrentUser(JSON.parse(raw));
      } catch {
        localStorage.removeItem("currentUser");
      }
    }
    const token = localStorage.getItem("authToken");
    if (token) {
      getCurrentUser(token)
        .then((res) => {
          if (res.user) {
            setCurrentUser(res.user);
            localStorage.setItem("currentUser", JSON.stringify(res.user));
          }
        })
        .catch(() => {});
    }
  }, []);

  const handleLogout = async () => {
    const token = localStorage.getItem("authToken");
    setLoggingOut(true);
    try {
      if (token) await logout(token);
    } catch {
      // still clear local session
    } finally {
      localStorage.removeItem("authToken");
      localStorage.removeItem("currentUser");
      router.replace("/login");
    }
  };

  return (
    <div className="page-stack">
      <section className="surface-card">
        <h1>Welcome, {currentUser?.nickname || "user"}</h1>
        <p className="helper-text">
          Your id: <strong>{currentUser?.id ?? "—"}</strong> — share it for invites, DMs, and follows.
        </p>
        <div className="dashboard-details">
          <p>
            <strong>Email:</strong> {currentUser?.email || "—"}
          </p>
          <p>
            <strong>Name:</strong>{" "}
            {[currentUser?.first_name, currentUser?.last_name].filter(Boolean).join(" ") || "—"}
          </p>
          <p>
            <strong>Followers / following:</strong> {currentUser?.followers_count ?? 0} /{" "}
            {currentUser?.following_count ?? 0}
          </p>
        </div>
        <div className="button-row">
          <Link className="button button-secondary" href="/dashboard/feed">
            Open feed
          </Link>
          <Link className="button button-secondary" href="/dashboard/notifications">
            Notifications
          </Link>
        </div>
        <button
          className="logout-button"
          type="button"
          onClick={handleLogout}
          disabled={loggingOut}
        >
          {loggingOut ? "Logging out…" : "Log out"}
        </button>
      </section>
    </div>
  );
}
