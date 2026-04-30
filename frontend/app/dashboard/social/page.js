"use client";

import { useState } from "react";
import {
  acceptFollowRequest,
  declineFollowRequest,
  followUser,
  getStoredToken,
  removeFollower,
  unfollowUser,
} from "../../../lib/api";

export default function SocialPage() {
  const token = getStoredToken();
  const [userId, setUserId] = useState("");
  const [followId, setFollowId] = useState("");
  const [followerId, setFollowerId] = useState("");
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);

  const run = async (fn) => {
    setBusy(true);
    setError("");
    setMessage("");
    try {
      const res = await fn();
      setMessage(res.message || "Done.");
      if (res.status) setMessage((m) => `${m} (${res.status})`);
    } catch (e) {
      setError(e.message);
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="page-stack">
      <section className="surface-card">
        <h1>Follow & connections</h1>
        <p className="helper-text">
          Follow actions use numeric user ids. Accept/decline use the follow <strong>request id</strong> (often from a notification), not the user id.
        </p>

        <div className="form-grid">
          <h2 className="h2-inline">Follow / unfollow</h2>
          <label>Target user id</label>
          <input
            type="number"
            min={1}
            value={userId}
            onChange={(e) => setUserId(e.target.value)}
            placeholder="e.g. 2"
          />
          <div className="button-row">
            <button
              type="button"
              disabled={busy}
              onClick={() =>
                run(() => followUser(token, Number(userId)))
              }
            >
              Follow
            </button>
            <button
              type="button"
              className="button-secondary"
              disabled={busy}
              onClick={() =>
                run(() => unfollowUser(token, Number(userId)))
              }
            >
              Unfollow
            </button>
          </div>
        </div>

        <div className="form-grid">
          <h2 className="h2-inline">Incoming requests</h2>
          <label>Follow request id</label>
          <input
            type="number"
            min={1}
            value={followId}
            onChange={(e) => setFollowId(e.target.value)}
          />
          <div className="button-row">
            <button
              type="button"
              disabled={busy}
              onClick={() =>
                run(() => acceptFollowRequest(token, Number(followId)))
              }
            >
              Accept
            </button>
            <button
              type="button"
              className="button-secondary"
              disabled={busy}
              onClick={() =>
                run(() => declineFollowRequest(token, Number(followId)))
              }
            >
              Decline
            </button>
          </div>
        </div>

        <div className="form-grid">
          <h2 className="h2-inline">Remove a follower</h2>
          <label>Follower user id</label>
          <input
            type="number"
            min={1}
            value={followerId}
            onChange={(e) => setFollowerId(e.target.value)}
          />
          <button
            type="button"
            disabled={busy}
            onClick={() =>
              run(() => removeFollower(token, Number(followerId)))
            }
          >
            Remove follower
          </button>
        </div>

        {error ? <p className="error-message">{error}</p> : null}
        {message ? <p className="success-message">{message}</p> : null}
      </section>
    </div>
  );
}
