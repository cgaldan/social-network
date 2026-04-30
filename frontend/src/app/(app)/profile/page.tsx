"use client";

import { FormEvent, useEffect, useState } from "react";
import { useAuth } from "@/components/AuthProvider";
import { FormMessage, SelectField, TextArea, TextField } from "@/components/forms";
import { api } from "@/lib/api";
import { dateInputToISOString, formatDate, toDateInputValue } from "@/lib/format";

export default function ProfilePage() {
  const { logout, refreshUser, user } = useAuth();
  const [message, setMessage] = useState("");
  const [tone, setTone] = useState<"error" | "success">("success");
  const [loading, setLoading] = useState(false);
  const [form, setForm] = useState({
    email: "",
    first_name: "",
    last_name: "",
    date_of_birth: "",
    nickname: "",
    gender: "",
    avatar_path: "",
    about_me: "",
    is_public: "true",
  });

  useEffect(() => {
    if (!user) {
      return;
    }

    setForm({
      email: user.email || "",
      first_name: user.first_name || "",
      last_name: user.last_name || "",
      date_of_birth: toDateInputValue(user.date_of_birth),
      nickname: user.nickname || "",
      gender: user.gender || "",
      avatar_path: user.avatar_path || "",
      about_me: user.about_me || "",
      is_public: String(user.is_public),
    });
  }, [user]);

  function update(name: keyof typeof form, value: string) {
    setForm((current) => ({ ...current, [name]: value }));
  }

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setMessage("");

    try {
      await api.updateMe({
        ...form,
        date_of_birth: dateInputToISOString(form.date_of_birth),
        is_public: form.is_public === "true",
      });
      await refreshUser();
      setTone("success");
      setMessage("Profile updated.");
    } catch (error) {
      setTone("error");
      setMessage(error instanceof Error ? error.message : "Update failed");
    } finally {
      setLoading(false);
    }
  }

  async function onDelete() {
    if (!confirm("Delete your account? This cannot be undone.")) {
      return;
    }

    try {
      await api.deleteMe();
      await logout();
    } catch (error) {
      setTone("error");
      setMessage(error instanceof Error ? error.message : "Delete failed");
    }
  }

  return (
    <div className="grid gap-6 lg:grid-cols-[1fr_20rem]">
      <section className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
        <h1 className="text-3xl font-bold text-slate-950">Profile settings</h1>
        <form className="mt-6 grid gap-5 md:grid-cols-2" onSubmit={onSubmit}>
          <TextField
            label="First name"
            name="first_name"
            onChange={(value) => update("first_name", value)}
            required
            value={form.first_name}
          />
          <TextField
            label="Last name"
            name="last_name"
            onChange={(value) => update("last_name", value)}
            required
            value={form.last_name}
          />
          <TextField
            label="Email"
            name="email"
            onChange={(value) => update("email", value)}
            required
            type="email"
            value={form.email}
          />
          <TextField
            label="Nickname"
            name="nickname"
            onChange={(value) => update("nickname", value)}
            required
            value={form.nickname}
          />
          <TextField
            label="Date of birth"
            name="date_of_birth"
            onChange={(value) => update("date_of_birth", value)}
            required
            type="date"
            value={form.date_of_birth}
          />
          <TextField
            label="Gender"
            name="gender"
            onChange={(value) => update("gender", value)}
            value={form.gender}
          />
          <div className="md:col-span-2">
            <TextField
              label="Avatar URL"
              name="avatar_path"
              onChange={(value) => update("avatar_path", value)}
              value={form.avatar_path}
            />
          </div>
          <SelectField
            label="Profile visibility"
            name="is_public"
            onChange={(value) => update("is_public", value)}
            options={[
              { label: "Public", value: "true" },
              { label: "Private", value: "false" },
            ]}
            value={form.is_public}
          />
          <div className="md:col-span-2">
            <TextArea
              label="About me"
              name="about_me"
              onChange={(value) => update("about_me", value)}
              value={form.about_me}
            />
          </div>
          <div className="grid gap-4 md:col-span-2">
            <FormMessage message={message} tone={tone} />
            <button
              className="rounded-xl bg-sky-600 px-4 py-3 font-semibold text-white transition hover:bg-sky-700"
              disabled={loading}
              type="submit"
            >
              {loading ? "Saving..." : "Save profile"}
            </button>
          </div>
        </form>
      </section>
      <aside className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
        <h2 className="text-xl font-bold text-slate-950">Account</h2>
        <dl className="mt-5 grid gap-4 text-sm">
          <div>
            <dt className="font-semibold text-slate-500">Followers</dt>
            <dd className="text-slate-950">{user?.followers_count ?? 0}</dd>
          </div>
          <div>
            <dt className="font-semibold text-slate-500">Following</dt>
            <dd className="text-slate-950">{user?.following_count ?? 0}</dd>
          </div>
          <div>
            <dt className="font-semibold text-slate-500">Joined</dt>
            <dd className="text-slate-950">{formatDate(user?.created_at)}</dd>
          </div>
          <div>
            <dt className="font-semibold text-slate-500">Last seen</dt>
            <dd className="text-slate-950">{formatDate(user?.last_seen)}</dd>
          </div>
        </dl>
        <button
          className="mt-6 w-full rounded-xl border border-red-200 px-4 py-3 font-semibold text-red-600 transition hover:bg-red-50"
          onClick={onDelete}
          type="button"
        >
          Delete account
        </button>
      </aside>
    </div>
  );
}
