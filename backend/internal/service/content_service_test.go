package service

import (
	"social-network/internal/domain"
	"testing"
	"time"
)

func TestContentService_CreatePost(t *testing.T) {
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

	tests := []struct {
		name        string
		postData    domain.CreatePostRequest
		expectError bool
	}{
		{
			name: "valid post",
			postData: domain.CreatePostRequest{
				Title:    "Test Post",
				Content:  "This is a test post content",
				Category: "general",
			},
			expectError: false,
		},
		{
			name: "empty title",
			postData: domain.CreatePostRequest{
				Title:    "",
				Content:  "This is a test post content",
				Category: "general",
			},
			expectError: true,
		},
		{
			name: "content too short",
			postData: domain.CreatePostRequest{
				Title:    "Test Post",
				Content:  "Short",
				Category: "general",
			},
			expectError: true,
		},
		{
			name: "title too short",
			postData: domain.CreatePostRequest{
				Title:    "Hi",
				Content:  "This is a valid content with enough characters",
				Category: "general",
			},
			expectError: true,
		},
		{
			name: "empty category",
			postData: domain.CreatePostRequest{
				Title:    "Test Post",
				Content:  "This is a valid content with enough characters",
				Category: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := services.Content.CreatePost(userID, tt.postData)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if post == nil {
				t.Fatal("Expected post but got nil")
			}

			if post.Title != tt.postData.Title {
				t.Errorf("Expected title %s, got %s", tt.postData.Title, post.Title)
			}

			if post.Content != tt.postData.Content {
				t.Errorf("Expected content %s, got %s", tt.postData.Content, post.Content)
			}

			if post.Category != tt.postData.Category {
				t.Errorf("Expected category %s, got %s", tt.postData.Category, post.Category)
			}

			if post.UserID != userID {
				t.Errorf("Expected user ID %d, got %d", userID, post.UserID)
			}
		})
	}
}

func TestContentService_UpdateAndDeletePost(t *testing.T) {
	services := SetupTestServices(t)

	userID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "owner@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "owner",
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
		Title:    "Original",
		Content:  "Original body with more than ten characters for the validator.",
		Category: "general",
	})

	updated, err := services.Content.UpdatePost(userID, post.ID, domain.UpdatePostRequest{
		Title:    "Updated title here",
		Content:  "Updated body with more than ten characters for the validator to pass checks.",
		Category: "tech",
	})
	if err != nil {
		t.Fatalf("UpdatePost: %v", err)
	}
	if updated.Title != "Updated title here" {
		t.Errorf("title: %s", updated.Title)
	}

	_, err = services.Content.UpdatePost(otherID, post.ID, domain.UpdatePostRequest{
		Title:    "Hacked title for sure long enough",
		Content:  "Hacked body with more than ten characters for the validator to pass all checks now.",
		Category: "general",
	})
	if err == nil {
		t.Error("expected error updating another user's post")
	}

	if err := services.Content.DeletePost(userID, post.ID); err != nil {
		t.Fatalf("DeletePost: %v", err)
	}

	if err := services.Content.DeletePost(userID, post.ID); err == nil {
		t.Error("expected error deleting same post again")
	}
}
