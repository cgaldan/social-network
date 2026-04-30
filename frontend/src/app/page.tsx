import Link from "next/link";

export default function HomePage() {
  return (
    <main className="mx-auto flex min-h-screen max-w-5xl flex-col justify-center px-6 py-16">
      <section className="rounded-3xl border border-slate-200 bg-white/90 p-8 shadow-xl shadow-slate-200/70 md:p-12">
        <p className="text-sm font-semibold uppercase tracking-[0.35em] text-sky-600">
          Social Network
        </p>
        <h1 className="mt-5 max-w-3xl text-4xl font-bold tracking-tight text-slate-950 md:text-6xl">
          Connect with your people, posts, groups, messages, and events.
        </h1>
        <p className="mt-5 max-w-2xl text-lg leading-8 text-slate-600">
          This frontend talks directly to the Go backend API on port 8000 and
          uses the backend session token for protected REST and WebSocket calls.
        </p>
        <div className="mt-8 flex flex-wrap gap-3">
          <Link
            className="rounded-full bg-sky-600 px-5 py-3 text-sm font-semibold text-white shadow-lg shadow-sky-200 transition hover:bg-sky-700"
            href="/login"
          >
            Log in
          </Link>
          <Link
            className="rounded-full border border-slate-300 px-5 py-3 text-sm font-semibold text-slate-800 transition hover:border-sky-400 hover:text-sky-700"
            href="/register"
          >
            Create account
          </Link>
        </div>
      </section>
    </main>
  );
}
