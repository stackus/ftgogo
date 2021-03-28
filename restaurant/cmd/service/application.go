package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/stackus/ftgogo/restaurant/internal/adapters"
	"github.com/stackus/ftgogo/restaurant/internal/application/commands"
	"github.com/stackus/ftgogo/restaurant/internal/application/queries"
	"serviceapis"
	"serviceapis/restaurantapi"
	"shared-go/applications"
	"shared-go/web"
)

// To regenerate the web server api use the following generate command
//go:generate oapi-codegen -config oapi-codegen.cfg.yaml openapi.yaml

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateRestaurant commands.CreateRestaurantHandler
}

type Queries struct {
	GetRestaurant queries.GetRestaurantHandler
}

func initApplication(svc *applications.Service) error {
	serviceapis.RegisterTypes()

	restaurantRepo := adapters.NewRestaurantPostgresRepository(svc.PgConn)
	restaurantPublisher := adapters.NewRestaurantPublisher(svc.Publisher)

	application := Application{
		Commands: Commands{
			CreateRestaurant: commands.NewCreateRestaurantHandler(restaurantRepo, restaurantPublisher),
		},
		Queries: Queries{
			GetRestaurant: queries.NewGetRestaurantHandler(restaurantRepo),
		},
	}

	svc.WebServer.Mount(svc.Cfg.Web.ApiPath, func(r chi.Router) http.Handler {
		return HandlerFromMux(NewWebHandlers(application), r)
	})

	return nil
}

type WebHandlers struct{ app Application }

func NewWebHandlers(app Application) WebHandlers { return WebHandlers{app: app} }

func (h WebHandlers) CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	request := CreateRestaurantJSONRequestBody{}

	if err := render.Decode(r, &request); err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	restaurantID, err := h.app.Commands.CreateRestaurant.Handle(r.Context(), commands.CreateRestaurant{
		Name:      request.Name,
		Address:   request.Address,
		MenuItems: request.Menu.MenuItems,
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, RestaurantIDResponse{Id: restaurantID})
}

func (h WebHandlers) GetRestaurant(w http.ResponseWriter, r *http.Request, restaurantID RestaurantID) {
	rid := string(restaurantID)

	restaurant, err := h.app.Queries.GetRestaurant.Handle(r.Context(), queries.GetRestaurant{RestaurantID: rid})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, RestaurantResponse{
		Address: restaurant.Address,
		Id:      restaurant.RestaurantID,
		Menu: struct {
			MenuItems []restaurantapi.MenuItem `json:"menu_items"`
		}{
			MenuItems: restaurant.MenuItems,
		},
		Name: restaurant.Name,
	})
}

type RpcHandlers struct {
	app Application
	restaurantapi.UnimplementedRestaurantServiceServer
}

var _ restaurantapi.RestaurantServiceServer = (*RpcHandlers)(nil)

func newRpcHandlers(app Application) RpcHandlers {
	return RpcHandlers{app: app}
}

func (h RpcHandlers) CreateRestaurant(ctx context.Context, request *restaurantapi.CreateRestaurantRequest) (*restaurantapi.CreateRestaurantResponse, error) {
	menuItems := make([]restaurantapi.MenuItem, 0)
	for _, item := range request.Menu.MenuItems {
		menuItems = append(menuItems, restaurantapi.MenuItem{
			ID:    item.ID,
			Name:  item.Name,
			Price: int(item.Price),
		})
	}

	restaurantID, err := h.app.Commands.CreateRestaurant.Handle(ctx, commands.CreateRestaurant{
		Name: request.Name,
		Address: restaurantapi.Address{
			Street1: request.Address.Street1,
			Street2: request.Address.Street2,
			City:    request.Address.City,
			State:   request.Address.State,
			Zip:     request.Address.Zip,
		},
		MenuItems: menuItems,
	})
	if err != nil {
		return nil, err
	}

	return &restaurantapi.CreateRestaurantResponse{RestaurantID: restaurantID}, nil
}

func (h RpcHandlers) GetRestaurant(ctx context.Context, request *restaurantapi.GetRestaurantRequest) (*restaurantapi.GetRestaurantResponse, error) {
	restaurant, err := h.app.Queries.GetRestaurant.Handle(ctx, queries.GetRestaurant{RestaurantID: request.RestaurantID})
	if err != nil {
		return nil, err
	}

	menuItems := make([]*restaurantapi.GetRestaurantResponseMenuItem, 0, len(restaurant.MenuItems))
	for _, item := range restaurant.MenuItems {
		menuItems = append(menuItems, &restaurantapi.GetRestaurantResponseMenuItem{
			ID:    item.ID,
			Name:  item.Name,
			Price: int64(item.Price),
		})
	}

	return &restaurantapi.GetRestaurantResponse{
		RestaurantID: restaurant.ID(),
		Name:         restaurant.Name,
		Address: &restaurantapi.GetRestaurantResponseAddress{
			Street1: restaurant.Address.Street1,
			Street2: restaurant.Address.Street2,
			City:    restaurant.Address.City,
			State:   restaurant.Address.State,
			Zip:     restaurant.Address.Zip,
		},
		Menu: &restaurantapi.GetRestaurantResponseMenu{MenuItems: menuItems},
	}, nil
}
