"use client";

import { useEffect, useId, useRef, useState } from "react";
import { useAuth } from "@/components/AuthProvider";
import { Avatar } from "@/components/Avatar";
import { api, ApiError, resolveMediaUrl } from "@/lib/api";
import type { CreatePostPayload, PrivacyLevel, UserSummary } from "@/types/api";

interface Props {
  onSubmit: (payload: CreatePostPayload) => Promise<void>;
  submitLabel?: string;
  inGroup?: boolean;
}

export function PostForm({ onSubmit, submitLabel = "Post", inGroup = false }: Props) {
  const { user } = useAuth();
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [category, setCategory] = useState("general");
  const [privacy, setPrivacy] = useState<PrivacyLevel>("public");
  const [mediaUrl, setMediaUrl] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [uploading, setUploading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const fileInputId = useId();

  const [followers, setFollowers] = useState<UserSummary[]>([]);
  const [audience, setAudience] = useState<Set<number>>(new Set());
  const [loadingFollowers, setLoadingFollowers] = useState(false);

  useEffect(() => {
    if (privacy !== "private" || inGroup || !user) return;
    if (followers.length > 0) return;
    let cancelled = false;
    setLoadingFollowers(true);
    api.users
      .followers(user.id, { limit: 100 })
      .then((res) => {
        if (!cancelled) setFollowers(res.users ?? []);
      })
      .catch(() => {})
      .finally(() => {
        if (!cancelled) setLoadingFollowers(false);
      });
    return () => {
      cancelled = true;
    };
  }, [privacy, inGroup, user, followers.length]);

  const reset = () => {
    setTitle("");
    setContent("");
    setMediaUrl("");
    setAudience(new Set());
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  const toggleAudience = (id: number) => {
    setAudience((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const handle = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSubmitting(true);
    try {
      await onSubmit({
        title,
        content,
        category,
        privacy_level: privacy,
        media_url: mediaUrl || undefined,
        audience: privacy === "private" && !inGroup ? Array.from(audience) : undefined,
      });
      reset();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to create post");
    } finally {
      setSubmitting(false);
    }
  };

  const onFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setError(null);
    setUploading(true);
    try {
      const res = await api.uploads.create(file);
      setMediaUrl(res.url);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Upload failed");
      if (fileInputRef.current) fileInputRef.current.value = "";
    } finally {
      setUploading(false);
    }
  };

  const clearImage = () => {
    setMediaUrl("");
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  const previewSrc = resolveMediaUrl(mediaUrl);
  const showAudiencePicker = privacy === "private" && !inGroup;

  return (
    <form onSubmit={handle} className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
      <input
        type="text"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="Title (min 3 characters)"
        required
        minLength={3}
        className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
      />
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder="What's on your mind? (min 10 characters)"
        rows={3}
        required
        minLength={10}
        className="mt-3 w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
      />

      <div className="mt-3">
        <label htmlFor={fileInputId} className="block text-xs font-medium text-slate-600">Image (optional, max 5MB)</label>
        <input
          id={fileInputId}
          ref={fileInputRef}
          type="file"
          accept="image/jpeg,image/png,image/gif,image/webp"
          onChange={onFileChange}
          disabled={uploading}
          className="mt-1 block w-full text-sm text-slate-700 file:mr-3 file:rounded-lg file:border-0 file:bg-indigo-50 file:px-3 file:py-1.5 file:text-sm file:font-medium file:text-indigo-700 hover:file:bg-indigo-100 disabled:opacity-60"
        />
        {uploading ? (
          <p className="mt-1 text-xs text-slate-500">Uploading…</p>
        ) : previewSrc ? (
          <div className="mt-2 flex items-start gap-3">
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img src={previewSrc} alt="preview" className="h-24 w-24 rounded-lg object-cover ring-1 ring-slate-200" />
            <button
              type="button"
              onClick={clearImage}
              className="text-xs text-slate-500 hover:text-red-600"
            >
              Remove
            </button>
          </div>
        ) : null}
      </div>

      <div className="mt-3 grid grid-cols-2 gap-3">
        <input
          type="text"
          value={category}
          onChange={(e) => setCategory(e.target.value)}
          placeholder="Category"
          required
          className="rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
        />
        <select
          value={privacy}
          onChange={(e) => setPrivacy(e.target.value as PrivacyLevel)}
          className="rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
        >
          <option value="public">Public</option>
          <option value="almost_private">Followers only</option>
          <option value="private">Private</option>
        </select>
      </div>

      {showAudiencePicker ? (
        <div className="mt-3 rounded-lg border border-slate-200 bg-slate-50 p-3">
          <p className="text-xs font-medium text-slate-700">
            Pick the followers who can see this post
          </p>
          <p className="mt-0.5 text-[11px] text-slate-500">
            Leave empty for visible to you only.
          </p>
          {loadingFollowers ? (
            <p className="mt-2 text-xs text-slate-500">Loading followers…</p>
          ) : followers.length === 0 ? (
            <p className="mt-2 text-xs text-slate-500">
              You don&apos;t have any followers yet.
            </p>
          ) : (
            <ul className="mt-2 max-h-48 space-y-1 overflow-y-auto">
              {followers.map((f) => {
                const checked = audience.has(f.id);
                return (
                  <li key={f.id}>
                    <label className="flex cursor-pointer items-center gap-2 rounded-md px-2 py-1 text-sm hover:bg-white">
                      <input
                        type="checkbox"
                        checked={checked}
                        onChange={() => toggleAudience(f.id)}
                        className="h-4 w-4 rounded border-slate-300 text-indigo-600 focus:ring-indigo-500"
                      />
                      <Avatar
                        src={f.avatar_path}
                        firstName={f.first_name}
                        lastName={f.last_name}
                        nickname={f.nickname}
                        size={24}
                      />
                      <span className="truncate text-slate-700">
                        {f.first_name} {f.last_name}{" "}
                        <span className="text-slate-400">@{f.nickname}</span>
                      </span>
                    </label>
                  </li>
                );
              })}
            </ul>
          )}
          {followers.length > 0 ? (
            <p className="mt-2 text-[11px] text-slate-500">
              {audience.size} of {followers.length} selected
            </p>
          ) : null}
        </div>
      ) : null}

      {error ? (
        <p className="mt-3 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p>
      ) : null}

      <div className="mt-3 flex justify-end">
        <button
          type="submit"
          disabled={submitting || uploading}
          className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm transition hover:bg-indigo-700 disabled:opacity-60"
        >
          {submitting ? "Posting…" : submitLabel}
        </button>
      </div>
    </form>
  );
}
