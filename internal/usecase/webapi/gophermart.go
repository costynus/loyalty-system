package webapi

import (
	"fmt"
	"net/http"

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


func (w *GophermartWebAPI) GetOrderInfo(orderNumber string) (entity.Order, error) {
    var order entity.Order
    resp, err := w.client.
        R().
        SetResult(order).
        Get("/api/orders/" + orderNumber)
    if err != nil {
        return entity.Order{}, fmt.Errorf("WebAPI - GetOrderInfo - w.client.R().Get: %w", err)
    }
    switch resp.StatusCode() {
    case http.StatusInternalServerError:
        return entity.Order{}, ErrInternalServerError
    case http.StatusTooManyRequests:
        return entity.Order{}, ErrTooManyRequests
    default:
        return order, nil
    }
}
