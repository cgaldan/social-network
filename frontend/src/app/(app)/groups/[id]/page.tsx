"use client";

import Link from "next/link";
import { useCallback, useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { PostCard } from "@/components/PostCard";
import { PostForm } from "@/components/PostForm";
import { UserPicker } from "@/components/UserPicker";
import { api, ApiError } from "@/lib/api";
import { formatDateTime } from "@/lib/format";
import type { CreatePostPayload, Group, GroupEvent, Post, RsvpResponse, UserSummary } from "@/types/api";

export default function GroupDetailPage() {
  const params = useParams<{ id: string }>();
  const groupId = Number(params.id);

  const [group, setGroup] = useState<Group | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [events, setEvents] = useState<GroupEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [tab, setTab] = useState<"posts" | "events">("posts");
  const [joining, setJoining] = useState(false);

  const [eventTitle, setEventTitle] = useState("");
  const [eventDesc, setEventDesc] = useState("");
  const [eventDate, setEventDate] = useState("");

  const [inviteStatus, setInviteStatus] = useState<string | null>(null);

  const load = useCallback(async () => {
    if (!Number.isFinite(groupId)) return;
    setLoading(true);
    setError(null);
    try {
      const groupRes = await api.groups.get(groupId);
      setGroup(groupRes.group);
      if (!groupRes.group.is_member) {
        setPosts([]);
        setEvents([]);
        return;
      }
      const [postsRes, eventsRes] = await Promise.all([
        api.groups.posts(groupId, { limit: 50 }),
        api.groups.events(groupId, { limit: 50 }),
      ]);
      setPosts(postsRes.posts ?? []);
      setEvents(eventsRes.events ?? []);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to load group");
    } finally {
      setLoading(false);
    }
  }, [groupId]);

  useEffect(() => {
    void load();
  }, [load]);

  const onJoinFromDetail = async () => {
    if (!group) return;
    setJoining(true);
    try {
      await api.groups.join(group.id);
      setGroup({ ...group, is_pending: true });
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to send join request");
    } finally {
      setJoining(false);
    }
  };

  const onCreatePost = async (payload: CreatePostPayload) => {
    const res = await api.groups.createPost(groupId, payload);
    setPosts((prev) => [res.post, ...prev]);
  };

  const onCreateEvent = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const startsAt = new Date(eventDate).toISOString();
      const res = await api.groups.createEvent(groupId, {
        title: eventTitle,
        description: eventDesc,
        starts_at: startsAt,
      });
      setEvents((prev) => [res.event, ...prev]);
      setEventTitle("");
      setEventDesc("");
      setEventDate("");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to create event");
    }
  };

  const onRsvp = async (eventId: number, response: RsvpResponse) => {
    setError(null);
    setEvents((prev) =>
      prev.map((ev) => {
        if (ev.id !== eventId) return ev;
        const next = { ...ev };
        if (next.my_response === "going") next.going_count = Math.max(0, next.going_count - 1);
        else if (next.my_response === "maybe") next.maybe_count = Math.max(0, next.maybe_count - 1);
        else if (next.my_response === "not_going")
          next.not_going_count = Math.max(0, next.not_going_count - 1);
        if (response === "going") next.going_count += 1;
        else if (response === "maybe") next.maybe_count += 1;
        else next.not_going_count += 1;
        next.my_response = response;
        return next;
      }),
    );
    try {
      await api.groups.rsvp(groupId, eventId, response);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to RSVP");
      void load();
    }
  };

  const onInvite = async (target: UserSummary) => {
    setInviteStatus(null);
    try {
      await api.groups.invite(groupId, target.id);
      setInviteStatus(`Invitation sent to @${target.nickname}`);
    } catch (err) {
      setInviteStatus(err instanceof ApiError ? err.message : "Failed to send invite");
    }
  };

  if (loading) {
    return <p className="text-sm text-slate-500">Loading…</p>;
  }

  if (!group) {
    return (
      <div className="space-y-4">
        <p className="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">
          {error ?? "Group not found"}
        </p>
        <Link href="/groups" className="text-sm text-indigo-600 hover:underline">
          ← All groups
        </Link>
      </div>
    );
  }

  if (!group.is_member) {
    return (
      <div className="space-y-6">
        <Link href="/groups" className="text-sm text-indigo-600 hover:underline">
          ← All groups
        </Link>
        <div className="rounded-xl border border-slate-200 bg-white p-6 shadow-sm">
          <h1 className="text-2xl font-semibold text-slate-900">{group.name}</h1>
          <p className="mt-2 text-sm text-slate-700">{group.description}</p>
          <p className="mt-3 text-xs text-slate-400">
            You&apos;re not a member of this group. Posts and events are only visible to members.
          </p>
          <div className="mt-4">
            {group.is_pending ? (
              <span className="inline-block rounded-lg border border-slate-300 bg-slate-50 px-3 py-1.5 text-sm font-medium text-slate-500">
                Request pending
              </span>
            ) : (
              <button
                onClick={onJoinFromDetail}
                disabled={joining}
                className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm transition hover:bg-indigo-700 disabled:opacity-60"
              >
                {joining ? "Sending…" : "Request to join"}
              </button>
            )}
          </div>
          {error ? (
            <p className="mt-3 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p>
          ) : null}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <Link href="/groups" className="text-sm text-indigo-600 hover:underline">
        ← All groups
      </Link>

      <h1 className="text-2xl font-semibold text-slate-900">{group.name}</h1>
      <p className="text-sm text-slate-600">{group.description}</p>

      <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
        <p className="mb-2 text-xs font-medium uppercase tracking-wide text-slate-500">
          Invite a user
        </p>
        <UserPicker onSelect={onInvite} placeholder="Search users to invite…" />
        {inviteStatus ? <p className="mt-2 text-xs text-slate-500">{inviteStatus}</p> : null}
      </div>

      <div className="flex gap-2 border-b border-slate-200">
        <TabButton active={tab === "posts"} onClick={() => setTab("posts")}>
          Posts
        </TabButton>
        <TabButton active={tab === "events"} onClick={() => setTab("events")}>
          Events
        </TabButton>
      </div>

      {error ? (
        <p className="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p>
      ) : null}

      {tab === "posts" ? (
        <div className="space-y-4">
          <PostForm onSubmit={onCreatePost} submitLabel="Post to group" inGroup />
          {posts.length === 0 ? (
            <p className="rounded-xl border border-dashed border-slate-300 bg-white p-6 text-center text-sm text-slate-500">
              No posts in this group yet.
            </p>
          ) : (
            posts.map((p) => <PostCard key={p.id} post={p} />)
          )}
        </div>
      ) : null}

      {tab === "events" ? (
        <div className="space-y-4">
          <form
            onSubmit={onCreateEvent}
            className="space-y-3 rounded-xl border border-slate-200 bg-white p-5 shadow-sm"
          >
            <h2 className="text-base font-semibold text-slate-900">Create event</h2>
            <input
              type="text"
              value={eventTitle}
              onChange={(e) => setEventTitle(e.target.value)}
              placeholder="Title"
              required
              className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            />
            <textarea
              value={eventDesc}
              onChange={(e) => setEventDesc(e.target.value)}
              rows={2}
              placeholder="Description"
              required
              className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            />
            <input
              type="datetime-local"
              value={eventDate}
              onChange={(e) => setEventDate(e.target.value)}
              required
              className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            />
            <div className="flex justify-end">
              <button
                type="submit"
                className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm transition hover:bg-indigo-700"
              >
                Create event
              </button>
            </div>
          </form>

          {events.length === 0 ? (
            <p className="rounded-xl border border-dashed border-slate-300 bg-white p-6 text-center text-sm text-slate-500">
              No events scheduled.
            </p>
          ) : (
            <ul className="space-y-3">
              {events.map((ev) => (
                <li
                  key={ev.id}
                  className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm"
                >
                  <h3 className="text-base font-semibold text-slate-900">{ev.title}</h3>
                  <p className="text-xs text-slate-500">{formatDateTime(ev.starts_at)}</p>
                  <p className="mt-2 whitespace-pre-wrap text-sm text-slate-700">
                    {ev.description}
                  </p>
                  <div className="mt-3 flex flex-wrap items-center gap-2">
                    {(
                      [
                        { value: "going" as RsvpResponse, label: "Going", count: ev.going_count },
                        { value: "maybe" as RsvpResponse, label: "Maybe", count: ev.maybe_count },
                        {
                          value: "not_going" as RsvpResponse,
                          label: "Not going",
                          count: ev.not_going_count,
                        },
                      ]
                    ).map(({ value, label, count }) => {
                      const active = ev.my_response === value;
                      return (
                        <button
                          key={value}
                          onClick={() => onRsvp(ev.id, value)}
                          className={`flex items-center gap-1.5 rounded-lg border px-3 py-1.5 text-xs font-medium transition ${
                            active
                              ? "border-indigo-600 bg-indigo-600 text-white shadow-sm"
                              : "border-slate-300 bg-white text-slate-700 hover:bg-slate-50"
                          }`}
                        >
                          <span>{label}</span>
                          <span
                            className={`rounded-full px-1.5 text-[10px] font-semibold ${
                              active ? "bg-white/20 text-white" : "bg-slate-100 text-slate-600"
                            }`}
                          >
                            {count}
                          </span>
                        </button>
                      );
                    })}
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>
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
