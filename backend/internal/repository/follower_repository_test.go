package repository

import (
	"testing"
	"time"
)

func TestFollowerRepository_CreateFollower(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	followerRepo := repos.Follower

	userID1, err := userRepo.CreateUser("user1@example.com", "hashedpass1", "User", "One", time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC), "user1", "male", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user 1: %v", err)
	}

	userID2, err := userRepo.CreateUser("user2@example.com", "hashedpass2", "User", "Two", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "user2", "female", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user 2: %v", err)
	}

	followerID, err := followerRepo.CreateFollower(int(userID1), int(userID2), "pending")
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
	followerRepo := repos.Follower

	userID1, _ := userRepo.CreateUser("user1@example.com", "hashedpass1", "User", "One", time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC), "user1", "male", "", "", false)
	userID2, _ := userRepo.CreateUser("user2@example.com", "hashedpass2", "User", "Two", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "user2", "female", "", "", false)

	followerID, _ := followerRepo.CreateFollower(int(userID1), int(userID2), "pending")

	follower, err := followerRepo.GetFollowerByID(int(followerID))
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

func TestFollowerRepository_GetFollowersByUserID(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	followerRepo := repos.Follower

	userID1, _ := userRepo.CreateUser("user1@example.com", "hashedpass1", "User", "One", time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC), "user1", "male", "", "", false)
	userID2, _ := userRepo.CreateUser("user2@example.com", "hashedpass2", "User", "Two", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "user2", "female", "", "", false)
	userID3, _ := userRepo.CreateUser("user3@example.com", "hashedpass3", "User", "Three", time.Date(1985, time.June, 1, 0, 0, 0, 0, time.UTC), "user3", "male", "", "", false)

	followerRepo.CreateFollower(int(userID1), int(userID3), "accepted")
	followerRepo.CreateFollower(int(userID2), int(userID3), "pending")

	followers, err := followerRepo.GetFollowersByUserID(int(userID3), 10, 0)
	if err != nil {
		t.Fatalf("Failed to get followers: %v", err)
	}

	if len(followers) != 2 {
		t.Errorf("Expected 2 followers, got %d", len(followers))
	}
}

func TestFollowerRepository_GetFollowingByUserID(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	followerRepo := repos.Follower

	userID1, _ := userRepo.CreateUser("user1@example.com", "hashedpass1", "User", "One", time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC), "user1", "male", "", "", false)
	userID2, _ := userRepo.CreateUser("user2@example.com", "hashedpass2", "User", "Two", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "user2", "female", "", "", false)
	userID3, _ := userRepo.CreateUser("user3@example.com", "hashedpass3", "User", "Three", time.Date(1985, time.June, 1, 0, 0, 0, 0, time.UTC), "user3", "male", "", "", false)

	followerRepo.CreateFollower(int(userID1), int(userID2), "accepted")
	followerRepo.CreateFollower(int(userID1), int(userID3), "pending")

	following, err := followerRepo.GetFollowingByUserID(int(userID1), 10, 0)
	if err != nil {
		t.Fatalf("Failed to get following: %v", err)
	}

	if len(following) != 2 {
		t.Errorf("Expected 2 following, got %d", len(following))
	}
}

func TestFollowerRepository_UpdateFollowerStatus(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	followerRepo := repos.Follower

	userID1, _ := userRepo.CreateUser("user1@example.com", "hashedpass1", "User", "One", time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC), "user1", "male", "", "", false)
	userID2, _ := userRepo.CreateUser("user2@example.com", "hashedpass2", "User", "Two", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "user2", "female", "", "", false)

	followerID, _ := followerRepo.CreateFollower(int(userID1), int(userID2), "pending")

	err := followerRepo.UpdateFollowerStatus(int(followerID), "accepted")
	if err != nil {
		t.Fatalf("Failed to update follower status: %v", err)
	}

	follower, _ := followerRepo.GetFollowerByID(int(followerID))
	if follower.Status != "accepted" {
		t.Errorf("Expected status 'accepted', got '%s'", follower.Status)
	}
}

func TestFollowerRepository_DeleteFollower(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	followerRepo := repos.Follower

	userID1, _ := userRepo.CreateUser("user1@example.com", "hashedpass1", "User", "One", time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC), "user1", "male", "", "", false)
	userID2, _ := userRepo.CreateUser("user2@example.com", "hashedpass2", "User", "Two", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "user2", "female", "", "", false)

	followerID, _ := followerRepo.CreateFollower(int(userID1), int(userID2), "pending")

	err := followerRepo.DeleteFollower(int(followerID))
	if err != nil {
		t.Fatalf("Failed to delete follower: %v", err)
	}

	_, err = followerRepo.GetFollowerByID(int(followerID))
	if err == nil {
		t.Error("Expected error when getting deleted follower")
	}
}

func TestFollowerRepository_FollowExists(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	followerRepo := repos.Follower

	userID1, _ := userRepo.CreateUser("user1@example.com", "hashedpass1", "User", "One", time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC), "user1", "male", "", "", false)
	userID2, _ := userRepo.CreateUser("user2@example.com", "hashedpass2", "User", "Two", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "user2", "female", "", "", false)

	followerRepo.CreateFollower(int(userID1), int(userID2), "pending")

	exists, err := followerRepo.FollowExists(int(userID1), int(userID2))
	if err != nil {
		t.Fatalf("Failed to check if follow exists: %v", err)
	}

	if !exists {
		t.Error("Expected follow relationship to exist")
	}

	exists, _ = followerRepo.FollowExists(int(userID2), int(userID1))
	if exists {
		t.Error("Expected follow relationship to not exist")
	}
}

func TestFollowerRepository_GetFollowStatus(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	followerRepo := repos.Follower

	userID1, _ := userRepo.CreateUser("user1@example.com", "hashedpass1", "User", "One", time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC), "user1", "male", "", "", false)
	userID2, _ := userRepo.CreateUser("user2@example.com", "hashedpass2", "User", "Two", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "user2", "female", "", "", false)

	followerRepo.CreateFollower(int(userID1), int(userID2), "accepted")

	status, err := followerRepo.GetFollowStatus(int(userID1), int(userID2))
	if err != nil {
		t.Fatalf("Failed to get follow status: %v", err)
	}

	if status != "accepted" {
		t.Errorf("Expected status 'accepted', got '%s'", status)
	}
}
