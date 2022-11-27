package entity

import (
	"time"

	"github.com/shopspring/decimal"
)


type (
    Withdrawal struct {
        Order string `json:"order"`
        Sum decimal.Decimal `json:"sum"`
    }

    Balance struct {
        Current decimal.Decimal `json:"current" db:"balance"`
        Withdraw decimal.Decimal `json:"withdrawn" db:"withdrawal"`
    }

    Withdraw struct {
        Order string `json:"order" db:"order_number"`
        Sum decimal.Decimal `json:"sum" db:"sum_number"`
        ProcessedAt time.Time `json:"processed_at" db:"updated_at"` 
    }
)
