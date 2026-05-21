"use client";

import { useRef, useState } from "react";
import { Avatar } from "@/components/Avatar";
import { api, ApiError, resolveMediaUrl } from "@/lib/api";

interface Props {
  value: string;
  onChange: (path: string) => void;
  firstName?: string;
  lastName?: string;
  nickname?: string;
  label?: string;
}

export function AvatarUpload({
  value,
  onChange,
  firstName,
  lastName,
  nickname,
  label = "Profile photo (optional)",
}: Props) {
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const onFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setError(null);
    setUploading(true);
    try {
      const res = await api.uploads.create(file);
      onChange(res.url);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Upload failed");
      if (inputRef.current) inputRef.current.value = "";
    } finally {
      setUploading(false);
    }
  };

  const clear = () => {
    onChange("");
    if (inputRef.current) inputRef.current.value = "";
  };

  return (
    <div>
      <span className="text-sm font-medium text-slate-700">{label}</span>
      <div className="mt-1 flex items-center gap-3">
        <Avatar
          src={resolveMediaUrl(value)}
          firstName={firstName}
          lastName={lastName}
          nickname={nickname}
          size={56}
        />
        <input
          ref={inputRef}
          type="file"
          accept="image/jpeg,image/png,image/gif,image/webp"
          onChange={onFileChange}
          disabled={uploading}
          className="block text-sm text-slate-700 file:mr-3 file:rounded-lg file:border-0 file:bg-indigo-50 file:px-3 file:py-1.5 file:text-sm file:font-medium file:text-indigo-700 hover:file:bg-indigo-100 disabled:opacity-60"
        />
        {value ? (
          <button
            type="button"
            onClick={clear}
            className="text-xs text-slate-500 hover:text-red-600"
          >
            Remove
          </button>
        ) : null}
      </div>
      {uploading ? <p className="mt-1 text-xs text-slate-500">Uploading…</p> : null}
      {error ? <p className="mt-1 text-xs text-red-600">{error}</p> : null}
    </div>
  );
}
