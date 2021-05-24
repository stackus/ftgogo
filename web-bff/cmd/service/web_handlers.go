package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/web-bff/internal/application/commands"
	"github.com/stackus/ftgogo/web-bff/internal/application/queries"
	"shared-go/web"
)

// To regenerate the API and a Chi server use the following generate command
//go:generate oapi-codegen -config oapi-codegen.cfg.yaml openapi.yaml
//
// The generated Chi server components are not used to construct the server it
// is instead used to verify our web handlers cover the API completely.

type webHandlers struct {
	app     Application
	jwtAuth *jwtauth.JWTAuth
}

const jwtAudience = "web"
const consumerCtxKey = "consumerID"

func newWebHandlers(app Application) webHandlers {
	// Generate a new set of keys with each start; for reasons
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	return webHandlers{
		app:     app,
		jwtAuth: jwtauth.New(jwa.RS256.String(), privateKey, privateKey.PublicKey),
	}
}

func (h webHandlers) Mount(r chi.Router) http.Handler {
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

		r.Route("/consumers", func(r chi.Router) {
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
			r.Route("/{orderID}", func(r chi.Router) {
				r.Get("/", h.withOrderID(h.GetOrder))
				r.Put("/revise", h.withOrderID(h.ReviseOrder))
				r.Put("/cancel", h.withOrderID(h.CancelOrder))
			})
		})

		r.Route("/restaurants", func(r chi.Router) {
			r.Route("/{restaurantID}", func(r chi.Router) {
				r.Get("/", h.withRestaurantID(h.GetRestaurant))
			})
		})
	})

	return r
}

func (h webHandlers) decodeClaimsIntoContext(next http.Handler) http.Handler {
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
func (h webHandlers) consumerID(ctx context.Context) string {
	v := ctx.Value(consumerCtxKey)
	switch s := v.(type) {
	case string:
		return s
	default:
		return ""
	}
}

func (h webHandlers) SignInConsumer(w http.ResponseWriter, r *http.Request) {
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

func (h webHandlers) RegisterConsumer(w http.ResponseWriter, r *http.Request) {
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

func (h webHandlers) GetConsumer(w http.ResponseWriter, r *http.Request) {
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

func (h webHandlers) AddConsumerAddress(w http.ResponseWriter, r *http.Request) {
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

func (h webHandlers) withConsumerAddressID(next func(http.ResponseWriter, *http.Request, ConsumerAddressID)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r, ConsumerAddressID(chi.URLParam(r, "consumerAddressID")))
	}
}

func (h webHandlers) RemoveConsumerAddress(w http.ResponseWriter, r *http.Request, consumerAddressID ConsumerAddressID) {
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

func (h webHandlers) GetConsumerAddress(w http.ResponseWriter, r *http.Request, consumerAddressID ConsumerAddressID) {
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

func (h webHandlers) UpdateConsumerAddress(w http.ResponseWriter, r *http.Request, consumerAddressID ConsumerAddressID) {
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

func (h webHandlers) CreateOrder(w http.ResponseWriter, r *http.Request) {
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

func (h webHandlers) withOrderID(next func(http.ResponseWriter, *http.Request, OrderID)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r, OrderID(chi.URLParam(r, "orderID")))
	}
}

func (h webHandlers) GetOrder(w http.ResponseWriter, r *http.Request, orderID OrderID) {
	panic("implement me")
}

func (h webHandlers) CancelOrder(w http.ResponseWriter, r *http.Request, orderID OrderID) {
	panic("implement me")
}

func (h webHandlers) ReviseOrder(w http.ResponseWriter, r *http.Request, orderID OrderID) {
	panic("implement me")
}

func (h webHandlers) withRestaurantID(next func(http.ResponseWriter, *http.Request, RestaurantID)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r, RestaurantID(chi.URLParam(r, "restaurantID")))
	}
}

func (h webHandlers) GetRestaurant(w http.ResponseWriter, r *http.Request, restaurantID RestaurantID) {
	panic("implement me")
}
