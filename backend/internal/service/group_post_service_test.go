package service

import (
	"social-network/internal/domain"
	"testing"
	"time"
)

func TestContentService_CreateGroupPost(t *testing.T) {
	services := SetupTestServices(t)

	creatorID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "creator@example.com",
		Password:    "password123",
		FirstName:   "Creator",
		LastName:    "User",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "creator",
		Gender:      "male",
		IsPublic:    true,
	})
	outsiderID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "outsider@example.com",
		Password:    "password123",
		FirstName:   "Outsider",
		LastName:    "User",
		DateOfBirth: time.Now().AddDate(-30, 0, 0),
		Nickname:    "outsider",
		Gender:      "female",
		IsPublic:    true,
	})

	group, err := services.Group.CreateGroup(&domain.Group{
		CreatorID:   creatorID,
		Title:       "Test Group",
		Description: "This is a test group.",
	})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	t.Run("member creates group post", func(t *testing.T) {
		post, err := services.Content.CreateGroupPost(creatorID, group.ID, domain.CreatePostRequest{
			Title:    "Group Post",
			Content:  "This is a group post with enough content",
			Category: "general",
		})
		if err != nil {
			t.Fatalf("Failed to create group post: %v", err)
		}
		if post.GroupID != group.ID {
			t.Errorf("Expected group ID %d, got %d", group.ID, post.GroupID)
		}
		if post.UserID != creatorID {
			t.Errorf("Expected user ID %d, got %d", creatorID, post.UserID)
		}
	})

	t.Run("non-member cannot create group post", func(t *testing.T) {
		_, err := services.Content.CreateGroupPost(outsiderID, group.ID, domain.CreatePostRequest{
			Title:    "Outsider Post",
			Content:  "This should not be created in the group",
			Category: "general",
		})
		if err == nil {
			t.Error("Expected error when non-member creates group post")
		}
	})

	t.Run("group posts excluded from global list", func(t *testing.T) {
		_, err := services.Content.CreatePost(creatorID, domain.CreatePostRequest{
			Title:    "Global Post",
			Content:  "This is a global post with enough content",
			Category: "general",
		})
		if err != nil {
			t.Fatalf("Failed to create global post: %v", err)
		}

		posts, err := services.Post.ListPosts("", 10, 0)
		if err != nil {
			t.Fatalf("Failed to list posts: %v", err)
		}
		for _, post := range posts {
			if post.GroupID != 0 {
				t.Errorf("Expected only non-group posts in global list, got post with group ID %d", post.GroupID)
			}
		}
	})
}

func TestPostService_ListPostsByGroupID(t *testing.T) {
	services := SetupTestServices(t)

	creatorID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "creator@example.com",
		Password:    "password123",
		FirstName:   "Creator",
		LastName:    "User",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "creator",
		Gender:      "male",
		IsPublic:    true,
	})
	outsiderID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "outsider@example.com",
		Password:    "password123",
		FirstName:   "Outsider",
		LastName:    "User",
		DateOfBirth: time.Now().AddDate(-30, 0, 0),
		Nickname:    "outsider",
		Gender:      "female",
		IsPublic:    true,
	})

	group, err := services.Group.CreateGroup(&domain.Group{
		CreatorID:   creatorID,
		Title:       "Test Group",
		Description: "This is a test group.",
	})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	for i, title := range []string{"Group Post 1", "Group Post 2", "Group Post 3"} {
		_, err := services.Content.CreateGroupPost(creatorID, group.ID, domain.CreatePostRequest{
			Title:    title,
			Content:  "Content of group post number with enough text",
			Category: "general",
		})
		if err != nil {
			t.Fatalf("Failed to create group post %d: %v", i, err)
		}
	}

	t.Run("member lists group posts", func(t *testing.T) {
		posts, err := services.Post.ListPostsByGroupID(creatorID, group.ID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to list group posts: %v", err)
		}
		if len(posts) != 3 {
			t.Errorf("Expected 3 group posts, got %d", len(posts))
		}
	})

	t.Run("non-member cannot list group posts", func(t *testing.T) {
		_, err := services.Post.ListPostsByGroupID(outsiderID, group.ID, 10, 0)
		if err == nil {
			t.Error("Expected error when non-member lists group posts")
		}
	})
}

