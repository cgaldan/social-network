package service

import (
	"social-network/internal/domain"
	"testing"
	"time"
)

func TestPostService_GetPostByID(t *testing.T) {
	services := SetupTestServices(t)

	userID := CreateTestUser(t, services, "testuser", "test@example.com", "password123", "John", "Doe", time.Now().AddDate(-25, 0, 0), "male")

	post := CreateTestPost(t, services, userID, domain.CreatePostRequest{
		Title:    "Test Post",
		Content:  "This is a test post content with enough characters",
		Category: "general",
	})

	t.Run("get existing post", func(t *testing.T) {
		retrievedPost, err := services.Post.GetPostByID(post.ID)
		if err != nil {
			t.Fatalf("Failed to get post: %v", err)
		}

		if retrievedPost == nil {
			t.Fatal("Expected post but got nil")
		}

		if retrievedPost.ID != post.ID {
			t.Errorf("Expected post ID %d, got %d", post.ID, retrievedPost.ID)
		}

		if retrievedPost.Title != post.Title {
			t.Errorf("Expected title %s, got %s", post.Title, retrievedPost.Title)
		}
	})

	t.Run("get non-existent post", func(t *testing.T) {
		_, err := services.Post.GetPostByID(99999)
		if err == nil {
			t.Error("Expected error for non-existent post")
		}
	})
}

func TestPostService_ListPosts(t *testing.T) {
	services := SetupTestServices(t)

	userID := CreateTestUser(t, services, "testuser", "test@example.com", "password123", "John", "Doe", time.Now().AddDate(-25, 0, 0), "male")

	posts := []domain.CreatePostRequest{
		{Title: "Post 1", Content: "Content 1 with enough characters", Category: "general"},
		{Title: "Post 2", Content: "Content 2 with enough characters", Category: "tech"},
		{Title: "Post 3", Content: "Content 3 with enough characters", Category: "general"},
	}

	for _, postData := range posts {
		CreateTestPost(t, services, userID, postData)
	}

	t.Run("list all posts", func(t *testing.T) {
		retrievedPosts, err := services.Post.ListPosts("", 10, 0)
		if err != nil {
			t.Fatalf("Failed to list posts: %v", err)
		}

		if len(retrievedPosts) != 3 {
			t.Errorf("Expected 3 posts, got %d", len(retrievedPosts))
		}
	})

	t.Run("list posts by category", func(t *testing.T) {
		retrievedPosts, err := services.Post.ListPosts("general", 10, 0)
		if err != nil {
			t.Fatalf("Failed to list posts by category: %v", err)
		}

		if len(retrievedPosts) != 2 {
			t.Errorf("Expected 2 general posts, got %d", len(retrievedPosts))
		}

		for _, post := range retrievedPosts {
			if post.Category != "general" {
				t.Errorf("Expected category 'general', got '%s'", post.Category)
			}
		}
	})

	t.Run("list posts with pagination", func(t *testing.T) {
		retrievedPosts, err := services.Post.ListPosts("", 2, 0)
		if err != nil {
			t.Fatalf("Failed to list posts with pagination: %v", err)
		}

		if len(retrievedPosts) != 2 {
			t.Errorf("Expected 2 posts with limit 2, got %d", len(retrievedPosts))
		}
	})
}

func TestPostService_GetPostsByUserID(t *testing.T) {
	services := SetupTestServices(t)

	user1ID := CreateTestUser(t, services, "user1", "user1@example.com", "password123", "John", "Doe", time.Now().AddDate(-25, 0, 0), "male")
	user2ID := CreateTestUser(t, services, "user2", "user2@example.com", "password123", "Jane", "Smith", time.Now().AddDate(-30, 0, 0), "female")

	CreateTestPost(t, services, user1ID, domain.CreatePostRequest{
		Title: "User1 Post 1", Content: "Content 1 with enough characters for validation", Category: "general",
	})

	CreateTestPost(t, services, user1ID, domain.CreatePostRequest{
		Title: "User1 Post 2", Content: "Content 2 with enough characters for validation", Category: "tech",
	})

	CreateTestPost(t, services, user2ID, domain.CreatePostRequest{
		Title: "User2 Post", Content: "Content 3 with enough characters for validation", Category: "general",
	})

	t.Run("get posts by user ID", func(t *testing.T) {
		posts, err := services.Post.GetPostsByUserID(user1ID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to get posts by user ID: %v", err)
		}

		if len(posts) != 2 {
			t.Errorf("Expected 2 posts for user1, got %d", len(posts))
		}

		for _, post := range posts {
			if post.UserID != user1ID {
				t.Errorf("Expected user ID %d, got %d", user1ID, post.UserID)
			}
		}
	})
}
