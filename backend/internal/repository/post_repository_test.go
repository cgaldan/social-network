package repository

import (
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

	postID, err := postRepo.CreatePost(int(userID), "Test Post", "This is a test post content.", "General", "public", "")
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
	postID, _ := postRepo.CreatePost(int(userID), "Test Post", "This is a test post content.", "General", "public", "")

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
	postRepo.CreatePost(int(userID), "Test Post 1", "Content 1", "General", "public", "")
	postRepo.CreatePost(int(userID), "Test Post 2", "Content 2", "General", "public", "")
	postRepo.CreatePost(int(userID), "Test Post 3", "Content 3", "General", "private", "")
	postRepo.CreatePost(int(userID), "Test Post 4", "Content 4", "General", "almost-private", "")

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

	postRepo.CreatePost(int(userID1), "User 1 Post", "Content for user 1", "General", "public", "")
	postRepo.CreatePost(int(userID2), "User 2 Post", "Content for user 2", "General", "public", "")

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

func TestPostRepository_PostExists(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	postRepo := repos.Post

	userID, _ := userRepo.CreateUser("test@example.com", "hashedpass", "John", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "testuser", "male", "", "", false)
	postID, _ := postRepo.CreatePost(int(userID), "Test Post", "This is a test post content.", "General", "public", "")

	exists, err := postRepo.PostExists(int(postID))
	if err != nil {
		t.Fatalf("Failed to check if post exists: %v", err)
	}

	if !exists {
		t.Error("Expected post to exist")
	}
}
