package service

import (
	"social-network/internal/domain"
	"testing"
	"time"
)

func TestConversationService_CreateDirectConversation(t *testing.T) {
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

	_, err := services.Follow.FollowUser(domain.FollowRequest{
		FollowerID: user1ID,
		FolloweeID: user2ID,
	})
	if err != nil {
		t.Fatalf("Failed to create follow relationship: %v", err)
	}

	services.Conversation.CreateDirectConversation(domain.DirectConversationRequest{
		SenderID:   user1ID,
		ReceiverID: user2ID,
	})

	tests := []struct {
		name        string
		userID1     int
		userID2     int
		expectError bool
	}{
		{
			name:        "valid conversation",
			userID1:     user1ID,
			userID2:     user2ID,
			expectError: false,
		},
		{
			name:        "same user conversation",
			userID1:     user1ID,
			userID2:     user1ID,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conversation, err := services.Conversation.CreateDirectConversation(domain.DirectConversationRequest{
				SenderID:   tt.userID1,
				ReceiverID: tt.userID2,
			})

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if conversation == nil {
					t.Error("Expected conversation but got nil")
				} else {
					if conversation.ID == 0 {
						t.Error("Expected conversation with valid ID")
					}
					if conversation.Type != "private" {
						t.Errorf("Expected private conversation, got type: %s", conversation.Type)
					}
				}
			}
		})
	}
}
