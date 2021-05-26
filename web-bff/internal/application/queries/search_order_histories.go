package queries

import (
	"time"

	"github.com/stackus/ftgogo/serviceapis/orderapi"
	"github.com/stackus/ftgogo/web-bff/internal/domain"
)

type SearchOrderHistories struct {
	ConsumerID string
	Filter     *struct {
		Keywords *[]string
		Since    *time.Time
		Status   *orderapi.OrderState
	}
	Next  *string
	Limit *int
}

type SearchOrderHistoriesHandler struct {
	repo domain.OrderHistoryRepository
}
