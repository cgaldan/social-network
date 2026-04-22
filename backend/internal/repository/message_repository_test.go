package repository

import (
	"social-network/internal/domain"
	"testing"
	"time"
)

func TestMessageRepository_SendMessage(t *testing.T) {
	repos := SetupTestDB(t)
	msgRepo := repos.Message
	userRepo := repos.User
	conversationRepo := repos.Conversation

	userID1, _ := userRepo.CreateUser(
		"sender@example.com",
		"hashedpass1",
		"User",
		"One",
		time.Now().AddDate(-25, 0, 0),
		"sender",
		"male",
		"",
		"",
		false,
	)

	userID2, _ := userRepo.CreateUser(
		"recipient@example.com",
		"hashedpass2",
		"User",
		"Two",
		time.Now().AddDate(-25, 0, 0),
		"recipient",
		"female",
		"",
		"",
		false,
	)

	conversation, err := conversationRepo.CreateDirectConversation(int(userID1), int(userID2))
	if err != nil {
		t.Fatalf("unexpected error creating conversation: %v", err)
	}
	if conversation == nil {
		t.Fatal("expected conversation but got nil")
	}

	id, err := msgRepo.CreateMessage(&domain.Message{
		ConversationID: int(conversation.ID),
		SenderID:       int(userID1),
		Content:        "Hello, World!",
	})
	if err != nil {
		t.Fatalf("unexpected error creating message: %v", err)
	}
	if id == 0 {
		t.Error("expected non-zero message ID")
	}

	message, err := msgRepo.GetMessageByID(int(id))
	if err != nil {
		t.Fatalf("unexpected error retrieving message: %v", err)
	}
	if message == nil {
		t.Fatal("expected message but got nil")
	}
	if message.Content != "Hello, World!" {
		t.Errorf("expected content 'Hello, World!', got '%s'", message.Content)
	}
	if message.SenderID != int(userID1) {
		t.Errorf("expected sender ID %d, got %d", userID1, message.SenderID)
	}
	if message.ConversationID != int(conversation.ID) {
		t.Errorf("expected conversation ID %d, got %d", conversation.ID, message.ConversationID)
	}
}
