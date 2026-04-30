"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";
import { useAuth } from "@/components/AuthProvider";
import { FormMessage, SelectField, TextArea, TextField } from "@/components/forms";
import { api, ApiError } from "@/lib/api";
import { dateInputToISOString } from "@/lib/format";

export default function RegisterPage() {
  const router = useRouter();
  const { setSession } = useAuth();
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);
  const [form, setForm] = useState({
    email: "",
    password: "",
    first_name: "",
    last_name: "",
    date_of_birth: "",
    nickname: "",
    gender: "",
    avatar_path: "",
    about_me: "",
    is_public: "true",
  });

  function update(name: keyof typeof form, value: string) {
    setForm((current) => ({ ...current, [name]: value }));
  }

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setMessage("");

    try {
      const response = await api.register({
        ...form,
        date_of_birth: dateInputToISOString(form.date_of_birth),
        is_public: form.is_public === "true",
      });

      if (!response.token || !response.user) {
        throw new ApiError("Registration response did not include a session.");
      }

      setSession(response.token, response.user);
      router.push("/feed");
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Registration failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="mx-auto flex min-h-screen max-w-3xl flex-col justify-center px-6 py-12">
      <form
        className="rounded-3xl border border-slate-200 bg-white p-8 shadow-xl shadow-slate-200/70"
        onSubmit={onSubmit}
      >
        <p className="text-sm font-semibold uppercase tracking-[0.3em] text-sky-600">
          Join the network
        </p>
        <h1 className="mt-3 text-3xl font-bold text-slate-950">
          Create account
        </h1>
        <div className="mt-8 grid gap-5 md:grid-cols-2">
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
            label="Password"
            name="password"
            onChange={(value) => update("password", value)}
            required
            type="password"
            value={form.password}
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
            <TextField
              label="Avatar URL"
              name="avatar_path"
              onChange={(value) => update("avatar_path", value)}
              value={form.avatar_path}
            />
          </div>
          <div className="md:col-span-2">
            <TextArea
              label="About me"
              name="about_me"
              onChange={(value) => update("about_me", value)}
              value={form.about_me}
            />
          </div>
        </div>
        <div className="mt-6 grid gap-4">
          <FormMessage message={message} tone="error" />
          <button
            className="rounded-xl bg-sky-600 px-4 py-3 font-semibold text-white transition hover:bg-sky-700"
            disabled={loading}
            type="submit"
          >
            {loading ? "Creating account..." : "Create account"}
          </button>
        </div>
        <p className="mt-6 text-center text-sm text-slate-600">
          Already have an account?{" "}
          <Link className="font-semibold text-sky-700" href="/login">
            Log in
          </Link>
        </p>
      </form>
    </main>
  );
}
