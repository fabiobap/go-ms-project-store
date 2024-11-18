package dto

type NewOrderItemDTO struct {
	Amount    int32
	Quantity  int32
	OrderId   uint64
	ProductId uint64
}
