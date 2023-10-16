package types

import (
	"time"

	db "github.com/mas-wig/simple-api/db/sqlc"
)

// NOTE: Kita sudah punya service yang digenerate dari SQLC tapi kita butuh custom struct sign user
// Untuk mengfilter data yang di return dari postgresql mengunakan UserResponse struct dengan tag omitemty

type SignInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
}

func FilteredUserResponse(user db.User) UserResponse {
	return UserResponse{
		CreatedAt: user.UpdatedAt,
		UpdatedAt: user.CreatedAt,
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
	}
}
