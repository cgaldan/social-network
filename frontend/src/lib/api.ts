import { getToken } from "./auth";
import type {
  AuthResponse,
  Comment,
  Conversation,
  FollowRequestSummary,
  ConversationSummary,
  CreateCommentPayload,
  CreateEventPayload,
  CreateGroupPayload,
  CreatePostPayload,
  Group,
  GroupEvent,
  LoginPayload,
  MeResponse,
  Message,
  Notification,
  Post,
  RegisterPayload,
  Rsvp,
  RsvpResponse,
  UpdatePostPayload,
  User,
  UserProfile,
  UserSummary,
} from "@/types/api";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8000";

export function resolveMediaUrl(path: string | undefined | null): string | undefined {
  if (!path) return undefined;
  if (/^https?:\/\//i.test(path)) return path;
  if (path.startsWith("/")) return API_URL + path;
  return API_URL + "/" + path;
}

export class ApiError extends Error {
  status: number;
  body: unknown;
  constructor(status: number, message: string, body: unknown) {
    super(message);
    this.status = status;
    this.body = body;
  }
}

interface RequestOptions {
  method?: string;
  body?: unknown;
  query?: object;
  auth?: boolean;
}

export interface PageOptions {
  limit?: number;
  offset?: number;
}

export type Paginated<K extends string, T> = {
  success: boolean;
  message: string;
  has_more: boolean;
} & { [P in K]: T[] };

async function request<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const { method = "GET", body, query, auth = true } = opts;
  const url = new URL(API_URL + path);
  if (query) {
    for (const [k, v] of Object.entries(query)) {
      if (v !== undefined && v !== null && v !== "") {
        url.searchParams.set(k, String(v));
      }
    }
  }
  const headers: Record<string, string> = {};
  if (body !== undefined) headers["Content-Type"] = "application/json";
  if (auth) {
    const token = getToken();
    if (token) headers["Authorization"] = token;
  }

  const res = await fetch(url.toString(), {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });

  const text = await res.text();
  let data: unknown = null;
  if (text) {
    try {
      data = JSON.parse(text);
    } catch {
      data = text;
    }
  }

  if (!res.ok) {
    const message =
      (data && typeof data === "object" && "message" in data && typeof (data as { message: unknown }).message === "string"
        ? (data as { message: string }).message
        : null) ?? res.statusText;
    throw new ApiError(res.status, message, data);
  }

  if (
    data &&
    typeof data === "object" &&
    "success" in data &&
    (data as { success: unknown }).success === false
  ) {
    const message =
      "message" in data && typeof (data as { message: unknown }).message === "string"
        ? (data as { message: string }).message
        : "Request failed";
    throw new ApiError(res.status, message, data);
  }

  return data as T;
}

async function uploadFile(file: File): Promise<{ success: boolean; message: string; url: string }> {
  const form = new FormData();
  form.append("file", file);

  const headers: Record<string, string> = {};
  const token = getToken();
  if (token) headers["Authorization"] = token;

  const res = await fetch(API_URL + "/api/uploads", {
    method: "POST",
    headers,
    body: form,
  });

  const text = await res.text();
  let data: unknown = null;
  if (text) {
    try {
      data = JSON.parse(text);
    } catch {
      data = text;
    }
  }

  const message =
    data && typeof data === "object" && "message" in data && typeof (data as { message: unknown }).message === "string"
      ? (data as { message: string }).message
      : res.statusText;

  if (!res.ok || (data && typeof data === "object" && "success" in data && (data as { success: unknown }).success === false)) {
    throw new ApiError(res.status, message || "Upload failed", data);
  }

  return data as { success: boolean; message: string; url: string };
}

