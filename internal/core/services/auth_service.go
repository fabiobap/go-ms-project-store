package services

import (
	"fmt"
	"net/http"

	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/go-ms-project-store/internal/core/domain"
	"github.com/go-ms-project-store/internal/core/enums"
	"github.com/go-ms-project-store/internal/core/ports"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/helpers"
)

type DefaultAuthService struct {
	repo ports.AuthRepository
}

func (s DefaultAuthService) Login(req dto.NewLoginRequest) (*dto.TokenResponse, *errs.AppError) {
	login := domain.NewLogin(req)

	user, err := s.repo.Login(login)
	if err != nil {
		if err.Code == http.StatusUnauthorized {
			return nil, err
		}
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	atAbility := []string{string(enums.AccessTokenAbility)}
	atDto := dto.NewTokenDTO{
		UserID:    uint64(user.Id),
		Name:      string(enums.AccessToken),
		Abilities: atAbility,
		ExpiresAt: helpers.GetAccessTokenExpiry(),
	}

	rtAbility := []string{string(enums.RefreshTokenAbility)}
	rtDto := dto.NewTokenDTO{
		UserID:    uint64(user.Id),
		Name:      string(enums.RefreshToken),
		Abilities: rtAbility,
		ExpiresAt: helpers.GetRefreshTokenExpiry(),
	}

	acToken := domain.NewToken(atDto)
	ac, err := s.repo.CreateAccessToken(acToken)
	if err != nil {
		if err.Code == http.StatusUnprocessableEntity {
			return nil, err
		}
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	rToken := domain.NewToken(rtDto)
	rt, err := s.repo.CreateRefreshToken(rToken)
	if err != nil {
		if err.Code == http.StatusUnprocessableEntity {
			return nil, err
		}
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	res := dto.TokenResponse{
		AccessToken:  fmt.Sprintf("%d|%s", ac.ID, ac.Token),
		RefreshToken: fmt.Sprintf("%d|%s", rt.ID, rt.Token),
		ExpiresIn:    60,
		TokenType:    "Bearer",
	}

	return &res, nil
}

func (s DefaultAuthService) Logout(user_id uint64) *errs.AppError {
	err := s.repo.Logout(user_id)
	if err != nil {
		return err
	}

	return nil
}

func (s DefaultAuthService) Me(user_id uint64) (*domain.User, *errs.AppError) {
	user, err := s.repo.UserRepo().FindById(user_id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s DefaultAuthService) RefreshToken(user_id uint64) (*dto.TokenResponse, *errs.AppError) {

	err := s.repo.RevokeAccessToken(user_id)
	if err != nil {
		return nil, err
	}

	atAbility := []string{string(enums.AccessTokenAbility)}
	atDto := dto.NewTokenDTO{
		UserID:    uint64(user_id),
		Name:      string(enums.AccessToken),
		ExpiresAt: helpers.GetAccessTokenExpiry(),
		Abilities: atAbility,
	}

	acToken := domain.NewToken(atDto)
	ac, err := s.repo.CreateAccessToken(acToken)
	if err != nil {
		if err.Code == http.StatusUnprocessableEntity {
			return nil, err
		}
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	res := dto.TokenResponse{
		AccessToken: fmt.Sprintf("%d|%s", ac.ID, ac.Token),
		ExpiresIn:   60,
		TokenType:   "Bearer",
	}

	return &res, nil
}

func (s DefaultAuthService) Register(dto dto.NewUserRegisterRequest) (*domain.User, *errs.AppError) {
	newUser := domain.NewUserRegister(dto)

	user, err := s.repo.Register(newUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewAuthService(repository ports.AuthRepository) DefaultAuthService {
	return DefaultAuthService{repo: repository}
}
