"use client";

import { FormEvent, useState } from "react";
import { FormMessage, TextArea, TextField } from "@/components/forms";
import { api } from "@/lib/api";
import type { Conversation, Message } from "@/types/api";

export default function MessagesPage() {
  const [receiverId, setReceiverId] = useState("");
  const [conversationId, setConversationId] = useState("");
  const [content, setContent] = useState("");
  const [conversation, setConversation] = useState<Conversation | null>(null);
  const [lastMessage, setLastMessage] = useState<Message | null>(null);
  const [message, setMessage] = useState("");
  const [tone, setTone] = useState<"error" | "success">("success");

  async function createConversation(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    try {
      const response = await api.createDirectConversation({
        receiver_id: Number(receiverId),
      });
      setConversation(response.conversation ?? null);
      if (response.conversation) {
        setConversationId(String(response.conversation.id));
      }
      setTone("success");
      setMessage(response.message || "Conversation ready.");
    } catch (error) {
      setTone("error");
      setMessage(error instanceof Error ? error.message : "Conversation failed");
    }
  }

  async function sendMessage(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    try {
      const response = await api.sendMessage({
        conversation_id: Number(conversationId),
        content,
      });
      setLastMessage(response.msg ?? null);
      setContent("");
      setTone("success");
      setMessage(response.message || "Message sent.");
    } catch (error) {
      setTone("error");
      setMessage(error instanceof Error ? error.message : "Message failed");
    }
  }

  return (
    <div className="grid gap-6 lg:grid-cols-[24rem_1fr]">
      <aside className="grid content-start gap-6">
        <form
          className="grid gap-4 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"
          onSubmit={createConversation}
        >
          <h1 className="text-2xl font-bold text-slate-950">Direct chat</h1>
          <p className="text-sm text-slate-600">
            Create or reuse a direct conversation by receiver user ID.
          </p>
          <TextField
            label="Receiver user ID"
            name="receiver_id"
            onChange={setReceiverId}
            required
            type="number"
            value={receiverId}
          />
          <button className="rounded-xl bg-sky-600 px-4 py-3 font-semibold text-white" type="submit">
            Start conversation
          </button>
        </form>
        <form
          className="grid gap-4 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"
          onSubmit={sendMessage}
        >
          <h2 className="text-xl font-bold text-slate-950">Send message</h2>
          <TextField
            label="Conversation ID"
            name="conversation_id"
            onChange={setConversationId}
            required
            type="number"
            value={conversationId}
          />
          <TextArea
            label="Message"
            name="content"
            onChange={setContent}
            required
            value={content}
          />
          <button className="rounded-xl bg-sky-600 px-4 py-3 font-semibold text-white" type="submit">
            Send
          </button>
        </form>
      </aside>
      <section className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
        <h2 className="text-3xl font-bold text-slate-950">Messages</h2>
        <p className="mt-3 text-slate-600">
          The backend currently supports creating direct conversations and
          sending messages. It does not expose conversation lists or message
          history routes yet, so this page shows the latest action result.
        </p>
        <div className="mt-6">
          <FormMessage message={message} tone={tone} />
        </div>
        {conversation ? (
          <div className="mt-6 rounded-xl bg-slate-50 p-4">
            <p className="text-sm font-semibold text-slate-500">
              Current conversation
            </p>
            <p className="mt-1 text-slate-950">
              #{conversation.id} / {conversation.name || conversation.type}
            </p>
          </div>
        ) : null}
        {lastMessage ? (
          <div className="mt-4 rounded-xl bg-sky-50 p-4">
            <p className="text-sm font-semibold text-sky-700">Last sent</p>
            <p className="mt-1 text-slate-950">{lastMessage.content}</p>
          </div>
        ) : null}
      </section>
    </div>
  );
}
