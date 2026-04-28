package service

import (
	"social-network/internal/domain"
	"testing"
	"time"
)

func TestAuthService_Register(t *testing.T) {
	services := SetupTestServices(t)

	tests := []struct {
		name        string
		userData    domain.RegisterRequest
		expectError bool
	}{
		{
			name: "valid registration",
			userData: domain.RegisterRequest{
				Email:       "test@example.com",
				Password:    "password123",
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-25, 0, 0),
				Nickname:    "testuser",
				Gender:      "male",
			},
			expectError: false,
		},
		{
			name: "duplicate nickname",
			userData: domain.RegisterRequest{
				Email:       "different@example.com",
				Password:    "password123",
				FirstName:   "Jane",
				LastName:    "Smith",
				DateOfBirth: time.Now().AddDate(-30, 0, 0),
				Nickname:    "testuser",
				Gender:      "female",
				AvatarPath:  "https://example.com/avatar.jpg",
				AboutMe:     "I am a test user.",
			},
			expectError: true,
		},
		{
			name: "invalid email",
			userData: domain.RegisterRequest{
				Email:       "invalid-email",
				Password:    "password123",
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-25, 0, 0),
				Nickname:    "user2",
				Gender:      "male",
			},
			expectError: true,
		},
		{
			name: "weak password",
			userData: domain.RegisterRequest{
				Email:       "user3@example.com",
				Password:    "123",
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-25, 0, 0),
				Nickname:    "user3",
				Gender:      "male",
			},
			expectError: true,
		},
		{
			name: "underage user",
			userData: domain.RegisterRequest{
				Email:       "user4@example.com",
				Password:    "password123",
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-12, 0, 0),
				Nickname:    "user4",
				Gender:      "male",
			},
			expectError: true,
		},
		{
			name: "invalid AboutMe",
			userData: domain.RegisterRequest{
				Email:       "user5@example.com",
				Password:    "password123",
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-25, 0, 0),
				Nickname:    "user5",
				Gender:      "male",
				AboutMe:     "This is a very long about me section that exceeds the maximum allowed length of 500 characters. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, sessionID, err := services.Auth.Register(tt.userData)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if user == nil {
				t.Fatal("Expected user but got nil")
			}

			if user.Nickname != tt.userData.Nickname {
				t.Errorf("Expected nickname %s, got %s", tt.userData.Nickname, user.Nickname)
			}

			if user.Email != tt.userData.Email {
				t.Errorf("Expected email %s, got %s", tt.userData.Email, user.Email)
			}

			if user.FirstName != tt.userData.FirstName {
				t.Errorf("Expected first name %s, got %s", tt.userData.FirstName, user.FirstName)
			}

			if user.LastName != tt.userData.LastName {
				t.Errorf("Expected last name %s, got %s", tt.userData.LastName, user.LastName)
			}

			if len(user.AboutMe) > 500 {
				t.Errorf("Expected about me to be at most 500 characters, got %d", len(user.AboutMe))
			}

			if sessionID == "" {
				t.Error("Expected non-empty session ID")
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
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
		loginData   domain.LoginRequest
		expectError bool
	}{
		{
			name: "valid login with nickname",
			loginData: domain.LoginRequest{
				Identifier: "testuser",
				Password:   "password123",
			},
			expectError: false,
		},
		{
			name: "valid login with email",
			loginData: domain.LoginRequest{
				Identifier: "test@example.com",
				Password:   "password123",
			},
			expectError: false,
		},
		{
			name: "invalid password",
			loginData: domain.LoginRequest{
				Identifier: "testuser",
				Password:   "wrongpassword",
			},
			expectError: true,
		},
		{
			name: "non-existent user",
			loginData: domain.LoginRequest{
				Identifier: "nonexistent",
				Password:   "password123",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, sessionID, err := services.Auth.Login(tt.loginData)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if user == nil {
				t.Fatal("Expected user but got nil")
			}

			if user.ID != userID {
				t.Errorf("Expected user ID %d, got %d", userID, user.ID)
			}

			if sessionID == "" {
				t.Error("Expected non-empty session ID")
			}
		})
	}
}

