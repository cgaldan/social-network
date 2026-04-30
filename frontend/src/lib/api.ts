import { getToken } from "@/lib/auth";
import type {
  ApiResponse,
  AuthResponse,
  CommentResponse,
  ConversationResponse,
  CreateCommentRequest,
  CreateGroupEventRequest,
  CreateGroupRequest,
  CreatePostRequest,
  DirectConversationRequest,
  FollowResponse,
  GroupEventsResponse,
  GroupEventResponse,
  GroupEventRSVPRequest,
  GroupEventRSVPResponse,
  GroupResponse,
  GroupsResponse,
  InviteToGroupRequest,
  JoinGroupRequest,
  LoginRequest,
  MessageResponse,
  NotificationResponse,
  NotificationsResponse,
  NotificationUnreadCountResponse,
  PostResponse,
  PostsResponse,
  RegisterRequest,
  SendMessageRequest,
  UpdatePostRequest,
  UpdateUserRequest,
} from "@/types/api";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8000";

type RequestOptions = {
  method?: "GET" | "POST" | "PUT" | "DELETE";
  body?: unknown;
  token?: string | null;
  auth?: boolean;
};

export class ApiError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "ApiError";
  }
}

function buildUrl(path: string) {
  if (path.startsWith("http")) {
    return path;
  }

  return `${API_URL}${path}`;
}

export async function apiRequest<T extends ApiResponse>(
  path: string,
  options: RequestOptions = {},
): Promise<T> {
  const token = options.token ?? (options.auth === false ? null : getToken());
  const headers = new Headers();

  headers.set("Content-Type", "application/json");
  if (token) {
    headers.set("Authorization", token);
  }

  const response = await fetch(buildUrl(path), {
    method: options.method ?? "GET",
    headers,
    body: options.body ? JSON.stringify(options.body) : undefined,
  });

  const data = (await response.json().catch(() => ({
    success: false,
    message: "Invalid server response",
  }))) as T;

  if (!response.ok || data.success === false) {
    throw new ApiError(data.message || "Request failed");
  }

  return data;
}

