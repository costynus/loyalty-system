package webapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/costynus/loyalty-system/internal/entity"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type GophermartWebAPI struct {
	client *resty.Client
}

var ErrTooManyRequests = errors.New("too many requests")
var ErrInternalServerError = errors.New("internal server error")

func New(client *resty.Client) *GophermartWebAPI {
	return &GophermartWebAPI{
		client: client,
	}
}

func (w *GophermartWebAPI) GetOrderInfo(orderNumber string) (entity.Order, time.Duration, error) {
	var order entity.OrderAddappter
	resp, err := w.client.
		R().
		SetResult(&order).
		Get("/api/orders/" + orderNumber)
	if err != nil {
		return entity.Order{}, 0, fmt.Errorf("WebAPI - GetOrderInfo - w.client.R().Get: %w", err)
	}
	switch resp.StatusCode() {
	case http.StatusInternalServerError:
		return entity.Order{}, 0, ErrInternalServerError
	case http.StatusTooManyRequests:
		timeout, err := time.ParseDuration(resp.Header().Get("Retry-After") + "s")
		if err != nil {
			return entity.Order{}, 0, fmt.Errorf("WebAPI - GetOrderInfo - time.ParseDuration: %w", err)
		}
		return entity.Order{}, timeout, ErrTooManyRequests
	default:
		return entity.Order{
			Number:  order.Number,
			Status:  order.Status,
			Accrual: order.Accrual,
		}, 0, nil
	}
}
