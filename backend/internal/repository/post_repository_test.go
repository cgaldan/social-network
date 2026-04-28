package repository

import (
	"social-network/internal/domain"
	"testing"
	"time"
)

func TestPostRepository_CreatePost(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	postRepo := repos.Post

	userID, err := userRepo.CreateUser("test@example.com", "hashedpass", "John", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "testuser", "male", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	postID, err := postRepo.CreatePost(int(userID), "Test Post", "This is a test post content.", "General", "public", "", 0)
	if err != nil {
		t.Fatalf("Failed to create post: %v", err)
	}

	if postID == 0 {
		t.Error("Expected non-zero post ID")
	}
}

func TestPostRepository_GetPostByID(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	postRepo := repos.Post

	userID, _ := userRepo.CreateUser("test@example.com", "hashedpass", "John", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "testuser", "male", "", "", false)
	postID, _ := postRepo.CreatePost(int(userID), "Test Post", "This is a test post content.", "General", "public", "", 0)

	post, err := postRepo.GetPostByID(int(postID))
	if err != nil {
		t.Fatalf("Failed to get post: %v", err)
	}

	if post.Title != "Test Post" {
		t.Errorf("Expected title 'Test Post', got '%s'", post.Title)
	}
	if post.Content != "This is a test post content." {
		t.Errorf("Expected content 'This is a test post content.', got '%s'", post.Content)
	}
	if post.PrivacyLevel != "public" {
		t.Errorf("Expected privacy level 'public', got '%s'", post.PrivacyLevel)
	}
}

func TestPostRepository_ListPosts(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	postRepo := repos.Post

	userID, _ := userRepo.CreateUser("test@example.com", "hashedpass", "John", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "testuser", "male", "", "", false)
	postRepo.CreatePost(int(userID), "Test Post 1", "Content 1", "General", "public", "", 0)
	postRepo.CreatePost(int(userID), "Test Post 2", "Content 2", "General", "public", "", 0)
	postRepo.CreatePost(int(userID), "Test Post 3", "Content 3", "General", "private", "", 0)
	postRepo.CreatePost(int(userID), "Test Post 4", "Content 4", "General", "almost-private", "", 0)

	posts, err := postRepo.ListPosts("", 10, 0)
	if err != nil {
		t.Fatalf("Failed to list posts: %v", err)
	}

	if len(posts) != 2 {
		t.Errorf("Expected 2 public posts, got %d", len(posts))
	}
}

func TestPostRepository_GetPostsByUserID(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	postRepo := repos.Post

	userID1, _ := userRepo.CreateUser("user1@example.com", "hashedpass1", "User", "One", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "user1", "male", "", "", false)
	userID2, _ := userRepo.CreateUser("user2@example.com", "hashedpass2", "User", "Two", time.Date(1985, time.June, 1, 0, 0, 0, 0, time.UTC), "user2", "female", "", "", false)

	postRepo.CreatePost(int(userID1), "User 1 Post", "Content for user 1", "General", "public", "", 0)
	postRepo.CreatePost(int(userID2), "User 2 Post", "Content for user 2", "General", "public", "", 0)

	posts, err := postRepo.GetPostsByUserID(int(userID1), 10, 0)
	if err != nil {
		t.Fatalf("Failed to get posts by user ID: %v", err)
	}

	if len(posts) != 1 {
		t.Errorf("Expected 1 post for user 1, got %d", len(posts))
	}
	if posts[0].Title != "User 1 Post" {
		t.Errorf("Expected title 'User 1 Post', got '%s'", posts[0].Title)
	}
}

func TestPostRepository_ListPostsByGroupID(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	postRepo := repos.Post
	groupRepo := repos.Group
	convRepo := repos.Conversation

	userID, _ := userRepo.CreateUser("test@example.com", "hashedpass", "John", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "testuser", "male", "", "", false)

	conversation, err := convRepo.CreateGroupConversation("Test Group", int(userID))
	if err != nil {
		t.Fatalf("Failed to create conversation: %v", err)
	}
	groupID, err := groupRepo.CreateGroup(&domain.Group{
		CreatorID:      int(userID),
		Title:          "Test Group",
		Description:    "Test description",
		ConversationID: conversation.ID,
	})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	postRepo.CreatePost(int(userID), "Group Post 1", "Content 1 of group post", "general", "public", "", int(groupID))
	postRepo.CreatePost(int(userID), "Group Post 2", "Content 2 of group post", "general", "public", "", int(groupID))
	postRepo.CreatePost(int(userID), "Global Post", "Content of global post", "general", "public", "", 0)

	groupPosts, err := postRepo.ListPostsByGroupID(int(groupID), 10, 0)
	if err != nil {
		t.Fatalf("Failed to list group posts: %v", err)
	}
	if len(groupPosts) != 2 {
		t.Errorf("Expected 2 group posts, got %d", len(groupPosts))
	}
	for _, post := range groupPosts {
		if post.GroupID != int(groupID) {
			t.Errorf("Expected post group ID %d, got %d", groupID, post.GroupID)
		}
	}

	globalPosts, err := postRepo.ListPosts("", 10, 0)
	if err != nil {
		t.Fatalf("Failed to list posts: %v", err)
	}
	if len(globalPosts) != 1 {
		t.Errorf("Expected 1 global post, got %d", len(globalPosts))
	}
	if globalPosts[0].Title != "Global Post" {
		t.Errorf("Expected global post title 'Global Post', got '%s'", globalPosts[0].Title)
	}

	userPosts, err := postRepo.GetPostsByUserID(int(userID), 10, 0)
	if err != nil {
		t.Fatalf("Failed to get posts by user ID: %v", err)
	}
	if len(userPosts) != 1 {
		t.Errorf("Expected 1 non-group post for user, got %d", len(userPosts))
	}
}

func TestPostRepository_GetPostByID_GroupID(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	postRepo := repos.Post
	groupRepo := repos.Group
	convRepo := repos.Conversation

	userID, _ := userRepo.CreateUser("test@example.com", "hashedpass", "John", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "testuser", "male", "", "", false)

	conversation, _ := convRepo.CreateGroupConversation("Test Group", int(userID))
	groupID, _ := groupRepo.CreateGroup(&domain.Group{
		CreatorID:      int(userID),
		Title:          "Test Group",
		Description:    "Test description",
		ConversationID: conversation.ID,
	})

	groupPostID, _ := postRepo.CreatePost(int(userID), "Group Post", "This is a group post content", "general", "public", "", int(groupID))
	globalPostID, _ := postRepo.CreatePost(int(userID), "Global Post", "This is a global post content", "general", "public", "", 0)

	groupPost, err := postRepo.GetPostByID(int(groupPostID))
	if err != nil {
		t.Fatalf("Failed to get group post: %v", err)
	}
	if groupPost.GroupID != int(groupID) {
		t.Errorf("Expected group ID %d, got %d", groupID, groupPost.GroupID)
	}

	globalPost, err := postRepo.GetPostByID(int(globalPostID))
	if err != nil {
		t.Fatalf("Failed to get global post: %v", err)
	}
	if globalPost.GroupID != 0 {
		t.Errorf("Expected global post group ID 0, got %d", globalPost.GroupID)
	}

	groupPosts, err := postRepo.ListPostsByGroupID(int(groupID), 10, 0)
	if err != nil {
		t.Fatalf("Failed to list group posts: %v", err)
	}
	if len(groupPosts) != 1 {
		t.Errorf("Expected 1 group post, got %d", len(groupPosts))
	}
	if groupPosts[0].Title != "Group Post" {
		t.Errorf("Expected group post title 'Group Post', got '%s'", groupPosts[0].Title)
	}

	globalPosts, err := postRepo.ListPosts("", 10, 0)
	if err != nil {
		t.Fatalf("Failed to list posts: %v", err)
	}
	if len(globalPosts) != 1 {
		t.Errorf("Expected 1 global post, got %d", len(globalPosts))
	}
	if globalPosts[0].Title != "Global Post" {
		t.Errorf("Expected global post title 'Global Post', got '%s'", globalPosts[0].Title)
	}
}

func TestPostRepository_PostExists(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	postRepo := repos.Post

	userID, _ := userRepo.CreateUser("test@example.com", "hashedpass", "John", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "testuser", "male", "", "", false)
	postID, _ := postRepo.CreatePost(int(userID), "Test Post", "This is a test post content.", "General", "public", "", 0)

	exists, err := postRepo.PostExists(int(postID))
	if err != nil {
		t.Fatalf("Failed to check if post exists: %v", err)
	}

	if !exists {
		t.Error("Expected post to exist")
	}
}

func TestPostRepository_UpdatePost(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	postRepo := repos.Post

	userID, _ := userRepo.CreateUser("test@example.com", "hashedpass", "John", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "testuser", "male", "", "", false)
	postID, _ := postRepo.CreatePost(int(userID), "Old Title", "This is a test post content for update.", "General", "public", "", 0)

	err := postRepo.UpdatePost(int(userID), int(postID), "New Title", "This is a test post content for update new body text here.", "Tech", "private", "")
	if err != nil {
		t.Fatalf("UpdatePost: %v", err)
	}

	post, err := postRepo.GetPostByID(int(postID))
	if err != nil {
		t.Fatalf("GetPostByID: %v", err)
	}
	if post.Title != "New Title" {
		t.Errorf("title: got %q", post.Title)
	}
	if post.PrivacyLevel != "private" {
		t.Errorf("privacy: got %q", post.PrivacyLevel)
	}
}

func TestPostRepository_DeletePost(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	postRepo := repos.Post

	userID, _ := userRepo.CreateUser("test@example.com", "hashedpass", "John", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "testuser", "male", "", "", false)
	postID, _ := postRepo.CreatePost(int(userID), "To Delete", "This is a test post content for delete action.", "General", "public", "", 0)

	if err := postRepo.DeletePost(int(userID), int(postID)); err != nil {
		t.Fatalf("DeletePost: %v", err)
	}
	_, err := postRepo.GetPostByID(int(postID))
	if err == nil {
		t.Error("expected post to be gone")
	}
}
