package handlers

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/customer-web/internal/application"
	"github.com/stackus/ftgogo/customer-web/internal/application/commands"
	"github.com/stackus/ftgogo/customer-web/internal/application/queries"
	"github.com/stackus/ftgogo/customer-web/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
	"shared-go/web"
)

// To regenerate the API and a Chi server use the following generate command
//go:generate oapi-codegen -config oapi-codegen.cfg.yaml openapi.yaml
//
// The generated Chi server components are not used to construct the server it
// is instead used to verify our web handlers cover the API completely.

type WebHandlers struct {
	app     application.Service
	jwtAuth *jwtauth.JWTAuth
}

var _ ServerInterface = (*WebHandlers)(nil)

const jwtAudience = "web"
const consumerCtxKey = "consumerID"

func NewWebHandlers(app application.Service) WebHandlers {
	// Generate a new set of keys with each start; for reasons
	// Can use this until scaling the customer-web is added to the demo
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	return WebHandlers{
		app:     app,
		jwtAuth: jwtauth.New(jwa.RS256.String(), privateKey, privateKey.PublicKey),
	}
}

func (h WebHandlers) Mount(r chi.Router) http.Handler {
	// Public Routes
	r.Post("/register", h.RegisterConsumer)
	r.Post("/signin", h.SignInConsumer)

	// Protected Routes
	r.Group(func(r chi.Router) {
		// JWT Session Authentication
		r.Use(
			jwtauth.Verifier(h.jwtAuth),
			jwtauth.Authenticator,
			h.decodeClaimsIntoContext,
		)

		r.Route("/consumer", func(r chi.Router) {
			r.Get("/", h.GetConsumer)
		})

		r.Route("/addresses", func(r chi.Router) {
			r.Post("/", h.AddConsumerAddress)
			r.Route("/{consumerAddressID}", func(r chi.Router) {
				r.Get("/", h.withConsumerAddressID(h.GetConsumerAddress))
				r.Put("/", h.withConsumerAddressID(h.UpdateConsumerAddress))
				r.Delete("/", h.withConsumerAddressID(h.RemoveConsumerAddress))
			})
		})

		r.Route("/orders", func(r chi.Router) {
			r.Post("/", h.CreateOrder)
			r.Get("/", h.withSearchOrdersParams(h.SearchOrders))
			r.Route("/{orderID}", func(r chi.Router) {
				r.Get("/", h.withOrderID(h.GetOrder))
				r.Put("/revise", h.withOrderID(h.ReviseOrder))
				r.Put("/cancel", h.withOrderID(h.CancelOrder))
			})
		})

		r.Route("/restaurants", func(r chi.Router) {
			// TODO Search Restaurants Endpoint
			// r.Get("/", h.withSearchRestaurantsParams(h.SearchRestaurants))
			r.Route("/{restaurantID}", func(r chi.Router) {
				r.Get("/", h.withRestaurantID(h.GetRestaurant))
			})
		})
	})

	return r
}

