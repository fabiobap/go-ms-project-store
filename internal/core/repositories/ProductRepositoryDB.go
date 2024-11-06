package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-ms-project-store/internal/core/domain"
	"github.com/go-ms-project-store/internal/pkg/db"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/helpers"
	"github.com/go-ms-project-store/internal/pkg/logger"
	"github.com/go-ms-project-store/internal/pkg/pagination"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
)

type ProductRepositoryDB struct {
	client   *sqlx.DB
	verifier *db.FieldVerifier
}

func (rdb ProductRepositoryDB) Create(p domain.Product) (*domain.Product, *errs.AppError) {
	var finalSlug string
	var nameExists *domain.Product

	nameExists, _ = rdb.FindByName(p.Name)
	if nameExists != nil {
		logger.Error("Error while creating new product, name already exists")
		return nil, errs.NewValidationError("name", "The name has already been taken")
	}

	if p.Slug != "" {
		// If slug is provided in the request, use it
		finalSlug = slug.Make(p.Slug)
	} else {
		// Generate slug from name
		finalSlug = slug.Make(p.Name)
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
	insertQuery := `INSERT INTO products 
		(name, 
		slug, 
		category_id, 
		description, 
		amount, 
		uuid, 
		created_at, 
		updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	res, sqlxErr := rdb.client.Exec(insertQuery, p.Name, finalSlug, p.CategoryId, p.Description, p.Amount, p.UUID, p.CreatedAt, p.UpdatedAt)
	if sqlxErr != nil {
		logger.Error("Error while creating new product " + sqlxErr.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	id, sqlxErr := res.LastInsertId()
	if sqlxErr != nil {
		logger.Error("Error while getting last insert id for new product " + sqlxErr.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	p.Id = id
	p.Slug = finalSlug

	return &p, nil
}

func (rdb ProductRepositoryDB) Delete(id int) *errs.AppError {
	query := `DELETE FROM products WHERE id = ?`

	result, err := rdb.client.ExecContext(context.Background(), query, id)
	if err != nil {
		logger.Error("Error while deleting product: " + err.Error())
		return errs.NewUnexpectedError("unexpected database error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("Error getting rows affected: " + err.Error())
		return errs.NewUnexpectedError("unexpected database error")
	}

	if rowsAffected == 0 {
		return errs.NewNotFoundError("Product not found")
	}

	return nil
}

func (rdb ProductRepositoryDB) FindById(id int) (*domain.Product, *errs.AppError) {
	var err error

	// Prepare query
	query := `
    SELECT 
        p.id,
        p.uuid,
        p.name as name,
        p.slug,
        p.category_id, 
        p.description, 
        p.amount,
        p.image,
        p.created_at,
        p.updated_at,
        c.id AS category_id,
        c.name AS category_name,
        c.slug AS category_slug,
		c.created_at AS category_created_at,
		c.updated_at AS category_updated_at
    FROM products p
    LEFT JOIN categories c ON p.category_id = c.id
    WHERE p.id = ?`

	result := make(map[string]interface{})
	err = rdb.client.QueryRowx(query, id).MapScan(result)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("product not found")
		}
		logger.Error("Error while scanning product row " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	product := &domain.Product{}

	// Handle UUID
	uuidBytes, ok := result["uuid"].([]uint8)
	if !ok {
		logger.Error(fmt.Sprintf("Unexpected UUID type: %T", result["uuid"]))
		return nil, errs.NewUnexpectedError("unexpected UUID format")
	}
	if len(uuidBytes) == 16 {
		product.UUID, err = uuid.FromBytes(uuidBytes)
	} else if len(uuidBytes) == 36 {
		product.UUID, err = uuid.ParseBytes(uuidBytes)
	} else {
		err = fmt.Errorf("unexpected UUID byte length: %d", len(uuidBytes))
	}
	if err != nil {
		logger.Error("Error processing UUID: " + err.Error())
		return nil, errs.NewUnexpectedError("error processing UUID")
	}

	// Handle other fields with type assertions
	product.Id, _ = result["id"].(int64)
	product.Name = helpers.DBByteToString(result["name"])
	product.Slug = helpers.DBByteToString(result["slug"])
	product.CategoryId, _ = result["category_id"].(int64)
	product.Description = helpers.DBByteToString(result["description"])
	if amount, ok := result["amount"].(int64); ok {
		product.Amount = int32(amount)
	}
	product.Image = helpers.DBByteToString(result["image"])
	product.CreatedAt, _ = result["created_at"].(time.Time)
	product.UpdatedAt, _ = result["updated_at"].(time.Time)

	product.Category.Id, _ = result["category_id"].(int64)
	product.Category.Name = helpers.DBByteToString(result["category_name"])
	product.Category.Slug = helpers.DBByteToString(result["category_slug"])
	product.Category.CreatedAt, _ = result["category_created_at"].(time.Time)
	product.Category.UpdatedAt, _ = result["category_updated_at"].(time.Time)

	return product, nil
}

func (rdb ProductRepositoryDB) FindAll(filter pagination.DataDBFilter) (domain.Products, int64, *errs.AppError) {
	var total int64
	products := domain.Products{}

	countQuery := `SELECT COUNT(*) FROM products`

	err := rdb.client.Get(&total, countQuery)
	if err != nil {
		logger.Error("Error while counting product table " + err.Error())
		return nil, 0, errs.NewUnexpectedError("unexpected database error")
	}

	query := fmt.Sprintf(`
	SELECT 
        p.id,
        p.uuid,
        p.name,
        p.slug,
        p.category_id, 
        p.description, 
        p.amount,
        p.image,
        p.created_at,
        p.updated_at,
        c.id AS category_id,
        c.name AS category_name,
        c.slug AS category_slug,
		c.created_at AS category_created_at,
		c.updated_at AS category_updated_at
    FROM products p
    LEFT JOIN categories c ON p.category_id = c.id
    ORDER BY %s %s
    LIMIT ? OFFSET ?
    `,
		filter.OrderBy,
		filter.OrderDir)

	// Calculate offset
	offset := (filter.Page - 1) * filter.PerPage

	rows, err := rdb.client.Queryx(
		query,
		filter.PerPage,
		offset,
	)

	if err != nil {
		logger.Error("Error while querying product table " + err.Error())
		return nil, 0, errs.NewUnexpectedError("unexpected database error")
	}
	defer rows.Close()

	for rows.Next() {
		var result map[string]interface{}
		err = rows.MapScan(result)
		if err != nil {
			logger.Error("Error while scanning product row " + err.Error())
			return nil, 0, errs.NewUnexpectedError("unexpected database error")
		}

		product := domain.Product{
			Id:          result["id"].(int64),
			UUID:        result["uuid"].(uuid.UUID),
			Name:        result["name"].(string),
			Slug:        result["slug"].(string),
			CategoryId:  result["category_id"].(int64),
			Description: result["description"].(string),
			Amount:      result["amount"].(int32),
			Image:       result["image"].(string),
			CreatedAt:   result["created_at"].(time.Time),
			UpdatedAt:   result["updated_at"].(time.Time),
			Category: domain.Category{
				Id:        result["category_id"].(int64),
				Name:      result["category_name"].(string),
				Slug:      result["category_slug"].(string),
				CreatedAt: result["category_created_at"].(time.Time),
				UpdatedAt: result["category_updated_at"].(time.Time),
			},
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		logger.Error("Error after iterating over product rows " + err.Error())
		return nil, 0, errs.NewUnexpectedError("unexpected database error")
	}

	return products, total, nil
}

func (rdb ProductRepositoryDB) FindByName(name string) (*domain.Product, *errs.AppError) {

	// Prepare query
	query := `SELECT
		id,
		uuid,
		name,
		slug,
		category_id, 
		description, 
		amount,
		image,
		created_at,
		updated_at
	FROM products
	WHERE name = ?
    `

	var product domain.Product

	err := rdb.client.Get(&product, query, name)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Product not found")
		} else {
			logger.Error("Error while querying product table " + err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}
	}

	return &product, nil
}

func (rdb ProductRepositoryDB) FindBySlug(slug string) (*domain.Product, *errs.AppError) {

	// Prepare query
	query := `SELECT
		id,
		uuid,
		name,
		slug,
		category_id, 
		description, 
		amount,
		image,
		created_at,
		updated_at
	FROM products
	WHERE slug = ?
    `

	var product domain.Product

	err := rdb.client.Get(&product, query, slug)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Product not found")
		} else {
			logger.Error("Error while querying product table " + err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}
	}

	return &product, nil
}

func (rdb ProductRepositoryDB) Update(p domain.Product) (*domain.Product, *errs.AppError) {
	var err error

	// First, check if the product exists
	existingProduct, errPkg := rdb.FindById(int(p.Id))
	if errPkg != nil {
		return nil, errs.NewNotFoundError("Product not found")
	}

	// Verify name uniqueness
	if err := rdb.verifier.VerifyUniqueField("name", p.Name, p.Id); err != nil {
		return nil, err
	}

	// Verify slug uniqueness
	if err := rdb.verifier.VerifyUniqueField("slug", p.Slug, p.Id); err != nil {
		return nil, err
	}

	updateQuery := `UPDATE products 
		SET 
		name = ?, 
		slug = ?, 
		category_id = ?, 
		description = ?, 
		amount = ? 
		WHERE id = ?`

	result, err := rdb.client.Exec(updateQuery, p.Name, p.Slug, p.CategoryId, p.Description, p.Amount, p.Id)
	if err != nil {
		logger.Error("Error while updating product: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("Error getting rows affected: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	if rowsAffected == 0 {
		return existingProduct, nil
	}

	// Fetch the updated product
	updatedProduct, errPkg := rdb.FindById(int(p.Id))
	if errPkg != nil {
		return nil, errs.NewUnexpectedError("Error fetching updated product")
	}

	return updatedProduct, nil
}

func NewProductRepositoryDB(dbClient *sqlx.DB) ProductRepositoryDB {
	return ProductRepositoryDB{
		client: dbClient,
		verifier: &db.FieldVerifier{
			DB:        dbClient,
			TableName: "products",
		},
	}
}
