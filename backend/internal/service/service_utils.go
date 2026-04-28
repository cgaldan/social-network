package service

import "social-network/internal/domain"

func displayUserName(user *domain.User) string {
	if user.Nickname != "" {
		return user.Nickname
	}
	if user.FirstName != "" || user.LastName != "" {
		name := user.FirstName
		if user.LastName != "" {
			if name != "" {
				name += " "
			}
			name += user.LastName
		}
		return name
	}
	return ""
}
