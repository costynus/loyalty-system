package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type (
    Order struct {
        Number string `json:"number" db:"order_number"`
        Status string `json:"status" db:"status"`
        Accrual decimal.Decimal `json:"accreal" db:"accrual"`
        UploadedAt time.Time`json:"uploaded_at" db:"uploaded_at"`
    }
)
