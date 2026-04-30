"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import {
  createGroupEvent,
  createGroupPost,
  getGroupPosts,
  listGroupEvents,
  setGroupEventRSVP,
  getStoredToken,
} from "../../../../lib/api";

const PRIVACY = [
  { value: "public", label: "Public" },
  { value: "almost_private", label: "Almost private" },
  { value: "private", label: "Private" },
];

export default function GroupDetailPage() {
  const params = useParams();
  const groupId = Number(params.id);
  const token = getStoredToken();
  const [posts, setPosts] = useState([]);
  const [events, setEvents] = useState([]);
  const [error, setError] = useState("");
  const [postForm, setPostForm] = useState({
    title: "",
    content: "",
    category: "general",
    privacy_level: "public",
  });
  const [eventForm, setEventForm] = useState({
    title: "",
    description: "",
    starts_at: "",
  });

  const load = async () => {
    setError("");
    try {
      const [pr, er] = await Promise.all([
        getGroupPosts(token, groupId, { limit: 40, offset: 0 }),
        listGroupEvents(token, groupId, { limit: 40, offset: 0 }),
      ]);
      setPosts(pr.posts ?? []);
      setEvents(er.events ?? []);
    } catch (e) {
      setError(e.message);
    }
  };

  useEffect(() => {
    if (Number.isFinite(groupId) && groupId > 0) load();
  }, [groupId]);

  const submitPost = async (e) => {
    e.preventDefault();
    setError("");
    try {
      await createGroupPost(token, groupId, postForm);
      setPostForm((f) => ({ ...f, title: "", content: "" }));
      await load();
    } catch (e) {
      setError(e.message);
    }
  };

  const submitEvent = async (e) => {
    e.preventDefault();
    setError("");
    try {
      await createGroupEvent(token, groupId, {
        title: eventForm.title,
        description: eventForm.description,
        starts_at: new Date(eventForm.starts_at).toISOString(),
      });
      setEventForm({ title: "", description: "", starts_at: "" });
      await load();
    } catch (e) {
      setError(e.message);
    }
  };

  const rsvp = async (eventId, response) => {
    setError("");
    try {
      await setGroupEventRSVP(token, groupId, eventId, response);
      await load();
    } catch (e) {
      setError(e.message);
    }
  };

  return (
    <div className="page-stack">
      <nav className="breadcrumb">
        <Link href="/dashboard/groups">Groups</Link>
        <span aria-hidden="true"> / </span>
        <span>Group #{groupId}</span>
      </nav>

      {error ? <p className="error-message">{error}</p> : null}

      <section className="surface-card">
        <h1>Group posts</h1>
        <form className="stack-form" onSubmit={submitPost}>
          <label>Title</label>
          <input
            value={postForm.title}
            onChange={(e) => setPostForm((f) => ({ ...f, title: e.target.value }))}
            required
          />
          <label>Content</label>
          <textarea
            rows={3}
            value={postForm.content}
            onChange={(e) => setPostForm((f) => ({ ...f, content: e.target.value }))}
            required
          />
          <label>Category</label>
          <input
            value={postForm.category}
            onChange={(e) => setPostForm((f) => ({ ...f, category: e.target.value }))}
          />
          <label>Privacy</label>
          <select
            value={postForm.privacy_level}
            onChange={(e) =>
              setPostForm((f) => ({ ...f, privacy_level: e.target.value }))
            }
          >
            {PRIVACY.map((p) => (
              <option key={p.value} value={p.value}>
                {p.label}
              </option>
            ))}
          </select>
          <button type="submit">Post to group</button>
        </form>

        <ul className="post-list">
          {posts.map((p) => (
            <li key={p.id} className="post-item">
              <Link href={`/dashboard/posts/${p.id}`} className="post-title-link">
                {p.title}
              </Link>
              <p className="post-meta">
                {p.author} · {new Date(p.created_at).toLocaleString()}
              </p>
              <p className="post-excerpt">{p.content}</p>
            </li>
          ))}
        </ul>
      </section>

      <section className="surface-card">
        <h1>Events</h1>
        <form className="stack-form" onSubmit={submitEvent}>
          <label>Title</label>
          <input
            value={eventForm.title}
            onChange={(e) => setEventForm((f) => ({ ...f, title: e.target.value }))}
            required
          />
          <label>Description</label>
          <textarea
            rows={2}
            value={eventForm.description}
            onChange={(e) =>
              setEventForm((f) => ({ ...f, description: e.target.value }))
            }
            required
          />
          <label>Starts at</label>
          <input
            type="datetime-local"
            value={eventForm.starts_at}
            onChange={(e) =>
              setEventForm((f) => ({ ...f, starts_at: e.target.value }))
            }
            required
          />
          <button type="submit">Create event</button>
        </form>

        <ul className="post-list">
          {events.map((ev) => (
            <li key={ev.id} className="post-item">
              <strong>{ev.title}</strong>
              <p className="post-meta">
                {new Date(ev.starts_at).toLocaleString()}
              </p>
              <p className="post-excerpt">{ev.description}</p>
              <div className="button-row">
                <button type="button" onClick={() => rsvp(ev.id, "going")}>
                  Going
                </button>
                <button
                  type="button"
                  className="button-secondary"
                  onClick={() => rsvp(ev.id, "not_going")}
                >
                  Not going
                </button>
              </div>
            </li>
          ))}
        </ul>
      </section>
    </div>
  );
}