export const api = {
  auth: {
    register: (payload: RegisterPayload) =>
      request<AuthResponse>("/api/auth/register", { method: "POST", body: payload, auth: false }),
    login: (payload: LoginPayload) =>
      request<AuthResponse>("/api/auth/login", { method: "POST", body: payload, auth: false }),
    logout: () => request<{ success: boolean; message: string }>("/api/auth/logout", { method: "POST" }),
    me: () => request<MeResponse>("/api/auth/me"),
    update: (payload: Partial<RegisterPayload>) =>
      request<MeResponse>("/api/auth/me", { method: "PUT", body: payload }),
    remove: () => request<{ success: boolean; message: string }>("/api/auth/me", { method: "DELETE" }),
  },
  posts: {
    list: (params: { category?: string } & PageOptions = {}) =>
      request<Paginated<"posts", Post>>("/api/posts", { query: params }),
    create: (payload: CreatePostPayload) =>
      request<{ success: boolean; message: string; post: Post }>("/api/posts", {
        method: "POST",
        body: payload,
      }),
    get: (id: number) =>
      request<{ success: boolean; message: string; post: Post }>(`/api/posts/${id}`),
    update: (id: number, payload: UpdatePostPayload) =>
      request<{ success: boolean; message: string; post: Post }>(`/api/posts/${id}`, {
        method: "PUT",
        body: payload,
      }),
    remove: (id: number) =>
      request<{ success: boolean; message: string }>(`/api/posts/${id}`, { method: "DELETE" }),
  },
  comments: {
    list: (postId: number, params: PageOptions = {}) =>
      request<Paginated<"comments", Comment>>(`/api/posts/${postId}/comments`, {
        query: params,
      }),
    create: (postId: number, payload: CreateCommentPayload) =>
      request<{ success: boolean; message: string; comment: Comment }>(
        `/api/posts/${postId}/comments`,
        { method: "POST", body: payload },
      ),
    update: (postId: number, commentId: number, payload: Partial<CreateCommentPayload>) =>
      request<{ success: boolean; message: string; comment: Comment }>(
        `/api/posts/${postId}/comments/${commentId}`,
        { method: "PUT", body: payload },
      ),
    remove: (postId: number, commentId: number) =>
      request<{ success: boolean; message: string }>(
        `/api/posts/${postId}/comments/${commentId}`,
        { method: "DELETE" },
      ),
  },
  follow: {
    requests: (direction: "incoming" | "outgoing" = "incoming", params: PageOptions = {}) =>
      request<Paginated<"requests", FollowRequestSummary>>("/api/follow/requests", {
        query: { direction, ...params },
      }),
    follow: (userId: number) =>
      request<{ success: boolean; message: string; status: string }>(`/api/follow/${userId}`, {
        method: "POST",
      }),
    accept: (requestId: number) =>
      request<{ success: boolean; message: string; status: string }>(
        `/api/follow/${requestId}/accept`,
        { method: "POST" },
      ),
    decline: (requestId: number) =>
      request<{ success: boolean; message: string; status: string }>(
        `/api/follow/${requestId}/decline`,
        { method: "POST" },
      ),
    unfollow: (userId: number) =>
      request<{ success: boolean; message: string; status: string }>(
        `/api/follow/${userId}/unfollow`,
        { method: "POST" },
      ),
    remove: (userId: number) =>
      request<{ success: boolean; message: string; status: string }>(
        `/api/follow/${userId}/remove`,
        { method: "POST" },
      ),
  },
  users: {
    list: (params: { q?: string } & PageOptions = {}) =>
      request<Paginated<"users", UserSummary>>("/api/users", { query: params }),
    get: (id: number) =>
      request<{ success: boolean; message: string; user: UserProfile }>(`/api/users/${id}`),
    posts: (id: number, params: PageOptions = {}) =>
      request<Paginated<"posts", Post>>(`/api/users/${id}/posts`, { query: params }),
    followers: (id: number, params: PageOptions = {}) =>
      request<Paginated<"users", UserSummary>>(`/api/users/${id}/followers`, { query: params }),
    following: (id: number, params: PageOptions = {}) =>
      request<Paginated<"users", UserSummary>>(`/api/users/${id}/following`, { query: params }),
  },
  uploads: {
    create: (file: File) => uploadFile(file),
  },
  conversations: {
    list: (params: PageOptions = {}) =>
      request<Paginated<"conversations", ConversationSummary>>("/api/conversations", {
        query: params,
      }),
    direct: (receiverId: number) =>
      request<{ success: boolean; message: string; conversation: Conversation }>(
        "/api/conversations/direct",
        { method: "POST", body: { receiver_id: receiverId } },
      ),
    messages: (
      conversationId: number,
      params: { limit?: number; before_id?: number } = {},
    ) =>
      request<Paginated<"messages", Message>>(
        `/api/conversations/${conversationId}/messages`,
        { query: params },
      ),
    markRead: (conversationId: number) =>
      request<{ success: boolean; message: string }>(
        `/api/conversations/${conversationId}/read`,
        { method: "POST" },
      ),
  },
  messages: {
    send: (conversationId: number, content: string) =>
      request<{ success: boolean; message: string; msg: Message }>("/api/messages", {
        method: "POST",
        body: { conversation_id: conversationId, content },
      }),
  },
  groups: {
    list: (params: PageOptions = {}) =>
      request<Paginated<"groups", Group>>("/api/groups", { query: params }),
    create: (payload: CreateGroupPayload) =>
      request<{ success: boolean; message: string; group: Group }>("/api/groups", {
        method: "POST",
        body: payload,
      }),
    get: (id: number) =>
      request<{ success: boolean; message: string; group: Group }>(`/api/groups/${id}`),
    updateAvatar: (id: number, avatarPath: string) =>
      request<{ success: boolean; message: string; group: Group }>(
        `/api/groups/${id}/avatar`,
        { method: "PUT", body: { avatar_path: avatarPath } },
      ),
    join: (groupId: number) =>
      request<{ success: boolean; message: string }>("/api/groups/join", {
        method: "POST",
        body: { group_id: groupId },
      }),
    acceptJoin: (requestId: number) =>
      request<{ success: boolean; message: string }>(`/api/groups/join/${requestId}/accept`, {
        method: "POST",
      }),
    declineJoin: (requestId: number) =>
      request<{ success: boolean; message: string }>(`/api/groups/join/${requestId}/decline`, {
        method: "POST",
      }),
    invite: (groupId: number, inviteeId: number) =>
      request<{ success: boolean; message: string }>("/api/groups/invitations", {
        method: "POST",
        body: { group_id: groupId, invitee_id: inviteeId },
      }),
    acceptInvite: (id: number) =>
      request<{ success: boolean; message: string }>(`/api/groups/invitations/${id}/accept`, {
        method: "POST",
      }),
    declineInvite: (id: number) =>
      request<{ success: boolean; message: string }>(`/api/groups/invitations/${id}/decline`, {
        method: "POST",
      }),
    posts: (groupId: number, params: PageOptions = {}) =>
      request<Paginated<"posts", Post>>(`/api/groups/${groupId}/posts`, { query: params }),
    createPost: (groupId: number, payload: CreatePostPayload) =>
      request<{ success: boolean; message: string; post: Post }>(`/api/groups/${groupId}/posts`, {
        method: "POST",
        body: payload,
      }),
    events: (groupId: number, params: PageOptions = {}) =>
      request<Paginated<"events", GroupEvent>>(`/api/groups/${groupId}/events`, {
        query: params,
      }),
    createEvent: (groupId: number, payload: CreateEventPayload) =>
      request<{ success: boolean; message: string; event: GroupEvent }>(
        `/api/groups/${groupId}/events`,
        { method: "POST", body: payload },
      ),
    rsvp: (groupId: number, eventId: number, response: RsvpResponse) =>
      request<{ success: boolean; message: string; rsvp: Rsvp }>(
        `/api/groups/${groupId}/events/${eventId}/rsvp`,
        { method: "POST", body: { response } },
      ),
  },
  notifications: {
    list: (params: PageOptions = {}) =>
      request<Paginated<"notifications", Notification>>("/api/notifications", { query: params }),
    unreadCount: () =>
      request<{ success: boolean; message: string; unread_count: number }>(
        "/api/notifications/unread-count",
      ),
    markRead: (id: number) =>
      request<{ success: boolean; message: string }>(`/api/notifications/${id}/read`, {
        method: "POST",
      }),
    markAllRead: () =>
      request<{ success: boolean; message: string }>("/api/notifications/read-all", {
        method: "POST",
      }),
  },
};

export type { User };
