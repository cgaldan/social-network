"use client";

import Link from "next/link";
import { useCallback, useState } from "react";
import { Avatar } from "@/components/Avatar";
import { UserPicker } from "@/components/UserPicker";
import { InfiniteScrollSentinel } from "@/components/InfiniteScrollSentinel";
import { usePaginatedList } from "@/hooks/usePaginatedList";
import { api, ApiError } from "@/lib/api";
import { formatRelative } from "@/lib/format";
import type { FollowRequestSummary, UserSummary } from "@/types/api";

type Tab = "incoming" | "outgoing";

const PAGE_SIZE = 20;

export default function FollowPage() {
  const [tab, setTab] = useState<Tab>("incoming");
  const [error, setError] = useState<string | null>(null);
  const [followStatus, setFollowStatus] = useState<string | null>(null);

  const fetcher = useCallback(
    async ({ limit, offset }: { limit: number; offset: number }) => {
      const res = await api.follow.requests(tab, { limit, offset });
      return { items: res.requests ?? [], hasMore: res.has_more };
    },
    [tab],
  );
  const { items: requests, hasMore, loading, error: fetchError, loadMore, setItems, reset } =
    usePaginatedList<FollowRequestSummary>(fetcher, PAGE_SIZE);

  const onFollow = async (target: UserSummary) => {
    setFollowStatus(null);
    try {
      const res = await api.follow.follow(target.id);
      setFollowStatus(`Followed @${target.nickname} — ${res.status}`);
      if (res.status === "pending" && tab === "outgoing") {
        reset();
      }
    } catch (err) {
      setFollowStatus(err instanceof ApiError ? err.message : "Follow failed");
    }
  };

  const act = async (requestId: number, kind: "accept" | "decline" | "cancel") => {
    try {
      if (kind === "accept") await api.follow.accept(requestId);
      else if (kind === "decline") await api.follow.decline(requestId);
      else if (kind === "cancel") {
        const req = requests.find((r) => r.id === requestId);
        if (req) await api.follow.unfollow(req.user.id);
      }
      setItems((prev) => prev.filter((r) => r.id !== requestId));
    } catch (err) {
      setError(err instanceof ApiError ? err.message : `Failed to ${kind}`);
    }
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold text-slate-900">Follow</h1>

      <section className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
        <h2 className="text-base font-semibold text-slate-900">Follow a user</h2>
        <p className="mt-1 text-xs text-slate-500">
          Search and pick a user. Public profiles auto-accept; private ones go to your outgoing list.
        </p>
        <div className="mt-3">
          <UserPicker onSelect={onFollow} placeholder="Search users by name or @nickname…" />
        </div>
        {followStatus ? <p className="mt-2 text-sm text-slate-600">{followStatus}</p> : null}
      </section>

      <div className="flex gap-2 border-b border-slate-200">
        <TabButton active={tab === "incoming"} onClick={() => setTab("incoming")}>
          Incoming
        </TabButton>
        <TabButton active={tab === "outgoing"} onClick={() => setTab("outgoing")}>
          Outgoing
        </TabButton>
      </div>

      {(error || fetchError) ? (
        <p className="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{error ?? fetchError}</p>
      ) : null}

      {loading && requests.length === 0 ? (
        <p className="text-sm text-slate-500">Loading…</p>
      ) : requests.length === 0 ? (
        <p className="rounded-xl border border-dashed border-slate-300 bg-white p-6 text-center text-sm text-slate-500">
          {tab === "incoming"
            ? "No one is waiting on you. Pending requests from people who tried to follow you will show here."
            : "You haven't sent any pending follow requests."}
        </p>
      ) : (
        <ul className="space-y-3">
          {requests.map((req) => (
            <li
              key={req.id}
              className="flex items-center justify-between rounded-xl border border-slate-200 bg-white p-4 shadow-sm"
            >
              <Link
                href={`/users/${req.user.id}`}
                className="flex items-center gap-3 min-w-0 rounded-md hover:bg-slate-50"
              >
                <Avatar
                  src={req.user.avatar_path}
                  firstName={req.user.first_name}
                  lastName={req.user.last_name}
                  nickname={req.user.nickname}
                  size={40}
                />
                <div className="min-w-0">
                  <p className="text-sm font-medium text-slate-900 truncate hover:text-indigo-600">
                    {req.user.first_name} {req.user.last_name}
                  </p>
                  <p className="text-xs text-slate-500 truncate">
                    @{req.user.nickname} · {formatRelative(req.created_at)}
                  </p>
                </div>
              </Link>
              <div className="flex gap-2 shrink-0">
                {tab === "incoming" ? (
                  <>
                    <button
                      onClick={() => act(req.id, "accept")}
                      className="rounded-lg bg-indigo-600 px-3 py-1.5 text-xs font-medium text-white shadow-sm hover:bg-indigo-700"
                    >
                      Accept
                    </button>
                    <button
                      onClick={() => act(req.id, "decline")}
                      className="rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-xs font-medium text-slate-700 hover:bg-slate-50"
                    >
                      Decline
                    </button>
                  </>
                ) : (
                  <button
                    onClick={() => act(req.id, "cancel")}
                    className="rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-xs font-medium text-slate-700 hover:bg-slate-50"
                  >
                    Cancel
                  </button>
                )}
              </div>
            </li>
          ))}
        </ul>
      )}

      {requests.length > 0 && hasMore ? (
        <InfiniteScrollSentinel onIntersect={loadMore} enabled={!loading} />
      ) : null}
      {loading && requests.length > 0 ? (
        <p className="text-center text-xs text-slate-400">Loading more…</p>
      ) : null}
    </div>
  );
}

function TabButton({
  active,
  onClick,
  children,
}: {
  active: boolean;
  onClick: () => void;
  children: React.ReactNode;
}) {
  return (
    <button
      onClick={onClick}
      className={`-mb-px border-b-2 px-3 py-2 text-sm font-medium transition ${
        active
          ? "border-indigo-600 text-indigo-600"
          : "border-transparent text-slate-600 hover:text-slate-900"
      }`}
    >
      {children}
    </button>
  );
}
