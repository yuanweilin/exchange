package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderType string

type OrderBook struct {
	Buy       []Order `json:"buy"`
	Sell      []Order `json:"sell"`
	TotalBuy  []Order `json:"total_buy"`
	TotalSell []Order `json:"total_sell"`
}

const (
	OrderTypeBuy  OrderType = "buy"
	OrderTypeSell OrderType = "sell"
)

type Order struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Amount    decimal.Decimal `json:"amount"`
	Price     decimal.Decimal `json:"price"`
	Type      OrderType       `gorm:"type:varchar(10)" json:"type"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
