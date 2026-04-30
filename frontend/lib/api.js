const FALLBACK_API_BASE_URL = "http://localhost:8000";

export const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL || FALLBACK_API_BASE_URL;

export function getStoredToken() {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("authToken");
}

/** Backend expects the raw session id in Authorization (not a Bearer prefix). */
export function authHeader(token) {
  const t = token ?? getStoredToken();
  if (!t) throw new Error("Not authenticated");
  return { Authorization: t };
}

async function parseJsonResponse(response) {
  let payload = null;
  try {
    payload = await response.json();
  } catch {
    throw new Error("Unable to parse server response.");
  }
  if (!response.ok || payload.success === false) {
    throw new Error(payload?.message || "Request failed.");
  }
  return payload;
}

export function getWebSocketUrl(token) {
  const t = token ?? getStoredToken();
  if (!t) throw new Error("Not authenticated");
  const wsBase = API_BASE_URL.replace(/^http/, "ws");
  return `${wsBase}/ws?token=${encodeURIComponent(t)}`;
}

// --- Auth (public) ---

export async function login({ identifier, password }) {
  const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ identifier, password }),
  });
  return parseJsonResponse(response);
}

export async function register({
  email,
  password,
  firstName,
  lastName,
  dateOfBirth,
  nickname,
  gender,
  aboutMe = "",
  isPublic = true,
  avatarPath = "",
}) {
  const response = await fetch(`${API_BASE_URL}/api/auth/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      email,
      password,
      first_name: firstName,
      last_name: lastName,
      date_of_birth: `${dateOfBirth}T00:00:00Z`,
      nickname,
      gender,
      about_me: aboutMe,
      is_public: isPublic,
      avatar_path: avatarPath,
    }),
  });
  return parseJsonResponse(response);
}

export async function logout(token) {
  const response = await fetch(`${API_BASE_URL}/api/auth/logout`, {
    method: "POST",
    headers: { ...authHeader(token), "Content-Type": "application/json" },
  });
  return parseJsonResponse(response);
}

export async function getCurrentUser(token) {
  const response = await fetch(`${API_BASE_URL}/api/auth/me`, {
    headers: authHeader(token),
  });
  return parseJsonResponse(response);
}

export async function updateCurrentUser(token, body) {
  const response = await fetch(`${API_BASE_URL}/api/auth/me`, {
    method: "PUT",
    headers: { ...authHeader(token), "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  return parseJsonResponse(response);
}

export async function deleteCurrentUser(token) {
  const response = await fetch(`${API_BASE_URL}/api/auth/me`, {
    method: "DELETE",
    headers: { ...authHeader(token), "Content-Type": "application/json" },
  });
  return parseJsonResponse(response);
}

// --- Posts ---

export async function getPosts({ category = "", limit = 20, offset = 0 } = {}) {
  const params = new URLSearchParams();
  if (category) params.set("category", category);
  params.set("limit", String(limit));
  params.set("offset", String(offset));
  const response = await fetch(
    `${API_BASE_URL}/api/posts?${params.toString()}`,
  );
  return parseJsonResponse(response);
}

export async function createPost(token, body) {
  const response = await fetch(`${API_BASE_URL}/api/posts`, {
    method: "POST",
    headers: { ...authHeader(token), "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  return parseJsonResponse(response);
}

export async function getPost(token, postId) {
  const response = await fetch(`${API_BASE_URL}/api/posts/${postId}`, {
    headers: authHeader(token),
  });
  return parseJsonResponse(response);
}

export async function updatePost(token, postId, body) {
  const response = await fetch(`${API_BASE_URL}/api/posts/${postId}`, {
    method: "PUT",
    headers: { ...authHeader(token), "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  return parseJsonResponse(response);
}

export async function deletePost(token, postId) {
  const response = await fetch(`${API_BASE_URL}/api/posts/${postId}`, {
    method: "DELETE",
    headers: { ...authHeader(token), "Content-Type": "application/json" },
  });
  return parseJsonResponse(response);
}

// --- Comments ---

export async function createComment(token, postId, body) {
  const response = await fetch(
    `${API_BASE_URL}/api/posts/${postId}/comments`,
    {
      method: "POST",
      headers: { ...authHeader(token), "Content-Type": "application/json" },
      body: JSON.stringify(body),
    },
  );
  return parseJsonResponse(response);
}

export async function updateComment(token, postId, commentId, body) {
  const response = await fetch(
    `${API_BASE_URL}/api/posts/${postId}/comments/${commentId}`,
    {
      method: "PUT",
      headers: { ...authHeader(token), "Content-Type": "application/json" },
      body: JSON.stringify(body),
    },
  );
  return parseJsonResponse(response);
}

export async function deleteComment(token, postId, commentId) {
  const response = await fetch(
    `${API_BASE_URL}/api/posts/${postId}/comments/${commentId}`,
    {
      method: "DELETE",
      headers: { ...authHeader(token), "Content-Type": "application/json" },
    },
  );
  return parseJsonResponse(response);
}

// --- Follow ---

export async function followUser(token, userId) {
  const response = await fetch(`${API_BASE_URL}/api/follow/${userId}`, {
    method: "POST",
    headers: { ...authHeader(token), "Content-Type": "application/json" },
  });
  return parseJsonResponse(response);
}

export async function acceptFollowRequest(token, followId) {
  const response = await fetch(
    `${API_BASE_URL}/api/follow/${followId}/accept`,
    {
      method: "POST",
      headers: { ...authHeader(token), "Content-Type": "application/json" },
    },
  );
  return parseJsonResponse(response);
}

export async function declineFollowRequest(token, followId) {
  const response = await fetch(
    `${API_BASE_URL}/api/follow/${followId}/decline`,
    {
      method: "POST",
      headers: { ...authHeader(token), "Content-Type": "application/json" },
    },
  );
  return parseJsonResponse(response);
}

export async function unfollowUser(token, userId) {
  const response = await fetch(
    `${API_BASE_URL}/api/follow/${userId}/unfollow`,
    {
      method: "POST",
      headers: { ...authHeader(token), "Content-Type": "application/json" },
    },
  );
  return parseJsonResponse(response);
}

export async function removeFollower(token, followerId) {
  const response = await fetch(
    `${API_BASE_URL}/api/follow/${followerId}/remove`,
    {
      method: "POST",
      headers: { ...authHeader(token), "Content-Type": "application/json" },
    },
  );
  return parseJsonResponse(response);
}

// --- Conversations & messages ---

export async function createDirectConversation(token, receiverId) {
  const response = await fetch(`${API_BASE_URL}/api/conversations/direct`, {
    method: "POST",
    headers: { ...authHeader(token), "Content-Type": "application/json" },
    body: JSON.stringify({ receiver_id: receiverId }),
  });
  return parseJsonResponse(response);
}

export async function sendMessage(token, { conversationId, content }) {
  const response = await fetch(`${API_BASE_URL}/api/messages`, {
    method: "POST",
    headers: { ...authHeader(token), "Content-Type": "application/json" },
    body: JSON.stringify({
      conversation_id: conversationId,
      content,
    }),
  });
  return parseJsonResponse(response);
}

// --- Groups ---

export async function listGroups(token, { limit = 20, offset = 0 } = {}) {
  const params = new URLSearchParams({
    limit: String(limit),
    offset: String(offset),
  });
  const response = await fetch(
    `${API_BASE_URL}/api/groups?${params.toString()}`,
    { headers: authHeader(token) },
  );
  return parseJsonResponse(response);
}

export async function createGroup(token, { title, description, conversationId = 0 }) {
  const response = await fetch(`${API_BASE_URL}/api/groups`, {
    method: "POST",
    headers: { ...authHeader(token), "Content-Type": "application/json" },
    body: JSON.stringify({
      title,
      description,
      conversation_id: conversationId,
    }),
  });
  return parseJsonResponse(response);
}

export async function joinGroup(token, groupId) {
  const response = await fetch(`${API_BASE_URL}/api/groups/join`, {
    method: "POST",
    headers: { ...authHeader(token), "Content-Type": "application/json" },
    body: JSON.stringify({ group_id: groupId }),
  });
  return parseJsonResponse(response);
}

export async function acceptGroupJoinRequest(token, requestId) {
  const response = await fetch(
    `${API_BASE_URL}/api/groups/join/${requestId}/accept`,
    {
      method: "POST",
      headers: { ...authHeader(token), "Content-Type": "application/json" },
    },
  );
  return parseJsonResponse(response);
}

export async function declineGroupJoinRequest(token, requestId) {
  const response = await fetch(
    `${API_BASE_URL}/api/groups/join/${requestId}/decline`,
    {
      method: "POST",
      headers: { ...authHeader(token), "Content-Type": "application/json" },
    },
  );
  return parseJsonResponse(response);
}

export async function inviteToGroup(token, { groupId, inviteeId }) {
  const response = await fetch(`${API_BASE_URL}/api/groups/invitations`, {
    method: "POST",
    headers: { ...authHeader(token), "Content-Type": "application/json" },
    body: JSON.stringify({
      group_id: groupId,
      invitee_id: inviteeId,
    }),
  });
  return parseJsonResponse(response);
}

export async function acceptGroupInvitation(token, invitationId) {
  const response = await fetch(
    `${API_BASE_URL}/api/groups/invitations/${invitationId}/accept`,
    {
      method: "POST",
      headers: { ...authHeader(token), "Content-Type": "application/json" },
    },
  );
  return parseJsonResponse(response);
}

export async function declineGroupInvitation(token, invitationId) {
  const response = await fetch(
    `${API_BASE_URL}/api/groups/invitations/${invitationId}/decline`,
    {
      method: "POST",
      headers: { ...authHeader(token), "Content-Type": "application/json" },
    },
  );
  return parseJsonResponse(response);
}

export async function getGroupPosts(token, groupId, { limit = 20, offset = 0 } = {}) {
  const params = new URLSearchParams({
    limit: String(limit),
    offset: String(offset),
  });
  const response = await fetch(
    `${API_BASE_URL}/api/groups/${groupId}/posts?${params.toString()}`,
    { headers: authHeader(token) },
  );
  return parseJsonResponse(response);
}

export async function createGroupPost(token, groupId, body) {
  const response = await fetch(
    `${API_BASE_URL}/api/groups/${groupId}/posts`,
    {
      method: "POST",
      headers: { ...authHeader(token), "Content-Type": "application/json" },
      body: JSON.stringify(body),
    },
  );
  return parseJsonResponse(response);
}

export async function listGroupEvents(token, groupId, { limit = 20, offset = 0 } = {}) {
  const params = new URLSearchParams({
    limit: String(limit),
    offset: String(offset),
  });
  const response = await fetch(
    `${API_BASE_URL}/api/groups/${groupId}/events?${params.toString()}`,
    { headers: authHeader(token) },
  );
  return parseJsonResponse(response);
}

export async function createGroupEvent(token, groupId, body) {
  const response = await fetch(
    `${API_BASE_URL}/api/groups/${groupId}/events`,
    {
      method: "POST",
      headers: { ...authHeader(token), "Content-Type": "application/json" },
      body: JSON.stringify(body),
    },
  );
  return parseJsonResponse(response);
}

export async function setGroupEventRSVP(token, groupId, eventId, response) {
  const res = await fetch(
    `${API_BASE_URL}/api/groups/${groupId}/events/${eventId}/rsvp`,
    {
      method: "POST",
      headers: { ...authHeader(token), "Content-Type": "application/json" },
      body: JSON.stringify({ response }),
    },
  );
  return parseJsonResponse(res);
}

// --- Notifications ---

export async function listNotifications(token, { limit = 30, offset = 0 } = {}) {
  const params = new URLSearchParams({
    limit: String(limit),
    offset: String(offset),
  });
  const response = await fetch(
    `${API_BASE_URL}/api/notifications?${params.toString()}`,
    { headers: authHeader(token) },
  );
  return parseJsonResponse(response);
}

export async function getUnreadNotificationCount(token) {
  const response = await fetch(
    `${API_BASE_URL}/api/notifications/unread-count`,
    { headers: authHeader(token) },
  );
  return parseJsonResponse(response);
}

export async function markNotificationRead(token, notificationId) {
  const response = await fetch(
    `${API_BASE_URL}/api/notifications/${notificationId}/read`,
    {
      method: "POST",
      headers: { ...authHeader(token), "Content-Type": "application/json" },
    },
  );
  return parseJsonResponse(response);
}

export async function markAllNotificationsRead(token) {
  const response = await fetch(
    `${API_BASE_URL}/api/notifications/read-all`,
    {
      method: "POST",
      headers: { ...authHeader(token), "Content-Type": "application/json" },
    },
  );
  return parseJsonResponse(response);
}
