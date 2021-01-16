package main

import (
	"net/http"

	"github.com/go-chi/chi"
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
