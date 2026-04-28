package service

import (
	"social-network/internal/domain"
	"testing"
	"time"
)

func TestCommentService_CreateComment(t *testing.T) {
	services := SetupTestServices(t)

	userID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "test@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "testuser",
		Gender:      "male",
		IsPublic:    true,
	})

	post := CreateTestPost(t, services, userID, domain.CreatePostRequest{
		Title:    "Test Post",
		Content:  "This is a test post content",
		Category: "general",
	})

	tests := []struct {
		name        string
		commentData domain.CreateCommentRequest
		expectError bool
	}{
		{
			name: "valid comment",
			commentData: domain.CreateCommentRequest{
				Content: "This is a valid comment",
			},
			expectError: false,
		},
		{
			name: "empty content",
			commentData: domain.CreateCommentRequest{
				Content: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comment, err := services.Comment.CreateComment(userID, post.ID, tt.commentData)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if comment == nil {
				t.Fatal("Expected comment but got nil")
			}

			if comment.Content != tt.commentData.Content {
				t.Errorf("Expected content %s, got %s", tt.commentData.Content, comment.Content)
			}

			if comment.UserID != userID {
				t.Errorf("Expected user ID %d, got %d", userID, comment.UserID)
			}

			if comment.PostID != post.ID {
				t.Errorf("Expected post ID %d, got %d", post.ID, comment.PostID)
			}

			retrievedComments, err := services.Comment.GetCommentsByPostID(userID, post.ID)
			if err != nil {
				t.Fatalf("Failed to get comments by post ID: %v", err)
			}

			if len(retrievedComments) == 0 {
				t.Fatal("Expected at least one comment from post retrieval")
			}
		})
	}
}

func TestCommentService_GetCommentsByPostID(t *testing.T) {
	services := SetupTestServices(t)

	userID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "test@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "testuser",
		Gender:      "male",
		IsPublic:    true,
	})

	post := CreateTestPost(t, services, userID, domain.CreatePostRequest{
		Title:    "Test Post",
		Content:  "This is a test post content",
		Category: "general",
	})

	comments := []string{
		"First comment",
		"Second comment",
		"Third comment",
	}

	for _, content := range comments {
		CreateTestComment(t, services, userID, post.ID, domain.CreateCommentRequest{
			Content: content,
		})
	}

	t.Run("get comments by post ID", func(t *testing.T) {
		retrievedComments, err := services.Comment.GetCommentsByPostID(userID, post.ID)
		if err != nil {
			t.Fatalf("Failed to get comments by post ID: %v", err)
		}

		if len(retrievedComments) != 3 {
			t.Errorf("Expected 3 comments, got %d", len(retrievedComments))
		}

		expectedContents := []string{"First comment", "Second comment", "Third comment"}
		for i, comment := range retrievedComments {
			if comment.Content != expectedContents[i] {
				t.Errorf("Expected comment content '%s', got '%s'", expectedContents[i], comment.Content)
			}
			if comment.PostID != post.ID {
				t.Errorf("Expected post ID %d, got %d", post.ID, comment.PostID)
			}
		}
	})

	t.Run("get comments for non-existent post", func(t *testing.T) {
		comments, err := services.Comment.GetCommentsByPostID(userID, 99999)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(comments) != 0 {
			t.Errorf("Expected 0 comments for non-existent post, got %d", len(comments))
		}
	})
}

func TestCommentService_GetCommentsByUserID(t *testing.T) {
	services := SetupTestServices(t)

	user1ID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "user1@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "user1",
		Gender:      "male",
		IsPublic:    true,
	})
	user2ID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "user2@example.com",
		Password:    "password123",
		FirstName:   "Jane",
		LastName:    "Smith",
		DateOfBirth: time.Now().AddDate(-30, 0, 0),
		Nickname:    "user2",
		Gender:      "female",
		IsPublic:    true,
	})

	post1 := CreateTestPost(t, services, user1ID, domain.CreatePostRequest{
		Title: "Post 1", Content: "Content 1 with enough characters", Category: "general",
	})

	post2 := CreateTestPost(t, services, user2ID, domain.CreatePostRequest{
		Title: "Post 2", Content: "Content 2 with enough characters", Category: "general",
	})

	CreateTestComment(t, services, user1ID, post1.ID, domain.CreateCommentRequest{
		Content: "User1 comment on post1",
	})

	CreateTestComment(t, services, user1ID, post2.ID, domain.CreateCommentRequest{
		Content: "User1 comment on post2",
	})

	CreateTestComment(t, services, user2ID, post1.ID, domain.CreateCommentRequest{
		Content: "User2 comment on post1",
	})

	t.Run("get comments by user ID", func(t *testing.T) {
		comments, err := services.Comment.GetCommentsByUserID(user1ID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to get comments by user ID: %v", err)
		}

		if len(comments) != 2 {
			t.Errorf("Expected 2 comments for user1, got %d", len(comments))
		}

		for _, comment := range comments {
			if comment.UserID != user1ID {
				t.Errorf("Expected user ID %d, got %d", user1ID, comment.UserID)
			}
		}
	})

	t.Run("get comments by user ID with pagination", func(t *testing.T) {
		comments, err := services.Comment.GetCommentsByUserID(user1ID, 1, 0)
		if err != nil {
			t.Fatalf("Failed to get comments by user ID with pagination: %v", err)
		}

		if len(comments) != 1 {
			t.Errorf("Expected 1 comment with limit 1, got %d", len(comments))
		}
	})
}

func TestCommentService_UpdateAndDeleteComment(t *testing.T) {
	services := SetupTestServices(t)

	userID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "author@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "author",
		Gender:      "male",
		IsPublic:    true,
	})

	otherID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "other@example.com",
		Password:    "password123",
		FirstName:   "Jane",
		LastName:    "Doe",
		DateOfBirth: time.Now().AddDate(-30, 0, 0),
		Nickname:    "other",
		Gender:      "female",
		IsPublic:    true,
	})

	post := CreateTestPost(t, services, userID, domain.CreatePostRequest{
		Title:    "Post",
		Content:  "This is a test post content for comment update and delete",
		Category: "general",
	})

	comment := CreateTestComment(t, services, userID, post.ID, domain.CreateCommentRequest{
		Content: "original comment text here",
	})

	updated, err := services.Comment.UpdateComment(userID, post.ID, comment.ID, domain.UpdateCommentRequest{
		Content: "replaced comment text for the test case validation",
	})
	if err != nil {
		t.Fatalf("UpdateComment: %v", err)
	}
	if updated.Content != "replaced comment text for the test case validation" {
		t.Errorf("content: %s", updated.Content)
	}

	_, err = services.Comment.UpdateComment(otherID, post.ID, comment.ID, domain.UpdateCommentRequest{
		Content: "hijack attempt with enough characters in this field",
	})
	if err == nil {
		t.Error("expected error updating another user's comment")
	}

	_, err = services.Comment.UpdateComment(userID, 99999, comment.ID, domain.UpdateCommentRequest{
		Content: "wrong post id but enough characters in the body here",
	})
	if err == nil {
		t.Error("expected error when post id does not match comment's post")
	}

	if err := services.Comment.DeleteComment(userID, post.ID, comment.ID); err != nil {
		t.Fatalf("DeleteComment: %v", err)
	}
	if err := services.Comment.DeleteComment(userID, post.ID, comment.ID); err == nil {
		t.Error("expected error deleting the same comment again")
	}
}
