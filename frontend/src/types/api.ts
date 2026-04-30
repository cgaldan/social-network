export type ISODateString = string;

export type User = {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  date_of_birth: ISODateString;
  nickname: string;
  gender: string;
  avatar_path: string;
  about_me: string;
  following_count: number;
  followers_count: number;
  is_online: boolean;
  is_public: boolean;
  created_at: ISODateString;
  last_seen: ISODateString;
};

export type Post = {
  id: number;
  user_id: number;
  group_id?: number;
  title: string;
  content: string;
  category: string;
  privacy_level: "public" | "almost_private" | "private" | string;
  media_url?: string;
  like_count: number;
  comment_count: number;
  created_at: ISODateString;
  updated_at: ISODateString;
  author: string;
};

export type Comment = {
  id: number;
  post_id: number;
  user_id: number;
  content: string;
  media_url?: string;
  created_at: ISODateString;
  updated_at: ISODateString;
  author: string;
};

export type Conversation = {
  id: number;
  name: string;
  type: string;
  created_at: ISODateString;
};

export type Message = {
  id: number;
  conversation_id: number;
  sender_id: number;
  content: string;
  created_at: ISODateString;
};

export type Group = {
  id: number;
  creator_id: number;
  name: string;
  description: string;
  conversation_id: number;
  created_at: ISODateString;
};

export type GroupEvent = {
  id: number;
  group_id: number;
  creator_id: number;
  title: string;
  description: string;
  starts_at: ISODateString;
  created_at: ISODateString;
};

export type GroupEventRSVP = {
  id: number;
  event_id: number;
  user_id: number;
  response: string;
  created_at: ISODateString;
  updated_at: ISODateString;
};

export type Notification = {
  id: number;
  recipient_id: number;
  actor_id?: number;
  type: string;
  title: string;
  body: string;
  entity_type?: string;
  entity_id?: number;
  action_url?: string;
  metadata?: string;
  read_at?: ISODateString;
  created_at: ISODateString;
};

export type RegisterRequest = {
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
};

export type LoginRequest = {
  identifier: string;
  password: string;
};

export type UpdateUserRequest = Omit<RegisterRequest, "password">;

export type CreatePostRequest = {
  title: string;
  content: string;
  category: string;
  privacy_level: string;
  media_url?: string;
};

export type UpdatePostRequest = CreatePostRequest;

export type CreateCommentRequest = {
  content: string;
  media_url?: string;
};

export type DirectConversationRequest = {
  receiver_id: number;
};

export type SendMessageRequest = {
  conversation_id: number;
  content: string;
};

export type CreateGroupRequest = {
  title: string;
  description: string;
  conversation_id?: number;
};

export type JoinGroupRequest = {
  group_id: number;
};

export type InviteToGroupRequest = {
  group_id: number;
  invitee_id: number;
};

export type CreateGroupEventRequest = {
  title: string;
  description: string;
  starts_at: string;
};

export type GroupEventRSVPRequest = {
  response: string;
};

export type ApiResponse = {
  success: boolean;
  message?: string;
};

export type AuthResponse = ApiResponse & {
  user?: User;
  token?: string;
};

export type PostsResponse = ApiResponse & {
  posts?: Post[];
};

export type PostResponse = ApiResponse & {
  post?: Post;
};

export type CommentResponse = ApiResponse & {
  comment?: Comment;
};

export type FollowResponse = ApiResponse & {
  status?: string;
};

export type ConversationResponse = ApiResponse & {
  conversation?: Conversation;
};

export type MessageResponse = ApiResponse & {
  msg?: Message;
};

export type GroupsResponse = ApiResponse & {
  groups?: Group[];
};

export type GroupResponse = ApiResponse & {
  group?: Group;
};

export type GroupEventsResponse = ApiResponse & {
  events?: GroupEvent[];
};

export type GroupEventResponse = ApiResponse & {
  event?: GroupEvent;
};

export type GroupEventRSVPResponse = ApiResponse & {
  rsvp?: GroupEventRSVP;
};

export type NotificationsResponse = ApiResponse & {
  notifications?: Notification[];
};

export type NotificationResponse = ApiResponse & {
  notification?: Notification;
};

export type NotificationUnreadCountResponse = ApiResponse & {
  unread_count: number;
};