export const api = {
  register: (body: RegisterRequest) =>
    apiRequest<AuthResponse>("/api/auth/register", {
      method: "POST",
      body,
      auth: false,
    }),
  login: (body: LoginRequest) =>
    apiRequest<AuthResponse>("/api/auth/login", {
      method: "POST",
      body,
      auth: false,
    }),
  logout: () =>
    apiRequest<AuthResponse>("/api/auth/logout", {
      method: "POST",
    }),
  me: () => apiRequest<AuthResponse>("/api/auth/me"),
  updateMe: (body: UpdateUserRequest) =>
    apiRequest<AuthResponse>("/api/auth/me", {
      method: "PUT",
      body,
    }),
  deleteMe: () =>
    apiRequest<AuthResponse>("/api/auth/me", {
      method: "DELETE",
    }),
  listPosts: (params?: { category?: string; limit?: number; offset?: number }) => {
    const search = new URLSearchParams();
    if (params?.category) search.set("category", params.category);
    if (params?.limit) search.set("limit", String(params.limit));
    if (params?.offset) search.set("offset", String(params.offset));

    return apiRequest<PostsResponse>(`/api/posts?${search.toString()}`, {
      auth: false,
    });
  },
  createPost: (body: CreatePostRequest) =>
    apiRequest<PostResponse>("/api/posts", {
      method: "POST",
      body,
    }),
  getPost: (id: number) => apiRequest<PostResponse>(`/api/posts/${id}`),
  updatePost: (id: number, body: UpdatePostRequest) =>
    apiRequest<PostResponse>(`/api/posts/${id}`, {
      method: "PUT",
      body,
    }),
  deletePost: (id: number) =>
    apiRequest<PostResponse>(`/api/posts/${id}`, {
      method: "DELETE",
    }),
  createComment: (postId: number, body: CreateCommentRequest) =>
    apiRequest<CommentResponse>(`/api/posts/${postId}/comments`, {
      method: "POST",
      body,
    }),
  updateComment: (
    postId: number,
    commentId: number,
    body: CreateCommentRequest,
  ) =>
    apiRequest<CommentResponse>(`/api/posts/${postId}/comments/${commentId}`, {
      method: "PUT",
      body,
    }),
  deleteComment: (postId: number, commentId: number) =>
    apiRequest<CommentResponse>(`/api/posts/${postId}/comments/${commentId}`, {
      method: "DELETE",
    }),
  followUser: (id: number) =>
    apiRequest<FollowResponse>(`/api/follow/${id}`, { method: "POST" }),
  acceptFollow: (id: number) =>
    apiRequest<FollowResponse>(`/api/follow/${id}/accept`, { method: "POST" }),
  declineFollow: (id: number) =>
    apiRequest<FollowResponse>(`/api/follow/${id}/decline`, { method: "POST" }),
  unfollowUser: (id: number) =>
    apiRequest<FollowResponse>(`/api/follow/${id}/unfollow`, {
      method: "POST",
    }),
  removeFollower: (id: number) =>
    apiRequest<FollowResponse>(`/api/follow/${id}/remove`, {
      method: "POST",
    }),
  createDirectConversation: (body: DirectConversationRequest) =>
    apiRequest<ConversationResponse>("/api/conversations/direct", {
      method: "POST",
      body,
    }),
  sendMessage: (body: SendMessageRequest) =>
    apiRequest<MessageResponse>("/api/messages", {
      method: "POST",
      body,
    }),
  listGroups: (params?: { limit?: number; offset?: number }) => {
    const search = new URLSearchParams();
    if (params?.limit) search.set("limit", String(params.limit));
    if (params?.offset) search.set("offset", String(params.offset));

    return apiRequest<GroupsResponse>(`/api/groups?${search.toString()}`);
  },
  createGroup: (body: CreateGroupRequest) =>
    apiRequest<GroupResponse>("/api/groups", {
      method: "POST",
      body,
    }),
  joinGroup: (body: JoinGroupRequest) =>
    apiRequest<GroupResponse>("/api/groups/join", {
      method: "POST",
      body,
    }),
  inviteToGroup: (body: InviteToGroupRequest) =>
    apiRequest<GroupResponse>("/api/groups/invitations", {
      method: "POST",
      body,
    }),
  acceptGroupJoin: (id: number) =>
    apiRequest<GroupResponse>(`/api/groups/join/${id}/accept`, {
      method: "POST",
    }),
  declineGroupJoin: (id: number) =>
    apiRequest<GroupResponse>(`/api/groups/join/${id}/decline`, {
      method: "POST",
    }),
  acceptGroupInvitation: (id: number) =>
    apiRequest<GroupResponse>(`/api/groups/invitations/${id}/accept`, {
      method: "POST",
    }),
  declineGroupInvitation: (id: number) =>
    apiRequest<GroupResponse>(`/api/groups/invitations/${id}/decline`, {
      method: "POST",
    }),
  listGroupPosts: (groupId: number) =>
    apiRequest<PostsResponse>(`/api/groups/${groupId}/posts`),
  createGroupPost: (groupId: number, body: CreatePostRequest) =>
    apiRequest<PostResponse>(`/api/groups/${groupId}/posts`, {
      method: "POST",
      body,
    }),
  listGroupEvents: (groupId: number) =>
    apiRequest<GroupEventsResponse>(`/api/groups/${groupId}/events`),
  createGroupEvent: (groupId: number, body: CreateGroupEventRequest) =>
    apiRequest<GroupEventResponse>(`/api/groups/${groupId}/events`, {
      method: "POST",
      body,
    }),
  rsvpGroupEvent: (
    groupId: number,
    eventId: number,
    body: GroupEventRSVPRequest,
  ) =>
    apiRequest<GroupEventRSVPResponse>(
      `/api/groups/${groupId}/events/${eventId}/rsvp`,
      {
        method: "POST",
        body,
      },
    ),
  listNotifications: () =>
    apiRequest<NotificationsResponse>("/api/notifications"),
  unreadNotifications: () =>
    apiRequest<NotificationUnreadCountResponse>(
      "/api/notifications/unread-count",
    ),
  markNotificationRead: (id: number) =>
    apiRequest<NotificationResponse>(`/api/notifications/${id}/read`, {
      method: "POST",
    }),
  markAllNotificationsRead: () =>
    apiRequest<NotificationsResponse>("/api/notifications/read-all", {
      method: "POST",
    }),
};
