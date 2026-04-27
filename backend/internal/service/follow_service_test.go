package service

import (
	"social-network/internal/domain"
	"testing"
	"time"
)

func TestFollowService_FollowUser_Success_PublicUser(t *testing.T) {
	services := SetupTestServices(t)

	user1ID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "follower@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Follower",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "follower_user",
		Gender:      "male",
		IsPublic:    true,
	})

	user2ID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "public@example.com",
		Password:    "password123",
		FirstName:   "Jane",
		LastName:    "Public",
		DateOfBirth: time.Now().AddDate(-30, 0, 0),
		Nickname:    "public_user",
		Gender:      "female",
		IsPublic:    true,
	})

	status, err := services.Follow.FollowUser(domain.FollowRequest{
		FollowerID: user1ID,
		FolloweeID: user2ID,
	})

	if err != nil {
		t.Fatalf("Expected no error when following a public user, got: %v", err)
	}

	if status != "accepted" {
		t.Errorf("Expected status 'accepted' for public user, got: %s", status)
	}
}

func TestFollowService_FollowUser_Success_PrivateUser(t *testing.T) {
	services := SetupTestServices(t)

	user1ID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "follower2@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Follower",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "follower_user2",
		Gender:      "male",
		IsPublic:    true,
	})

	user2ID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "private@example.com",
		Password:    "password123",
		FirstName:   "Jane",
		LastName:    "Private",
		DateOfBirth: time.Now().AddDate(-30, 0, 0),
		Nickname:    "private_user",
		Gender:      "female",
		IsPublic:    false,
	})

	status, err := services.Follow.FollowUser(domain.FollowRequest{
		FollowerID: user1ID,
		FolloweeID: user2ID,
	})

	if err != nil {
		t.Fatalf("Expected no error when following a private user, got: %v", err)
	}

	if status != "pending" {
		t.Errorf("Expected status 'pending' for private user, got: %s", status)
	}
}

func TestFollowService_FollowUser_ReopensDeclinedRequest(t *testing.T) {
	services := SetupTestServices(t)

	followService := services.Follow.(*FollowService)

	user1ID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "follower3@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Follower",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "follower_user3",
		Gender:      "male",
		IsPublic:    true,
	})

	user2ID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "private2@example.com",
		Password:    "password123",
		FirstName:   "Jane",
		LastName:    "Private",
		DateOfBirth: time.Now().AddDate(-30, 0, 0),
		Nickname:    "private_user2",
		Gender:      "female",
		IsPublic:    false,
	})

	_, err := services.Follow.FollowUser(domain.FollowRequest{
		FollowerID: user1ID,
		FolloweeID: user2ID,
	})
	if err != nil {
		t.Fatalf("Expected initial follow request to succeed, got: %v", err)
	}

	existingFollow, err := followService.followRepo.GetFollowByUsers(user1ID, user2ID)
	if err != nil {
		t.Fatalf("Failed to get follow relationship after initial request: %v", err)
	}
	if existingFollow == nil {
		t.Fatal("Expected follow relationship to exist after initial request")
	}

	err = followService.followRepo.UpdateFollowStatus(existingFollow.ID, FollowStatusRejected)
	if err != nil {
		t.Fatalf("Failed to decline follow relationship for setup: %v", err)
	}

	status, err := services.Follow.FollowUser(domain.FollowRequest{
		FollowerID: user1ID,
		FolloweeID: user2ID,
	})
	if err != nil {
		t.Fatalf("Expected re-follow after declined request to succeed, got: %v", err)
	}
	if status != FollowStatusPending {
		t.Fatalf("Expected re-follow status to be pending, got: %s", status)
	}

	updatedFollow, err := followService.followRepo.GetFollowByUsers(user1ID, user2ID)
	if err != nil {
		t.Fatalf("Failed to get follow relationship after re-follow: %v", err)
	}
	if updatedFollow == nil {
		t.Fatal("Expected follow relationship to exist after re-follow")
	}
	if updatedFollow.ID != existingFollow.ID {
		t.Fatalf("Expected same follow row to be reused, old=%d new=%d", existingFollow.ID, updatedFollow.ID)
	}
	if updatedFollow.Status != FollowStatusPending {
		t.Fatalf("Expected reused follow row status to be pending, got: %s", updatedFollow.Status)
	}
}

func TestFollowService_UnfollowUser_Success(t *testing.T) {
	services := SetupTestServices(t)

	user1ID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "follower4@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Follower",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "follower_user4",
		Gender:      "male",
		IsPublic:    true,
	})

	user2ID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "public2@example.com",
		Password:    "password123",
		FirstName:   "Jane",
		LastName:    "Public",
		DateOfBirth: time.Now().AddDate(-30, 0, 0),
		Nickname:    "public_user2",
		Gender:      "female",
		IsPublic:    true,
	})

	_, err := services.Follow.FollowUser(domain.FollowRequest{
		FollowerID: user1ID,
		FolloweeID: user2ID,
	})
	if err != nil {
		t.Fatalf("Expected initial follow request to succeed, got: %v", err)
	}

	err = services.Follow.UnfollowUser(domain.UnfollowRequest{
		FollowerID: user1ID,
		FolloweeID: user2ID,
	})
	if err != nil {
		t.Fatalf("Expected unfollow to succeed, got: %v", err)
	}

	err = services.Follow.UnfollowUser(domain.UnfollowRequest{
		FollowerID: user1ID,
		FolloweeID: user2ID,
	})
	if err == nil {
		t.Fatalf("Expected error when unfollowing a user that is not followed, got: %v", err)
	}
	if err.Error() != "there is no follow relationship between you and this user" {
		t.Fatalf("Expected error message 'there is no follow relationship between you and this user', got: %v", err.Error())
	}
}

func TestFollowService_RemoveFollower_Success(t *testing.T) {
	services := SetupTestServices(t)

	user1ID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "follower5@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Follower",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "follower_user5",
		Gender:      "male",
		IsPublic:    true,
	})

	user2ID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "public3@example.com",
		Password:    "password123",
		FirstName:   "Jane",
		LastName:    "Public",
		DateOfBirth: time.Now().AddDate(-30, 0, 0),
		Nickname:    "public_user3",
		Gender:      "female",
		IsPublic:    true,
	})

	_, err := services.Follow.FollowUser(domain.FollowRequest{
		FollowerID: user1ID,
		FolloweeID: user2ID,
	})
	if err != nil {
		t.Fatalf("Expected initial follow request to succeed, got: %v", err)
	}

	err = services.Follow.RemoveFollower(domain.RemoveFollowerRequest{
		FolloweeID: user2ID,
		FollowerID: user1ID,
	})
	if err != nil {
		t.Fatalf("Expected remove follower to succeed, got: %v", err)
	}

	err = services.Follow.RemoveFollower(domain.RemoveFollowerRequest{
		FolloweeID: user2ID,
		FollowerID: user1ID,
	})
	if err == nil {
		t.Fatalf("Expected error when removing a follower that is not followed, got: %v", err)
	}
	if err.Error() != "this user is not following you" {
		t.Fatalf("Expected error message 'this user is not following you', got: %v", err.Error())
	}
}
