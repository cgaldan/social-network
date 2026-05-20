package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"social-network/internal/domain"
	"social-network/internal/repository"
	"social-network/packages/logger"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    repository.UserRepositoryInterface
	sessionRepo repository.SessionRepositoryInterface
	followRepo  repository.FollowRepositoryInterface
	logger      *logger.Logger
}

func NewAuthService(userRepo repository.UserRepositoryInterface, sessionRepo repository.SessionRepositoryInterface, followRepo repository.FollowRepositoryInterface, logger *logger.Logger) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		followRepo:  followRepo,
		logger:      logger,
	}
}

func (s *AuthService) Register(registrationData domain.RegisterRequest) (*domain.User, string, error) {
	autoNickname := strings.TrimSpace(registrationData.Nickname) == ""
	if autoNickname {
		registrationData.Nickname = baseNickname(registrationData.FirstName, registrationData.LastName)
	}

	if err := s.validateRegistrationData(registrationData); err != nil {
		return nil, "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registrationData.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", "error", err)
		return nil, "", fmt.Errorf("failed to process password")
	}

	var userID int64
	nickname := registrationData.Nickname
	for attempts := 0; ; attempts++ {
		userID, err = s.userRepo.CreateUser(
			registrationData.Email,
			string(hashedPassword),
			registrationData.FirstName,
			registrationData.LastName,
			registrationData.DateOfBirth,
			nickname,
			registrationData.Gender,
			registrationData.AvatarPath,
			registrationData.AboutMe,
			registrationData.IsPublic,
		)
		if err == nil {
			break
		}
		if !strings.Contains(err.Error(), "UNIQUE constraint failed") {
			s.logger.Error("Failed to create user", "error", err)
			return nil, "", fmt.Errorf("failed to create user")
		}
		if !autoNickname || attempts >= 5 || strings.Contains(err.Error(), "users.email") {
			return nil, "", fmt.Errorf("nickname or email already in use")
		}
		nickname = baseNickname(registrationData.FirstName, registrationData.LastName) + randomNicknameSuffix()
	}
	registrationData.Nickname = nickname

	user, err := s.userRepo.GetUserByID(int(userID))
	if err != nil {
		s.logger.Error("Failed to retrieve created user", "error", err)
		return nil, "", fmt.Errorf("failed to retrieve user after creation")
	}

	sessionID, err := s.createSession(int(userID))
	if err != nil {
		s.logger.Error("Failed to create session for new user", "error", err)
		return nil, "", fmt.Errorf("failed to create session")
	}

	s.logger.Info("User registered successfully", "userID", userID, "nickname", registrationData.Nickname)
	return user, sessionID, nil
}

func (s *AuthService) Login(loginData domain.LoginRequest) (*domain.User, string, error) {
	if loginData.Identifier == "" || loginData.Password == "" {
		return nil, "", fmt.Errorf("identifier and password are required")
	}

	user, passwordHash, err := s.userRepo.GetUserByIdentifier(loginData.Identifier)
	if err != nil {
		s.logger.Debug("Login failed - user not found", "identifier", loginData.Identifier)
		return nil, "", fmt.Errorf("invalid identifier or password")
	}

	// REMEMBER TO REMOVE THIS BEFORE SUBMIT
	// if user.Nickname != "cgaldan" && user.Nickname != "cmarkos" && user.Nickname != "testuser" && user.Nickname != "jane" {
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(loginData.Password)); err != nil {
		s.logger.Debug("Login failed - password mismatch", "identifier", loginData.Identifier)
		return nil, "", fmt.Errorf("invalid identifier or password")
	}
	// }

	if err := s.userRepo.UpdateLastSeen(user.ID); err != nil {
		s.logger.Error("Failed to update last seen for user", "userID", user.ID, "error", err)
	}

	sessionID, err := s.createSession(user.ID)
	if err != nil {
		s.logger.Error("Failed to create session for logged in user", "error", err)
		return nil, "", fmt.Errorf("failed to create session")
	}

	s.logger.Info("User logged in successfully", "userID", user.ID, "nickname", user.Nickname)
	return user, sessionID, nil
}

func (s *AuthService) Logout(sessionID string) error {
	if err := s.sessionRepo.DeleteSession(sessionID); err != nil {
		s.logger.Error("Failed to delete session during logout", "sessionID", sessionID, "error", err)
		return fmt.Errorf("failed to logout")
	}

	s.logger.Info("User logged out successfully", "sessionID", sessionID)
	return nil
}

func (s *AuthService) ValidateSession(sessionID string) (*domain.User, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}

	session, err := s.sessionRepo.GetSessionBySessionID(sessionID)
	if err != nil {
		s.logger.Debug("Session validation failed - session not found", "sessionID", sessionID)
		return nil, fmt.Errorf("invalid session")
	}

	user, err := s.userRepo.GetUserByID(session.UserID)
	if err != nil {
		s.logger.Error("Failed to retrieve user for session validation", "userID", session.UserID, "error", err)
		return nil, fmt.Errorf("invalid session")
	}

	return user, nil
}

