"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import {
  acceptGroupInvitation,
  acceptGroupJoinRequest,
  createGroup,
  declineGroupInvitation,
  declineGroupJoinRequest,
  inviteToGroup,
  joinGroup,
  listGroups,
  getStoredToken,
} from "../../../lib/api";

export default function GroupsPage() {
  const token = getStoredToken();
  const [groups, setGroups] = useState([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);
  const [createForm, setCreateForm] = useState({ title: "", description: "" });
  const [joinId, setJoinId] = useState("");
  const [inviteForm, setInviteForm] = useState({ groupId: "", inviteeId: "" });
  const [invitationId, setInvitationId] = useState("");
  const [joinRequestId, setJoinRequestId] = useState("");
  const [msg, setMsg] = useState("");

  const load = async () => {
    setLoading(true);
    setError("");
    try {
      const res = await listGroups(token, { limit: 50, offset: 0 });
      setGroups(res.groups ?? []);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const handleCreate = async (e) => {
    e.preventDefault();
    setMsg("");
    setError("");
    try {
      await createGroup(token, {
        title: createForm.title,
        description: createForm.description,
      });
      setCreateForm({ title: "", description: "" });
      setMsg("Group created.");
      await load();
    } catch (e) {
      setError(e.message);
    }
  };

  const handleJoin = async (e) => {
    e.preventDefault();
    setMsg("");
    setError("");
    try {
      await joinGroup(token, Number(joinId));
      setMsg("Join request sent (if applicable).");
    } catch (e) {
      setError(e.message);
    }
  };

  const handleInvite = async (e) => {
    e.preventDefault();
    setMsg("");
    setError("");
    try {
      await inviteToGroup(token, {
        groupId: Number(inviteForm.groupId),
        inviteeId: Number(inviteForm.inviteeId),
      });
      setMsg("Invitation sent.");
    } catch (e) {
      setError(e.message);
    }
  };

  return (
    <div className="page-stack">
      <section className="surface-card">
        <h1>Groups</h1>
        <p className="helper-text">
          Create groups, request to join, invite others, and resolve invitations or join requests by id (from notifications or admins).
        </p>

        <form className="stack-form" onSubmit={handleCreate}>
          <h2 className="h2-inline">Create group</h2>
          <label>Title</label>
          <input
            value={createForm.title}
            onChange={(e) =>
              setCreateForm((f) => ({ ...f, title: e.target.value }))
            }
            required
          />
          <label>Description</label>
          <textarea
            rows={2}
            value={createForm.description}
            onChange={(e) =>
              setCreateForm((f) => ({ ...f, description: e.target.value }))
            }
            required
          />
          <button type="submit">Create</button>
        </form>

        <form className="stack-form" onSubmit={handleJoin}>
          <h2 className="h2-inline">Request to join</h2>
          <label>Group id</label>
          <input
            type="number"
            min={1}
            value={joinId}
            onChange={(e) => setJoinId(e.target.value)}
            required
          />
          <button type="submit">Submit join request</button>
        </form>

        <form className="stack-form" onSubmit={handleInvite}>
          <h2 className="h2-inline">Invite user</h2>
          <label>Group id</label>
          <input
            type="number"
            min={1}
            value={inviteForm.groupId}
            onChange={(e) =>
              setInviteForm((f) => ({ ...f, groupId: e.target.value }))
            }
            required
          />
          <label>Invitee user id</label>
          <input
            type="number"
            min={1}
            value={inviteForm.inviteeId}
            onChange={(e) =>
              setInviteForm((f) => ({ ...f, inviteeId: e.target.value }))
            }
            required
          />
          <button type="submit">Send invitation</button>
        </form>

        <div className="form-grid">
          <h2 className="h2-inline">Invitation id</h2>
          <input
            type="number"
            min={1}
            value={invitationId}
            onChange={(e) => setInvitationId(e.target.value)}
            placeholder="from notification"
          />
          <div className="button-row">
            <button
              type="button"
              onClick={async () => {
                setError("");
                try {
                  await acceptGroupInvitation(token, Number(invitationId));
                  setMsg("Invitation accepted.");
                  await load();
                } catch (e) {
                  setError(e.message);
                }
              }}
            >
              Accept invite
            </button>
            <button
              type="button"
              className="button-secondary"
              onClick={async () => {
                setError("");
                try {
                  await declineGroupInvitation(token, Number(invitationId));
                  setMsg("Invitation declined.");
                } catch (e) {
                  setError(e.message);
                }
              }}
            >
              Decline invite
            </button>
          </div>
        </div>

        <div className="form-grid">
          <h2 className="h2-inline">Join request id (for moderators)</h2>
          <input
            type="number"
            min={1}
            value={joinRequestId}
            onChange={(e) => setJoinRequestId(e.target.value)}
          />
          <div className="button-row">
            <button
              type="button"
              onClick={async () => {
                setError("");
                try {
                  await acceptGroupJoinRequest(token, Number(joinRequestId));
                  setMsg("Join request accepted.");
                } catch (e) {
                  setError(e.message);
                }
              }}
            >
              Accept join
            </button>
            <button
              type="button"
              className="button-secondary"
              onClick={async () => {
                setError("");
                try {
                  await declineGroupJoinRequest(token, Number(joinRequestId));
                  setMsg("Join request declined.");
                } catch (e) {
                  setError(e.message);
                }
              }}
            >
              Decline join
            </button>
          </div>
        </div>

        {error ? <p className="error-message">{error}</p> : null}
        {msg ? <p className="success-message">{msg}</p> : null}
      </section>

      <section className="surface-card">
        <div className="toolbar">
          <h2 className="h2-inline">All groups</h2>
          <button type="button" className="button-text" onClick={load}>
            Refresh
          </button>
        </div>
        {loading ? (
          <p className="helper-text">Loading…</p>
        ) : groups.length === 0 ? (
          <p className="helper-text">No groups yet.</p>
        ) : (
          <ul className="post-list">
            {groups.map((g) => (
              <li key={g.id} className="post-item">
                <Link href={`/dashboard/groups/${g.id}`} className="post-title-link">
                  {g.name}
                </Link>
                <p className="post-meta">
                  id {g.id} · creator {g.creator_id}
                </p>
                <p className="post-excerpt">{g.description}</p>
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}
