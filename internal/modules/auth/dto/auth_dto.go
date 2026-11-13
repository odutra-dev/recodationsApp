package dto

import usersDto "github.com/dutra/recodationsApp/internal/modules/users/dto"

// LoginInput defines credentials required to authenticate.
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse returns the generated JWT token and user profile.
type LoginResponse struct {
	AccessToken string                `json:"accessToken"`
	TokenType   string                `json:"tokenType"`
	ExpiresIn   int64                 `json:"expiresIn"`
	User        usersDto.UserResponse `json:"user"`
}
