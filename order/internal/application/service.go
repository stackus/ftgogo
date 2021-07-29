package application

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/application/commands"
	"github.com/stackus/ftgogo/order/internal/application/ports"
	"github.com/stackus/ftgogo/order/internal/application/queries"
	"github.com/stackus/ftgogo/order/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

type ServiceApplication interface {
	CreateOrder(ctx context.Context, cmd commands.CreateOrder) (string, error)
	ApproveOrder(ctx context.Context, cmd commands.ApproveOrder) error
	RejectOrder(ctx context.Context, cmd commands.RejectOrder) error
	BeginCancelOrder(ctx context.Context, cmd commands.BeginCancelOrder) error
	UndoCancelOrder(ctx context.Context, cmd commands.UndoCancelOrder) error
	ConfirmCancelOrder(ctx context.Context, cmd commands.ConfirmCancelOrder) error
	BeginReviseOrder(ctx context.Context, cmd commands.BeginReviseOrder) (int, error)
	UndoReviseOrder(ctx context.Context, cmd commands.UndoReviseOrder) error
	ConfirmReviseOrder(ctx context.Context, cmd commands.ConfirmReviseOrder) error
	StartCreateOrderSaga(ctx context.Context, cmd commands.StartCreateOrderSaga) error
	StartCancelOrderSaga(ctx context.Context, cmd commands.StartCancelOrderSaga) (orderapi.OrderState, error)
	StartReviseOrderSaga(ctx context.Context, cmd commands.StartReviseOrderSaga) (orderapi.OrderState, error)
	CreateRestaurant(ctx context.Context, cmd commands.CreateRestaurant) error
	ReviseRestaurantMenu(ctx context.Context, cmd commands.ReviseRestaurantMenu) error
	GetOrder(ctx context.Context, query queries.GetOrder) (*domain.Order, error)
	GetRestaurant(ctx context.Context, query queries.GetRestaurant) (*domain.Restaurant, error)
}

type Service struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateOrder          commands.CreateOrderHandler
	ApproveOrder         commands.ApproveOrderHandler
	RejectOrder          commands.RejectOrderHandler
	BeginCancelOrder     commands.BeginCancelOrderHandler
	UndoCancelOrder      commands.UndoCancelOrderHandler
	ConfirmCancelOrder   commands.ConfirmCancelOrderHandler
	BeginReviseOrder     commands.BeginReviseOrderHandler
	UndoReviseOrder      commands.UndoReviseOrderHandler
	ConfirmReviseOrder   commands.ConfirmReviseOrderHandler
	StartCreateOrderSaga commands.StartCreateOrderSagaHandler
	StartCancelOrderSaga commands.StartCancelOrderSagaHandler
	StartReviseOrderSaga commands.StartReviseOrderSagaHandler
	CreateRestaurant     commands.CreateRestaurantHandler
	ReviseRestaurantMenu commands.ReviseRestaurantMenuHandler
}

type Queries struct {
	GetOrder      queries.GetOrderHandler
	GetRestaurant queries.GetRestaurantHandler
}

func NewServiceApplication(
	orderRepo ports.OrderRepository, restaurantRepo ports.RestaurantRepository,
	createOrderSaga ports.CreateOrderSaga, cancelOrderSaga ports.CancelOrderSaga, reviseOrderSaga ports.ReviseOrderSaga,
	ordersPlacedCounter, ordersApprovedCounter, ordersRejectedCounter ports.Counter,
) *Service {
	return &Service{
		Commands: Commands{
			CreateOrder:          commands.NewCreateOrderHandler(orderRepo, restaurantRepo),
			ApproveOrder:         commands.NewApproveOrderHandler(orderRepo, ordersApprovedCounter),
			RejectOrder:          commands.NewRejectOrderHandler(orderRepo, ordersRejectedCounter),
			BeginCancelOrder:     commands.NewBeginCancelOrderHandler(orderRepo),
			UndoCancelOrder:      commands.NewUndoCancelOrderHandler(orderRepo),
			ConfirmCancelOrder:   commands.NewConfirmCancelOrderHandler(orderRepo),
			BeginReviseOrder:     commands.NewBeginReviseOrderHandler(orderRepo),
			UndoReviseOrder:      commands.NewUndoReviseOrderHandler(orderRepo),
			ConfirmReviseOrder:   commands.NewConfirmReviseOrderHandler(orderRepo),
			StartCreateOrderSaga: commands.NewStartCreateOrderSagaHandler(createOrderSaga, ordersPlacedCounter),
			StartCancelOrderSaga: commands.NewStartCancelOrderSagaHandler(orderRepo, cancelOrderSaga),
			StartReviseOrderSaga: commands.NewStartReviseOrderSagaHandler(orderRepo, reviseOrderSaga),
			CreateRestaurant:     commands.NewCreateRestaurantHandler(restaurantRepo),
			ReviseRestaurantMenu: commands.NewReviseRestaurantMenuHandler(restaurantRepo),
		},
		Queries: Queries{
			GetOrder:      queries.NewGetOrderHandler(orderRepo),
			GetRestaurant: queries.NewGetRestaurantHandler(restaurantRepo),
		},
	}
}

