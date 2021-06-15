package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/stackus/ftgogo/serviceapis/orderapi"
	"github.com/stackus/ftgogo/store-web/internal/application"
	"github.com/stackus/ftgogo/store-web/internal/application/commands"
	"github.com/stackus/ftgogo/store-web/internal/application/queries"
	"github.com/stackus/ftgogo/store-web/internal/domain"
	"shared-go/web"
)

// To regenerate the API and a Chi server use the following generate command
//go:generate oapi-codegen -config oapi-codegen.cfg.yaml openapi.yaml
//
// The generated Chi server components are not used to construct the server it
// is instead used to verify our web handlers cover the API completely.

type WebHandlers struct {
	app application.Service
}

var _ ServerInterface = (*WebHandlers)(nil)

const jwtAudience = "web"
const consumerCtxKey = "consumerID"

func NewWebHandlers(app application.Service) WebHandlers {
	return WebHandlers{
		app: app,
	}
}

func (h WebHandlers) Mount(r chi.Router) http.Handler {
	r.Route("/accounts", func(r chi.Router) {
		r.Route("/{accountID}", func(r chi.Router) {
			r.Get("/", h.withAccountID(h.GetAccount))
			r.Put("/disable", h.withAccountID(h.DisableAccount))
			r.Put("/enable", h.withAccountID(h.EnableAccount))
		})
	})

	r.Route("/consumers", func(r chi.Router) {
		r.Route("/{consumerID}", func(r chi.Router) {
			r.Get("/", h.withConsumerID(h.GetConsumer))
		})
	})

	r.Route("/couriers", func(r chi.Router) {
		r.Route("/{courierID}", func(r chi.Router) {
			r.Get("/availability", h.withCourierID(h.SetCourierAvailability))
		})
	})

	r.Route("/deliveries", func(r chi.Router) {
		r.Route("/{deliveryID}", func(r chi.Router) {
			r.Get("/", h.withDeliveryID(h.GetDeliveryHistory))
		})
	})

	r.Route("/orders", func(r chi.Router) {
		r.Route("/{orderID}", func(r chi.Router) {
			r.Get("/", h.withOrderID(h.GetOrder))
			r.Put("/cancel", h.withOrderID(h.CancelOrder))
		})
	})

	r.Route("/restaurants", func(r chi.Router) {
		r.Route("/{restaurantID}", func(r chi.Router) {
			r.Get("/", h.withRestaurantID(h.GetRestaurant))
		})
	})

	return r
}

func (h WebHandlers) withAccountID(next func(http.ResponseWriter, *http.Request, AccountID)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r, AccountID(chi.URLParam(r, "accountID")))
	}
}

func (h WebHandlers) GetAccount(w http.ResponseWriter, r *http.Request, accountID AccountID) {
	account, err := h.app.Queries.GetAccount.Handle(r.Context(), queries.GetAccount{
		AccountID: string(accountID),
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, AccountResponse{
		Account: h.toAccountJson(account),
	})
}

func (h WebHandlers) DisableAccount(w http.ResponseWriter, r *http.Request, accountID AccountID) {
	err := h.app.Commands.DisableAccount.Handle(r.Context(), commands.DisableAccount{
		AccountID: string(accountID),
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, AccountIDResponse{Id: string(accountID)})
}

func (h WebHandlers) EnableAccount(w http.ResponseWriter, r *http.Request, accountID AccountID) {
	err := h.app.Commands.EnableAccount.Handle(r.Context(), commands.EnableAccount{
		AccountID: string(accountID),
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, AccountIDResponse{Id: string(accountID)})
}

func (h WebHandlers) withConsumerID(next func(http.ResponseWriter, *http.Request, ConsumerID)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r, ConsumerID(chi.URLParam(r, "consumerID")))
	}
}

func (h WebHandlers) GetConsumer(w http.ResponseWriter, r *http.Request, consumerID ConsumerID) {
	consumer, err := h.app.Queries.GetConsumer.Handle(r.Context(), queries.GetConsumer{
		ConsumerID: string(consumerID),
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, ConsumerResponse{
		ConsumerId: consumer.ConsumerID,
		Name:       consumer.Name,
	})
}

func (h WebHandlers) withCourierID(next func(http.ResponseWriter, *http.Request, CourierID)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r, CourierID(chi.URLParam(r, "courierID")))
	}
}

