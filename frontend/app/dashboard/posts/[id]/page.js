"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import {
  createComment,
  deleteComment,
  deletePost,
  getPost,
  getStoredToken,
  updateComment,
  updatePost,
} from "../../../../lib/api";

const PRIVACY = [
  { value: "public", label: "Public" },
  { value: "almost_private", label: "Almost private" },
  { value: "private", label: "Private" },
];

export default function PostDetailPage() {
  const params = useParams();
  const router = useRouter();
  const postId = Number(params.id);
  const [post, setPost] = useState(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);
  const [editMode, setEditMode] = useState(false);
  const [editForm, setEditForm] = useState({});
  const [commentText, setCommentText] = useState("");
  const [localComments, setLocalComments] = useState([]);
  const [currentUser, setCurrentUser] = useState(null);

  const token = getStoredToken();

  useEffect(() => {
    const raw = localStorage.getItem("currentUser");
    if (raw) {
      try {
        setCurrentUser(JSON.parse(raw));
      } catch {
        /* ignore */
      }
    }
  }, []);

  const load = async () => {
    if (!Number.isFinite(postId) || postId <= 0) {
      setError("Invalid post.");
      setLoading(false);
      return;
    }
    setLoading(true);
    setError("");
    try {
      const res = await getPost(token, postId);
      setPost(res.post);
      setEditForm({
        title: res.post.title,
        content: res.post.content,
        category: res.post.category,
        privacy_level: res.post.privacy_level,
        media_url: res.post.media_url || "",
      });
    } catch (e) {
      setError(e.message);
      setPost(null);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, [postId]);

  const savePost = async (e) => {
    e.preventDefault();
    try {
      await updatePost(token, postId, {
        title: editForm.title,
        content: editForm.content,
        category: editForm.category,
        privacy_level: editForm.privacy_level,
        ...(editForm.media_url?.trim()
          ? { media_url: editForm.media_url.trim() }
          : {}),
      });
      setEditMode(false);
      await load();
    } catch (e) {
      setError(e.message);
    }
  };

  const handleDeletePost = async () => {
    if (!confirm("Delete this post?")) return;
    try {
      await deletePost(token, postId);
      router.push("/dashboard/feed");
    } catch (e) {
      setError(e.message);
    }
  };

  const submitComment = async (e) => {
    e.preventDefault();
    if (!commentText.trim()) return;
    try {
      const res = await createComment(token, postId, { content: commentText.trim() });
      setCommentText("");
      if (res.comment) {
        setLocalComments((prev) => [...prev, res.comment]);
      }
    } catch (e) {
      setError(e.message);
    }
  };

  const isOwner = currentUser && post && currentUser.id === post.user_id;

  if (loading) {
    return (
      <div className="page-stack">
        <p className="helper-text">Loading post…</p>
      </div>
    );
  }

  if (!post) {
    return (
      <div className="page-stack">
        <p className="error-message">{error || "Post not found."}</p>
        <Link href="/dashboard/feed">Back to feed</Link>
      </div>
    );
  }

  return (
    <div className="page-stack">
      <nav className="breadcrumb">
        <Link href="/dashboard/feed">Feed</Link>
        <span aria-hidden="true"> / </span>
        <span>Post #{post.id}</span>
      </nav>

      <section className="surface-card">
        {error ? <p className="error-message">{error}</p> : null}

        {!editMode ? (
          <>
            <h1>{post.title}</h1>
            <p className="post-meta">
              by {post.author} · {post.category} · {post.privacy_level} ·{" "}
              {new Date(post.created_at).toLocaleString()}
            </p>
            <div className="post-body">{post.content}</div>
            {post.media_url ? (
              <p>
                <a href={post.media_url} target="_blank" rel="noreferrer">
                  Media link
                </a>
              </p>
            ) : null}
            {isOwner ? (
              <div className="button-row">
                <button type="button" onClick={() => setEditMode(true)}>
                  Edit
                </button>
                <button type="button" className="danger-outline" onClick={handleDeletePost}>
                  Delete
                </button>
              </div>
            ) : null}
          </>
        ) : (
          <form className="stack-form" onSubmit={savePost}>
            <label>Title</label>
            <input
              value={editForm.title}
              onChange={(e) => setEditForm((f) => ({ ...f, title: e.target.value }))}
            />
            <label>Content</label>
            <textarea
              rows={5}
              value={editForm.content}
              onChange={(e) => setEditForm((f) => ({ ...f, content: e.target.value }))}
            />
            <label>Category</label>
            <input
              value={editForm.category}
              onChange={(e) => setEditForm((f) => ({ ...f, category: e.target.value }))}
            />
            <label>Privacy</label>
            <select
              value={editForm.privacy_level}
              onChange={(e) =>
                setEditForm((f) => ({ ...f, privacy_level: e.target.value }))
              }
            >
              {PRIVACY.map((p) => (
                <option key={p.value} value={p.value}>
                  {p.label}
                </option>
              ))}
            </select>
            <label>Media URL</label>
            <input
              value={editForm.media_url}
              onChange={(e) => setEditForm((f) => ({ ...f, media_url: e.target.value }))}
            />
            <div className="button-row">
              <button type="submit">Save</button>
              <button type="button" className="button-secondary" onClick={() => setEditMode(false)}>
                Cancel
              </button>
            </div>
          </form>
        )}
      </section>

      <section className="surface-card">
        <h2>Comments ({localComments.length} this session)</h2>
        <p className="helper-text small">
          There is no GET comments API; only comments you add in this session are listed below for convenience.
        </p>
        <ul className="comment-list">
          {localComments.map((c) => (
            <li key={c.id} className="comment-item">
              <strong>{c.author}</strong> · {new Date(c.created_at).toLocaleString()}
              <p>{c.content}</p>
              {currentUser && c.user_id === currentUser.id ? (
                <CommentActions
                  token={token}
                  postId={postId}
                  comment={c}
                  onUpdated={(updated) =>
                    setLocalComments((prev) =>
                      prev.map((x) => (x.id === updated.id ? updated : x)),
                    )
                  }
                  onDeleted={(id) =>
                    setLocalComments((prev) => prev.filter((x) => x.id !== id))
                  }
                />
              ) : null}
            </li>
          ))}
        </ul>

        <form className="stack-form" onSubmit={submitComment}>
          <label htmlFor="cnew">Add comment</label>
          <textarea
            id="cnew"
            rows={2}
            value={commentText}
            onChange={(e) => setCommentText(e.target.value)}
            placeholder="Write a comment…"
          />
          <button type="submit">Post comment</button>
        </form>
      </section>
    </div>
  );
}

function CommentActions({ token, postId, comment, onUpdated, onDeleted }) {
  const [editing, setEditing] = useState(false);
  const [text, setText] = useState(comment.content);

  const save = async () => {
    const res = await updateComment(token, postId, comment.id, { content: text });
    if (res.comment) onUpdated(res.comment);
    setEditing(false);
  };

  const remove = async () => {
    if (!confirm("Delete this comment?")) return;
    await deleteComment(token, postId, comment.id);
    onDeleted(comment.id);
  };

  if (editing) {
    return (
      <div className="comment-actions">
        <textarea rows={2} value={text} onChange={(e) => setText(e.target.value)} />
        <button type="button" onClick={save}>
          Save
        </button>
        <button type="button" className="button-text" onClick={() => setEditing(false)}>
          Cancel
        </button>
      </div>
    );
  }

  return (
    <div className="comment-actions">
      <button type="button" className="button-text" onClick={() => setEditing(true)}>
        Edit
      </button>
      <button type="button" className="button-text danger-text" onClick={remove}>
        Delete
      </button>
    </div>
  );
}
