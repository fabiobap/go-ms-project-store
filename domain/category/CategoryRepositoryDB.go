package domain

import (
	"github.com/go-ms-project-store/errs"
	"github.com/go-ms-project-store/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type CategoryRepositoryDB struct {
	client *sqlx.DB
}

func (rdb CategoryRepositoryDB) FindAll() ([]Category, *errs.AppError) {
	var query string
	var err error
	categories := make([]Category, 0)

	query = "SELECT id, name, slug, created_at, updated_at from categories"
	err = rdb.client.Select(&categories, query)

	if err != nil {
		logger.Error("Error while querying category table " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	return categories, nil
}

func NewCategoryRepositoryDB(dbClient *sqlx.DB) CategoryRepositoryDB {
	return CategoryRepositoryDB{client: dbClient}
}
