package repository

import (
	"testing"
	"time"
)

func TestFollowerRepository_CreateFollower(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	followRepo := repos.Follow

	userID1, err := userRepo.CreateUser("user1@example.com", "hashedpass1", "User", "One", time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC), "user1", "male", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user 1: %v", err)
	}

	userID2, err := userRepo.CreateUser("user2@example.com", "hashedpass2", "User", "Two", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "user2", "female", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user 2: %v", err)
	}

	followerID, err := followRepo.CreateFollow(int(userID1), int(userID2), "pending")
	if err != nil {
		t.Fatalf("Failed to create follower: %v", err)
	}

	if followerID == 0 {
		t.Error("Expected non-zero follower ID")
	}
}

func TestFollowerRepository_GetFollowerByID(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	followerRepo := repos.Follow

	userID1, _ := userRepo.CreateUser("user1@example.com", "hashedpass1", "User", "One", time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC), "user1", "male", "", "", false)
	userID2, _ := userRepo.CreateUser("user2@example.com", "hashedpass2", "User", "Two", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "user2", "female", "", "", false)

	followerID, _ := followerRepo.CreateFollow(int(userID1), int(userID2), "pending")

	follower, err := followerRepo.GetFollowByID(int(followerID))
	if err != nil {
		t.Fatalf("Failed to get follower: %v", err)
	}

	if follower.FollowerID != int(userID1) {
		t.Errorf("Expected follower_id %d, got %d", userID1, follower.FollowerID)
	}
	if follower.FollowingID != int(userID2) {
		t.Errorf("Expected following_id %d, got %d", userID2, follower.FollowingID)
	}
	if follower.Status != "pending" {
		t.Errorf("Expected status 'pending', got '%s'", follower.Status)
	}
}

func TestFollowerRepository_GetFollowRequestsByFollowingID(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	followRepo := repos.Follow

	userID1, _ := userRepo.CreateUser("user1@example.com", "hashedpass1", "User", "One", time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC), "user1", "male", "", "", false)
	userID2, _ := userRepo.CreateUser("user2@example.com", "hashedpass2", "User", "Two", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "user2", "female", "", "", false)
	userID3, _ := userRepo.CreateUser("user3@example.com", "hashedpass3", "User", "Three", time.Date(1985, time.June, 1, 0, 0, 0, 0, time.UTC), "user3", "male", "", "", false)

	followRepo.CreateFollow(int(userID1), int(userID3), "accepted")
	followRepo.CreateFollow(int(userID2), int(userID3), "pending")

	followRequests, err := followRepo.GetFollowRequestsByFollowingID(int(userID3), 10, 0)
	if err != nil {
		t.Fatalf("Failed to get follow requests: %v", err)
	}

	if len(followRequests) != 2 {
		t.Errorf("Expected 2 follow requests, got %d", len(followRequests))
	}
	if followRequests[0].FollowerID != int(userID1) {
		t.Errorf("Expected follower ID %d, got %d", userID1, followRequests[0].FollowerID)
	}
	if followRequests[0].FollowingID != int(userID3) {
		t.Errorf("Expected following ID %d, got %d", userID3, followRequests[0].FollowingID)
	}
	if followRequests[0].Status != "pending" {
		t.Errorf("Expected status 'pending', got '%s'", followRequests[0].Status)
	}
	if followRequests[1].FollowerID != int(userID2) {
		t.Errorf("Expected follower ID %d, got %d", userID2, followRequests[1].FollowerID)
	}
	if followRequests[1].FollowingID != int(userID3) {
		t.Errorf("Expected following ID %d, got %d", userID3, followRequests[1].FollowingID)
	}
	if followRequests[1].Status != "pending" {
		t.Errorf("Expected status 'pending', got '%s'", followRequests[1].Status)
	}
}

