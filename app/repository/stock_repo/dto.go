package stockrepo

type AvailableProductStockResponse struct {
	ProductID      int64 `json:"product_id"`
	AvailableStock int64 `json:"available_stock"`
}
