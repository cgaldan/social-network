"use client";

import Link from "next/link";
import { useCallback, useEffect, useRef, useState } from "react";
import { useAuth } from "@/components/AuthProvider";
import { useWebSocket } from "@/components/WebSocketProvider";
import { Avatar } from "@/components/Avatar";
import { EmojiPicker } from "@/components/EmojiPicker";
import { useMessageCount } from "@/components/MessageCountProvider";
import { UserPicker } from "@/components/UserPicker";
import { api, ApiError } from "@/lib/api";
import { formatRelative } from "@/lib/format";
import type {
  ConversationParticipant,
  ConversationSummary,
  Message,
  MessageCreatedPayload,
  UserSummary,
} from "@/types/api";

export default function MessagesPage() {
  const { user } = useAuth();
  const { on, connected } = useWebSocket();
  const { decrement: decrementMessages } = useMessageCount();
  const [conversations, setConversations] = useState<ConversationSummary[]>([]);
  const [activeId, setActiveId] = useState<number | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [loadingList, setLoadingList] = useState(true);
  const [loadingMessages, setLoadingMessages] = useState(false);
  const [draft, setDraft] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [sending, setSending] = useState(false);
  const scrollerRef = useRef<HTMLDivElement>(null);

  const loadConversations = useCallback(async () => {
    setLoadingList(true);
    try {
      const res = await api.conversations.list({ limit: 50 });
      setConversations(res.conversations ?? []);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to load conversations");
    } finally {
      setLoadingList(false);
    }
  }, []);

  useEffect(() => {
    void loadConversations();
  }, [loadConversations]);

  useEffect(() => {
    if (!activeId) {
      setMessages([]);
      return;
    }
    let cancelled = false;
    setLoadingMessages(true);
    api.conversations
      .messages(activeId, { limit: 100 })
      .then((res) => {
        if (!cancelled) setMessages(res.messages ?? []);
      })
      .catch((err) => {
        if (!cancelled)
          setError(err instanceof ApiError ? err.message : "Failed to load messages");
      })
      .finally(() => {
        if (!cancelled) setLoadingMessages(false);
      });
    return () => {
      cancelled = true;
    };
  }, [activeId]);

  useEffect(() => {
    if (scrollerRef.current) {
      scrollerRef.current.scrollTop = scrollerRef.current.scrollHeight;
    }
  }, [messages]);

  useEffect(() => {
    return on((msg) => {
      if (msg.type !== "message.created") return;
      const payload = msg.payload as MessageCreatedPayload;
      const incoming = payload.message;

      if (incoming.conversation_id === activeId) {
        setMessages((prev) => {
          if (prev.some((m) => m.id === incoming.id)) return prev;
          return [...prev, incoming];
        });
        if (incoming.sender_id !== user?.id) {
          void api.conversations.markRead(activeId).catch(() => {});
          decrementMessages(1);
        }
      }

      setConversations((prev) => {
        const idx = prev.findIndex((c) => c.id === incoming.conversation_id);
        if (idx === -1) {
          void loadConversations();
          return prev;
        }
        const updated: ConversationSummary = {
          ...prev[idx],
          last_message: incoming,
          unread_count:
            incoming.conversation_id === activeId || incoming.sender_id === user?.id
              ? prev[idx].unread_count
              : prev[idx].unread_count + 1,
        };
        const next = [updated, ...prev.slice(0, idx), ...prev.slice(idx + 1)];
        return next;
      });
    });
  }, [on, activeId, user?.id, loadConversations]);

  const onSelectConversation = (id: number) => {
    setActiveId(id);
    const previousUnread = conversations.find((c) => c.id === id)?.unread_count ?? 0;
    setConversations((prev) =>
      prev.map((c) => (c.id === id ? { ...c, unread_count: 0 } : c)),
    );
    if (previousUnread > 0) decrementMessages(previousUnread);
    void api.conversations.markRead(id).catch(() => {});
  };

  const onStartConversation = async (target: UserSummary) => {
    if (!target || target.id === user?.id) return;
    setError(null);
    try {
      const res = await api.conversations.direct(target.id);
      await loadConversations();
      setActiveId(res.conversation.id);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to open conversation");
    }
  };

  const send = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!activeId || !draft.trim()) return;
    setSending(true);
    try {
      const res = await api.messages.send(activeId, draft);
      setMessages((prev) =>
        prev.some((m) => m.id === res.msg.id) ? prev : [...prev, res.msg],
      );
      setDraft("");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to send");
    } finally {
      setSending(false);
    }
  };

  const active = conversations.find((c) => c.id === activeId) ?? null;

  return (
    <div className="space-y-4">
      <header className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold text-slate-900">Messages</h1>
        <span className="flex items-center gap-2 text-xs text-slate-500">
          <span className={`h-2 w-2 rounded-full ${connected ? "bg-emerald-500" : "bg-slate-300"}`} />
          {connected ? "Live" : "Offline"}
        </span>
      </header>

      <section className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
        <p className="mb-2 text-xs font-medium uppercase tracking-wide text-slate-500">
          Start a conversation
        </p>
        <UserPicker onSelect={onStartConversation} placeholder="Search users by name…" />
      </section>

      {error ? (
        <p className="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p>
      ) : null}

      <div className="grid gap-4 md:grid-cols-[260px_1fr]">
        <aside className="rounded-xl border border-slate-200 bg-white shadow-sm">
          <header className="border-b border-slate-200 px-4 py-2 text-xs font-medium uppercase tracking-wide text-slate-500">
            Inbox
          </header>
          {loadingList ? (
            <p className="px-4 py-3 text-sm text-slate-500">Loading…</p>
          ) : conversations.length === 0 ? (
            <p className="px-4 py-3 text-sm text-slate-500">No conversations yet.</p>
          ) : (
            <ul className="max-h-[480px] overflow-y-auto">
              {conversations.map((c) => {
                const isActive = c.id === activeId;
                const other = c.participants.find((p) => p.user_id !== user?.id);
                const title = c.type === "group" ? c.name || "Group" : other
                  ? `${other.first_name} ${other.last_name}`
                  : "Direct chat";
                return (
                  <li key={c.id}>
                    <button
                      onClick={() => onSelectConversation(c.id)}
                      className={`flex w-full items-center gap-3 px-3 py-2 text-left transition ${
                        isActive ? "bg-indigo-50" : "hover:bg-slate-50"
                      }`}
                    >
                      <Avatar
                        src={other?.avatar_path}
                        firstName={other?.first_name}
                        lastName={other?.last_name}
                        nickname={other?.nickname ?? c.name}
                        size={36}
                      />
                      <div className="min-w-0 flex-1">
                        <div className="flex items-center justify-between gap-2">
                          <p className="truncate text-sm font-medium text-slate-900">
                            {title}
                          </p>
                          {c.last_message ? (
                            <span className="shrink-0 text-[10px] text-slate-400">
                              {formatRelative(c.last_message.created_at)}
                            </span>
                          ) : null}
                        </div>
                        <p className="truncate text-xs text-slate-500">
                          {c.last_message?.content ?? "No messages yet"}
                        </p>
                      </div>
                      {c.unread_count > 0 ? (
                        <span className="rounded-full bg-indigo-600 px-1.5 py-0.5 text-[10px] font-medium text-white">
                          {c.unread_count}
                        </span>
                      ) : null}
                    </button>
                  </li>
                );
              })}
            </ul>
          )}
        </aside>

        <section className="rounded-xl border border-slate-200 bg-white shadow-sm">
          {active ? (
            <ConversationPane
              conversation={active}
              messages={messages}
              loading={loadingMessages}
              draft={draft}
              sending={sending}
              currentUserId={user?.id}
              onDraftChange={setDraft}
              onSend={send}
              scrollerRef={scrollerRef}
            />
          ) : (
            <div className="flex h-64 items-center justify-center p-6 text-center text-sm text-slate-500">
              Select a conversation or search for a user to start chatting.
            </div>
          )}
        </section>
      </div>
    </div>
  );
}

