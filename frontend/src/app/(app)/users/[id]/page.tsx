"use client";

import { useCallback, useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { Avatar } from "@/components/Avatar";
import { AvatarUpload } from "@/components/AvatarUpload";
import { PostCard } from "@/components/PostCard";
import { useAuth } from "@/components/AuthProvider";
import { api, ApiError } from "@/lib/api";
import { formatDate } from "@/lib/format";
import type { FollowStatus, Post, RegisterPayload, UserProfile, UserSummary } from "@/types/api";

type Tab = "posts" | "followers" | "following";

export default function UserProfilePage() {
  const params = useParams<{ id: string }>();
  const id = Number(params.id);
  const router = useRouter();
  const { user: me, updateUser, logout } = useAuth();

  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [acting, setActing] = useState(false);

  const [tab, setTab] = useState<Tab>("posts");
  const [posts, setPosts] = useState<Post[]>([]);
  const [followers, setFollowers] = useState<UserSummary[]>([]);
  const [following, setFollowing] = useState<UserSummary[]>([]);
  const [tabLoading, setTabLoading] = useState(false);
  const [tabError, setTabError] = useState<string | null>(null);

  const [editing, setEditing] = useState(false);

  const loadProfile = useCallback(async () => {
    if (!Number.isFinite(id)) return;
    setLoading(true);
    setError(null);
    try {
      const res = await api.users.get(id);
      setProfile(res.user);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to load profile");
    } finally {
      setLoading(false);
    }
  }, [id]);

  useEffect(() => {
    void loadProfile();
  }, [loadProfile]);

  useEffect(() => {
    if (!profile?.can_view) return;
    let cancelled = false;
    setTabLoading(true);
    setTabError(null);

    const fetcher = (() => {
      if (tab === "posts") return api.users.posts(id, { limit: 50 });
      if (tab === "followers") return api.users.followers(id, { limit: 100 });
      return api.users.following(id, { limit: 100 });
    })();

    fetcher
      .then((res) => {
        if (cancelled) return;
        if (tab === "posts") setPosts((res as { posts: Post[] }).posts ?? []);
        else if (tab === "followers") setFollowers((res as { users: UserSummary[] }).users ?? []);
        else setFollowing((res as { users: UserSummary[] }).users ?? []);
      })
      .catch((err) => {
        if (cancelled) return;
        setTabError(err instanceof ApiError ? err.message : "Failed to load");
      })
      .finally(() => {
        if (!cancelled) setTabLoading(false);
      });

    return () => {
      cancelled = true;
    };
  }, [tab, id, profile?.can_view]);

  const updateFollowState = (status: FollowStatus, deltaFollowers = 0) => {
    setProfile((prev) =>
      prev
        ? {
            ...prev,
            follow_status: status,
            followers_count: Math.max(0, prev.followers_count + deltaFollowers),
            can_view: prev.is_public || status === "self" || status === "following",
          }
        : prev,
    );
  };

  const onFollow = async () => {
    if (!profile) return;
    setActing(true);
    try {
      const res = await api.follow.follow(profile.id);
      const status: FollowStatus = res.status === "accepted" ? "following" : "pending";
      const delta = res.status === "accepted" ? 1 : 0;
      updateFollowState(status, delta);
      if (status === "following") void loadProfile();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Follow failed");
    } finally {
      setActing(false);
    }
  };

  const onUnfollow = async () => {
    if (!profile) return;
    if (!confirm(`Unfollow @${profile.nickname}?`)) return;
    setActing(true);
    try {
      await api.follow.unfollow(profile.id);
      updateFollowState("none", profile.follow_status === "following" ? -1 : 0);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Unfollow failed");
    } finally {
      setActing(false);
    }
  };

  const onMessage = async () => {
    if (!profile) return;
    try {
      await api.conversations.direct(profile.id);
      router.push("/messages");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to open conversation");
    }
  };

  const onTogglePrivacy = async () => {
    if (!profile || profile.follow_status !== "self") return;
    setActing(true);
    try {
      const res = await api.auth.update({ is_public: !profile.is_public });
      updateUser(res.user);
      setProfile((prev) => (prev ? { ...prev, is_public: res.user.is_public } : prev));
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to update privacy");
    } finally {
      setActing(false);
    }
  };

  const onSaveEdit = async (payload: Partial<RegisterPayload>) => {
    try {
      const res = await api.auth.update(payload);
      updateUser(res.user);
      setEditing(false);
      void loadProfile();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to update profile");
    }
  };

  const onDeleteAccount = async () => {
    if (!profile || profile.follow_status !== "self") return;
    if (!confirm("Permanently delete your account? This cannot be undone.")) return;
    try {
      await api.auth.remove();
    } catch {
      // continue locally
    }
    await logout();
    router.replace("/login");
  };

  if (loading) return <p className="text-sm text-slate-500">Loading…</p>;

  if (!profile) {
    return (
      <div className="space-y-4">
        <p className="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">
          {error ?? "User not found"}
        </p>
        <Link href="/feed" className="text-sm text-indigo-600 hover:underline">
          ← Back to feed
        </Link>
      </div>
    );
  }

  const isSelf = profile.follow_status === "self";
  const followBtn = (() => {
    if (isSelf) return null;
    if (profile.follow_status === "following") {
      return (
        <button
          onClick={onUnfollow}
          disabled={acting}
          className="rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 hover:bg-slate-50 disabled:opacity-60"
        >
          Following
        </button>
      );
    }
    if (profile.follow_status === "pending") {
      return (
        <button
          onClick={onUnfollow}
          disabled={acting}
          className="rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 hover:bg-slate-50 disabled:opacity-60"
        >
          Requested
        </button>
      );
    }
    return (
      <button
        onClick={onFollow}
        disabled={acting}
        className="rounded-lg bg-indigo-600 px-3 py-1.5 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 disabled:opacity-60"
      >
        Follow
      </button>
    );
  })();

  return (
    <div className="space-y-6">
      <Link href="/feed" className="text-sm text-indigo-600 hover:underline">
        ← Back
      </Link>

      <section className="rounded-xl border border-slate-200 bg-white p-6 shadow-sm">
        <div className="flex items-start gap-4">
          <Avatar
            src={profile.avatar_path}
            firstName={profile.first_name}
            lastName={profile.last_name}
            nickname={profile.nickname}
            size={72}
          />
          <div className="flex-1 min-w-0">
            <h1 className="text-xl font-semibold text-slate-900 truncate">
              {profile.first_name} {profile.last_name}
            </h1>
            <p className="text-sm text-slate-500">@{profile.nickname}</p>
            <p className="mt-1 text-xs text-slate-400">
              Joined {formatDate(profile.created_at)}
              {" · "}
              {profile.is_public ? "Public profile" : "Private profile"}
              {profile.is_online ? " · online" : ""}
            </p>
          </div>
          <div className="flex shrink-0 flex-wrap items-start justify-end gap-2">
            {isSelf ? (
              <>
                <button
                  onClick={onTogglePrivacy}
                  disabled={acting}
                  className="rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 hover:bg-slate-50 disabled:opacity-60"
                >
                  {profile.is_public ? "Make private" : "Make public"}
                </button>
                <button
                  onClick={() => setEditing((v) => !v)}
                  className="rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 hover:bg-slate-50"
                >
                  {editing ? "Cancel" : "Edit"}
                </button>
              </>
            ) : (
              <>
                {followBtn}
                <button
                  onClick={onMessage}
                  className="rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 hover:bg-slate-50"
                >
                  Message
                </button>
              </>
            )}
          </div>
        </div>

        {profile.can_view ? (
          <>
            <dl className="mt-5 grid grid-cols-1 gap-x-6 gap-y-2 text-sm sm:grid-cols-2">
              {profile.email ? (
                <Field label="Email" value={profile.email} />
              ) : null}
              {profile.date_of_birth ? (
                <Field label="Date of birth" value={formatDate(profile.date_of_birth)} />
              ) : null}
              {profile.gender ? <Field label="Gender" value={profile.gender} /> : null}
            </dl>

            {profile.about_me ? (
              <div className="mt-4">
                <p className="text-xs font-medium uppercase tracking-wide text-slate-500">
                  About
                </p>
                <p className="mt-1 whitespace-pre-wrap text-sm text-slate-700">
                  {profile.about_me}
                </p>
              </div>
            ) : null}

            <div className="mt-4 flex gap-6 text-sm">
              <button
                onClick={() => setTab("followers")}
                className="hover:text-indigo-600"
              >
                <span className="font-semibold text-slate-900">
                  {profile.followers_count}
                </span>{" "}
                <span className="text-slate-500">Followers</span>
              </button>
              <button
                onClick={() => setTab("following")}
                className="hover:text-indigo-600"
              >
                <span className="font-semibold text-slate-900">
                  {profile.following_count}
                </span>{" "}
                <span className="text-slate-500">Following</span>
              </button>
            </div>
          </>
        ) : null}
      </section>

      {error ? (
        <p className="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p>
      ) : null}

      {isSelf && editing ? (
        <EditForm initial={me!} onSave={onSaveEdit} />
      ) : null}

      {profile.can_view ? (
        <>
          <div className="flex gap-2 border-b border-slate-200">
            <TabButton active={tab === "posts"} onClick={() => setTab("posts")}>
              Posts
            </TabButton>
            <TabButton active={tab === "followers"} onClick={() => setTab("followers")}>
              Followers
            </TabButton>
            <TabButton active={tab === "following"} onClick={() => setTab("following")}>
              Following
            </TabButton>
          </div>

          {tabError ? (
            <p className="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{tabError}</p>
          ) : null}

          {tabLoading ? (
            <p className="text-sm text-slate-500">Loading…</p>
          ) : tab === "posts" ? (
            posts.length === 0 ? (
              <p className="rounded-xl border border-dashed border-slate-300 bg-white p-6 text-center text-sm text-slate-500">
                No posts yet.
              </p>
            ) : (
              <div className="space-y-4">
                {posts.map((p) => (
                  <PostCard key={p.id} post={p} />
                ))}
              </div>
            )
          ) : (
            <UserList users={tab === "followers" ? followers : following} emptyText={tab === "followers" ? "No followers yet." : "Not following anyone yet."} />
          )}
        </>
      ) : (
        <p className="rounded-xl border border-dashed border-slate-300 bg-white p-6 text-center text-sm text-slate-500">
          This profile is private. Follow to see their posts, followers, and details.
        </p>
      )}

      {isSelf ? (
        <div className="rounded-xl border border-red-200 bg-red-50 p-4">
          <h3 className="text-sm font-semibold text-red-900">Danger zone</h3>
          <p className="mt-1 text-xs text-red-700">
            Deleting your account is permanent.
          </p>
          <button
            onClick={onDeleteAccount}
            className="mt-3 rounded-lg border border-red-300 bg-white px-3 py-1.5 text-sm font-medium text-red-700 hover:bg-red-100"
          >
            Delete account
          </button>
        </div>
      ) : null}
    </div>
  );
}

function Field({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <dt className="text-xs font-medium uppercase tracking-wide text-slate-500">{label}</dt>
      <dd className="mt-0.5 text-slate-800">{value}</dd>
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

function UserList({ users, emptyText }: { users: UserSummary[]; emptyText: string }) {
  if (users.length === 0) {
    return (
      <p className="rounded-xl border border-dashed border-slate-300 bg-white p-6 text-center text-sm text-slate-500">
        {emptyText}
      </p>
    );
  }
  return (
    <ul className="space-y-2">
      {users.map((u) => (
        <li key={u.id}>
          <Link
            href={`/users/${u.id}`}
            className="flex items-center gap-3 rounded-xl border border-slate-200 bg-white p-3 shadow-sm hover:bg-slate-50"
          >
            <Avatar
              src={u.avatar_path}
              firstName={u.first_name}
              lastName={u.last_name}
              nickname={u.nickname}
              size={40}
            />
            <div className="min-w-0">
              <p className="truncate text-sm font-medium text-slate-900">
                {u.first_name} {u.last_name}
              </p>
              <p className="truncate text-xs text-slate-500">
                @{u.nickname}
                {u.is_online ? " · online" : ""}
              </p>
            </div>
          </Link>
        </li>
      ))}
    </ul>
  );
}

function EditForm({
  initial,
  onSave,
}: {
  initial: {
    email: string;
    first_name: string;
    last_name: string;
    nickname: string;
    gender: string;
    avatar_path: string;
    about_me: string;
    is_public: boolean;
  };
  onSave: (payload: Partial<RegisterPayload>) => Promise<void>;
}) {
  const [first_name, setFirst] = useState(initial.first_name);
  const [last_name, setLast] = useState(initial.last_name);
  const [nickname, setNickname] = useState(initial.nickname);
  const [email, setEmail] = useState(initial.email);
  const [avatar_path, setAvatar] = useState(initial.avatar_path);
  const [about_me, setAbout] = useState(initial.about_me);
  const [saving, setSaving] = useState(false);

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      await onSave({ first_name, last_name, nickname, email, avatar_path, about_me });
    } finally {
      setSaving(false);
    }
  };

  return (
    <form
      onSubmit={submit}
      className="space-y-3 rounded-xl border border-slate-200 bg-white p-5 shadow-sm"
    >
      <div className="grid grid-cols-2 gap-3">
        <Input label="First name" value={first_name} onChange={setFirst} />
        <Input label="Last name" value={last_name} onChange={setLast} />
      </div>
      <Input label="Nickname" value={nickname} onChange={setNickname} />
      <Input label="Email" type="email" value={email} onChange={setEmail} />
      <AvatarUpload
        value={avatar_path}
        onChange={setAvatar}
        firstName={first_name}
        lastName={last_name}
        nickname={nickname}
        label="Profile photo"
      />
      <label className="block">
        <span className="text-sm font-medium text-slate-700">About</span>
        <textarea
          value={about_me}
          onChange={(e) => setAbout(e.target.value)}
          rows={3}
          className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
        />
      </label>
      <div className="flex justify-end">
        <button
          type="submit"
          disabled={saving}
          className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm transition hover:bg-indigo-700 disabled:opacity-60"
        >
          {saving ? "Saving…" : "Save"}
        </button>
      </div>
    </form>
  );
}

function Input({
  label,
  value,
  onChange,
  type = "text",
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  type?: string;
}) {
  return (
    <label className="block">
      <span className="text-sm font-medium text-slate-700">{label}</span>
      <input
        type={type}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
      />
    </label>
  );
}