func (h WebHandlers) decodeClaimsIntoContext(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			render.Render(w, r, web.NewErrorResponse(err))
			return
		}
		var consumerID string
		if subject, exists := claims[jwt.SubjectKey]; !exists {
			render.Render(w, r, web.NewErrorResponse(errors.Wrap(errors.ErrUnauthenticated, "missing claims subject")))
			return
		} else {
			switch s := subject.(type) {
			case string:
				consumerID = s
			case []string:
				if len(s) == 0 {
					render.Render(w, r, web.NewErrorResponse(errors.Wrap(errors.ErrUnauthenticated, "invalid claims subject")))
					return
				}
				consumerID = s[0]
			default:
				render.Render(w, r, web.NewErrorResponse(errors.Wrap(errors.ErrUnauthenticated, "invalid claims subject")))
				return
			}
		}
		ctx := context.WithValue(r.Context(), consumerCtxKey, consumerID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// Setting and fetching the consumerID into and from claims is strictly the responsibility of these web handlers!
//
// The application (commands, queries), the domain, and the adapters don't know and should not know where any
// session values might come from.
func (h WebHandlers) consumerID(ctx context.Context) string {
	v := ctx.Value(consumerCtxKey)
	switch s := v.(type) {
	case string:
		return s
	default:
		return ""
	}
}

func (h WebHandlers) SignInConsumer(w http.ResponseWriter, r *http.Request) {
	// Simple JWT authentication
	request := SignInConsumerJSONRequestBody{}
	if err := render.Decode(r, &request); err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	// For Demo: "login" with any known consumerID
	consumer, err := h.app.Queries.GetConsumer.Handle(r.Context(), queries.GetConsumer{
		ConsumerID: request.ConsumerId,
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	var token string
	_, token, err = h.jwtAuth.Encode(map[string]interface{}{
		jwt.SubjectKey:    consumer.ConsumerID,
		jwt.AudienceKey:   jwtAudience,
		jwt.ExpirationKey: jwtauth.ExpireIn(24 * time.Hour),
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(errors.Wrap(errors.ErrInternal, "could not authenticate you at this time")))
		return
	}

	render.Respond(w, r, SignInResponse{Token: token})
}

func (h WebHandlers) RegisterConsumer(w http.ResponseWriter, r *http.Request) {
	request := RegisterConsumerJSONRequestBody{}
	if err := render.Decode(r, &request); err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	consumerID, err := h.app.Commands.RegisterConsumer.Handle(r.Context(), commands.RegisterConsumer{
		Name: request.Name,
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, ConsumerIDResponse{Id: consumerID})
}

func (h WebHandlers) GetConsumer(w http.ResponseWriter, r *http.Request) {
	consumer, err := h.app.Queries.GetConsumer.Handle(r.Context(), queries.GetConsumer{
		ConsumerID: h.consumerID(r.Context()),
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

func (h WebHandlers) AddConsumerAddress(w http.ResponseWriter, r *http.Request) {
	request := AddConsumerAddressJSONRequestBody{}
	if err := render.Decode(r, &request); err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	consumerID := h.consumerID(r.Context())
	err := h.app.Commands.AddConsumerAddress.Handle(r.Context(), commands.AddConsumerAddress{
		ConsumerID: consumerID,
		AddressID:  request.Name,
		Address:    &request.Address,
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, ConsumerAddressIDResponse{
		AddressId:  request.Name,
		ConsumerId: consumerID,
	})
}

func (h WebHandlers) withConsumerAddressID(next func(http.ResponseWriter, *http.Request, ConsumerAddressID)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r, ConsumerAddressID(chi.URLParam(r, "consumerAddressID")))
	}
}

func (h WebHandlers) RemoveConsumerAddress(w http.ResponseWriter, r *http.Request, consumerAddressID ConsumerAddressID) {
	consumerID := h.consumerID(r.Context())
	err := h.app.Commands.RemoveConsumerAddress.Handle(r.Context(), commands.RemoveConsumerAddress{
		ConsumerID: consumerID,
		AddressID:  string(consumerAddressID),
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Status(r, http.StatusCreated)
}

func (h WebHandlers) GetConsumerAddress(w http.ResponseWriter, r *http.Request, consumerAddressID ConsumerAddressID) {
	consumerID := h.consumerID(r.Context())
	address, err := h.app.Queries.GetConsumerAddress.Handle(r.Context(), queries.GetConsumerAddress{
		ConsumerID: consumerID,
		AddressID:  string(consumerAddressID),
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, ConsumerAddressResponse{
		Address: *address,
	})
}

func (h WebHandlers) UpdateConsumerAddress(w http.ResponseWriter, r *http.Request, consumerAddressID ConsumerAddressID) {
	request := UpdateConsumerAddressJSONRequestBody{}
	if err := render.Decode(r, &request); err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	consumerID := h.consumerID(r.Context())
	err := h.app.Commands.UpdateConsumerAddress.Handle(r.Context(), commands.UpdateConsumerAddress{
		ConsumerID: consumerID,
		AddressID:  string(consumerAddressID),
		Address:    &request.Address,
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, ConsumerAddressIDResponse{
		AddressId:  string(consumerAddressID),
		ConsumerId: consumerID,
	})
}

func (h WebHandlers) CreateOrder(w http.ResponseWriter, r *http.Request) {
	request := CreateOrderJSONRequestBody{}
	if err := render.Decode(r, &request); err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	orderID, err := h.app.Commands.CreateOrder.Handle(r.Context(), commands.CreateOrder{
		ConsumerID:   request.ConsumerId,
		RestaurantID: request.RestaurantId,
		AddressID:    request.AddressId,
		LineItems:    request.LineItems,
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, OrderIDResponse{Id: orderID})
}

func (h WebHandlers) withOrderID(next func(http.ResponseWriter, *http.Request, OrderID)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r, OrderID(chi.URLParam(r, "orderID")))
	}
}

func (h WebHandlers) GetOrder(w http.ResponseWriter, r *http.Request, orderID OrderID) {
	consumerID := h.consumerID(r.Context())
	order, err := h.app.Queries.GetOrder.Handle(r.Context(), queries.GetOrder{
		OrderID:    string(orderID),
		ConsumerID: consumerID,
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
	consumerID := h.consumerID(r.Context())
	status, err := h.app.Commands.CancelOrder.Handle(r.Context(), commands.CancelOrder{
		ConsumerID: consumerID,
		OrderID:    string(orderID),
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, OrderStatusResponse{
		Status: h.toOrderStateJson(status),
	})
}

func (h WebHandlers) ReviseOrder(w http.ResponseWriter, r *http.Request, orderID OrderID) {
	consumerID := h.consumerID(r.Context())
	status, err := h.app.Commands.ReviseOrder.Handle(r.Context(), commands.ReviseOrder{
		ConsumerID: consumerID,
		OrderID:    string(orderID),
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, OrderStatusResponse{
		Status: h.toOrderStateJson(status),
	})
}

func (h WebHandlers) withSearchOrdersParams(next func(w http.ResponseWriter, r *http.Request, params SearchOrdersParams)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		// FRAGILE: Copied from generated code
		// Parameter object where we will unmarshal all parameters from the context
		var params SearchOrdersParams

		err = runtime.BindQueryParameter("deepObject", true, false, "filter", r.URL.Query(), &params.Filter)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid format for parameter filter: %s", err), http.StatusBadRequest)
			return
		}

		err = runtime.BindQueryParameter("form", true, false, "next", r.URL.Query(), &params.Next)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid format for parameter next: %s", err), http.StatusBadRequest)
			return
		}

		err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid format for parameter limit: %s", err), http.StatusBadRequest)
			return
		}

		next(w, r, params)
	}
}

func (h WebHandlers) SearchOrders(w http.ResponseWriter, r *http.Request, params SearchOrdersParams) {
	consumerID := h.consumerID(r.Context())

	filters, next, limit := h.fromSearchOrdersParamsJson(params)

	results, err := h.app.Queries.SearchOrders.Handle(r.Context(), queries.SearchOrders{
		ConsumerID: consumerID,
		Filters:    filters,
		Next:       next,
		Limit:      limit,
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, SearchOrdersResponse{
		Next:   results.Next,
		Orders: nil,
	})
}

func (h WebHandlers) withRestaurantID(next func(http.ResponseWriter, *http.Request, RestaurantID)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r, RestaurantID(chi.URLParam(r, "restaurantID")))
	}
}

func (h WebHandlers) GetRestaurant(w http.ResponseWriter, r *http.Request, restaurantID RestaurantID) {
	panic("implement me")
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

func (h WebHandlers) fromOrderStateJson(orderState OrderState) orderapi.OrderState {
	switch orderState {
	case OrderStateApprovalPending:
		return orderapi.ApprovalPending
	case OrderStateApproved:
		return orderapi.Approved
	case OrderStateCancelPending:
		return orderapi.CancelPending
	case OrderStateCancelled:
		return orderapi.Cancelled
	case OrderStateRevisionPending:
		return orderapi.RevisionPending
	case OrderStateRejected:
		return orderapi.Rejected
	default:
		return orderapi.UnknownOrderState
	}
}

func (h WebHandlers) fromSearchOrdersParamsJson(params SearchOrdersParams) (*domain.SearchOrdersFilters, string, int) {
	keywords := []string{}
	if params.Filter.Keywords != nil {
		keywords = *params.Filter.Keywords
	}
	since := time.Time{}
	if params.Filter.Since != nil {
		since = *params.Filter.Since
	}
	status := orderapi.UnknownOrderState
	if params.Filter.Status != nil {
		status = h.fromOrderStateJson(*params.Filter.Status)
	}

	filters := &domain.SearchOrdersFilters{
		Keywords: keywords,
		Since:    since,
		Status:   status,
	}
	next := ""
	if params.Next != nil {
		next = string(*params.Next)
	}
	limit := 0
	if params.Limit != nil {
		limit = int(*params.Limit)
	}
	return filters, next, limit
}

func (h WebHandlers) toOrderDetailJson(order *domain.OrderHistory) OrderDetail {
	return OrderDetail{
		CreatedAt:      order.CreatedAt,
		OrderId:        order.OrderID,
		RestaurantId:   order.RestaurantID,
		RestaurantName: order.RestaurantName,
		Status:         h.toOrderStateJson(order.Status),
	}
}
