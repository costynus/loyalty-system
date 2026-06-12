package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type (
    Order struct {
        Number string `json:"number"`
        Status string `json:"status"`
        Accrual decimal.Decimal `json:"accreal"`
        UploadedAt time.Time`json:"uploaded_at"`
    }
)
