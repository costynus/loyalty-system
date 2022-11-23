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
        Current decimal.Decimal `json:"current"`
        Withdraw decimal.Decimal `json:"withdraw"`
    }

    Withdraw struct {
        Order string `json:"order"`
        Sum decimal.Decimal `json:"sum"`
        ProcessedAt time.Time `json:"processed_at"` 
    }
)
