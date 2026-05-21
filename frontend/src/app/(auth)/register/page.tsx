"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useMemo, useRef, useState } from "react";
import { useAuth } from "@/components/AuthProvider";
import { Avatar } from "@/components/Avatar";
import { api, ApiError } from "@/lib/api";

interface FormState {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  date_of_birth: string;
  nickname: string;
  gender: string;
  avatar_path: string;
  about_me: string;
  is_public: boolean;
}

const initial: FormState = {
  email: "",
  password: "",
  first_name: "",
  last_name: "",
  date_of_birth: "",
  nickname: "",
  gender: "",
  avatar_path: "",
  about_me: "",
  is_public: true,
};

export default function RegisterPage() {
  const { register, updateUser } = useAuth();
  const router = useRouter();
  const [form, setForm] = useState<FormState>(initial);
  const [avatarFile, setAvatarFile] = useState<File | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const update = <K extends keyof FormState>(key: K, value: FormState[K]) =>
    setForm((prev) => ({ ...prev, [key]: value }));

  const previewUrl = useMemo(() => (avatarFile ? URL.createObjectURL(avatarFile) : ""), [avatarFile]);

  const onPickAvatar = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] ?? null;
    setAvatarFile(file);
  };

  const clearAvatar = () => {
    setAvatarFile(null);
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSubmitting(true);
    try {
      const dob = new Date(form.date_of_birth);
      await register({
        ...form,
        date_of_birth: dob.toISOString(),
      });
      if (avatarFile) {
        try {
          const uploaded = await api.uploads.create(avatarFile);
          const res = await api.auth.update({ avatar_path: uploaded.url });
          updateUser(res.user);
        } catch {
          // non-blocking: account is created; user can set avatar later in profile
        }
      }
      router.replace("/feed");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Registration failed");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold text-slate-900">Create your account</h1>
        <p className="mt-1 text-sm text-slate-500">Join the community</p>
      </div>

      <form onSubmit={onSubmit} className="space-y-4">
        <div className="grid grid-cols-2 gap-3">
          <Field label="First name" value={form.first_name} onChange={(v) => update("first_name", v)} required />
          <Field label="Last name" value={form.last_name} onChange={(v) => update("last_name", v)} required />
        </div>
        <Field
          label="Nickname (optional)"
          value={form.nickname}
          onChange={(v) => update("nickname", v)}
          placeholder="Leave blank to use your name"
        />
        <Field label="Email" type="email" value={form.email} onChange={(v) => update("email", v)} required />
        <Field label="Password" type="password" value={form.password} onChange={(v) => update("password", v)} required />
        <Field label="Date of birth" type="date" value={form.date_of_birth} onChange={(v) => update("date_of_birth", v)} required />
        <label className="block">
          <span className="text-sm font-medium text-slate-700">Gender</span>
          <select
            value={form.gender}
            onChange={(e) => update("gender", e.target.value)}
            required
            className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
          >
            <option value="">Select…</option>
            <option value="male">Male</option>
            <option value="female">Female</option>
            <option value="other">Other</option>
          </select>
        </label>
        <div>
          <span className="text-sm font-medium text-slate-700">Profile photo (optional)</span>
          <div className="mt-1 flex items-center gap-3">
            <Avatar
              src={previewUrl || undefined}
              firstName={form.first_name}
              lastName={form.last_name}
              nickname={form.nickname}
              size={56}
            />
            <input
              ref={fileInputRef}
              type="file"
              accept="image/jpeg,image/png,image/gif,image/webp"
              onChange={onPickAvatar}
              className="block text-sm text-slate-700 file:mr-3 file:rounded-lg file:border-0 file:bg-indigo-50 file:px-3 file:py-1.5 file:text-sm file:font-medium file:text-indigo-700 hover:file:bg-indigo-100"
            />
            {avatarFile ? (
              <button
                type="button"
                onClick={clearAvatar}
                className="text-xs text-slate-500 hover:text-red-600"
              >
                Remove
              </button>
            ) : null}
          </div>
          <p className="mt-1 text-[11px] text-slate-500">
            Uploaded after your account is created.
          </p>
        </div>
        <label className="block">
          <span className="text-sm font-medium text-slate-700">About you</span>
          <textarea
            value={form.about_me}
            onChange={(e) => update("about_me", e.target.value)}
            rows={3}
            className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
          />
        </label>
        <label className="flex items-center gap-2 text-sm text-slate-700">
          <input
            type="checkbox"
            checked={form.is_public}
            onChange={(e) => update("is_public", e.target.checked)}
            className="h-4 w-4 rounded border-slate-300 text-indigo-600 focus:ring-indigo-500"
          />
          Public profile (anyone can follow without approval)
        </label>

        {error ? <p className="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p> : null}

        <button
          type="submit"
          disabled={submitting}
          className="w-full rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm transition hover:bg-indigo-700 disabled:opacity-60"
        >
          {submitting ? "Creating account…" : "Create account"}
        </button>
      </form>

      <p className="text-center text-sm text-slate-500">
        Already have an account?{" "}
        <Link href="/login" className="font-medium text-indigo-600 hover:underline">
          Sign in
        </Link>
      </p>
    </div>
  );
}

function Field({
  label,
  value,
  onChange,
  type = "text",
  required,
  placeholder,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  type?: string;
  required?: boolean;
  placeholder?: string;
}) {
  return (
    <label className="block">
      <span className="text-sm font-medium text-slate-700">{label}</span>
      <input
        type={type}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        required={required}
        placeholder={placeholder}
        className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
      />
    </label>
  );
}
