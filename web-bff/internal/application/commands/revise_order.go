package commands

import (
	"context"

	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
	"github.com/stackus/ftgogo/web-bff/internal/domain"
)

type ReviseOrder struct {
	ConsumerID        string
	OrderID           string
	RevisedQuantities commonapi.MenuItemQuantities
}

type ReviseOrderHandler struct {
	repo domain.OrderRepository
}

func NewReviseOrderHandler(repo domain.OrderRepository) ReviseOrderHandler {
	return ReviseOrderHandler{repo: repo}
}

func (h ReviseOrderHandler) Handle(ctx context.Context, cmd ReviseOrder) (orderapi.OrderState, error) {
	order, err := h.repo.Find(ctx, domain.FindOrder{
		OrderID: cmd.OrderID,
	})
	if err != nil {
		return orderapi.UnknownOrderState, err
	}

	if order.ConsumerID != cmd.ConsumerID {
		// being opaque intentionally; Could also send a permission denied error
		return orderapi.UnknownOrderState, errors.Wrap(errors.ErrNotFound, "order not found")
	}

	return h.repo.Revise(ctx, domain.ReviseOrder{
		OrderID:           cmd.OrderID,
		RevisedQuantities: cmd.RevisedQuantities,
	})
}
