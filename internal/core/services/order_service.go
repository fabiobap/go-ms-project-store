package services

import (
	"net/http"
	"time"

	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/go-ms-project-store/internal/core/domain"
	"github.com/go-ms-project-store/internal/core/ports"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/google/uuid"
)

type DefaultOrderService struct {
	repo ports.OrderRepository
}

func (s DefaultOrderService) CreateOrder(req dto.NewOrderRequest, user_id uint64) (*domain.Order, *errs.AppError) {
	productUUIDs := make([]string, len(req.Products))

	// Extract UUIDs from the request
	for i, product := range req.Products {
		productUUIDs[i] = product.ID
	}

	// Get products from database using WhereIn
	products, err := s.repo.ProductRepo().WhereIn(productUUIDs)
	if err != nil {
		return nil, err
	}

	var total_amount int32
	orderItems := make([]domain.OrderItem, 0, len(req.Products))

	for _, dbProduct := range products {
		// Find the quantity from request for this product
		for _, reqProduct := range req.Products {
			if dbProduct.UUID.String() == reqProduct.ID {
				total_amount += dbProduct.Amount * int32(reqProduct.Quantity)

				orderItem := domain.OrderItem{
					ProductId: uint64(dbProduct.Id),
					Quantity:  int32(reqProduct.Quantity),
					Amount:    dbProduct.Amount,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				orderItems = append(orderItems, orderItem)
				break
			}
		}
	}

	order := domain.Order{
		UUID:       uuid.New(),
		Status:     "pending",
		Amount:     total_amount,
		UserId:     user_id,
		ExternalId: uuid.New().String(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		OrderItems: orderItems,
	}

	newOrder, err := s.repo.Create(order)
	if err != nil {
		if err.Code != http.StatusUnprocessableEntity {
			return nil, errs.NewUnexpectedError("unexpected database error")
		} else {
			return nil, err
		}
	}

	return newOrder, nil
}

func NewOrderService(repository ports.OrderRepository) DefaultOrderService {
	return DefaultOrderService{repo: repository}
}
