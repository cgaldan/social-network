"use client";

import { useParams } from "next/navigation";
import { FormEvent, useEffect, useMemo, useState } from "react";
import { PostCard } from "@/components/PostCard";
import { PostForm } from "@/components/PostForm";
import { FormMessage, SelectField, TextArea, TextField } from "@/components/forms";
import { api } from "@/lib/api";
import { formatDate } from "@/lib/format";
import type { GroupEvent, Post } from "@/types/api";

export default function GroupDetailPage() {
  const params = useParams<{ id: string }>();
  const groupId = useMemo(() => Number(params.id), [params.id]);
  const [posts, setPosts] = useState<Post[]>([]);
  const [events, setEvents] = useState<GroupEvent[]>([]);
  const [message, setMessage] = useState("");
  const [tone, setTone] = useState<"error" | "success">("success");
  const [eventForm, setEventForm] = useState({
    title: "",
    description: "",
    starts_at: "",
  });
  const [rsvp, setRsvp] = useState({ event_id: "", response: "going" });

  async function loadGroupData() {
    try {
      const [postResponse, eventResponse] = await Promise.all([
        api.listGroupPosts(groupId),
        api.listGroupEvents(groupId),
      ]);
      setPosts(postResponse.posts ?? []);
      setEvents(eventResponse.events ?? []);
    } catch (error) {
      setTone("error");
      setMessage(error instanceof Error ? error.message : "Could not load group");
    }
  }

  useEffect(() => {
    if (Number.isFinite(groupId)) {
      void loadGroupData();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [groupId]);

  async function createEvent(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    try {
      await api.createGroupEvent(groupId, {
        ...eventForm,
        starts_at: new Date(eventForm.starts_at).toISOString(),
      });
      setEventForm({ title: "", description: "", starts_at: "" });
      setTone("success");
      setMessage("Event created.");
      await loadGroupData();
    } catch (error) {
      setTone("error");
      setMessage(error instanceof Error ? error.message : "Event failed");
    }
  }

  async function submitRsvp(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    try {
      await api.rsvpGroupEvent(groupId, Number(rsvp.event_id), {
        response: rsvp.response,
      });
      setTone("success");
      setMessage("RSVP saved.");
    } catch (error) {
      setTone("error");
      setMessage(error instanceof Error ? error.message : "RSVP failed");
    }
  }

  return (
    <div className="grid gap-6 lg:grid-cols-[1fr_24rem]">
      <section className="grid content-start gap-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold text-slate-950">
            Group #{groupId}
          </h1>
          <button
            className="rounded-xl border border-slate-300 px-4 py-2 text-sm font-semibold text-slate-700"
            onClick={loadGroupData}
            type="button"
          >
            Refresh
          </button>
        </div>
        <FormMessage message={message} tone={tone} />
        <section>
          <h2 className="mb-4 text-2xl font-bold text-slate-950">Posts</h2>
          <div className="grid gap-4">
            {posts.map((post) => (
              <PostCard key={post.id} post={post} />
            ))}
            {posts.length === 0 ? (
              <p className="rounded-2xl border border-dashed border-slate-300 bg-white p-8 text-center text-slate-600">
                No group posts yet.
              </p>
            ) : null}
          </div>
        </section>
        <section>
          <h2 className="mb-4 text-2xl font-bold text-slate-950">Events</h2>
          <div className="grid gap-4">
            {events.map((groupEvent) => (
              <article
                className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"
                key={groupEvent.id}
              >
                <h3 className="text-xl font-bold text-slate-950">
                  {groupEvent.title}
                </h3>
                <p className="mt-2 text-slate-700">{groupEvent.description}</p>
                <p className="mt-4 text-sm text-slate-500">
                  Starts {formatDate(groupEvent.starts_at)}
                </p>
              </article>
            ))}
            {events.length === 0 ? (
              <p className="rounded-2xl border border-dashed border-slate-300 bg-white p-8 text-center text-slate-600">
                No events yet.
              </p>
            ) : null}
          </div>
        </section>
      </section>
      <aside className="grid content-start gap-6">
        <section>
          <h2 className="mb-4 text-xl font-bold text-slate-950">
            Create group post
          </h2>
          <PostForm
            onSubmit={async (body) => {
              await api.createGroupPost(groupId, body);
              await loadGroupData();
            }}
          />
        </section>
        <form
          className="grid gap-4 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"
          onSubmit={createEvent}
        >
          <h2 className="text-xl font-bold text-slate-950">Create event</h2>
          <TextField
            label="Title"
            name="title"
            onChange={(value) =>
              setEventForm((current) => ({ ...current, title: value }))
            }
            required
            value={eventForm.title}
          />
          <TextArea
            label="Description"
            name="description"
            onChange={(value) =>
              setEventForm((current) => ({ ...current, description: value }))
            }
            required
            value={eventForm.description}
          />
          <TextField
            label="Starts at"
            name="starts_at"
            onChange={(value) =>
              setEventForm((current) => ({ ...current, starts_at: value }))
            }
            required
            type="datetime-local"
            value={eventForm.starts_at}
          />
          <button className="rounded-xl bg-sky-600 px-4 py-3 font-semibold text-white" type="submit">
            Create event
          </button>
        </form>
        <form
          className="grid gap-4 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"
          onSubmit={submitRsvp}
        >
          <h2 className="text-xl font-bold text-slate-950">RSVP</h2>
          <TextField
            label="Event ID"
            name="event_id"
            onChange={(value) =>
              setRsvp((current) => ({ ...current, event_id: value }))
            }
            required
            type="number"
            value={rsvp.event_id}
          />
          <SelectField
            label="Response"
            name="response"
            onChange={(value) =>
              setRsvp((current) => ({ ...current, response: value }))
            }
            options={[
              { label: "Going", value: "going" },
              { label: "Not going", value: "not_going" },
              { label: "Maybe", value: "maybe" },
            ]}
            value={rsvp.response}
          />
          <button className="rounded-xl bg-sky-600 px-4 py-3 font-semibold text-white" type="submit">
            Save RSVP
          </button>
        </form>
      </aside>
    </div>
  );
}
