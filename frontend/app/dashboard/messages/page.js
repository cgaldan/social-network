"use client";

import { useState } from "react";
import {
  createDirectConversation,
  sendMessage,
  getStoredToken,
} from "../../../lib/api";

export default function MessagesPage() {
  const token = getStoredToken();
  const [receiverId, setReceiverId] = useState("");
  const [conversationId, setConversationId] = useState("");
  const [content, setContent] = useState("");
  const [error, setError] = useState("");
  const [msg, setMsg] = useState("");

  const startDm = async (e) => {
    e.preventDefault();
    setError("");
    setMsg("");
    try {
      const res = await createDirectConversation(token, Number(receiverId));
      const id = res.conversation?.id;
      if (id) setConversationId(String(id));
      setMsg(res.message || "Conversation ready.");
    } catch (e) {
      setError(e.message);
    }
  };

  const send = async (e) => {
    e.preventDefault();
    setError("");
    setMsg("");
    try {
      await sendMessage(token, {
        conversationId: Number(conversationId),
        content: content.trim(),
      });
      setContent("");
      setMsg("Message sent.");
    } catch (e) {
      setError(e.message);
    }
  };

  return (
    <div className="page-stack">
      <section className="surface-card">
        <h1>Messages</h1>
        <p className="helper-text">
          Start a direct conversation by receiver user id, then send messages using the returned conversation id. There is no messages history API in this backend yet.
        </p>

        <form className="stack-form" onSubmit={startDm}>
          <h2 className="h2-inline">Start or open DM</h2>
          <label>Receiver user id</label>
          <input
            type="number"
            min={1}
            value={receiverId}
            onChange={(e) => setReceiverId(e.target.value)}
            required
          />
          <button type="submit">Create direct conversation</button>
        </form>

        <form className="stack-form" onSubmit={send}>
          <h2 className="h2-inline">Send message</h2>
          <label>Conversation id</label>
          <input
            type="number"
            min={1}
            value={conversationId}
            onChange={(e) => setConversationId(e.target.value)}
            required
          />
          <label>Message</label>
          <textarea
            rows={3}
            value={content}
            onChange={(e) => setContent(e.target.value)}
            required
          />
          <button type="submit">Send</button>
        </form>

        {error ? <p className="error-message">{error}</p> : null}
        {msg ? <p className="success-message">{msg}</p> : null}
      </section>
    </div>
  );
}