func (s Service) CreateOrder(ctx context.Context, cmd commands.CreateOrder) (string, error) {
	return s.Commands.CreateOrder.Handle(ctx, cmd)
}

func (s Service) ApproveOrder(ctx context.Context, cmd commands.ApproveOrder) error {
	return s.Commands.ApproveOrder.Handle(ctx, cmd)
}

func (s Service) RejectOrder(ctx context.Context, cmd commands.RejectOrder) error {
	return s.Commands.RejectOrder.Handle(ctx, cmd)
}

func (s Service) BeginCancelOrder(ctx context.Context, cmd commands.BeginCancelOrder) error {
	return s.Commands.BeginCancelOrder.Handle(ctx, cmd)
}

func (s Service) UndoCancelOrder(ctx context.Context, cmd commands.UndoCancelOrder) error {
	return s.Commands.UndoCancelOrder.Handle(ctx, cmd)
}

func (s Service) ConfirmCancelOrder(ctx context.Context, cmd commands.ConfirmCancelOrder) error {
	return s.Commands.ConfirmCancelOrder.Handle(ctx, cmd)
}

func (s Service) BeginReviseOrder(ctx context.Context, cmd commands.BeginReviseOrder) (int, error) {
	return s.Commands.BeginReviseOrder.Handle(ctx, cmd)
}

func (s Service) UndoReviseOrder(ctx context.Context, cmd commands.UndoReviseOrder) error {
	return s.Commands.UndoReviseOrder.Handle(ctx, cmd)
}

func (s Service) ConfirmReviseOrder(ctx context.Context, cmd commands.ConfirmReviseOrder) error {
	return s.Commands.ConfirmReviseOrder.Handle(ctx, cmd)
}

func (s Service) StartCreateOrderSaga(ctx context.Context, cmd commands.StartCreateOrderSaga) error {
	return s.Commands.StartCreateOrderSaga.Handle(ctx, cmd)
}

func (s Service) StartCancelOrderSaga(ctx context.Context, cmd commands.StartCancelOrderSaga) (orderapi.OrderState, error) {
	return s.Commands.StartCancelOrderSaga.Handle(ctx, cmd)
}

func (s Service) StartReviseOrderSaga(ctx context.Context, cmd commands.StartReviseOrderSaga) (orderapi.OrderState, error) {
	return s.Commands.StartReviseOrderSaga.Handle(ctx, cmd)
}

func (s Service) CreateRestaurant(ctx context.Context, cmd commands.CreateRestaurant) error {
	return s.Commands.CreateRestaurant.Handle(ctx, cmd)
}

func (s Service) ReviseRestaurantMenu(ctx context.Context, cmd commands.ReviseRestaurantMenu) error {
	return s.Commands.ReviseRestaurantMenu.Handle(ctx, cmd)
}

func (s Service) GetOrder(ctx context.Context, query queries.GetOrder) (*domain.Order, error) {
	return s.Queries.GetOrder.Handle(ctx, query)
}

func (s Service) GetRestaurant(ctx context.Context, query queries.GetRestaurant) (*domain.Restaurant, error) {
	return s.Queries.GetRestaurant.Handle(ctx, query)
}
