package repository

import (
	"social-network/internal/domain"
	"testing"
)

func TestConversationRepository_CreateConversation(t *testing.T) {
	repos := SetupTestDB(t)
	conversationRepo := repos.Conversation

	conversationID, err := conversationRepo.CreateConversation(&domain.Conversation{
		Name: "Test Conversation",
		Type: "private",
	})
	if err != nil {
		t.Fatalf("unexpected error creating conversation: %v", err)
	}
	if conversationID == 0 {
		t.Error("expected non-zero conversation ID")
	}

	conversation, err := conversationRepo.GetConversationByID(int(conversationID))
	if err != nil {
		t.Fatalf("unexpected error retrieving conversation: %v", err)
	}
	if conversation == nil {
		t.Fatal("expected conversation but got nil")
	}
	if conversation.Name != "Test Conversation" {
		t.Errorf("expected name 'Test Conversation', got '%s'", conversation.Name)
	}
	if conversation.Type != "private" {
		t.Errorf("expected type 'private', got '%s'", conversation.Type)
	}
}
