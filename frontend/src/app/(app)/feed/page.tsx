"use client";

import { FormEvent, useEffect, useState } from "react";
import { PostCard } from "@/components/PostCard";
import { PostForm } from "@/components/PostForm";
import { FormMessage, TextField } from "@/components/forms";
import { api } from "@/lib/api";
import type { Post } from "@/types/api";

export default function FeedPage() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [category, setCategory] = useState("");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(true);

  async function loadPosts(nextCategory = category) {
    setLoading(true);
    setMessage("");

    try {
      const response = await api.listPosts({
        category: nextCategory || undefined,
        limit: 50,
      });
      setPosts(response.posts ?? []);
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Could not load posts");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void loadPosts("");
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  function filterPosts(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    void loadPosts(category);
  }

  return (
    <div className="grid gap-6 lg:grid-cols-[22rem_1fr]">
      <aside className="grid content-start gap-6">
        <section className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <h1 className="text-2xl font-bold text-slate-950">Create post</h1>
          <div className="mt-4">
            <PostForm onSubmit={async (body) => {
              await api.createPost(body);
              await loadPosts();
            }} />
          </div>
        </section>
        <form
          className="grid gap-4 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"
          onSubmit={filterPosts}
        >
          <TextField
            label="Filter by category"
            name="category"
            onChange={setCategory}
            placeholder="general"
            value={category}
          />
          <button
            className="rounded-xl border border-slate-300 px-4 py-3 font-semibold text-slate-700 transition hover:border-sky-400 hover:text-sky-700"
            type="submit"
          >
            Apply filter
          </button>
        </form>
      </aside>
      <section>
        <div className="mb-5 flex items-center justify-between">
          <h2 className="text-3xl font-bold text-slate-950">Feed</h2>
          <button
            className="rounded-xl border border-slate-300 px-4 py-2 text-sm font-semibold text-slate-700 transition hover:border-sky-400 hover:text-sky-700"
            onClick={() => loadPosts()}
            type="button"
          >
            Refresh
          </button>
        </div>
        <FormMessage message={message} tone="error" />
        {loading ? (
          <p className="mt-6 text-slate-600">Loading posts...</p>
        ) : posts.length === 0 ? (
          <p className="mt-6 rounded-2xl border border-dashed border-slate-300 bg-white p-8 text-center text-slate-600">
            No posts found.
          </p>
        ) : (
          <div className="mt-6 grid gap-4">
            {posts.map((post) => (
              <PostCard key={post.id} post={post} />
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
