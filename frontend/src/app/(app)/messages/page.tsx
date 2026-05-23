"use client";

import Link from "next/link";
import { useCallback, useEffect, useRef, useState } from "react";
import { useAuth } from "@/components/AuthProvider";
import { useWebSocket } from "@/components/WebSocketProvider";
import { Avatar } from "@/components/Avatar";
import { EmojiPicker } from "@/components/EmojiPicker";
import { useMessageCount } from "@/components/MessageCountProvider";
import { UserPicker } from "@/components/UserPicker";
import { InfiniteScrollSentinel } from "@/components/InfiniteScrollSentinel";
import { usePaginatedList } from "@/hooks/usePaginatedList";
import { api, ApiError } from "@/lib/api";
import { formatRelative } from "@/lib/format";
import type {
  ConversationParticipant,
  ConversationSummary,
  Message,
  MessageCreatedPayload,
  MessageDeletedPayload,
  MessageUpdatedPayload,
  UserSummary,
} from "@/types/api";

const CONVERSATIONS_PAGE_SIZE = 20;
const MESSAGES_PAGE_SIZE = 30;

export default function MessagesPage() {
  const { user } = useAuth();
  const { on, connected } = useWebSocket();
  const { decrement: decrementMessages } = useMessageCount();
  const [activeId, setActiveId] = useState<number | null>(null);
  const [draft, setDraft] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [sending, setSending] = useState(false);
  const [editingMessageId, setEditingMessageId] = useState<number | null>(null);
  const [editDraft, setEditDraft] = useState("");
  const [savingMessage, setSavingMessage] = useState(false);
  const scrollerRef = useRef<HTMLDivElement>(null);

  // --- Conversations sidebar (standard top-to-bottom infinite scroll) ---
  const conversationsFetcher = useCallback(
    async ({ limit, offset }: { limit: number; offset: number }) => {
      const res = await api.conversations.list({ limit, offset });
      return { items: res.conversations ?? [], hasMore: res.has_more };
    },
    [],
  );
  const {
    items: conversations,
    hasMore: hasMoreConversations,
    loading: loadingList,
    loadMore: loadMoreConversations,
    setItems: setConversations,
    reset: resetConversations,
  } = usePaginatedList<ConversationSummary>(
    conversationsFetcher,
    CONVERSATIONS_PAGE_SIZE,
    (c) => c.id,
  );

  // --- Messages inside the active conversation ---
  // Cursor-based pagination (before_id), not offset, because new messages
  // are inserted at the head all the time and offset would drift. We store
  // messages in chronological order (ASC: oldest first, newest last) so the
  // render is a plain .map() with no reversal anywhere.
  const [messages, setMessages] = useState<Message[]>([]);
  const [hasOlder, setHasOlder] = useState(false);
  const [loadingMessages, setLoadingMessages] = useState(false);
  const messagesRef = useRef<Message[]>([]);
  messagesRef.current = messages;

  // Scroll bookkeeping:
  //  - preLoadScrollRef captures pre-prepend scroll metrics so we can restore
  //    the visual position after older messages slot in at the top.
  //  - stuckToBottomRef tracks whether the user is reading the latest message
  //    (so we auto-scroll on new messages) or has scrolled up (so we leave
  //    them alone).
  const preLoadScrollRef = useRef<{ height: number; top: number } | null>(null);
  const stuckToBottomRef = useRef(true);

  // Initial load for the active conversation (or reset when switching).
  useEffect(() => {
    if (activeId == null) {
      setMessages([]);
      setHasOlder(false);
      return;
    }
    let cancelled = false;
    stuckToBottomRef.current = true;
    setLoadingMessages(true);
    api.conversations
      .messages(activeId, { limit: MESSAGES_PAGE_SIZE })
      .then((res) => {
        if (cancelled) return;
        // Backend returns DESC (newest first); flip to ASC for chronological render.
        const asc = [...(res.messages ?? [])].reverse();
        setMessages(asc);
        setHasOlder(res.has_more);
      })
      .catch((err) => {
        if (!cancelled) setError(err instanceof ApiError ? err.message : "Failed to load messages");
      })
      .finally(() => {
        if (!cancelled) setLoadingMessages(false);
      });
    return () => {
      cancelled = true;
    };
  }, [activeId]);

  // After messages change, manage scroll position.
  useEffect(() => {
    const el = scrollerRef.current;
    if (!el) return;
    const snap = preLoadScrollRef.current;
    if (snap) {
      el.scrollTop = el.scrollHeight - snap.height + snap.top;
      preLoadScrollRef.current = null;
      return;
    }
    if (stuckToBottomRef.current) {
      el.scrollTop = el.scrollHeight;
    }
  }, [messages]);

  const onMessagesScroll = (e: React.UIEvent<HTMLDivElement>) => {
    const el = e.currentTarget;
    const nearBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 80;
    stuckToBottomRef.current = nearBottom;
  };

  const loadOlderMessages = useCallback(() => {
    if (!scrollerRef.current || loadingMessages || !hasOlder || activeId == null) return;
    const oldest = messagesRef.current[0];
    if (!oldest) return;
    preLoadScrollRef.current = {
      height: scrollerRef.current.scrollHeight,
      top: scrollerRef.current.scrollTop,
    };
    setLoadingMessages(true);
    api.conversations
      .messages(activeId, { limit: MESSAGES_PAGE_SIZE, before_id: oldest.id })
      .then((res) => {
        const olderAsc = [...(res.messages ?? [])].reverse();
        setMessages((prev) => {
          // Defensive: drop any rows we already have (shouldn't happen with
          // cursor pagination, but bulletproof against duplicate-key crashes).
          const seen = new Set(prev.map((m) => m.id));
          const novel = olderAsc.filter((m) => !seen.has(m.id));
          return [...novel, ...prev];
        });
        setHasOlder(res.has_more);
      })
      .catch((err) => {
        preLoadScrollRef.current = null;
        setError(err instanceof ApiError ? err.message : "Failed to load older messages");
      })
      .finally(() => setLoadingMessages(false));
  }, [activeId, hasOlder, loadingMessages]);

  useEffect(() => {
    return on((msg) => {
      if (msg.type === "message.created") {
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
            resetConversations();
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
          return [updated, ...prev.slice(0, idx), ...prev.slice(idx + 1)];
        });
        return;
      }

      if (msg.type === "message.updated") {
        const payload = msg.payload as MessageUpdatedPayload;
        const updated = payload.message;
        if (updated.conversation_id === activeId) {
          setMessages((prev) =>
            prev.map((m) => (m.id === updated.id ? updated : m)),
          );
        }
        setConversations((prev) =>
          prev.map((c) =>
            c.id === updated.conversation_id && c.last_message?.id === updated.id
              ? { ...c, last_message: updated }
              : c,
          ),
        );
        return;
      }

      if (msg.type === "message.deleted") {
        const payload = msg.payload as MessageDeletedPayload;
        if (payload.conversation_id === activeId) {
          setMessages((prev) => prev.filter((m) => m.id !== payload.message_id));
        }
        setConversations((prev) =>
          prev.map((c) =>
            c.id === payload.conversation_id && c.last_message?.id === payload.message_id
              ? { ...c, last_message: undefined }
              : c,
          ),
        );
        return;
      }
    });
  }, [on, activeId, user?.id, resetConversations, setMessages, setConversations, decrementMessages]);

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
      // Insert (or move-to-top) optimistically using what we already know
      // from the picker, so the sidebar and chat header render immediately
      // instead of flashing empty while a refetch runs.
      setConversations((prev) => {
        const existing = prev.find((c) => c.id === res.conversation.id);
        if (existing) {
          return [existing, ...prev.filter((c) => c.id !== res.conversation.id)];
        }
        const summary: ConversationSummary = {
          id: res.conversation.id,
          type: res.conversation.type || "private",
          name: res.conversation.name || "",
          participants: [
            {
              user_id: target.id,
              nickname: target.nickname,
              first_name: target.first_name,
              last_name: target.last_name,
              avatar_path: target.avatar_path,
              is_online: target.is_online,
            },
          ],
          unread_count: 0,
          created_at: res.conversation.created_at ?? new Date().toISOString(),
        };
        return [summary, ...prev];
      });
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
      stuckToBottomRef.current = true;
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

  const startEditMessage = (m: Message) => {
    setEditingMessageId(m.id);
    setEditDraft(m.content);
  };

  const cancelEditMessage = () => {
    setEditingMessageId(null);
    setEditDraft("");
  };

  const saveEditMessage = async () => {
    if (editingMessageId == null || !editDraft.trim()) return;
    setSavingMessage(true);
    try {
      const res = await api.messages.update(editingMessageId, editDraft);
      setMessages((prev) => prev.map((m) => (m.id === res.msg.id ? res.msg : m)));
      setConversations((prev) =>
        prev.map((c) =>
          c.id === res.msg.conversation_id && c.last_message?.id === res.msg.id
            ? { ...c, last_message: res.msg }
            : c,
        ),
      );
      cancelEditMessage();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to update message");
    } finally {
      setSavingMessage(false);
    }
  };

  const deleteMessage = async (messageId: number) => {
    if (!confirm("Delete this message?")) return;
    try {
      await api.messages.remove(messageId);
      setMessages((prev) => prev.filter((m) => m.id !== messageId));
      setConversations((prev) =>
        prev.map((c) =>
          c.last_message?.id === messageId ? { ...c, last_message: undefined } : c,
        ),
      );
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to delete message");
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
          {loadingList && conversations.length === 0 ? (
            <p className="px-4 py-3 text-sm text-slate-500">Loading…</p>
          ) : conversations.length === 0 ? (
            <p className="px-4 py-3 text-sm text-slate-500">No conversations yet.</p>
          ) : (
            <ul className="max-h-[480px] overflow-y-auto">
              {conversations.map((c) => {
                const isActive = c.id === activeId;
                const isGroup = c.type === "group";
                const other = isGroup
                  ? undefined
                  : c.participants.find((p) => p.user_id !== user?.id);
                const title = isGroup
                  ? c.name || "Group"
                  : other
                    ? `${other.first_name} ${other.last_name}`
                    : "Direct chat";
                const avatarSrc = isGroup ? c.avatar_path : other?.avatar_path;
                const avatarNickname = isGroup ? c.name || "Group" : other?.nickname;
                return (
                  <li key={c.id}>
                    <button
                      onClick={() => onSelectConversation(c.id)}
                      className={`flex w-full items-center gap-3 px-3 py-2 text-left transition ${
                        isActive ? "bg-indigo-50" : "hover:bg-slate-50"
                      }`}
                    >
                      <Avatar
                        src={avatarSrc}
                        firstName={isGroup ? undefined : other?.first_name}
                        lastName={isGroup ? undefined : other?.last_name}
                        nickname={avatarNickname}
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
              {hasMoreConversations ? (
                <li>
                  <InfiniteScrollSentinel
                    onIntersect={loadMoreConversations}
                    enabled={!loadingList}
                  />
                </li>
              ) : null}
            </ul>
          )}
        </aside>

        <section className="rounded-xl border border-slate-200 bg-white shadow-sm">
          {active ? (
            <ConversationPane
              conversation={active}
              messages={messages}
              loading={loadingMessages}
              hasMoreOlder={hasOlder}
              onLoadOlder={loadOlderMessages}
              onScroll={onMessagesScroll}
              draft={draft}
              sending={sending}
              currentUserId={user?.id}
              onDraftChange={setDraft}
              onSend={send}
              scrollerRef={scrollerRef}
              editingMessageId={editingMessageId}
              editDraft={editDraft}
              savingMessage={savingMessage}
              onEditDraftChange={setEditDraft}
              onStartEditMessage={startEditMessage}
              onCancelEditMessage={cancelEditMessage}
              onSaveEditMessage={saveEditMessage}
              onDeleteMessage={deleteMessage}
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
  hasMoreOlder,
  onLoadOlder,
  onScroll,
  draft,
  sending,
  currentUserId,
  onDraftChange,
  onSend,
  scrollerRef,
  editingMessageId,
  editDraft,
  savingMessage,
  onEditDraftChange,
  onStartEditMessage,
  onCancelEditMessage,
  onSaveEditMessage,
  onDeleteMessage,
}: {
  conversation: ConversationSummary;
  messages: Message[];
  loading: boolean;
  hasMoreOlder: boolean;
  onLoadOlder: () => void;
  onScroll: (e: React.UIEvent<HTMLDivElement>) => void;
  draft: string;
  sending: boolean;
  currentUserId?: number;
  onDraftChange: (v: string) => void;
  onSend: (e: React.FormEvent) => void;
  scrollerRef: React.RefObject<HTMLDivElement | null>;
  editingMessageId: number | null;
  editDraft: string;
  savingMessage: boolean;
  onEditDraftChange: (v: string) => void;
  onStartEditMessage: (m: Message) => void;
  onCancelEditMessage: () => void;
  onSaveEditMessage: () => void;
  onDeleteMessage: (id: number) => void;
}) {
  const isGroup = conversation.type === "group";
  const other: ConversationParticipant | undefined = isGroup
    ? undefined
    : conversation.participants.find((p) => p.user_id !== currentUserId);
  const title = isGroup
    ? conversation.name || "Group"
    : other
      ? `${other.first_name} ${other.last_name}`
      : "Direct chat";
  const memberCount = isGroup ? conversation.participants.length : 0;

  const headerInner = (
    <>
      <Avatar
        src={isGroup ? conversation.avatar_path : other?.avatar_path}
        firstName={isGroup ? undefined : other?.first_name}
        lastName={isGroup ? undefined : other?.last_name}
        nickname={isGroup ? conversation.name || "Group" : other?.nickname}
        size={36}
      />
      <div>
        <p className="text-sm font-medium text-slate-900">{title}</p>
        {isGroup ? (
          <p className="text-xs text-slate-500">
            {memberCount} {memberCount === 1 ? "member" : "members"}
          </p>
        ) : other ? (
          <p className="text-xs text-slate-500">
            @{other.nickname}
            {other.is_online ? " · online" : " · offline"}
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

      <div
        ref={scrollerRef}
        onScroll={onScroll}
        className="flex-1 space-y-2 overflow-y-auto p-4"
      >
        {loading && messages.length === 0 ? (
          <p className="text-center text-sm text-slate-500">Loading messages…</p>
        ) : messages.length === 0 ? (
          <p className="text-center text-sm text-slate-500">No messages yet. Say hello.</p>
        ) : (
          <>
            {hasMoreOlder ? (
              <InfiniteScrollSentinel onIntersect={onLoadOlder} enabled={!loading} />
            ) : null}
            {loading && messages.length > 0 ? (
              <p className="text-center text-xs text-slate-400">Loading older messages…</p>
            ) : null}
            {messages.map((m) => {
              const mine = m.sender_id === currentUserId;
              const isEditing = editingMessageId === m.id;
              return (
                <div key={m.id} className={`flex ${mine ? "justify-end" : "justify-start"}`}>
                  <div className={`max-w-[75%] ${mine ? "items-end" : "items-start"} flex flex-col gap-1`}>
                    <div
                      className={`rounded-2xl px-4 py-2 text-sm shadow-sm ${
                        mine ? "bg-indigo-600 text-white" : "bg-slate-100 text-slate-900"
                      }`}
                    >
                      {isEditing ? (
                        <div className="flex flex-col gap-2">
                          <textarea
                            name="editMessageBody"
                            value={editDraft}
                            onChange={(e) => onEditDraftChange(e.target.value)}
                            rows={2}
                            className="w-full min-w-[200px] rounded-md border border-white/30 bg-white/10 px-2 py-1 text-sm text-white placeholder-indigo-200 focus:outline-none"
                          />
                          <div className="flex justify-end gap-2 text-[11px]">
                            <button
                              type="button"
                              onClick={onCancelEditMessage}
                              className="rounded-md bg-white/20 px-2 py-0.5 text-white hover:bg-white/30"
                            >
                              Cancel
                            </button>
                            <button
                              type="button"
                              onClick={onSaveEditMessage}
                              disabled={savingMessage || !editDraft.trim()}
                              className="rounded-md bg-white px-2 py-0.5 font-medium text-indigo-700 hover:bg-indigo-50 disabled:opacity-60"
                            >
                              {savingMessage ? "Saving…" : "Save"}
                            </button>
                          </div>
                        </div>
                      ) : (
                        <>
                          <p className="whitespace-pre-wrap">{m.content}</p>
                          <p
                            className={`mt-1 text-[10px] ${mine ? "text-indigo-200" : "text-slate-500"}`}
                          >
                            {formatRelative(m.created_at)}
                            {m.updated_at ? " · edited" : ""}
                          </p>
                        </>
                      )}
                    </div>
                    {mine && !isEditing ? (
                      <div className="flex gap-2 text-[10px] text-slate-400">
                        <button
                          type="button"
                          onClick={() => onStartEditMessage(m)}
                          className="hover:text-indigo-600"
                        >
                          Edit
                        </button>
                        <button
                          type="button"
                          onClick={() => onDeleteMessage(m.id)}
                          className="hover:text-red-600"
                        >
                          Delete
                        </button>
                      </div>
                    ) : null}
                  </div>
                </div>
              );
            })}
          </>
        )}
      </div>

      <form onSubmit={onSend} className="flex gap-2 border-t border-slate-200 p-3">
        <input
          type="text"
          name="messageBody"
          value={draft}
          onChange={(e) => onDraftChange(e.target.value)}
          placeholder="Type a message…"
          autoComplete="off"
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
