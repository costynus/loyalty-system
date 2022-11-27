package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type (
    Order struct {
        Number string `json:"number" db:"order_number"`
        Status string `json:"status" db:"status"`
        Accrual decimal.Decimal `json:"accrual,omitempty" db:"accrual"`
        UploadedAt time.Time`json:"uploaded_at" db:"uploaded_at"`
        UserID int `json:"-" db:"user_id"`
    }

    OrderAddappter struct {
        Number string `json:"order"`
        Status string `json:"status"`
        Accrual decimal.Decimal `json:"accrual"`
    }
)
