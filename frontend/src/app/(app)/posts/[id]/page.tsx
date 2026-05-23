"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { useAuth } from "@/components/AuthProvider";
import { PostCard } from "@/components/PostCard";
import { InfiniteScrollSentinel } from "@/components/InfiniteScrollSentinel";
import { usePaginatedList } from "@/hooks/usePaginatedList";
import { api, ApiError, resolveMediaUrl } from "@/lib/api";
import { formatRelative } from "@/lib/format";
import type { Comment, Post, PrivacyLevel } from "@/types/api";

const COMMENTS_PAGE_SIZE = 20;

export default function PostDetailPage() {
  const params = useParams<{ id: string }>();
  const id = Number(params.id);
  const router = useRouter();
  const { user } = useAuth();
  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [commentText, setCommentText] = useState("");
  const [commentMediaUrl, setCommentMediaUrl] = useState("");
  const [uploadingComment, setUploadingComment] = useState(false);
  const [sending, setSending] = useState(false);
  const commentFileRef = useRef<HTMLInputElement>(null);

  const [editingPost, setEditingPost] = useState(false);
  const [editTitle, setEditTitle] = useState("");
  const [editContent, setEditContent] = useState("");
  const [editCategory, setEditCategory] = useState("");
  const [editPrivacy, setEditPrivacy] = useState<PrivacyLevel>("public");
  const [editMediaUrl, setEditMediaUrl] = useState("");
  const [savingPost, setSavingPost] = useState(false);
  const [uploadingEdit, setUploadingEdit] = useState(false);
  const editFileRef = useRef<HTMLInputElement>(null);

  const [editingCommentId, setEditingCommentId] = useState<number | null>(null);
  const [editCommentText, setEditCommentText] = useState("");
  const [editCommentMediaUrl, setEditCommentMediaUrl] = useState("");
  const [savingComment, setSavingComment] = useState(false);
  const [uploadingEditComment, setUploadingEditComment] = useState(false);
  const editCommentFileRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (!Number.isFinite(id)) return;
    let cancelled = false;
    setLoading(true);
    api.posts
      .get(id)
      .then((res) => {
        if (!cancelled) setPost(res.post);
      })
      .catch((err) => {
        if (!cancelled) setError(err instanceof ApiError ? err.message : "Failed to load post");
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [id]);

  const commentsFetcher = useCallback(
    async ({ limit, offset }: { limit: number; offset: number }) => {
      if (!Number.isFinite(id)) return { items: [], hasMore: false };
      const res = await api.comments.list(id, { limit, offset });
      return { items: res.comments ?? [], hasMore: res.has_more };
    },
    [id],
  );
  const {
    items: comments,
    hasMore: hasMoreComments,
    loading: loadingComments,
    loadMore: loadMoreComments,
    setItems: setComments,
  } = usePaginatedList<Comment>(commentsFetcher, COMMENTS_PAGE_SIZE);

  const onCommentFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setError(null);
    setUploadingComment(true);
    try {
      const res = await api.uploads.create(file);
      setCommentMediaUrl(res.url);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Upload failed");
      if (commentFileRef.current) commentFileRef.current.value = "";
    } finally {
      setUploadingComment(false);
    }
  };

  const clearCommentImage = () => {
    setCommentMediaUrl("");
    if (commentFileRef.current) commentFileRef.current.value = "";
  };

  const onComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!commentText.trim()) return;
    setSending(true);
    try {
      const res = await api.comments.create(id, {
        content: commentText,
        media_url: commentMediaUrl || undefined,
      });
      setComments((prev: Comment[]) => [...prev, res.comment]);
      setCommentText("");
      clearCommentImage();
      if (post) setPost({ ...post, comment_count: post.comment_count + 1 });
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to send comment");
    } finally {
      setSending(false);
    }
  };

  const onDeletePost = async () => {
    if (!confirm("Delete this post?")) return;
    try {
      await api.posts.remove(id);
      router.replace("/feed");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to delete post");
    }
  };

  const onStartEditPost = () => {
    if (!post) return;
    setEditTitle(post.title);
    setEditContent(post.content);
    setEditCategory(post.category);
    setEditPrivacy(post.privacy_level);
    setEditMediaUrl(post.media_url ?? "");
    setEditingPost(true);
  };

  const onCancelEditPost = () => {
    setEditingPost(false);
    setError(null);
    if (editFileRef.current) editFileRef.current.value = "";
  };

  const onEditFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setError(null);
    setUploadingEdit(true);
    try {
      const res = await api.uploads.create(file);
      setEditMediaUrl(res.url);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Upload failed");
      if (editFileRef.current) editFileRef.current.value = "";
    } finally {
      setUploadingEdit(false);
    }
  };

  const onSavePost = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!post) return;
    setSavingPost(true);
    try {
      const res = await api.posts.update(id, {
        title: editTitle,
        content: editContent,
        category: editCategory,
        privacy_level: editPrivacy,
        media_url: editMediaUrl || undefined,
      });
      setPost(res.post);
      setEditingPost(false);
      if (editFileRef.current) editFileRef.current.value = "";
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to update post");
    } finally {
      setSavingPost(false);
    }
  };

  const onDeleteComment = async (commentId: number) => {
    if (!confirm("Delete this comment?")) return;
    try {
      await api.comments.remove(id, commentId);
      setComments((prev: Comment[]) => prev.filter((c) => c.id !== commentId));
      if (post) setPost({ ...post, comment_count: Math.max(0, post.comment_count - 1) });
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to delete comment");
    }
  };

  const onStartEditComment = (c: Comment) => {
    setEditingCommentId(c.id);
    setEditCommentText(c.content);
    setEditCommentMediaUrl(c.media_url ?? "");
  };

  const onCancelEditComment = () => {
    setEditingCommentId(null);
    setEditCommentText("");
    setEditCommentMediaUrl("");
    if (editCommentFileRef.current) editCommentFileRef.current.value = "";
  };

  const onEditCommentFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setError(null);
    setUploadingEditComment(true);
    try {
      const res = await api.uploads.create(file);
      setEditCommentMediaUrl(res.url);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Upload failed");
      if (editCommentFileRef.current) editCommentFileRef.current.value = "";
    } finally {
      setUploadingEditComment(false);
    }
  };

  const onSaveComment = async (commentId: number) => {
    if (!editCommentText.trim()) return;
    setSavingComment(true);
    try {
      const res = await api.comments.update(id, commentId, {
        content: editCommentText,
        media_url: editCommentMediaUrl || undefined,
      });
      setComments((prev: Comment[]) =>
        prev.map((c) => (c.id === commentId ? res.comment : c)),
      );
      onCancelEditComment();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to update comment");
    } finally {
      setSavingComment(false);
    }
  };

  if (loading) return <p className="text-sm text-slate-500">Loading…</p>;
  if (!post) {
    return (
      <div className="space-y-4">
        <p className="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">
          {error ?? "Post not found"}
        </p>
        <Link href="/feed" className="text-sm text-indigo-600 hover:underline">
          ← Back to feed
        </Link>
      </div>
    );
  }

  const isOwner = user?.id === post.user_id;

  return (
    <div className="space-y-6">
      <Link href="/feed" className="text-sm text-indigo-600 hover:underline">
        ← Back to feed
      </Link>

      {editingPost ? (
        <form
          onSubmit={onSavePost}
          className="space-y-3 rounded-xl border border-slate-200 bg-white p-5 shadow-sm"
        >
          <input
            type="text"
            value={editTitle}
            onChange={(e) => setEditTitle(e.target.value)}
            placeholder="Title (min 3 characters)"
            required
            minLength={3}
            className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
          />
          <textarea
            value={editContent}
            onChange={(e) => setEditContent(e.target.value)}
            rows={4}
            required
            minLength={10}
            placeholder="Content (min 10 characters)"
            className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
          />
          <div className="grid grid-cols-2 gap-3">
            <input
              type="text"
              value={editCategory}
              onChange={(e) => setEditCategory(e.target.value)}
              placeholder="Category"
              required
              className="rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            />
            <select
              value={editPrivacy}
              onChange={(e) => setEditPrivacy(e.target.value as PrivacyLevel)}
              className="rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            >
              <option value="public">Public</option>
              <option value="almost_private">Followers only</option>
              <option value="private">Private</option>
            </select>
          </div>
          <div>
            <label className="block text-xs font-medium text-slate-600">Image (optional)</label>
            <input
              ref={editFileRef}
              type="file"
              accept="image/jpeg,image/png,image/gif,image/webp"
              onChange={onEditFileChange}
              disabled={uploadingEdit}
              className="mt-1 block w-full text-sm text-slate-700 file:mr-3 file:rounded-lg file:border-0 file:bg-indigo-50 file:px-3 file:py-1.5 file:text-sm file:font-medium file:text-indigo-700 hover:file:bg-indigo-100 disabled:opacity-60"
            />
            {uploadingEdit ? (
              <p className="mt-1 text-xs text-slate-500">Uploading…</p>
            ) : editMediaUrl ? (
              <div className="mt-2 flex items-start gap-3">
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img
                  src={resolveMediaUrl(editMediaUrl)}
                  alt="preview"
                  className="h-24 w-24 rounded-lg object-cover ring-1 ring-slate-200"
                />
                <button
                  type="button"
                  onClick={() => {
                    setEditMediaUrl("");
                    if (editFileRef.current) editFileRef.current.value = "";
                  }}
                  className="text-xs text-slate-500 hover:text-red-600"
                >
                  Remove
                </button>
              </div>
            ) : null}
          </div>
          <div className="flex justify-end gap-2">
            <button
              type="button"
              onClick={onCancelEditPost}
              className="rounded-lg border border-slate-300 bg-white px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={savingPost || uploadingEdit}
              className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm transition hover:bg-indigo-700 disabled:opacity-60"
            >
              {savingPost ? "Saving…" : "Save changes"}
            </button>
          </div>
        </form>
      ) : (
        <PostCard post={post} />
      )}

      {isOwner && !editingPost ? (
        <div className="flex gap-3">
          <button
            onClick={onStartEditPost}
            className="text-sm text-indigo-600 hover:underline"
          >
            Edit post
          </button>
          <button
            onClick={onDeletePost}
            className="text-sm text-red-600 hover:underline"
          >
            Delete post
          </button>
        </div>
      ) : null}

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-slate-900">Comments</h2>

        <form onSubmit={onComment} className="space-y-2 rounded-lg border border-slate-200 bg-white p-3 shadow-sm">
          <div className="flex gap-2">
            <input
              type="text"
              value={commentText}
              onChange={(e) => setCommentText(e.target.value)}
              placeholder="Add a comment…"
              className="flex-1 rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            />
            <button
              type="submit"
              disabled={sending || uploadingComment || !commentText.trim()}
              className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm transition hover:bg-indigo-700 disabled:opacity-60"
            >
              Send
            </button>
          </div>

          <div>
            <input
              ref={commentFileRef}
              type="file"
              accept="image/jpeg,image/png,image/gif,image/webp"
              onChange={onCommentFileChange}
              disabled={uploadingComment}
              className="block text-xs text-slate-700 file:mr-3 file:rounded-lg file:border-0 file:bg-indigo-50 file:px-3 file:py-1 file:text-xs file:font-medium file:text-indigo-700 hover:file:bg-indigo-100 disabled:opacity-60"
            />
            {uploadingComment ? (
              <p className="mt-1 text-xs text-slate-500">Uploading…</p>
            ) : commentMediaUrl ? (
              <div className="mt-2 flex items-start gap-3">
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img
                  src={resolveMediaUrl(commentMediaUrl)}
                  alt="preview"
                  className="h-20 w-20 rounded-lg object-cover ring-1 ring-slate-200"
                />
                <button
                  type="button"
                  onClick={clearCommentImage}
                  className="text-xs text-slate-500 hover:text-red-600"
                >
                  Remove
                </button>
              </div>
            ) : null}
          </div>
        </form>

        {comments.length === 0 && !loadingComments ? (
          <p className="text-sm text-slate-500">No comments yet.</p>
        ) : (
          <ul className="space-y-3">
            {comments.map((c) => {
              const isEditing = editingCommentId === c.id;
              return (
                <li
                  key={c.id}
                  className="rounded-lg border border-slate-200 bg-white p-3 text-sm shadow-sm"
                >
                  <div className="flex items-center justify-between text-xs text-slate-500">
                    <Link
                      href={`/users/${c.user_id}`}
                      className="font-medium text-slate-700 hover:text-indigo-600 hover:underline"
                    >
                      @{c.author}
                    </Link>
                    <span>{formatRelative(c.created_at)}</span>
                  </div>

                  {isEditing ? (
                    <div className="mt-2 space-y-2">
                      <input
                        type="text"
                        value={editCommentText}
                        onChange={(e) => setEditCommentText(e.target.value)}
                        className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
                      />
                      <div>
                        <input
                          ref={editCommentFileRef}
                          type="file"
                          accept="image/jpeg,image/png,image/gif,image/webp"
                          onChange={onEditCommentFileChange}
                          disabled={uploadingEditComment}
                          className="block text-xs text-slate-700 file:mr-3 file:rounded-lg file:border-0 file:bg-indigo-50 file:px-3 file:py-1 file:text-xs file:font-medium file:text-indigo-700 hover:file:bg-indigo-100 disabled:opacity-60"
                        />
                        {uploadingEditComment ? (
                          <p className="mt-1 text-xs text-slate-500">Uploading…</p>
                        ) : editCommentMediaUrl ? (
                          <div className="mt-2 flex items-start gap-3">
                            {/* eslint-disable-next-line @next/next/no-img-element */}
                            <img
                              src={resolveMediaUrl(editCommentMediaUrl)}
                              alt="preview"
                              className="h-20 w-20 rounded-lg object-cover ring-1 ring-slate-200"
                            />
                            <button
                              type="button"
                              onClick={() => {
                                setEditCommentMediaUrl("");
                                if (editCommentFileRef.current)
                                  editCommentFileRef.current.value = "";
                              }}
                              className="text-xs text-slate-500 hover:text-red-600"
                            >
                              Remove
                            </button>
                          </div>
                        ) : null}
                      </div>
                      <div className="flex justify-end gap-2">
                        <button
                          type="button"
                          onClick={onCancelEditComment}
                          className="rounded-lg border border-slate-300 bg-white px-3 py-1 text-xs font-medium text-slate-700 hover:bg-slate-50"
                        >
                          Cancel
                        </button>
                        <button
                          type="button"
                          onClick={() => onSaveComment(c.id)}
                          disabled={savingComment || uploadingEditComment || !editCommentText.trim()}
                          className="rounded-lg bg-indigo-600 px-3 py-1 text-xs font-medium text-white shadow-sm transition hover:bg-indigo-700 disabled:opacity-60"
                        >
                          {savingComment ? "Saving…" : "Save"}
                        </button>
                      </div>
                    </div>
                  ) : (
                    <>
                      <p className="mt-1 whitespace-pre-wrap text-slate-700">{c.content}</p>
                      {c.media_url ? (
                        // eslint-disable-next-line @next/next/no-img-element
                        <img
                          src={resolveMediaUrl(c.media_url)}
                          alt=""
                          className="mt-2 max-h-72 rounded-lg object-cover"
                        />
                      ) : null}
                      {user?.id === c.user_id ? (
                        <div className="mt-2 flex gap-3">
                          <button
                            onClick={() => onStartEditComment(c)}
                            className="text-xs text-indigo-600 hover:underline"
                          >
                            Edit
                          </button>
                          <button
                            onClick={() => onDeleteComment(c.id)}
                            className="text-xs text-red-600 hover:underline"
                          >
                            Delete
                          </button>
                        </div>
                      ) : null}
                    </>
                  )}
                </li>
              );
            })}
          </ul>
        )}

        {comments.length > 0 && hasMoreComments ? (
          <InfiniteScrollSentinel onIntersect={loadMoreComments} enabled={!loadingComments} />
        ) : null}
        {loadingComments && comments.length > 0 ? (
          <p className="text-center text-xs text-slate-400">Loading more comments…</p>
        ) : null}
      </section>

      {error ? (
        <p className="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p>
      ) : null}
    </div>
  );
}
