package repositories

import (
	"github.com/go-ms-project-store/internal/core/domain"
	"github.com/go-ms-project-store/internal/pkg/db"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type OrderItemRepositoryDB struct {
	client   *sqlx.DB
	verifier *db.FieldVerifier
}

func (rdb OrderItemRepositoryDB) Create(o domain.OrderItem) (*domain.OrderItem, *errs.AppError) {
	insertQuery := `INSERT INTO order_items (amount, quantity, order_id, product_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`

	res, sqlxErr := rdb.client.Exec(insertQuery, o.Amount, o.Quantity, o.OrderId, o.ProductId, o.CreatedAt, o.UpdatedAt)
	if sqlxErr != nil {
		logger.Error("Error while creating new order item" + sqlxErr.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	id, sqlxErr := res.LastInsertId()
	if sqlxErr != nil {
		logger.Error("Error while getting last insert id for new order item " + sqlxErr.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	o.ID = uint64(id)

	return &o, nil
}

func NewOrderItemRepositoryDB(dbClient *sqlx.DB) OrderItemRepositoryDB {
	return OrderItemRepositoryDB{
		client: dbClient,
		verifier: &db.FieldVerifier{
			DB:        dbClient,
			TableName: "order_items",
		},
	}
}
