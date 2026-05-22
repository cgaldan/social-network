"use client";

import Link from "next/link";
import { useCallback, useState } from "react";
import { InfiniteScrollSentinel } from "@/components/InfiniteScrollSentinel";
import { usePaginatedList } from "@/hooks/usePaginatedList";
import { api, ApiError } from "@/lib/api";
import { formatDate } from "@/lib/format";
import type { Group } from "@/types/api";

const PAGE_SIZE = 20;

export default function GroupsPage() {
  const [error, setError] = useState<string | null>(null);

  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [creating, setCreating] = useState(false);

  const fetcher = useCallback(
    async ({ limit, offset }: { limit: number; offset: number }) => {
      const res = await api.groups.list({ limit, offset });
      return { items: res.groups ?? [], hasMore: res.has_more };
    },
    [],
  );
  const { items: groups, hasMore, loading, error: fetchError, loadMore, setItems } =
    usePaginatedList<Group>(fetcher, PAGE_SIZE);

  const onCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;
    setCreating(true);
    try {
      const res = await api.groups.create({ title, description, conversation_id: 0 });
      setItems((prev) => [res.group, ...prev]);
      setTitle("");
      setDescription("");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to create group");
    } finally {
      setCreating(false);
    }
  };

  const onJoin = async (id: number) => {
    try {
      await api.groups.join(id);
      setItems((prev) =>
        prev.map((g) => (g.id === id ? { ...g, is_pending: true } : g)),
      );
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to send join request");
    }
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold text-slate-900">Groups</h1>

      <form onSubmit={onCreate} className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
        <h2 className="text-base font-semibold text-slate-900">Create a group</h2>
        <input
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="Group name"
          required
          className="mt-3 w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
        />
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          rows={2}
          placeholder="Description"
          required
          className="mt-3 w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
        />
        <div className="mt-3 flex justify-end">
          <button
            type="submit"
            disabled={creating}
            className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm transition hover:bg-indigo-700 disabled:opacity-60"
          >
            {creating ? "Creating…" : "Create"}
          </button>
        </div>
      </form>

      {(error || fetchError) ? (
        <p className="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{error ?? fetchError}</p>
      ) : null}

      {loading && groups.length === 0 ? (
        <p className="text-sm text-slate-500">Loading…</p>
      ) : groups.length === 0 ? (
        <p className="rounded-xl border border-dashed border-slate-300 bg-white p-6 text-center text-sm text-slate-500">
          No groups yet.
        </p>
      ) : (
        <ul className="space-y-3">
          {groups.map((g) => (
            <li
              key={g.id}
              className="flex items-start justify-between rounded-xl border border-slate-200 bg-white p-4 shadow-sm"
            >
              <div className="min-w-0">
                {g.is_member ? (
                  <Link
                    href={`/groups/${g.id}`}
                    className="text-base font-semibold text-slate-900 hover:underline"
                  >
                    {g.name}
                  </Link>
                ) : (
                  <span className="text-base font-semibold text-slate-900">{g.name}</span>
                )}
                <p className="mt-1 text-sm text-slate-600">{g.description}</p>
                <p className="mt-1 text-xs text-slate-400">Created {formatDate(g.created_at)}</p>
              </div>
              <div className="shrink-0">
                {g.is_member ? (
                  <Link
                    href={`/groups/${g.id}`}
                    className="rounded-lg bg-indigo-600 px-3 py-1.5 text-xs font-medium text-white shadow-sm hover:bg-indigo-700"
                  >
                    Open
                  </Link>
                ) : g.is_pending ? (
                  <span className="inline-block rounded-lg border border-slate-300 bg-slate-50 px-3 py-1.5 text-xs font-medium text-slate-500">
                    Pending
                  </span>
                ) : (
                  <button
                    onClick={() => onJoin(g.id)}
                    className="rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-xs font-medium text-slate-700 hover:bg-slate-50"
                  >
                    Join
                  </button>
                )}
              </div>
            </li>
          ))}
        </ul>
      )}

      {groups.length > 0 && hasMore ? (
        <InfiniteScrollSentinel onIntersect={loadMore} enabled={!loading} />
      ) : null}
      {loading && groups.length > 0 ? (
        <p className="text-center text-xs text-slate-400">Loading more…</p>
      ) : null}
    </div>
  );
}
