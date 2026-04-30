"use client";

import { useParams, useRouter } from "next/navigation";
import { FormEvent, useEffect, useMemo, useState } from "react";
import { PostForm } from "@/components/PostForm";
import { FormMessage, TextArea, TextField } from "@/components/forms";
import { api } from "@/lib/api";
import { formatDate } from "@/lib/format";
import type { Post } from "@/types/api";

export default function PostDetailPage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const postId = useMemo(() => Number(params.id), [params.id]);
  const [post, setPost] = useState<Post | null>(null);
  const [message, setMessage] = useState("");
  const [tone, setTone] = useState<"error" | "success">("success");
  const [loading, setLoading] = useState(true);
  const [comment, setComment] = useState({ content: "", media_url: "" });
  const [editComment, setEditComment] = useState({
    comment_id: "",
    content: "",
    media_url: "",
  });

  async function loadPost() {
    setLoading(true);
    setMessage("");

    try {
      const response = await api.getPost(postId);
      setPost(response.post ?? null);
    } catch (error) {
      setTone("error");
      setMessage(error instanceof Error ? error.message : "Could not load post");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    if (Number.isFinite(postId)) {
      void loadPost();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [postId]);

  async function createComment(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    try {
      await api.createComment(postId, {
        content: comment.content,
        media_url: comment.media_url || undefined,
      });
      setComment({ content: "", media_url: "" });
      setTone("success");
      setMessage("Comment created.");
      await loadPost();
    } catch (error) {
      setTone("error");
      setMessage(error instanceof Error ? error.message : "Comment failed");
    }
  }

  async function updateComment(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    try {
      await api.updateComment(postId, Number(editComment.comment_id), {
        content: editComment.content,
        media_url: editComment.media_url || undefined,
      });
      setTone("success");
      setMessage("Comment updated.");
    } catch (error) {
      setTone("error");
      setMessage(error instanceof Error ? error.message : "Update failed");
    }
  }

  async function deleteComment() {
    try {
      await api.deleteComment(postId, Number(editComment.comment_id));
      setTone("success");
      setMessage("Comment deleted.");
      setEditComment({ comment_id: "", content: "", media_url: "" });
    } catch (error) {
      setTone("error");
      setMessage(error instanceof Error ? error.message : "Delete failed");
    }
  }

  async function deletePost() {
    if (!confirm("Delete this post?")) {
      return;
    }

    try {
      await api.deletePost(postId);
      router.push("/feed");
    } catch (error) {
      setTone("error");
      setMessage(error instanceof Error ? error.message : "Delete failed");
    }
  }

  if (loading) {
    return <p className="text-slate-600">Loading post...</p>;
  }

  return (
    <div className="grid gap-6 lg:grid-cols-[1fr_22rem]">
      <section className="grid gap-6">
        <FormMessage message={message} tone={tone} />
        {post ? (
          <article className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
            <p className="text-sm font-semibold uppercase tracking-wide text-sky-700">
              {post.category || "General"} / {post.privacy_level}
            </p>
            <h1 className="mt-3 text-4xl font-bold text-slate-950">
              {post.title}
            </h1>
            <p className="mt-4 whitespace-pre-wrap text-lg leading-8 text-slate-700">
              {post.content}
            </p>
            <footer className="mt-6 flex flex-wrap gap-4 text-sm text-slate-500">
              <span>By {post.author || `User #${post.user_id}`}</span>
              <span>{formatDate(post.created_at)}</span>
              <span>{post.comment_count} comments</span>
            </footer>
            <button
              className="mt-6 rounded-xl border border-red-200 px-4 py-3 font-semibold text-red-600 transition hover:bg-red-50"
              onClick={deletePost}
              type="button"
            >
              Delete post
            </button>
          </article>
        ) : (
          <p className="rounded-2xl border border-dashed border-slate-300 bg-white p-8 text-center text-slate-600">
            Post not found.
          </p>
        )}
        {post ? (
          <section>
            <h2 className="mb-4 text-2xl font-bold text-slate-950">
              Edit post
            </h2>
            <PostForm
              initialPost={post}
              onSubmit={async (body) => {
                await api.updatePost(postId, body);
                await loadPost();
                setTone("success");
                setMessage("Post updated.");
              }}
              submitLabel="Save post"
            />
          </section>
        ) : null}
      </section>
      <aside className="grid content-start gap-6">
        <form
          className="grid gap-4 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"
          onSubmit={createComment}
        >
          <h2 className="text-xl font-bold text-slate-950">Add comment</h2>
          <TextArea
            label="Comment"
            name="content"
            onChange={(value) =>
              setComment((current) => ({ ...current, content: value }))
            }
            required
            value={comment.content}
          />
          <TextField
            label="Media URL"
            name="media_url"
            onChange={(value) =>
              setComment((current) => ({ ...current, media_url: value }))
            }
            value={comment.media_url}
          />
          <button
            className="rounded-xl bg-sky-600 px-4 py-3 font-semibold text-white transition hover:bg-sky-700"
            type="submit"
          >
            Add comment
          </button>
          <p className="text-xs text-slate-500">
            The backend accepts comment create/edit/delete, but does not expose a
            comment listing route yet.
          </p>
        </form>
        <form
          className="grid gap-4 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"
          onSubmit={updateComment}
        >
          <h2 className="text-xl font-bold text-slate-950">Edit comment by ID</h2>
          <TextField
            label="Comment ID"
            name="comment_id"
            onChange={(value) =>
              setEditComment((current) => ({ ...current, comment_id: value }))
            }
            required
            type="number"
            value={editComment.comment_id}
          />
          <TextArea
            label="New content"
            name="content"
            onChange={(value) =>
              setEditComment((current) => ({ ...current, content: value }))
            }
            required
            value={editComment.content}
          />
          <TextField
            label="Media URL"
            name="media_url"
            onChange={(value) =>
              setEditComment((current) => ({ ...current, media_url: value }))
            }
            value={editComment.media_url}
          />
          <button
            className="rounded-xl bg-sky-600 px-4 py-3 font-semibold text-white transition hover:bg-sky-700"
            type="submit"
          >
            Update comment
          </button>
          <button
            className="rounded-xl border border-red-200 px-4 py-3 font-semibold text-red-600 transition hover:bg-red-50"
            disabled={!editComment.comment_id}
            onClick={deleteComment}
            type="button"
          >
            Delete comment
          </button>
        </form>
      </aside>
    </div>
  );
}