func (h WebHandlers) SetCourierAvailability(w http.ResponseWriter, r *http.Request, courierID CourierID) {
	request := SetCourierAvailabilityJSONRequestBody{}
	if err := render.Decode(r, &request); err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	err := h.app.Commands.SetCourierAvailability.Handle(r.Context(), commands.SetCourierAvailability{
		CourierID: string(courierID),
		Available: request.Available,
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, CourierAvailableResponse{
		Available: request.Available,
	})
}

func (h WebHandlers) withDeliveryID(next func(http.ResponseWriter, *http.Request, DeliveryID)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r, DeliveryID(chi.URLParam(r, "deliveryID")))
	}
}

func (h WebHandlers) GetDeliveryHistory(w http.ResponseWriter, r *http.Request, deliveryID DeliveryID) {
	history, err := h.app.Queries.GetDeliveryHistory.Handle(r.Context(), queries.GetDeliveryHistory{DeliveryID: string(deliveryID)})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, DeliveryHistoryResponse{
		DeliveryHistory: h.toDeliveryHistoryJson(history),
	})
}

func (h WebHandlers) withOrderID(next func(http.ResponseWriter, *http.Request, OrderID)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r, OrderID(chi.URLParam(r, "orderID")))
	}
}

func (h WebHandlers) GetOrder(w http.ResponseWriter, r *http.Request, orderID OrderID) {
	order, err := h.app.Queries.GetOrder.Handle(r.Context(), queries.GetOrder{
		OrderID: string(orderID),
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, OrderResponse{
		Order: h.toOrderJson(order),
	})
}

func (h WebHandlers) CancelOrder(w http.ResponseWriter, r *http.Request, orderID OrderID) {
	status, err := h.app.Commands.CancelOrder.Handle(r.Context(), commands.CancelOrder{
		OrderID: string(orderID),
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, OrderStatusResponse{
		Status: h.toOrderStateJson(status),
	})
}

func (h WebHandlers) CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (h WebHandlers) withRestaurantID(next func(http.ResponseWriter, *http.Request, RestaurantID)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r, RestaurantID(chi.URLParam(r, "restaurantID")))
	}
}

func (h WebHandlers) GetRestaurant(w http.ResponseWriter, r *http.Request, restaurantID RestaurantID) {
	panic("implement me")
}

func (h WebHandlers) toAccountJson(account *domain.Account) Account {
	return Account{
		Enabled:   account.Enabled,
		AccountId: account.AccountID,
	}
}

func (h WebHandlers) toDeliveryHistoryJson(history *domain.DeliveryHistory) DeliveryHistory {
	return DeliveryHistory{
		AssignedCourier: history.AssignedCourier,
		CourierActions:  history.CourierActions,
		Id:              history.ID,
		Status:          history.Status,
	}
}

func (h WebHandlers) toOrderJson(order *domain.Order) Order {
	return Order{
		OrderId:    order.OrderID,
		OrderTotal: order.Total,
		Status:     h.toOrderStateJson(order.Status),
	}
}

func (h WebHandlers) toOrderStateJson(orderState orderapi.OrderState) OrderState {
	switch orderState {
	case orderapi.ApprovalPending:
		return OrderStateApprovalPending
	case orderapi.Approved:
		return OrderStateApproved
	case orderapi.CancelPending:
		return OrderStateCancelPending
	case orderapi.Cancelled:
		return OrderStateCancelled
	case orderapi.RevisionPending:
		return OrderStateRevisionPending
	case orderapi.Rejected:
		return OrderStateRejected
	default:
		return OrderStateUnknown
	}
}
