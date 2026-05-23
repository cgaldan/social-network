export type PrivacyLevel = "public" | "almost_private" | "private";

export interface User {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  date_of_birth: string;
  nickname: string;
  gender: string;
  avatar_path: string;
  about_me: string;
  following_count: number;
  followers_count: number;
  is_online: boolean;
  is_public: boolean;
  created_at: string;
  last_seen: string;
}

export interface RegisterPayload {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  date_of_birth: string;
  nickname: string;
  gender: string;
  avatar_path: string;
  about_me: string;
  is_public: boolean;
}

export interface LoginPayload {
  identifier: string;
  password: string;
}

export interface AuthResponse {
  success: boolean;
  message: string;
  user: User;
  token: string;
}

export interface MeResponse {
  success: boolean;
  message: string;
  user: User;
}

export interface Post {
  id: number;
  user_id: number;
  group_id: number;
  title: string;
  content: string;
  category: string;
  privacy_level: PrivacyLevel;
  media_url?: string;
  like_count: number;
  comment_count: number;
  created_at: string;
  updated_at: string;
  author: string;
}

export interface CreatePostPayload {
  title: string;
  content: string;
  category: string;
  privacy_level: PrivacyLevel;
  media_url?: string;
  audience?: number[];
}

export type UpdatePostPayload = Partial<CreatePostPayload>;

export interface Comment {
  id: number;
  post_id: number;
  user_id: number;
  content: string;
  media_url?: string;
  created_at: string;
  updated_at: string;
  author: string;
}

export interface CreateCommentPayload {
  content: string;
  media_url?: string;
}

export interface Conversation {
  id: number;
  name: string;
  type: "direct" | "group" | "private";
  created_at: string;
}

export interface ConversationParticipant {
  user_id: number;
  nickname: string;
  first_name: string;
  last_name: string;
  avatar_path?: string;
  is_online: boolean;
}

export interface ConversationSummary {
  id: number;
  type: string;
  name: string;
  avatar_path?: string;
  participants: ConversationParticipant[];
  last_message?: Message;
  unread_count: number;
  created_at: string;
}

export interface Message {
  id: number;
  conversation_id: number;
  sender_id: number;
  content: string;
  created_at: string;
  updated_at?: string;
}

export interface MessageUpdatedPayload {
  message: Message;
}

export interface MessageDeletedPayload {
  conversation_id: number;
  message_id: number;
}

export interface UserSummary {
  id: number;
  nickname: string;
  first_name: string;
  last_name: string;
  avatar_path?: string;
  is_online: boolean;
  is_public: boolean;
}

export interface FollowRequestSummary {
  id: number;
  status: string;
  created_at: string;
  user: UserSummary;
}

export type FollowStatus = "self" | "following" | "pending" | "none";

export interface UserProfile {
  id: number;
  nickname: string;
  first_name: string;
  last_name: string;
  avatar_path?: string;
  is_online: boolean;
  is_public: boolean;
  created_at: string;
  follow_status: FollowStatus;
  can_view: boolean;
  followers_count: number;
  following_count: number;

  // Present only when can_view is true:
  email?: string;
  date_of_birth?: string;
  gender?: string;
  about_me?: string;
}

export interface MessageCreatedPayload {
  message: Message;
}

export interface Group {
  id: number;
  creator_id: number;
  name: string;
  description: string;
  conversation_id: number;
  avatar_path?: string;
  created_at: string;
  is_member: boolean;
  is_pending: boolean;
}

export interface CreateGroupPayload {
  title: string;
  description: string;
  conversation_id: number;
}

export interface GroupEvent {
  id: number;
  group_id: number;
  creator_id: number;
  title: string;
  description: string;
  starts_at: string;
  created_at: string;
  going_count: number;
  maybe_count: number;
  not_going_count: number;
  my_response?: RsvpResponse;
}

export interface CreateEventPayload {
  title: string;
  description: string;
  starts_at: string;
}

export type RsvpResponse = "going" | "not_going" | "maybe";

export interface Rsvp {
  id: number;
  event_id: number;
  user_id: number;
  response: RsvpResponse;
  created_at: string;
  updated_at: string;
}

export interface Notification {
  id: number;
  recipient_id: number;
  actor_id?: number;
  type: string;
  title: string;
  body: string;
  entity_type?: string;
  entity_id?: number;
  entity_status?: string;
  action_url?: string;
  metadata?: string | null;
  read_at?: string | null;
  created_at: string;
}

export interface ApiEnvelope {
  success: boolean;
  message: string;
}

export interface WsEnvelope<T = unknown> {
  type: string;
  payload: T;
}

export interface OnlineUser {
  user_id: number;
  nickname: string;
  is_online: boolean;
}

export interface UserStatusPayload {
  user_id: number;
  nickname: string;
  is_online: boolean;
}

export interface OnlineUsersPayload {
  users: OnlineUser[];
}

export interface NotificationCreatedPayload {
  notification: Notification;
}