func (s *AuthService) UpdateUser(userID int, data domain.UpdateUserRequest) (*domain.User, error) {
	current, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	merged := *current
	if data.Email != nil {
		merged.Email = *data.Email
	}
	if data.FirstName != nil {
		merged.FirstName = *data.FirstName
	}
	if data.LastName != nil {
		merged.LastName = *data.LastName
	}
	if data.DateOfBirth != nil {
		merged.DateOfBirth = *data.DateOfBirth
	}
	if data.Nickname != nil {
		merged.Nickname = *data.Nickname
	}
	if data.Gender != nil {
		merged.Gender = *data.Gender
	}
	if data.AvatarPath != nil {
		merged.AvatarPath = *data.AvatarPath
	}
	if data.AboutMe != nil {
		merged.AboutMe = *data.AboutMe
	}
	if data.IsPublic != nil {
		merged.IsPublic = *data.IsPublic
	}

	if err := s.validateCommonUserData(
		merged.Email,
		merged.FirstName,
		merged.LastName,
		merged.DateOfBirth,
		merged.Nickname,
		merged.Gender,
		merged.AboutMe,
	); err != nil {
		return nil, err
	}

	if err := s.userRepo.UpdateUser(
		userID,
		merged.Email,
		merged.FirstName,
		merged.LastName,
		merged.DateOfBirth,
		merged.Nickname,
		merged.Gender,
		merged.AvatarPath,
		merged.AboutMe,
		merged.IsPublic,
	); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, fmt.Errorf("nickname or email already in use")
		}
		s.logger.Error("Failed to update user", "error", err, "userID", userID)
		return nil, fmt.Errorf("failed to update user")
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		s.logger.Error("Failed to retrieve updated user", "error", err, "userID", userID)
		return nil, fmt.Errorf("failed to retrieve updated user")
	}

	return user, nil
}

func (s *AuthService) DeleteUser(userID int) error {
	if err := s.userRepo.DeleteUser(userID); err != nil {
		s.logger.Error("Failed to delete user", "error", err, "userID", userID)
		return fmt.Errorf("failed to delete user")
	}
	return nil
}

func (s *AuthService) GetUserProfile(viewerID, targetID int) (*domain.UserProfile, error) {
	user, err := s.userRepo.GetUserByID(targetID)
	if err != nil {
		return nil, err
	}

	followStatus := "none"
	if viewerID == targetID {
		followStatus = "self"
	} else if s.followRepo != nil {
		follow, err := s.followRepo.GetFollowByUsers(viewerID, targetID)
		if err != nil {
			s.logger.Error("Failed to get follow status", "error", err, "viewerID", viewerID, "targetID", targetID)
		} else if follow != nil {
			switch follow.Status {
			case "accepted":
				followStatus = "following"
			case "pending":
				followStatus = "pending"
			}
		}
	}

	canView := user.IsPublic || followStatus == "self" || followStatus == "following"

	profile := &domain.UserProfile{
		ID:             user.ID,
		Nickname:       user.Nickname,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		AvatarPath:     user.AvatarPath,
		IsOnline:       user.IsOnline,
		IsPublic:       user.IsPublic,
		CreatedAt:      user.CreatedAt,
		FollowStatus:   followStatus,
		CanView:        canView,
		FollowersCount: user.FollowersCount,
		FollowingCount: user.FollowingCount,
	}

	if canView {
		profile.Email = user.Email
		profile.DateOfBirth = user.DateOfBirth
		profile.Gender = user.Gender
		profile.AboutMe = user.AboutMe
	}

	return profile, nil
}

func (s *AuthService) CanViewUser(viewerID, targetID int) (bool, error) {
	if viewerID == targetID {
		return true, nil
	}
	user, err := s.userRepo.GetUserByID(targetID)
	if err != nil {
		return false, err
	}
	if user.IsPublic {
		return true, nil
	}
	if s.followRepo == nil {
		return false, nil
	}
	follow, err := s.followRepo.GetFollowByUsers(viewerID, targetID)
	if err != nil {
		return false, err
	}
	return follow != nil && follow.Status == "accepted", nil
}

// Helper functions

func (s *AuthService) validateRegistrationData(data domain.RegisterRequest) error {
	if data.Password == "" || len(data.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	return s.validateCommonUserData(
		data.Email,
		data.FirstName,
		data.LastName,
		data.DateOfBirth,
		data.Nickname,
		data.Gender,
		data.AboutMe,
	)
}

func (s *AuthService) validateCommonUserData(email, firstName, lastName string, dateOfBirth time.Time, nickname, gender, aboutMe string) error {
	if email == "" || !strings.Contains(email, "@") {
		return fmt.Errorf("valid email is required")
	}
	if firstName == "" {
		return fmt.Errorf("first name is required")
	}
	if lastName == "" {
		return fmt.Errorf("last name is required")
	}
	if dateOfBirth.IsZero() {
		return fmt.Errorf("date of birth is required")
	}
	now := time.Now()
	age := now.Year() - dateOfBirth.Year()
	if now.YearDay() < dateOfBirth.YearDay() {
		age--
	}
	if age < 13 || age > 120 {
		return fmt.Errorf("age must be between 13 and 120")
	}
	if nickname == "" || len(nickname) < 3 {
		return fmt.Errorf("nickname must be at least 3 characters")
	}
	if gender == "" {
		return fmt.Errorf("gender is required")
	}
	if aboutMe != "" && len(aboutMe) > 500 {
		return fmt.Errorf("about me must be less than 500 characters")
	}
	return nil
}

func (s *AuthService) createSession(userID int) (string, error) {
	sessionID := generateSessionID()
	expiresAt := time.Now().Add(24 * time.Hour)

	if err := s.sessionRepo.CreateSession(sessionID, userID, expiresAt); err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	return sessionID, nil
}

func generateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func baseNickname(firstName, lastName string) string {
	combined := strings.TrimSpace(firstName) + strings.TrimSpace(lastName)
	var b strings.Builder
	for _, r := range strings.ToLower(combined) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	if b.Len() == 0 {
		return "user"
	}
	return b.String()
}

func randomNicknameSuffix() string {
	var buf [2]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano()%10000)
	}
	n := int(buf[0])<<8 | int(buf[1])
	return fmt.Sprintf("%04d", n%10000)
}