function ConversationPane({
  conversation,
  messages,
  loading,
  draft,
  sending,
  currentUserId,
  onDraftChange,
  onSend,
  scrollerRef,
}: {
  conversation: ConversationSummary;
  messages: Message[];
  loading: boolean;
  draft: string;
  sending: boolean;
  currentUserId?: number;
  onDraftChange: (v: string) => void;
  onSend: (e: React.FormEvent) => void;
  scrollerRef: React.RefObject<HTMLDivElement | null>;
}) {
  const other: ConversationParticipant | undefined = conversation.participants.find(
    (p) => p.user_id !== currentUserId,
  );
  const title =
    conversation.type === "group"
      ? conversation.name || "Group"
      : other
        ? `${other.first_name} ${other.last_name}`
        : "Direct chat";

  const headerInner = (
    <>
      <Avatar
        src={other?.avatar_path}
        firstName={other?.first_name}
        lastName={other?.last_name}
        nickname={other?.nickname ?? conversation.name}
        size={36}
      />
      <div>
        <p className="text-sm font-medium text-slate-900">{title}</p>
        {other ? (
          <p className="text-xs text-slate-500">
            @{other.nickname}
            {other.is_online ? " · online" : ""}
          </p>
        ) : null}
      </div>
    </>
  );

  return (
    <div className="flex h-[560px] flex-col">
      {other ? (
        <Link
          href={`/users/${other.user_id}`}
          className="flex items-center gap-3 border-b border-slate-200 px-4 py-3 hover:bg-slate-50"
        >
          {headerInner}
        </Link>
      ) : (
        <header className="flex items-center gap-3 border-b border-slate-200 px-4 py-3">
          {headerInner}
        </header>
      )}

      <div ref={scrollerRef} className="flex-1 space-y-2 overflow-y-auto p-4">
        {loading ? (
          <p className="text-center text-sm text-slate-500">Loading messages…</p>
        ) : messages.length === 0 ? (
          <p className="text-center text-sm text-slate-500">No messages yet. Say hello.</p>
        ) : (
          messages.map((m) => {
            const mine = m.sender_id === currentUserId;
            return (
              <div key={m.id} className={`flex ${mine ? "justify-end" : "justify-start"}`}>
                <div
                  className={`max-w-[75%] rounded-2xl px-4 py-2 text-sm shadow-sm ${
                    mine ? "bg-indigo-600 text-white" : "bg-slate-100 text-slate-900"
                  }`}
                >
                  <p className="whitespace-pre-wrap">{m.content}</p>
                  <p
                    className={`mt-1 text-[10px] ${mine ? "text-indigo-200" : "text-slate-500"}`}
                  >
                    {formatRelative(m.created_at)}
                  </p>
                </div>
              </div>
            );
          })
        )}
      </div>

      <form onSubmit={onSend} className="flex gap-2 border-t border-slate-200 p-3">
        <input
          type="text"
          value={draft}
          onChange={(e) => onDraftChange(e.target.value)}
          placeholder="Type a message…"
          className="flex-1 rounded-lg border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
        />
        <EmojiPicker onPick={(emoji) => onDraftChange(draft + emoji)} />
        <button
          type="submit"
          disabled={sending || !draft.trim()}
          className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm transition hover:bg-indigo-700 disabled:opacity-60"
        >
          Send
        </button>
      </form>
    </div>
  );
}