func TestAuthService_ValidateSession(t *testing.T) {
	services := SetupTestServices(t)

	CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "test@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "testuser",
		Gender:      "male",
		IsPublic:    true,
	})

	_, sessionID, err := services.Auth.Login(domain.LoginRequest{
		Identifier: "testuser",
		Password:   "password123",
	})
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	t.Run("valid session", func(t *testing.T) {
		user, err := services.Auth.ValidateSession(sessionID)
		if err != nil {
			t.Fatalf("Failed to validate session: %v", err)
		}

		if user == nil {
			t.Fatal("Expected user but got nil")
		}

		if user.Nickname != "testuser" {
			t.Errorf("Expected nickname 'testuser', got '%s'", user.Nickname)
		}
	})

	t.Run("invalid session", func(t *testing.T) {
		_, err := services.Auth.ValidateSession("invalid-session-id")
		if err == nil {
			t.Error("Expected error for invalid session")
		}
	})
}

func TestAuthService_Logout(t *testing.T) {
	services := SetupTestServices(t)

	CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "test@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "testuser",
		Gender:      "male",
		IsPublic:    true,
	})

	_, sessionID, err := services.Auth.Login(domain.LoginRequest{
		Identifier: "testuser",
		Password:   "password123",
	})
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	t.Run("successful logout", func(t *testing.T) {
		err := services.Auth.Logout(sessionID)
		if err != nil {
			t.Fatalf("Failed to logout: %v", err)
		}

		_, err = services.Auth.ValidateSession(sessionID)
		if err == nil {
			t.Error("Expected error after logout")
		}
	})

	t.Run("logout with invalid session", func(t *testing.T) {
		err := services.Auth.Logout("invalid-session-id")
		if err != nil {
			t.Errorf("Logout with invalid session should not error, got: %v", err)
		}
	})
}

func TestAuthService_UpdateUser(t *testing.T) {
	services := SetupTestServices(t)

	userID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "before@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "beforeuser",
		Gender:      "male",
		IsPublic:    true,
	})

	updated, err := services.Auth.UpdateUser(userID, domain.UpdateUserRequest{
		Email:       "after@example.com",
		FirstName:   "Jane",
		LastName:    "Smith",
		DateOfBirth: time.Now().AddDate(-30, 0, 0),
		Nickname:    "afteruser",
		Gender:      "female",
		AvatarPath:  "/avatar.png",
		AboutMe:     "Updated bio",
		IsPublic:    false,
	})
	if err != nil {
		t.Fatalf("Expected update to succeed, got error: %v", err)
	}

	if updated.Email != "after@example.com" {
		t.Errorf("Expected updated email, got %s", updated.Email)
	}
	if updated.Nickname != "afteruser" {
		t.Errorf("Expected updated nickname, got %s", updated.Nickname)
	}

	_, err = services.Auth.UpdateUser(userID, domain.UpdateUserRequest{
		Email:       "invalid-email",
		FirstName:   "Jane",
		LastName:    "Smith",
		DateOfBirth: time.Now().AddDate(-30, 0, 0),
		Nickname:    "afteruser",
		Gender:      "female",
		IsPublic:    true,
	})
	if err == nil {
		t.Error("Expected validation error for invalid email")
	}
}

func TestAuthService_DeleteUser(t *testing.T) {
	services := SetupTestServices(t)

	userID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "deleteuser@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "deleteuser",
		Gender:      "male",
		IsPublic:    true,
	})

	if err := services.Auth.DeleteUser(userID); err != nil {
		t.Fatalf("Expected delete to succeed, got error: %v", err)
	}

	if err := services.Auth.DeleteUser(userID); err == nil {
		t.Error("Expected error when deleting user second time")
	}
}
