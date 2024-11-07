package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-ms-project-store/internal/core/domain"
	"github.com/go-ms-project-store/internal/core/enums"
	"github.com/go-ms-project-store/internal/pkg/db"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/logger"
	"github.com/go-ms-project-store/internal/pkg/pagination"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryDB struct {
	client   *sqlx.DB
	verifier *db.FieldVerifier
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func (rdb UserRepositoryDB) Delete(id string) *errs.AppError {
	query := `DELETE FROM users WHERE uuid = ?`

	result, err := rdb.client.ExecContext(context.Background(), query, id)
	if err != nil {
		logger.Error("Error while deleting user: " + err.Error())
		return errs.NewUnexpectedError("unexpected database error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("Error getting rows affected: " + err.Error())
		return errs.NewUnexpectedError("unexpected database error")
	}

	if rowsAffected == 0 {
		return errs.NewNotFoundError("User not found")
	}

	return nil
}

func (rdb UserRepositoryDB) FindById(id string) (*domain.User, *errs.AppError) {
	query := `
    SELECT 
        id,
        uuid,
        name,
        email,
        role_id, 
        email_verified_at, 
        created_at,
        updated_at
    FROM users
    WHERE uuid = ?`

	row := rdb.client.QueryRowx(query, id)
	return rdb.scanUser(row)
}

func (rdb UserRepositoryDB) FindAll(filter pagination.DataDBFilter, roleName string) (domain.Users, int64, *errs.AppError) {
	var total int64
	users := domain.Users{}

	// Base query for counting
	countQuery := "SELECT COUNT(*) FROM users u"

	// Base query for selecting users
	baseQuery := `
    SELECT 
        u.id,
        u.uuid,
        u.name,
        u.email,
        u.role_id, 
        u.email_verified_at, 
        u.created_at,
        u.updated_at
    FROM users u`

	// If roleName is provided, add it to the queries
	var args []interface{}
	if roleName != "" {
		countQuery += " JOIN roles ON u.role_id = roles.id WHERE roles.name = ?"
		baseQuery += " JOIN roles ON u.role_id = roles.id WHERE roles.name = ?"
		args = append(args, roleName)
	}

	// Execute count query
	err := rdb.client.Get(&total, countQuery, args...)
	if err != nil {
		logger.Error("Error while counting user table " + err.Error())
		return nil, 0, errs.NewUnexpectedError("unexpected database error")
	}

	if filter.OrderBy == "id" {
		filter.OrderBy = "u.id"
	}

	// Construct the final query
	query := fmt.Sprintf("%s ORDER BY %s %s LIMIT ? OFFSET ?",
		baseQuery,
		filter.OrderBy,
		filter.OrderDir)

	// Add pagination parameters to args
	offset := (filter.Page - 1) * filter.PerPage
	args = append(args, filter.PerPage, offset)

	// Execute the main query
	rows, err := rdb.client.Queryx(query, args...)
	if err != nil {
		logger.Error("Error while querying user table " + err.Error())
		return nil, 0, errs.NewUnexpectedError("unexpected database error")
	}
	defer rows.Close()

	for rows.Next() {
		user, err := rdb.scanUser(rows)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, *user)
	}

	if err = rows.Err(); err != nil {
		logger.Error("Error after iterating over user rows " + err.Error())
		return nil, 0, errs.NewUnexpectedError("unexpected database error")
	}

	return users, total, nil
}

func (rdb UserRepositoryDB) FindAllCustomers(filter pagination.DataDBFilter) (domain.Users, int64, *errs.AppError) {
	return rdb.FindAll(filter, string(enums.CustomerRole))
}

func (rdb UserRepositoryDB) FindAllAdmins(filter pagination.DataDBFilter) (domain.Users, int64, *errs.AppError) {
	return rdb.FindAll(filter, string(enums.AdminRole))
}

func NewUserRepositoryDB(dbClient *sqlx.DB) UserRepositoryDB {
	return UserRepositoryDB{
		client: dbClient,
		verifier: &db.FieldVerifier{
			DB:        dbClient,
			TableName: "users",
		},
	}
}

func (rdb UserRepositoryDB) processUser(user *domain.User, uuidBytes []byte) (*domain.User, *errs.AppError) {
	processedUUID, err := db.ProcessUUID(uuidBytes)
	if err != nil {
		return nil, errs.NewUnexpectedError("error processing UUID")
	}

	user.UUID = processedUUID
	return user, nil
}

func (rdb UserRepositoryDB) scanUser(s scanner) (*domain.User, *errs.AppError) {
	var user domain.User
	var uuidBytes []byte

	err := s.Scan(
		&user.Id,
		&uuidBytes,
		&user.Name,
		&user.Email,
		&user.RoleId,
		&user.EmailVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("user not found")
		}
		logger.Error("Error while scanning user row " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	return rdb.processUser(&user, uuidBytes)
}
