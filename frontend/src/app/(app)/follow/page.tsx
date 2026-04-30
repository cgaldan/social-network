"use client";

import { FormEvent, useState } from "react";
import { FormMessage, SelectField, TextField } from "@/components/forms";
import { api } from "@/lib/api";

const actions = [
  { label: "Follow or request follow", value: "follow" },
  { label: "Accept follow request", value: "accept" },
  { label: "Decline follow request", value: "decline" },
  { label: "Unfollow user", value: "unfollow" },
  { label: "Remove follower", value: "remove" },
];

export default function FollowPage() {
  const [userId, setUserId] = useState("");
  const [action, setAction] = useState("follow");
  const [message, setMessage] = useState("");
  const [tone, setTone] = useState<"error" | "success">("success");

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const id = Number(userId);

    try {
      const response =
        action === "accept"
          ? await api.acceptFollow(id)
          : action === "decline"
            ? await api.declineFollow(id)
            : action === "unfollow"
              ? await api.unfollowUser(id)
              : action === "remove"
                ? await api.removeFollower(id)
                : await api.followUser(id);

      setTone("success");
      setMessage(response.message || response.status || "Follow action complete.");
    } catch (error) {
      setTone("error");
      setMessage(error instanceof Error ? error.message : "Follow action failed");
    }
  }

  return (
    <section className="mx-auto max-w-2xl rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
      <h1 className="text-3xl font-bold text-slate-950">Follow controls</h1>
      <p className="mt-3 text-slate-600">
        The backend exposes follow actions by user ID. A user search/list route
        is not available yet, so this screen uses explicit IDs.
      </p>
      <form className="mt-6 grid gap-5" onSubmit={onSubmit}>
        <TextField
          label="User ID"
          name="user_id"
          onChange={setUserId}
          required
          type="number"
          value={userId}
        />
        <SelectField
          label="Action"
          name="action"
          onChange={setAction}
          options={actions}
          value={action}
        />
        <FormMessage message={message} tone={tone} />
        <button
          className="rounded-xl bg-sky-600 px-4 py-3 font-semibold text-white transition hover:bg-sky-700"
          type="submit"
        >
          Run action
        </button>
      </form>
    </section>
  );
}
