package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-ms-project-store/internal/core/domain"
	"github.com/go-ms-project-store/internal/core/enums"
	"github.com/go-ms-project-store/internal/core/ports"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/helpers"
	"github.com/go-ms-project-store/internal/pkg/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepositoryDB struct {
	client   *sqlx.DB
	userRepo ports.UserRepository
	roleRepo ports.RoleRepository
}

func (rdb AuthRepositoryDB) createToken(tokenType string, au domain.Token) (*domain.Token, *errs.AppError) {
	genToken, err := helpers.GenerateToken()
	if err != nil {
		logger.Error("Error while generating token " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	hashedToken := helpers.HashToken(genToken)

	query := `INSERT INTO personal_access_tokens 
    (tokenable_id, tokenable_type, name, token, abilities, expires_at, created_at, updated_at) 
    VALUES 
    (?, ?, ?, ?, ?, ?, ?, ?)`

	result, sqlxErr := rdb.client.Exec(
		query,
		au.UserID,
		"App\\Models\\User",
		tokenType,
		hashedToken,
		au.Abilities,
		au.ExpiresAt,
		au.CreatedAt,
		au.UpdatedAt)
	if sqlxErr != nil {
		logger.Error("Error while creating new token " + sqlxErr.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	// Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("Error while getting last insert id " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	au.ID = uint64(id)
	au.HashedToken = hashedToken
	au.Token = genToken

	return &au, nil
}

func (rdb AuthRepositoryDB) CreateAccessToken(au domain.Token) (*domain.Token, *errs.AppError) {
	return rdb.createToken(string(enums.AccessToken), au)
}

func (rdb AuthRepositoryDB) CreateRefreshToken(au domain.Token) (*domain.Token, *errs.AppError) {
	return rdb.createToken(string(enums.RefreshToken), au)
}

func (rdb AuthRepositoryDB) Login(au domain.AuthUser) (*domain.User, *errs.AppError) {
	// Prepare query
	query := `SELECT id, email, password from users where email = ?`
	// Execute query
	var user domain.User

	err := rdb.client.Get(&user, query, au.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("No user found with email: " + au.Email)
			return nil, errs.NewUnauthorizedError("Invalid credentials")
		} else {
			logger.Error("Error while querying users table " + err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(au.Password))
	if err != nil {
		logger.Error("Error while generating token " + err.Error())
		return nil, errs.NewUnauthorizedError("Invalid credentials")
	}

	return &user, nil
}

func (rdb AuthRepositoryDB) Register(au domain.UserRegister) (*domain.User, *errs.AppError) {

	// Check if user already exists
	query := `SELECT id, email, password from users where email = ?`
	// Execute query
	var user domain.User

	err := rdb.client.Get(&user, query, au.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Error("Error while querying users table " + err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}
	} else {
		logger.Error("User already exists")
		return nil, errs.NewUnexpectedError("User already exists")
	}

	//find role id by name
	role, errRole := rdb.roleRepo.FindByName(string(enums.CustomerRole))
	au.RoleId = uint64(role.Id)
	if errRole != nil {
		logger.Error("Error while finding role " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	// Prepare query
	query = `INSERT INTO users (name, email, password, role_id, uuid) VALUES (?, ?, ?, ?, ?)`
	// Execute query
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(au.Password), 12)
	if err != nil {
		logger.Error("Error while hashing password " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	result, err := rdb.client.Exec(query, au.Name, au.Email, string(hashedPassword), au.RoleId, au.UUID)
	if err != nil {
		logger.Error("Error while creating new user " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	// Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("Error while getting last insert id " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	newUser := domain.User{
		Id:        int64(id),
		Name:      au.Name,
		Email:     au.Email,
		Password:  string(hashedPassword),
		UUID:      au.UUID,
		CreatedAt: au.CreatedAt,
	}

	return &newUser, nil
}

func (rdb AuthRepositoryDB) Logout(id uint64) *errs.AppError {
	// Prepare query
	query := `DELETE FROM personal_access_tokens WHERE tokenable_id = ?`

	// Execute query
	_, err := rdb.client.Exec(query, id)
	if err != nil {
		logger.Error("Error while deleting tokens " + err.Error())
		return errs.NewUnexpectedError("unexpected database error")
	}

	return nil
}

func (rdb AuthRepositoryDB) RoleRepo() ports.RoleRepository {
	return rdb.roleRepo
}

func (rdb AuthRepositoryDB) UserRepo() ports.UserRepository {
	return rdb.userRepo
}

func (rdb AuthRepositoryDB) ValidateToken(fullToken string) (uint64, *errs.AppError) {
	tokenID, tokenString, err := helpers.ParseToken(fullToken)
	if err != nil {
		logger.Error("Error while parsing token " + err.Error())
		return 0, errs.NewUnauthorizedError("Invalid Token")
	}

	hashedToken := helpers.HashToken(tokenString)

	var userID uint64
	var expiresAt sql.NullTime
	var lastUsedAt sql.NullTime

	query := `
	SELECT 
		tokenable_id, 
		expires_at, 
		last_used_at 
	FROM personal_access_tokens 
	WHERE id = ? AND token = ?`
	// Execute query
	err = rdb.client.QueryRow(query, tokenID, hashedToken).Scan(&userID, &expiresAt, &lastUsedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error(fmt.Sprintf("No token found with id: %d", tokenID))
			return 0, errs.NewUnauthorizedError("Invalid Token")
		} else {
			logger.Error("Error while querying personal_access_tokens table " + err.Error())
			return 0, errs.NewUnexpectedError("unexpected database error")
		}
	}

	if expiresAt.Valid && expiresAt.Time.Before(time.Now()) {
		return 0, errs.NewUnauthorizedError("Token Expired")
	}

	query = `UPDATE personal_access_tokens SET last_used_at = ? WHERE id = ?`
	_, err = rdb.client.Exec(query, time.Now(), tokenID)
	if err != nil {
		logger.Error("Error while updating last_used_at " + err.Error())
		return 0, errs.NewUnexpectedError("unexpected database error")
	}

	return userID, nil
}

func NewAuthRepositoryDB(dbClient *sqlx.DB) AuthRepositoryDB {
	return AuthRepositoryDB{
		client:   dbClient,
		userRepo: NewUserRepositoryDB(dbClient),
		roleRepo: NewRoleRepositoryDB(dbClient),
	}
}
