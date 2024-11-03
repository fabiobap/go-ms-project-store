package repositories

import (
	"fmt"

	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/go-ms-project-store/internal/core/domain"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type CategoryRepositoryDB struct {
	client *sqlx.DB
}

func (rdb CategoryRepositoryDB) FindAll(filter dto.DataDBFilter) (domain.Categories, int64, *errs.AppError) {
	var total int64
	categories := domain.Categories{}

	countQuery := `SELECT COUNT(*) FROM categories`
	err := rdb.client.Get(&total, countQuery)
	if err != nil {
		logger.Error("Error while counting category table " + err.Error())
		return nil, 0, errs.NewUnexpectedError("unexpected database error")
	}

	query := fmt.Sprintf(`
	SELECT 
		id, 
		name, 
		slug, 
		created_at, 
		updated_at 
	FROM categories 
	ORDER BY %s %s
	LIMIT ? OFFSET ?
    `, filter.OrderBy,
		filter.OrderDir)
	// Calculate offset
	offset := (filter.Page - 1) * filter.PerPage

	err = rdb.client.Select(
		&categories,
		query,
		filter.PerPage,
		offset,
	)

	if err != nil {
		logger.Error("Error while querying category table " + err.Error())
		return nil, 0, errs.NewUnexpectedError("unexpected database error")
	}

	return categories, total, nil
}

func NewCategoryRepositoryDB(dbClient *sqlx.DB) CategoryRepositoryDB {
	return CategoryRepositoryDB{client: dbClient}
}