func TestFollowerRepository_GetFollowRequestsByFollowerID(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	followRepo := repos.Follow

	userID1, _ := userRepo.CreateUser("user1@example.com", "hashedpass1", "User", "One", time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC), "user1", "male", "", "", false)
	userID2, _ := userRepo.CreateUser("user2@example.com", "hashedpass2", "User", "Two", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "user2", "female", "", "", false)
	userID3, _ := userRepo.CreateUser("user3@example.com", "hashedpass3", "User", "Three", time.Date(1985, time.June, 1, 0, 0, 0, 0, time.UTC), "user3", "male", "", "", false)

	followRepo.CreateFollow(int(userID1), int(userID2), "accepted")
	followRepo.CreateFollow(int(userID1), int(userID3), "pending")

	followRequests, err := followRepo.GetFollowRequestsByFollowerID(int(userID1), 10, 0)
	if err != nil {
		t.Fatalf("Failed to get follow requests: %v", err)
	}

	if len(followRequests) != 2 {
		t.Errorf("Expected 2 follow requests, got %d", len(followRequests))
	}
}

func TestFollowerRepository_UpdateFollowerStatus(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	followRepo := repos.Follow

	userID1, _ := userRepo.CreateUser("user1@example.com", "hashedpass1", "User", "One", time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC), "user1", "male", "", "", false)
	userID2, _ := userRepo.CreateUser("user2@example.com", "hashedpass2", "User", "Two", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "user2", "female", "", "", false)

	followerID, _ := followRepo.CreateFollow(int(userID1), int(userID2), "pending")

	err := followRepo.UpdateFollowStatus(int(followerID), "accepted")
	if err != nil {
		t.Fatalf("Failed to update follow status: %v", err)
	}

	follower, _ := followRepo.GetFollowByID(int(followerID))
	if follower.Status != "accepted" {
		t.Errorf("Expected status 'accepted', got '%s'", follower.Status)
	}
}

func TestFollowerRepository_DeleteFollower(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	followRepo := repos.Follow

	userID1, _ := userRepo.CreateUser("user1@example.com", "hashedpass1", "User", "One", time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC), "user1", "male", "", "", false)
	userID2, _ := userRepo.CreateUser("user2@example.com", "hashedpass2", "User", "Two", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "user2", "female", "", "", false)

	followerID, _ := followRepo.CreateFollow(int(userID1), int(userID2), "pending")

	err := followRepo.DeleteFollow(int(followerID))
	if err != nil {
		t.Fatalf("Failed to delete follow: %v", err)
	}

	_, err = followRepo.GetFollowByID(int(followerID))
	if err == nil {
		t.Error("Expected error when getting deleted follower")
	}
}

func TestFollowerRepository_FollowExists(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	followRepo := repos.Follow

	userID1, _ := userRepo.CreateUser("user1@example.com", "hashedpass1", "User", "One", time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC), "user1", "male", "", "", false)
	userID2, _ := userRepo.CreateUser("user2@example.com", "hashedpass2", "User", "Two", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "user2", "female", "", "", false)

	followRepo.CreateFollow(int(userID1), int(userID2), "pending")

	exists, err := followRepo.FollowExists(int(userID1), int(userID2))
	if err != nil {
		t.Fatalf("Failed to check if follow exists: %v", err)
	}

	if !exists {
		t.Error("Expected follow relationship to exist")
	}

	exists, _ = followRepo.FollowExists(int(userID2), int(userID1))
	if exists {
		t.Error("Expected follow relationship to not exist")
	}
}

func TestFollowerRepository_GetFollowStatus(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	followRepo := repos.Follow

	userID1, _ := userRepo.CreateUser("user1@example.com", "hashedpass1", "User", "One", time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC), "user1", "male", "", "", false)
	userID2, _ := userRepo.CreateUser("user2@example.com", "hashedpass2", "User", "Two", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "user2", "female", "", "", false)

	followerID, _ := followRepo.CreateFollow(int(userID1), int(userID2), "accepted")

	status, err := followRepo.GetFollowStatusByFollowID(int(followerID))
	if err != nil {
		t.Fatalf("Failed to get follow status: %v", err)
	}

	if status != "accepted" {
		t.Errorf("Expected status 'accepted', got '%s'", status)
	}
}
