package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-ms-project-store/internal/core/domain"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/logger"
	"github.com/go-ms-project-store/internal/pkg/pagination"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
)

type CategoryRepositoryDB struct {
	client   *sqlx.DB
	verifier *FieldVerifier
}

type FieldVerifier struct {
	DB        *sqlx.DB
	TableName string
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

func (rdb CategoryRepositoryDB) Delete(id int) *errs.AppError {
	query := `DELETE FROM categories WHERE id = ?`

	result, err := rdb.client.ExecContext(context.Background(), query, id)
	if err != nil {
		logger.Error("Error while deleting category: " + err.Error())
		return errs.NewUnexpectedError("unexpected database error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("Error getting rows affected: " + err.Error())
		return errs.NewUnexpectedError("unexpected database error")
	}

	if rowsAffected == 0 {
		return errs.NewNotFoundError("Category not found")
	}

	return nil
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

func (rdb CategoryRepositoryDB) FindAll(filter pagination.DataDBFilter) (domain.Categories, int64, *errs.AppError) {
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

func (rdb CategoryRepositoryDB) Update(c domain.Category) (*domain.Category, *errs.AppError) {
	var err error

	// First, check if the category exists
	existingCategory, errPkg := rdb.FindById(int(c.Id))
	if errPkg != nil {
		return nil, errs.NewNotFoundError("Category not found")
	}

	// Verify name uniqueness
	if err := rdb.verifier.VerifyUniqueField("name", c.Name, c.Id); err != nil {
		return nil, err
	}

	// Verify slug uniqueness
	if err := rdb.verifier.VerifyUniqueField("slug", c.Slug, c.Id); err != nil {
		return nil, err
	}

	updateQuery := `UPDATE categories SET name = ?, slug = ? WHERE id = ?`
	result, err := rdb.client.Exec(updateQuery, c.Name, c.Slug, c.Id)
	if err != nil {
		logger.Error("Error while updating category: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("Error getting rows affected: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	if rowsAffected == 0 {
		return existingCategory, nil
	}

	// Fetch the updated category
	updatedCategory, errPkg := rdb.FindById(int(c.Id))
	if errPkg != nil {
		return nil, errs.NewUnexpectedError("Error fetching updated category")
	}

	return updatedCategory, nil
}

// VerifyUniqueField checks if a field value is unique in the table, excluding a specific ID
func (fv *FieldVerifier) VerifyUniqueField(fieldName, fieldValue string, excludeID int64) *errs.AppError {
	query := fmt.Sprintf("SELECT id FROM %s WHERE %s = ? AND id != ?", fv.TableName, fieldName)
	var existingID int64
	err := fv.DB.QueryRow(query, fieldValue, excludeID).Scan(&existingID)

	if err != nil && err != sql.ErrNoRows {
		logger.Error(fmt.Sprintf("Error checking for existing %s: %s", fieldName, err.Error()))
		return errs.NewUnexpectedError("Unexpected database error")
	}

	if err == nil {
		return errs.NewValidationError(fieldName, fmt.Sprintf("A record with this %s already exists", fieldName))
	}

	return nil
}

func NewCategoryRepositoryDB(dbClient *sqlx.DB) CategoryRepositoryDB {
	return CategoryRepositoryDB{
		client: dbClient,
		verifier: &FieldVerifier{
			DB:        dbClient,
			TableName: "categories",
		},
	}
}
