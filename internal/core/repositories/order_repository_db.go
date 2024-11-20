package repositories

import (
	"database/sql"
	"time"

	"github.com/go-ms-project-store/internal/core/domain"
	"github.com/go-ms-project-store/internal/pkg/db"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type OrderRepositoryDB struct {
	client   *sqlx.DB
	verifier *db.FieldVerifier
}

func (rdb OrderRepositoryDB) Create(o domain.Order) (*domain.Order, *errs.AppError) {
	insertQuery := `INSERT INTO orders (uuid, external_id, status, amount, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`

	res, sqlxErr := rdb.client.Exec(insertQuery, o.UUID, o.ExternalId, o.Status, o.Amount, o.UserId, o.CreatedAt, o.UpdatedAt)
	if sqlxErr != nil {
		logger.Error("Error while creating new order " + sqlxErr.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	id, sqlxErr := res.LastInsertId()
	if sqlxErr != nil {
		logger.Error("Error while getting last insert id for new order " + sqlxErr.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	o.ID = uint64(id)

	return &o, nil
}

func (rdb OrderRepositoryDB) FindById(id uint64) (*domain.Order, *errs.AppError) {
	query := `
        SELECT 
            o.id,
            o.uuid,
            o.external_id,
            o.status,
            o.amount,
            o.created_at,
            o.updated_at,
            u.id as user_id,
            u.uuid as user_uuid,
            u.name as user_name,
            u.email as user_email,
			u.created_at as user_created_at,
            oi.id as order_item_id,
            oi.quantity as order_item_quantity,
            oi.amount as order_item_amount,
            p.id as product_id,
            p.uuid as product_uuid,
            p.name as product_name,
            p.slug as product_slug,
            p.image as product_image,
            p.description as product_description,
            p.amount as product_amount,
            p.created_at as product_created_at
        FROM orders o
        LEFT JOIN users u ON o.user_id = u.id
        LEFT JOIN order_items oi ON oi.order_id = o.id
        LEFT JOIN products p ON oi.product_id = p.id
        WHERE o.id = ?
    `

	rows, err := rdb.client.Queryx(query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Order not found")
		}
		logger.Error("Error while querying order table " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}
	defer rows.Close()

	var order *domain.Order
	orderItems := make(map[uint64]domain.OrderItem)

	for rows.Next() {
		var row struct {
			ID               uint64         `db:"id"`
			UUIDBytes        []byte         `db:"uuid"`
			ExternalID       string         `db:"external_id"`
			Status           string         `db:"status"`
			Amount           int32          `db:"amount"`
			CreatedAt        time.Time      `db:"created_at"`
			UpdatedAt        time.Time      `db:"updated_at"`
			UserID           uint64         `db:"user_id"`
			UserUUIDBytes    []byte         `db:"user_uuid"`
			UserName         string         `db:"user_name"`
			UserEmail        string         `db:"user_email"`
			UserCreatedAt    time.Time      `db:"user_created_at"`
			OrderItemID      sql.NullInt64  `db:"order_item_id"`
			Quantity         sql.NullInt32  `db:"order_item_quantity"`
			ItemAmount       sql.NullInt32  `db:"order_item_amount"`
			ProductID        sql.NullInt64  `db:"product_id"`
			ProductUUIDBytes []byte         `db:"product_uuid"`
			ProductName      sql.NullString `db:"product_name"`
			ProductSlug      sql.NullString `db:"product_slug"`
			ProductImage     sql.NullString `db:"product_image"`
			ProductDesc      sql.NullString `db:"product_description"`
			ProductAmount    sql.NullInt32  `db:"product_amount"`
			ProductCreatedAt time.Time      `db:"product_created_at"`
		}

		if err := rows.StructScan(&row); err != nil {
			logger.Error("Error while scanning order row " + err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}

		orderUUID, err := db.ProcessUUID(row.UUIDBytes)
		if err != nil {
			return nil, errs.NewUnexpectedError("error processing UUID")
		}

		userUUID, err := db.ProcessUUID(row.UserUUIDBytes)
		if err != nil {
			return nil, errs.NewUnexpectedError("error processing UUID")
		}

		// Initialize order only once
		if order == nil {
			order = &domain.Order{
				ID:         row.ID,
				UUID:       orderUUID,
				ExternalId: row.ExternalID,
				Status:     row.Status,
				Amount:     row.Amount,
				CreatedAt:  row.CreatedAt,
				UpdatedAt:  row.UpdatedAt,
				User: domain.User{
					Id:        int64(row.UserID),
					UUID:      userUUID,
					Name:      row.UserName,
					Email:     row.UserEmail,
					CreatedAt: row.UserCreatedAt,
				},
			}
		}

		// Add order items if they exist
		if row.OrderItemID.Valid {
			orderItemID := uint64(row.OrderItemID.Int64)
			if _, exists := orderItems[orderItemID]; !exists {
				productUUID, err := db.ProcessUUID(row.ProductUUIDBytes)
				if err != nil {
					return nil, errs.NewUnexpectedError("error processing product UUID")
				}

				orderItems[orderItemID] = domain.OrderItem{
					ID:       orderItemID,
					Quantity: row.Quantity.Int32,
					Amount:   row.ItemAmount.Int32,
					Product: domain.Product{
						Id:          row.ProductID.Int64,
						UUID:        productUUID,
						Name:        row.ProductName.String,
						Description: row.ProductDesc.String,
						Amount:      row.ProductAmount.Int32,
						Image:       row.ProductImage.String,
						Slug:        row.ProductSlug.String,
						CreatedAt:   row.ProductCreatedAt,
					},
				}
			}
		}
	}

	// Convert order items map to slice
	var items []domain.OrderItem
	for _, item := range orderItems {
		items = append(items, item)
	}
	order.OrderItems = items

	if order == nil {
		return nil, errs.NewNotFoundError("Order not found")
	}

	return order, nil
}

func NewOrderRepositoryDB(dbClient *sqlx.DB) OrderRepositoryDB {
	return OrderRepositoryDB{
		client: dbClient,
		verifier: &db.FieldVerifier{
			DB:        dbClient,
			TableName: "orders",
		},
	}
}
