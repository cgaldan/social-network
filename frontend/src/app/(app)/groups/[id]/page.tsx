"use client";

import Link from "next/link";
import { useCallback, useEffect, useRef, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { useAuth } from "@/components/AuthProvider";
import { Avatar } from "@/components/Avatar";
import { PostCard } from "@/components/PostCard";
import { PostForm } from "@/components/PostForm";
import { UserPicker } from "@/components/UserPicker";
import { InfiniteScrollSentinel } from "@/components/InfiniteScrollSentinel";
import { usePaginatedList } from "@/hooks/usePaginatedList";
import { api, ApiError } from "@/lib/api";
import { formatDateTime } from "@/lib/format";
import type { CreatePostPayload, Group, GroupEvent, Post, RsvpResponse, UserSummary } from "@/types/api";

const PAGE_SIZE = 20;

export default function GroupDetailPage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const groupId = Number(params.id);
  const { user } = useAuth();

  const [group, setGroup] = useState<Group | null>(null);
  const [groupLoading, setGroupLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [tab, setTab] = useState<"posts" | "events">("posts");
  const [joining, setJoining] = useState(false);

  const [editingGroup, setEditingGroup] = useState(false);
  const [editTitle, setEditTitle] = useState("");
  const [editDescription, setEditDescription] = useState("");
  const [savingGroup, setSavingGroup] = useState(false);

  const [eventTitle, setEventTitle] = useState("");
  const [eventDesc, setEventDesc] = useState("");
  const [eventDate, setEventDate] = useState("");

  const [editingEventId, setEditingEventId] = useState<number | null>(null);
  const [editEventTitle, setEditEventTitle] = useState("");
  const [editEventDesc, setEditEventDesc] = useState("");
  const [editEventDate, setEditEventDate] = useState("");
  const [savingEvent, setSavingEvent] = useState(false);

  const [inviteStatus, setInviteStatus] = useState<
    { kind: "success" | "error"; message: string } | null
  >(null);
  const [uploadingAvatar, setUploadingAvatar] = useState(false);
  const avatarInputRef = useRef<HTMLInputElement>(null);

  const onPickGroupAvatar = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file || !group) return;
    setError(null);
    setUploadingAvatar(true);
    try {
      const uploaded = await api.uploads.create(file);
      const res = await api.groups.updateAvatar(group.id, uploaded.url);
      setGroup(res.group);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to update group avatar");
    } finally {
      setUploadingAvatar(false);
      if (avatarInputRef.current) avatarInputRef.current.value = "";
    }
  };

  useEffect(() => {
    if (!Number.isFinite(groupId)) return;
    let cancelled = false;
    setGroupLoading(true);
    api.groups
      .get(groupId)
      .then((res) => {
        if (!cancelled) setGroup(res.group);
      })
      .catch((err) => {
        if (!cancelled) setError(err instanceof ApiError ? err.message : "Failed to load group");
      })
      .finally(() => {
        if (!cancelled) setGroupLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [groupId]);

  const isMember = group?.is_member === true;

  const postsFetcher = useCallback(
    async ({ limit, offset }: { limit: number; offset: number }) => {
      if (!Number.isFinite(groupId) || !isMember) return { items: [], hasMore: false };
      const res = await api.groups.posts(groupId, { limit, offset });
      return { items: res.posts ?? [], hasMore: res.has_more };
    },
    [groupId, isMember],
  );
  const {
    items: posts,
    hasMore: hasMorePosts,
    loading: loadingPosts,
    loadMore: loadMorePosts,
    setItems: setPosts,
  } = usePaginatedList<Post>(postsFetcher, PAGE_SIZE);

  const eventsFetcher = useCallback(
    async ({ limit, offset }: { limit: number; offset: number }) => {
      if (!Number.isFinite(groupId) || !isMember) return { items: [], hasMore: false };
      const res = await api.groups.events(groupId, { limit, offset });
      return { items: res.events ?? [], hasMore: res.has_more };
    },
    [groupId, isMember],
  );
  const {
    items: events,
    hasMore: hasMoreEvents,
    loading: loadingEvents,
    loadMore: loadMoreEvents,
    setItems: setEvents,
    reset: resetEvents,
  } = usePaginatedList<GroupEvent>(eventsFetcher, PAGE_SIZE);

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

  const onStartEditEvent = (ev: GroupEvent) => {
    setEditingEventId(ev.id);
    setEditEventTitle(ev.title);
    setEditEventDesc(ev.description);
    const d = new Date(ev.starts_at);
    const pad = (n: number) => String(n).padStart(2, "0");
    setEditEventDate(
      `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`,
    );
  };

  const onCancelEditEvent = () => {
    setEditingEventId(null);
    setEditEventTitle("");
    setEditEventDesc("");
    setEditEventDate("");
  };

  const onSaveEditEvent = async (eventId: number) => {
    setSavingEvent(true);
    try {
      const startsAt = new Date(editEventDate).toISOString();
      const res = await api.groups.updateEvent(groupId, eventId, {
        title: editEventTitle,
        description: editEventDesc,
        starts_at: startsAt,
      });
      setEvents((prev) => prev.map((ev) => (ev.id === eventId ? res.event : ev)));
      onCancelEditEvent();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to update event");
    } finally {
      setSavingEvent(false);
    }
  };

  const onDeleteEvent = async (eventId: number) => {
    if (!confirm("Delete this event?")) return;
    try {
      await api.groups.removeEvent(groupId, eventId);
      setEvents((prev) => prev.filter((ev) => ev.id !== eventId));
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to delete event");
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
      resetEvents();
    }
  };

  const onInvite = async (target: UserSummary) => {
    setInviteStatus(null);
    try {
      await api.groups.invite(groupId, target.id);
      setInviteStatus({ kind: "success", message: `Invitation sent to @${target.nickname}` });
    } catch (err) {
      setInviteStatus({
        kind: "error",
        message: err instanceof ApiError ? err.message : "Failed to send invite",
      });
    }
  };

  const onStartEditGroup = () => {
    if (!group) return;
    setEditTitle(group.name);
    setEditDescription(group.description);
    setEditingGroup(true);
  };

  const onCancelEditGroup = () => {
    setEditingGroup(false);
    setError(null);
  };

  const onSaveGroup = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!group) return;
    setSavingGroup(true);
    try {
      const res = await api.groups.update(group.id, {
        title: editTitle,
        description: editDescription,
      });
      setGroup(res.group);
      setEditingGroup(false);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to update group");
    } finally {
      setSavingGroup(false);
    }
  };

  const onDeleteGroup = async () => {
    if (!group) return;
    if (!confirm("Delete this group? This cannot be undone.")) return;
    try {
      await api.groups.remove(group.id);
      router.replace("/groups");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to delete group");
    }
  };

  const onLeaveGroup = async () => {
    if (!group) return;
    if (!confirm("Leave this group?")) return;
    try {
      await api.groups.leave(group.id);
      router.replace("/groups");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to leave group");
    }
  };

  if (groupLoading) {
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

  const isCreator = user?.id === group.creator_id;

  return (
    <div className="space-y-6">
      <Link href="/groups" className="text-sm text-indigo-600 hover:underline">
        ← All groups
      </Link>

      <div className="flex items-start gap-4">
        <Avatar
          src={group.avatar_path}
          nickname={group.name}
          size={64}
        />
        <div className="flex-1">
          {editingGroup ? (
            <form onSubmit={onSaveGroup} className="space-y-3">
              <input
                type="text"
                name="editGroupTitle"
                value={editTitle}
                onChange={(e) => setEditTitle(e.target.value)}
                placeholder="Group name"
                required
                minLength={3}
                className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
              />
              <textarea
                name="editGroupDescription"
                value={editDescription}
                onChange={(e) => setEditDescription(e.target.value)}
                placeholder="Description"
                rows={3}
                required
                className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
              />
              <div className="flex justify-end gap-2">
                <button
                  type="button"
                  onClick={onCancelEditGroup}
                  className="rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 hover:bg-slate-50"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={savingGroup}
                  className="rounded-lg bg-indigo-600 px-3 py-1.5 text-sm font-medium text-white shadow-sm transition hover:bg-indigo-700 disabled:opacity-60"
                >
                  {savingGroup ? "Saving…" : "Save"}
                </button>
              </div>
            </form>
          ) : (
            <>
              <h1 className="text-2xl font-semibold text-slate-900">{group.name}</h1>
              <p className="text-sm text-slate-600">{group.description}</p>
            </>
          )}
          {isCreator && !editingGroup ? (
            <div className="mt-2">
              <input
                ref={avatarInputRef}
                type="file"
                name="groupAvatar"
                accept="image/jpeg,image/png,image/gif,image/webp"
                onChange={onPickGroupAvatar}
                disabled={uploadingAvatar}
                className="block text-xs text-slate-600 file:mr-2 file:rounded-md file:border-0 file:bg-indigo-50 file:px-2 file:py-1 file:text-xs file:font-medium file:text-indigo-700 hover:file:bg-indigo-100"
              />
              <p className="mt-1 text-[11px] text-slate-400">
                {uploadingAvatar
                  ? "Uploading…"
                  : group.avatar_path
                    ? "Pick a new image to replace the group avatar."
                    : "Add a group avatar so members can recognize this chat."}
              </p>
            </div>
          ) : null}
        </div>
      </div>

      {!editingGroup ? (
        <div className="flex flex-wrap gap-3 text-sm">
          <Link
            href={`/messages?conversation=${group.conversation_id}`}
            className="rounded-lg border border-indigo-300 bg-indigo-50 px-3 py-1.5 text-sm font-medium text-indigo-700 hover:bg-indigo-100"
          >
            💬 Open chat
          </Link>
          {isCreator ? (
            <>
              <button
                onClick={onStartEditGroup}
                className="rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 hover:bg-slate-50"
              >
                Edit group
              </button>
              <button
                onClick={onDeleteGroup}
                className="rounded-lg border border-red-300 bg-white px-3 py-1.5 text-sm font-medium text-red-600 hover:bg-red-50"
              >
                Delete group
              </button>
            </>
          ) : (
            <button
              onClick={onLeaveGroup}
              className="rounded-lg border border-red-300 bg-white px-3 py-1.5 text-sm font-medium text-red-600 hover:bg-red-50"
            >
              Leave group
            </button>
          )}
        </div>
      ) : null}

      <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
        <p className="mb-2 text-xs font-medium uppercase tracking-wide text-slate-500">
          Invite a user
        </p>
        <UserPicker onSelect={onInvite} placeholder="Search users to invite…" />
        {inviteStatus ? (
          <p
            className={
              inviteStatus.kind === "error"
                ? "mt-2 rounded-md bg-red-50 px-2 py-1 text-xs font-medium text-red-700"
                : "mt-2 text-xs text-emerald-600"
            }
          >
            {inviteStatus.message}
          </p>
        ) : null}
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
          {loadingPosts && posts.length === 0 ? (
            <p className="text-sm text-slate-500">Loading…</p>
          ) : posts.length === 0 ? (
            <p className="rounded-xl border border-dashed border-slate-300 bg-white p-6 text-center text-sm text-slate-500">
              No posts in this group yet.
            </p>
          ) : (
            posts.map((p) => <PostCard key={p.id} post={p} />)
          )}
          {posts.length > 0 && hasMorePosts ? (
            <InfiniteScrollSentinel onIntersect={loadMorePosts} enabled={!loadingPosts} />
          ) : null}
          {loadingPosts && posts.length > 0 ? (
            <p className="text-center text-xs text-slate-400">Loading more…</p>
          ) : null}
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
              name="eventTitle"
              value={eventTitle}
              onChange={(e) => setEventTitle(e.target.value)}
              placeholder="Title"
              required
              className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            />
            <textarea
              name="eventDescription"
              value={eventDesc}
              onChange={(e) => setEventDesc(e.target.value)}
              rows={2}
              placeholder="Description"
              required
              className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            />
            <input
              type="datetime-local"
              name="eventStartsAt"
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

          {loadingEvents && events.length === 0 ? (
            <p className="text-sm text-slate-500">Loading…</p>
          ) : events.length === 0 ? (
            <p className="rounded-xl border border-dashed border-slate-300 bg-white p-6 text-center text-sm text-slate-500">
              No events scheduled.
            </p>
          ) : (
            <ul className="space-y-3">
              {events.map((ev) => {
                const isEditing = editingEventId === ev.id;
                const isEventCreator = user?.id === ev.creator_id;
                return (
                  <li
                    key={ev.id}
                    className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm"
                  >
                    {isEditing ? (
                      <div className="space-y-3">
                        <input
                          type="text"
                          name="editEventTitle"
                          value={editEventTitle}
                          onChange={(e) => setEditEventTitle(e.target.value)}
                          placeholder="Title"
                          required
                          className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
                        />
                        <textarea
                          name="editEventDescription"
                          value={editEventDesc}
                          onChange={(e) => setEditEventDesc(e.target.value)}
                          rows={2}
                          placeholder="Description"
                          required
                          className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
                        />
                        <input
                          type="datetime-local"
                          name="editEventStartsAt"
                          value={editEventDate}
                          onChange={(e) => setEditEventDate(e.target.value)}
                          required
                          className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
                        />
                        <div className="flex justify-end gap-2">
                          <button
                            type="button"
                            onClick={onCancelEditEvent}
                            className="rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 hover:bg-slate-50"
                          >
                            Cancel
                          </button>
                          <button
                            type="button"
                            onClick={() => onSaveEditEvent(ev.id)}
                            disabled={savingEvent}
                            className="rounded-lg bg-indigo-600 px-3 py-1.5 text-sm font-medium text-white shadow-sm transition hover:bg-indigo-700 disabled:opacity-60"
                          >
                            {savingEvent ? "Saving…" : "Save"}
                          </button>
                        </div>
                      </div>
                    ) : (
                      <>
                        <div className="flex items-start justify-between gap-3">
                          <div>
                            <h3 className="text-base font-semibold text-slate-900">{ev.title}</h3>
                            <p className="text-xs text-slate-500">{formatDateTime(ev.starts_at)}</p>
                          </div>
                          {isEventCreator ? (
                            <div className="flex gap-3 text-xs">
                              <button
                                onClick={() => onStartEditEvent(ev)}
                                className="text-indigo-600 hover:underline"
                              >
                                Edit
                              </button>
                              <button
                                onClick={() => onDeleteEvent(ev.id)}
                                className="text-red-600 hover:underline"
                              >
                                Delete
                              </button>
                            </div>
                          ) : null}
                        </div>
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
                      </>
                    )}
                  </li>
                );
              })}
            </ul>
          )}
          {events.length > 0 && hasMoreEvents ? (
            <InfiniteScrollSentinel onIntersect={loadMoreEvents} enabled={!loadingEvents} />
          ) : null}
          {loadingEvents && events.length > 0 ? (
            <p className="text-center text-xs text-slate-400">Loading more…</p>
          ) : null}
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
