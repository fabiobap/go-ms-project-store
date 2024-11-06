package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-ms-project-store/internal/core/domain"
	"github.com/go-ms-project-store/internal/pkg/db"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/logger"
	"github.com/go-ms-project-store/internal/pkg/pagination"
	_ "github.com/go-sql-driver/mysql"
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
	query := `
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
        c.id,
        c.name,
        c.slug,
        c.created_at,
        c.updated_at
    FROM products p
    LEFT JOIN categories c ON p.category_id = c.id
    WHERE p.id = ?`

	row := rdb.client.QueryRowx(query, id)
	return rdb.scanProduct(row)
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

	if filter.OrderBy == "id" {
		filter.OrderBy = "p.id"
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
        c.id,
        c.name,
        c.slug,
        c.created_at,
        c.updated_at
    FROM products p
    LEFT JOIN categories c ON p.category_id = c.id
    ORDER BY %s %s
    LIMIT ? OFFSET ?
    `,
		filter.OrderBy,
		filter.OrderDir)

	offset := (filter.Page - 1) * filter.PerPage

	rows, err := rdb.client.Queryx(query, filter.PerPage, offset)
	if err != nil {
		logger.Error("Error while querying product table " + err.Error())
		return nil, 0, errs.NewUnexpectedError("unexpected database error")
	}
	defer rows.Close()

	for rows.Next() {
		product, err := rdb.scanProducts(rows)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, *product)
	}

	if err = rows.Err(); err != nil {
		logger.Error("Error after iterating over product rows " + err.Error())
		return nil, 0, errs.NewUnexpectedError("unexpected database error")
	}

	return products, total, nil
}

func (rdb ProductRepositoryDB) FindByName(name string) (*domain.Product, *errs.AppError) {
	return rdb.findByField("name", name)
}

func (rdb ProductRepositoryDB) FindBySlug(slug string) (*domain.Product, *errs.AppError) {
	return rdb.findByField("slug", slug)
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

func (rdb ProductRepositoryDB) findByField(field, value string) (*domain.Product, *errs.AppError) {
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
    WHERE ` + field + ` = ?`

	var product domain.Product

	err := rdb.client.Get(&product, query, value)

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

func (rdb ProductRepositoryDB) processProduct(product *domain.Product, category *domain.Category, uuidBytes []byte) (*domain.Product, *errs.AppError) {
	processedUUID, err := db.ProcessUUID(uuidBytes)
	if err != nil {
		return nil, errs.NewUnexpectedError("error processing UUID")
	}

	product.UUID = processedUUID
	product.Category = *category
	return product, nil
}

func (rdb ProductRepositoryDB) scanProduct(row *sqlx.Row) (*domain.Product, *errs.AppError) {
	var product domain.Product
	var category domain.Category
	var uuidBytes []byte

	err := row.Scan(
		&product.Id,
		&uuidBytes,
		&product.Name,
		&product.Slug,
		&product.CategoryId,
		&product.Description,
		&product.Amount,
		&product.Image,
		&product.CreatedAt,
		&product.UpdatedAt,
		&category.Id,
		&category.Name,
		&category.Slug,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("product not found")
		}
		logger.Error("Error while scanning product row " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	return rdb.processProduct(&product, &category, uuidBytes)
}

func (rdb ProductRepositoryDB) scanProducts(rows *sqlx.Rows) (*domain.Product, *errs.AppError) {
	var product domain.Product
	var category domain.Category
	var uuidBytes []byte

	err := rows.Scan(
		&product.Id,
		&uuidBytes,
		&product.Name,
		&product.Slug,
		&product.CategoryId,
		&product.Description,
		&product.Amount,
		&product.Image,
		&product.CreatedAt,
		&product.UpdatedAt,
		&category.Id,
		&category.Name,
		&category.Slug,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		logger.Error("Error while scanning product row " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	return rdb.processProduct(&product, &category, uuidBytes)
}
