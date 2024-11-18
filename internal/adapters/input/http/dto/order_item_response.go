package dto

type OrderItemResponse struct {
	Amount   string                `json:"amount"`
	Quantity int32                 `json:"quantity"`
	Product  ProductPublicResponse `json:"product"`
}
