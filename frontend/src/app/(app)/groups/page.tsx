"use client";

import Link from "next/link";
import { FormEvent, useEffect, useState } from "react";
import { FormMessage, TextArea, TextField } from "@/components/forms";
import { api } from "@/lib/api";
import { formatDate } from "@/lib/format";
import type { Group } from "@/types/api";

export default function GroupsPage() {
  const [groups, setGroups] = useState<Group[]>([]);
  const [message, setMessage] = useState("");
  const [tone, setTone] = useState<"error" | "success">("success");
  const [createForm, setCreateForm] = useState({ title: "", description: "" });
  const [joinGroupId, setJoinGroupId] = useState("");
  const [invite, setInvite] = useState({ group_id: "", invitee_id: "" });
  const [decisionId, setDecisionId] = useState("");

  async function loadGroups() {
    try {
      const response = await api.listGroups({ limit: 50 });
      setGroups(response.groups ?? []);
    } catch (error) {
      setTone("error");
      setMessage(error instanceof Error ? error.message : "Could not load groups");
    }
  }

  useEffect(() => {
    void loadGroups();
  }, []);

  async function createGroup(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    try {
      await api.createGroup(createForm);
      setCreateForm({ title: "", description: "" });
      setTone("success");
      setMessage("Group created.");
      await loadGroups();
    } catch (error) {
      setTone("error");
      setMessage(error instanceof Error ? error.message : "Create group failed");
    }
  }

  async function joinGroup(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    try {
      await api.joinGroup({ group_id: Number(joinGroupId) });
      setTone("success");
      setMessage("Join request sent.");
    } catch (error) {
      setTone("error");
      setMessage(error instanceof Error ? error.message : "Join failed");
    }
  }

  async function inviteToGroup(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    try {
      await api.inviteToGroup({
        group_id: Number(invite.group_id),
        invitee_id: Number(invite.invitee_id),
      });
      setTone("success");
      setMessage("Invitation sent.");
    } catch (error) {
      setTone("error");
      setMessage(error instanceof Error ? error.message : "Invite failed");
    }
  }

  async function decide(type: "join-accept" | "join-decline" | "invite-accept" | "invite-decline") {
    try {
      const id = Number(decisionId);
      if (type === "join-accept") await api.acceptGroupJoin(id);
      if (type === "join-decline") await api.declineGroupJoin(id);
      if (type === "invite-accept") await api.acceptGroupInvitation(id);
      if (type === "invite-decline") await api.declineGroupInvitation(id);
      setTone("success");
      setMessage("Decision saved.");
    } catch (error) {
      setTone("error");
      setMessage(error instanceof Error ? error.message : "Decision failed");
    }
  }

  return (
    <div className="grid gap-6 lg:grid-cols-[1fr_24rem]">
      <section>
        <div className="mb-5 flex items-center justify-between">
          <h1 className="text-3xl font-bold text-slate-950">Groups</h1>
          <button
            className="rounded-xl border border-slate-300 px-4 py-2 text-sm font-semibold text-slate-700 transition hover:border-sky-400 hover:text-sky-700"
            onClick={loadGroups}
            type="button"
          >
            Refresh
          </button>
        </div>
        <div className="grid gap-4">
          {groups.map((group) => (
            <article
              className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"
              key={group.id}
            >
              <h2 className="text-2xl font-bold text-slate-950">
                <Link className="hover:text-sky-700" href={`/groups/${group.id}`}>
                  {group.name}
                </Link>
              </h2>
              <p className="mt-2 text-slate-700">{group.description}</p>
              <p className="mt-4 text-sm text-slate-500">
                Created {formatDate(group.created_at)} by user #{group.creator_id}
              </p>
            </article>
          ))}
          {groups.length === 0 ? (
            <p className="rounded-2xl border border-dashed border-slate-300 bg-white p-8 text-center text-slate-600">
              No groups found.
            </p>
          ) : null}
        </div>
      </section>
      <aside className="grid content-start gap-6">
        <FormMessage message={message} tone={tone} />
        <form
          className="grid gap-4 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"
          onSubmit={createGroup}
        >
          <h2 className="text-xl font-bold text-slate-950">Create group</h2>
          <TextField
            label="Title"
            name="title"
            onChange={(value) =>
              setCreateForm((current) => ({ ...current, title: value }))
            }
            required
            value={createForm.title}
          />
          <TextArea
            label="Description"
            name="description"
            onChange={(value) =>
              setCreateForm((current) => ({ ...current, description: value }))
            }
            required
            value={createForm.description}
          />
          <button className="rounded-xl bg-sky-600 px-4 py-3 font-semibold text-white" type="submit">
            Create
          </button>
        </form>
        <form
          className="grid gap-4 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"
          onSubmit={joinGroup}
        >
          <h2 className="text-xl font-bold text-slate-950">Join group</h2>
          <TextField
            label="Group ID"
            name="group_id"
            onChange={setJoinGroupId}
            required
            type="number"
            value={joinGroupId}
          />
          <button className="rounded-xl bg-sky-600 px-4 py-3 font-semibold text-white" type="submit">
            Request join
          </button>
        </form>
        <form
          className="grid gap-4 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"
          onSubmit={inviteToGroup}
        >
          <h2 className="text-xl font-bold text-slate-950">Invite user</h2>
          <TextField
            label="Group ID"
            name="group_id"
            onChange={(value) =>
              setInvite((current) => ({ ...current, group_id: value }))
            }
            required
            type="number"
            value={invite.group_id}
          />
          <TextField
            label="Invitee user ID"
            name="invitee_id"
            onChange={(value) =>
              setInvite((current) => ({ ...current, invitee_id: value }))
            }
            required
            type="number"
            value={invite.invitee_id}
          />
          <button className="rounded-xl bg-sky-600 px-4 py-3 font-semibold text-white" type="submit">
            Send invite
          </button>
        </form>
        <section className="grid gap-4 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <h2 className="text-xl font-bold text-slate-950">Requests by ID</h2>
          <TextField
            label="Request or invitation ID"
            name="decision_id"
            onChange={setDecisionId}
            required
            type="number"
            value={decisionId}
          />
          <div className="grid grid-cols-2 gap-2">
            <button className="rounded-xl border border-slate-300 px-3 py-2 text-sm font-semibold" onClick={() => decide("join-accept")} type="button">
              Accept join
            </button>
            <button className="rounded-xl border border-slate-300 px-3 py-2 text-sm font-semibold" onClick={() => decide("join-decline")} type="button">
              Decline join
            </button>
            <button className="rounded-xl border border-slate-300 px-3 py-2 text-sm font-semibold" onClick={() => decide("invite-accept")} type="button">
              Accept invite
            </button>
            <button className="rounded-xl border border-slate-300 px-3 py-2 text-sm font-semibold" onClick={() => decide("invite-decline")} type="button">
              Decline invite
            </button>
          </div>
        </section>
      </aside>
    </div>
  );
}
