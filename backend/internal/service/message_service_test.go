package service

import (
	"social-network/internal/domain"
	"testing"
	"time"
)

func TestMessageService_SendMessage(t *testing.T) {
	services := SetupTestServices(t)

	senderID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "test@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "testuser",
		Gender:      "male",
		IsPublic:    true,
	})

	receiverID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "test2@example.com",
		Password:    "password123",
		FirstName:   "Jane",
		LastName:    "Smith",
		DateOfBirth: time.Now().AddDate(-30, 0, 0),
		Nickname:    "testuser2",
		Gender:      "female",
		IsPublic:    true,
	})

	var err error
	_, err = services.Follow.FollowUser(domain.FollowRequest{
		FollowerID: senderID,
		FolloweeID: receiverID,
	})
	if err != nil {
		t.Fatalf("Failed to create follow relationship: %v", err)
	}

	_, err = services.Follow.FollowUser(domain.FollowRequest{
		FollowerID: receiverID,
		FolloweeID: senderID,
	})
	if err != nil {
		t.Fatalf("Failed to create follow relationship: %v", err)
	}

	conversation := CreateTestDirectChat(t, services, senderID, receiverID)
	conversationID := conversation.ID

	tests := []struct {
		name           string
		conversationID int
		senderID       int
		content        string
		expectError    bool
	}{
		{
			name:           "valid message",
			senderID:       senderID,
			conversationID: conversationID,
			content:        "Hello, this is a test message!",
			expectError:    false,
		},
		{
			name:           "empty content",
			senderID:       senderID,
			conversationID: conversationID,
			content:        "",
			expectError:    true,
		},
		{
			name:           "message too long",
			senderID:       senderID,
			conversationID: conversationID,
			content:        string(make([]byte, 1001)),
			expectError:    true,
		},
		{
			name:           "non-existent sender",
			senderID:       99999,
			conversationID: conversationID,
			content:        "Message from ghost",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := services.Message.SendMessage(tt.senderID, tt.conversationID, tt.content)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if message == nil {
				t.Fatal("Expected message but got nil")
			}

			if message.SenderID != tt.senderID {
				t.Errorf("Expected sender ID %d, got %d", tt.senderID, message.SenderID)
			}

			if message.ConversationID != tt.conversationID {
				t.Errorf("Expected conversation ID %d, got %d", tt.conversationID, message.ConversationID)
			}

			if message.Content != tt.content {
				t.Errorf("Expected content %s, got %s", tt.content, message.Content)
			}
		})
	}
}
