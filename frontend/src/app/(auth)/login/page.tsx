"use client";

import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { FormEvent, Suspense, useState } from "react";
import { useAuth } from "@/components/AuthProvider";
import { FormMessage, TextField } from "@/components/forms";
import { api, ApiError } from "@/lib/api";

function LoginForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { setSession } = useAuth();
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setMessage("");

    try {
      const response = await api.login({ identifier, password });
      if (!response.token || !response.user) {
        throw new ApiError("Login response did not include a session.");
      }

      setSession(response.token, response.user);
      router.push(searchParams.get("next") || "/feed");
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Login failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="mx-auto flex min-h-screen max-w-md flex-col justify-center px-6 py-12">
      <form
        className="rounded-3xl border border-slate-200 bg-white p-8 shadow-xl shadow-slate-200/70"
        onSubmit={onSubmit}
      >
        <p className="text-sm font-semibold uppercase tracking-[0.3em] text-sky-600">
          Welcome back
        </p>
        <h1 className="mt-3 text-3xl font-bold text-slate-950">Log in</h1>
        <div className="mt-8 grid gap-5">
          <TextField
            label="Email or nickname"
            name="identifier"
            onChange={setIdentifier}
            required
            value={identifier}
          />
          <TextField
            label="Password"
            name="password"
            onChange={setPassword}
            required
            type="password"
            value={password}
          />
          <FormMessage message={message} tone="error" />
          <button
            className="rounded-xl bg-sky-600 px-4 py-3 font-semibold text-white transition hover:bg-sky-700"
            disabled={loading}
            type="submit"
          >
            {loading ? "Logging in..." : "Log in"}
          </button>
        </div>
        <p className="mt-6 text-center text-sm text-slate-600">
          New here?{" "}
          <Link className="font-semibold text-sky-700" href="/register">
            Create an account
          </Link>
        </p>
      </form>
    </main>
  );
}

export default function LoginPage() {
  return (
    <Suspense
      fallback={
        <main className="flex min-h-screen items-center justify-center text-slate-600">
          Loading login...
        </main>
      }
    >
      <LoginForm />
    </Suspense>
  );
}
