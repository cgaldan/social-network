"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import {
  createPost,
  getPosts,
  getStoredToken,
} from "../../../lib/api";

const PRIVACY = [
  { value: "public", label: "Public" },
  { value: "almost_private", label: "Almost private" },
  { value: "private", label: "Private" },
];

export default function FeedPage() {
  const [posts, setPosts] = useState([]);
  const [category, setCategory] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);
  const [form, setForm] = useState({
    title: "",
    content: "",
    category: "general",
    privacy_level: "public",
    media_url: "",
  });
  const [submitting, setSubmitting] = useState(false);

  const load = async () => {
    setLoading(true);
    setError("");
    try {
      const res = await getPosts({
        category: category.trim(),
        limit: 30,
        offset: 0,
      });
      setPosts(res.posts ?? []);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, [category]);

  const handleCreate = async (e) => {
    e.preventDefault();
    const token = getStoredToken();
    if (!token) return;
    setSubmitting(true);
    setError("");
    try {
      await createPost(token, {
        title: form.title,
        content: form.content,
        category: form.category,
        privacy_level: form.privacy_level,
        ...(form.media_url.trim() ? { media_url: form.media_url.trim() } : {}),
      });
      setForm((f) => ({
        ...f,
        title: "",
        content: "",
        media_url: "",
      }));
      await load();
    } catch (err) {
      setError(err.message);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="page-stack">
      <section className="surface-card">
        <h1>Feed</h1>
        <p className="helper-text">
          Public posts from the network. Open a post to comment (your new comments appear below after posting; the API does not expose a comment list yet).
        </p>

        <form className="stack-form" onSubmit={handleCreate}>
          <h2 className="h2-inline">New post</h2>
          <label htmlFor="ptitle">Title</label>
          <input
            id="ptitle"
            value={form.title}
            onChange={(e) => setForm((f) => ({ ...f, title: e.target.value }))}
            required
          />
          <label htmlFor="pcontent">Content</label>
          <textarea
            id="pcontent"
            rows={3}
            value={form.content}
            onChange={(e) => setForm((f) => ({ ...f, content: e.target.value }))}
            required
          />
          <label htmlFor="pcat">Category</label>
          <input
            id="pcat"
            value={form.category}
            onChange={(e) => setForm((f) => ({ ...f, category: e.target.value }))}
          />
          <label htmlFor="ppriv">Privacy</label>
          <select
            id="ppriv"
            value={form.privacy_level}
            onChange={(e) =>
              setForm((f) => ({ ...f, privacy_level: e.target.value }))
            }
          >
            {PRIVACY.map((p) => (
              <option key={p.value} value={p.value}>
                {p.label}
              </option>
            ))}
          </select>
          <label htmlFor="pmedia">Media URL (optional)</label>
          <input
            id="pmedia"
            value={form.media_url}
            onChange={(e) => setForm((f) => ({ ...f, media_url: e.target.value }))}
            placeholder="https://"
          />
          <button type="submit" disabled={submitting}>
            {submitting ? "Posting…" : "Publish"}
          </button>
        </form>
      </section>

      <section className="surface-card">
        <div className="toolbar">
          <label className="inline-label">
            Filter category{" "}
            <input
              className="narrow-input"
              value={category}
              onChange={(e) => setCategory(e.target.value)}
              placeholder="e.g. general"
            />
          </label>
          <button type="button" className="button-text" onClick={load}>
            Refresh
          </button>
        </div>
        {error ? <p className="error-message">{error}</p> : null}
        {loading ? (
          <p className="helper-text">Loading posts…</p>
        ) : posts.length === 0 ? (
          <p className="helper-text">No posts yet.</p>
        ) : (
          <ul className="post-list">
            {posts.map((p) => (
              <li key={p.id} className="post-item">
                <Link href={`/dashboard/posts/${p.id}`} className="post-title-link">
                  {p.title}
                </Link>
                <p className="post-meta">
                  by {p.author} · {p.category} · {p.privacy_level} ·{" "}
                  {new Date(p.created_at).toLocaleString()}
                </p>
                <p className="post-excerpt">{p.content}</p>
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}
