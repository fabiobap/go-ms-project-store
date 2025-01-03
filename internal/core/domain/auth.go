package domain

import (
	"time"

	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/google/uuid"
)

type Token struct {
	ID          uint64    `db:"id"`
	UserID      uint64    `db:"tokenable_id"`
	Name        string    `db:"name"`
	Token       string    `db:"token"`
	Abilities   []string  `db:"abilities"`
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

type UserRegister struct {
	Email     string
	Name      string
	Password  string
	RoleId    uint64
	UUID      uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewToken(dto dto.NewTokenDTO) Token {
	return Token{
		UserID:    dto.UserID,
		Name:      dto.Name,
		Abilities: dto.Abilities,
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

func NewUserRegister(dto dto.NewUserRegisterRequest) UserRegister {
	return UserRegister{
		Email:     dto.Email,
		Password:  dto.Password,
		Name:      dto.Name,
		UUID:      uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
