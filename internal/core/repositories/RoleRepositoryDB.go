package repositories

import (
	"database/sql"

	"github.com/go-ms-project-store/internal/core/domain"
	"github.com/go-ms-project-store/internal/pkg/db"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type RoleRepositoryDB struct {
	client   *sqlx.DB
	verifier *db.FieldVerifier
}

func (rdb RoleRepositoryDB) FindByName(name string) (*domain.Role, *errs.AppError) {

	// Prepare query
	query := `SELECT
		id,
		name,
		created_at,
		updated_at
	FROM roles
	WHERE name = ?
    `

	var role domain.Role

	err := rdb.client.Get(&role, query, name)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Role not found")
		} else {
			logger.Error("Error while querying role table " + err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}
	}

	return &role, nil
}

func NewRoleRepositoryDB(dbClient *sqlx.DB) RoleRepositoryDB {
	return RoleRepositoryDB{
		client: dbClient,
		verifier: &db.FieldVerifier{
			DB:        dbClient,
			TableName: "roles",
		},
	}
}
