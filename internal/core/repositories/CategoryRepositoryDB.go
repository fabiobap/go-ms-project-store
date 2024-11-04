package repositories

import (
	"database/sql"
	"fmt"

	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/go-ms-project-store/internal/core/domain"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gosimple/slug"
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
    `,
		filter.OrderBy,
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

func (rdb CategoryRepositoryDB) FindById(id int) (*domain.Category, *errs.AppError) {
	// Prepare query
	query := `SELECT
		id,
		name,
		slug,
		created_at,
		updated_at
	FROM categories
	WHERE id = ?
    `

	var category domain.Category

	err := rdb.client.Get(&category, query, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Category not found")
		} else {
			logger.Error("Error while querying category table " + err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}
	}

	return &category, nil
}

func (rdb CategoryRepositoryDB) Create(c domain.Category) (*domain.Category, *errs.AppError) {
	var finalSlug string
	var nameExists *domain.Category

	nameExists, _ = rdb.FindByName(c.Name)
	if nameExists != nil {
		logger.Error("Error while creating new category, name already exists")
		return nil, errs.NewValidationError("name", "The name has already been taken")
	}

	if c.Slug != "" {
		// If slug is provided in the request, use it
		finalSlug = slug.Make(c.Slug)
	} else {
		// Generate slug from name
		finalSlug = slug.Make(c.Name)
	}

	// Check if slug exists and increment if necessary
	baseSlug := finalSlug
	counter := 1
	for {
		// Try to find if the current slug exists
		existing, err := rdb.FindBySlug(finalSlug)
		if err != nil {
			// If error is because slug doesn't exist, we can use this slug
			break
		}
		if existing != nil {
			// Slug exists, try next increment
			finalSlug = fmt.Sprintf("%s-%d", baseSlug, counter)
			counter++
			continue
		}
		break
	}

	// Prepare query
	insertQuery := `INSERT INTO categories (name, slug, created_at, updated_at) VALUES (?, ?, ?, ?)`

	res, sqlxErr := rdb.client.Exec(insertQuery, c.Name, finalSlug, c.CreatedAt, c.UpdatedAt)
	if sqlxErr != nil {
		logger.Error("Error while creating new category " + sqlxErr.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	id, sqlxErr := res.LastInsertId()
	if sqlxErr != nil {
		logger.Error("Error while getting last insert id for new category " + sqlxErr.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	c.Id = id
	c.Slug = finalSlug

	return &c, nil
}

func (rdb CategoryRepositoryDB) Update(c domain.Category) (*domain.Category, *errs.AppError) {

	// Prepare query
	query := `UPDATE categories SET name = ?, slug = ?, updated_at = ? WHERE id = ?`

	_, err := rdb.client.Exec(query, c.Name, c.Slug, c.UpdatedAt, c.Id)
	if err != nil {
		logger.Error("Error while updating category " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	return &c, nil
}

func (rdb CategoryRepositoryDB) Delete(id int) *errs.AppError {

	// Prepare query
	query := `DELETE FROM categories WHERE id = ?`

	_, err := rdb.client.Exec(query, id)
	if err != nil {
		logger.Error("Error while deleting category " + err.Error())
		return errs.NewUnexpectedError("unexpected database error")
	}

	return nil
}

func (rdb CategoryRepositoryDB) FindBySlug(slug string) (*domain.Category, *errs.AppError) {

	// Prepare query
	query := `SELECT
		id,
		name,
		slug,
		created_at,
		updated_at
	FROM categories
	WHERE slug = ?
    `

	var category domain.Category

	err := rdb.client.Get(&category, query, slug)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Category not found")
		} else {
			logger.Error("Error while querying category table " + err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}
	}

	return &category, nil
}

func (rdb CategoryRepositoryDB) FindByName(name string) (*domain.Category, *errs.AppError) {

	// Prepare query
	query := `SELECT
		id,
		name,
		slug,
		created_at,
		updated_at
	FROM categories
	WHERE name = ?
    `

	var category domain.Category

	err := rdb.client.Get(&category, query, name)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Category not found")
		} else {
			logger.Error("Error while querying category table " + err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}
	}

	return &category, nil
}

func NewCategoryRepositoryDB(dbClient *sqlx.DB) CategoryRepositoryDB {
	return CategoryRepositoryDB{client: dbClient}
}
