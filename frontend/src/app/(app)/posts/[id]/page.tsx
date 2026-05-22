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
import type { Comment, Post } from "@/types/api";

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

      <PostCard post={post} />

      {isOwner ? (
        <button
          onClick={onDeletePost}
          className="text-sm text-red-600 hover:underline"
        >
          Delete post
        </button>
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
            {comments.map((c) => (
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
                  <button
                    onClick={() => onDeleteComment(c.id)}
                    className="mt-2 text-xs text-red-600 hover:underline"
                  >
                    Delete
                  </button>
                ) : null}
              </li>
            ))}
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
