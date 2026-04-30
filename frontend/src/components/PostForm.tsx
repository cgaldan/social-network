"use client";

import { FormEvent, useState } from "react";
import { FormMessage, SelectField, TextArea, TextField } from "@/components/forms";
import type { CreatePostRequest, Post } from "@/types/api";

const privacyOptions = [
  { label: "Public", value: "public" },
  { label: "Almost private", value: "almost_private" },
  { label: "Private", value: "private" },
];

export function PostForm({
  initialPost,
  onSubmit,
  submitLabel = "Publish post",
}: {
  initialPost?: Post;
  onSubmit: (body: CreatePostRequest) => Promise<void>;
  submitLabel?: string;
}) {
  const [form, setForm] = useState({
    title: initialPost?.title ?? "",
    content: initialPost?.content ?? "",
    category: initialPost?.category ?? "",
    privacy_level: initialPost?.privacy_level ?? "public",
    media_url: initialPost?.media_url ?? "",
  });
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");

  function update(name: keyof typeof form, value: string) {
    setForm((current) => ({ ...current, [name]: value }));
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setMessage("");

    try {
      await onSubmit({
        ...form,
        media_url: form.media_url || undefined,
      });
      if (!initialPost) {
        setForm({
          title: "",
          content: "",
          category: "",
          privacy_level: "public",
          media_url: "",
        });
      }
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Post action failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <form
      className="grid gap-4 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"
      onSubmit={handleSubmit}
    >
      <TextField
        label="Title"
        name="title"
        onChange={(value) => update("title", value)}
        required
        value={form.title}
      />
      <TextArea
        label="Content"
        name="content"
        onChange={(value) => update("content", value)}
        required
        value={form.content}
      />
      <div className="grid gap-4 md:grid-cols-3">
        <TextField
          label="Category"
          name="category"
          onChange={(value) => update("category", value)}
          placeholder="general"
          value={form.category}
        />
        <SelectField
          label="Privacy"
          name="privacy_level"
          onChange={(value) => update("privacy_level", value)}
          options={privacyOptions}
          value={form.privacy_level}
        />
        <TextField
          label="Media URL"
          name="media_url"
          onChange={(value) => update("media_url", value)}
          value={form.media_url}
        />
      </div>
      <FormMessage message={message} tone="error" />
      <button
        className="rounded-xl bg-sky-600 px-4 py-3 font-semibold text-white transition hover:bg-sky-700"
        disabled={loading}
        type="submit"
      >
        {loading ? "Working..." : submitLabel}
      </button>
    </form>
  );
}
