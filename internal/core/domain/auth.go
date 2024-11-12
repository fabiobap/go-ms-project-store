package domain

import (
	"time"

	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/go-ms-project-store/internal/pkg/errs"
)

type Token struct {
	ID          uint64    `db:"id"`
	UserID      uint64    `db:"tokenable_id"`
	Name        string    `db:"name"`
	Token       string    `db:"token"`
	Abilities   string    `db:"abilities"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	ExpiresAt   time.Time `db:"expires_at"`
	LastUsedAt  time.Time `db:"last_used_at"`
	HashedToken string
}

type AuthUser struct {
	Email    string
	Password string
}

type AuthRepository interface {
	CreateAccessToken(Token) (*Token, *errs.AppError)
	CreateRefreshToken(Token) (*Token, *errs.AppError)
	ValidateToken(string) (uint64, *errs.AppError)
	Login(AuthUser) (*User, *errs.AppError)
	Logout(uint64) *errs.AppError
	// Register(User) (*User, *errs.AppError)
}

func NewToken(dto dto.NewTokenDTO) Token {
	return Token{
		UserID:    dto.UserID,
		Name:      dto.Name,
		Abilities: "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: dto.ExpiresAt,
	}
}

func NewLogin(dto dto.NewLoginRequest) AuthUser {
	return AuthUser{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

// func (c Categories) ToDTO() []dto.CategoryResponse {
// 	dtos := make([]dto.CategoryResponse, len(c))
// 	for i, category := range c {
// 		dtos[i] = category.ToCategoryDTO()
// 	}
// 	return dtos
// }