func TestPostService_GetPostByID_GroupAccess(t *testing.T) {
	services := SetupTestServices(t)

	creatorID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "creator@example.com",
		Password:    "password123",
		FirstName:   "Creator",
		LastName:    "User",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "creator",
		Gender:      "male",
		IsPublic:    true,
	})
	outsiderID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "outsider@example.com",
		Password:    "password123",
		FirstName:   "Outsider",
		LastName:    "User",
		DateOfBirth: time.Now().AddDate(-30, 0, 0),
		Nickname:    "outsider",
		Gender:      "female",
		IsPublic:    true,
	})

	group, err := services.Group.CreateGroup(&domain.Group{
		CreatorID:   creatorID,
		Title:       "Test Group",
		Description: "This is a test group.",
	})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	groupPost, err := services.Content.CreateGroupPost(creatorID, group.ID, domain.CreatePostRequest{
		Title:    "Group Post",
		Content:  "This is a group post with enough content",
		Category: "general",
	})
	if err != nil {
		t.Fatalf("Failed to create group post: %v", err)
	}

	globalPost, err := services.Content.CreatePost(creatorID, domain.CreatePostRequest{
		Title:    "Global Post",
		Content:  "This is a global post with enough content",
		Category: "general",
	})
	if err != nil {
		t.Fatalf("Failed to create global post: %v", err)
	}

	t.Run("member can read group post", func(t *testing.T) {
		post, err := services.Post.GetPostByID(creatorID, groupPost.ID)
		if err != nil {
			t.Fatalf("Failed to get group post as member: %v", err)
		}
		if post.ID != groupPost.ID {
			t.Errorf("Expected post ID %d, got %d", groupPost.ID, post.ID)
		}
	})

	t.Run("non-member cannot read group post", func(t *testing.T) {
		_, err := services.Post.GetPostByID(outsiderID, groupPost.ID)
		if err == nil {
			t.Error("Expected error when non-member reads group post")
		}
	})

	t.Run("any user can read global post", func(t *testing.T) {
		post, err := services.Post.GetPostByID(outsiderID, globalPost.ID)
		if err != nil {
			t.Fatalf("Failed to get global post: %v", err)
		}
		if post.ID != globalPost.ID {
			t.Errorf("Expected post ID %d, got %d", globalPost.ID, post.ID)
		}
	})
}

func TestCommentService_GroupPostMembership(t *testing.T) {
	services := SetupTestServices(t)

	creatorID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "creator@example.com",
		Password:    "password123",
		FirstName:   "Creator",
		LastName:    "User",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "creator",
		Gender:      "male",
		IsPublic:    true,
	})
	outsiderID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "outsider@example.com",
		Password:    "password123",
		FirstName:   "Outsider",
		LastName:    "User",
		DateOfBirth: time.Now().AddDate(-30, 0, 0),
		Nickname:    "outsider",
		Gender:      "female",
		IsPublic:    true,
	})

	group, err := services.Group.CreateGroup(&domain.Group{
		CreatorID:   creatorID,
		Title:       "Test Group",
		Description: "This is a test group.",
	})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	groupPost, err := services.Content.CreateGroupPost(creatorID, group.ID, domain.CreatePostRequest{
		Title:    "Group Post",
		Content:  "This is a group post with enough content",
		Category: "general",
	})
	if err != nil {
		t.Fatalf("Failed to create group post: %v", err)
	}

	t.Run("member can comment on group post", func(t *testing.T) {
		comment, err := services.Comment.CreateComment(creatorID, groupPost.ID, domain.CreateCommentRequest{
			Content: "This is a member comment",
		})
		if err != nil {
			t.Fatalf("Failed to create comment as member: %v", err)
		}
		if comment.PostID != groupPost.ID {
			t.Errorf("Expected post ID %d, got %d", groupPost.ID, comment.PostID)
		}
	})

	t.Run("non-member cannot comment on group post", func(t *testing.T) {
		_, err := services.Comment.CreateComment(outsiderID, groupPost.ID, domain.CreateCommentRequest{
			Content: "This is an outsider comment",
		})
		if err == nil {
			t.Error("Expected error when non-member comments on group post")
		}
	})

	t.Run("member can read group post comments", func(t *testing.T) {
		comments, err := services.Comment.GetCommentsByPostID(creatorID, groupPost.ID)
		if err != nil {
			t.Fatalf("Failed to get comments as member: %v", err)
		}
		if len(comments) == 0 {
			t.Error("Expected at least one comment for group post")
		}
	})

	t.Run("non-member cannot read group post comments", func(t *testing.T) {
		_, err := services.Comment.GetCommentsByPostID(outsiderID, groupPost.ID)
		if err == nil {
			t.Error("Expected error when non-member reads group post comments")
		}
	})
}
