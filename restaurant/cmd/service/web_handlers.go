package main

import (
	"github.com/go-chi/render"
	"github.com/stackus/ftgogo/restaurant/internal/application/commands"
	"github.com/stackus/ftgogo/restaurant/internal/application/queries"
	"net/http"
	"serviceapis/restaurantapi"
	"shared-go/web"
)

// To regenerate the web server api use the following generate command
//go:generate oapi-codegen -config oapi-codegen.cfg.yaml openapi.yaml

type webHandlers struct{ app Application }

func newWebHandlers(app Application) webHandlers { return webHandlers{app: app} }

func (h webHandlers) CreateRestaurant(w http.ResponseWriter, r *http.Request) {
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

func (h webHandlers) GetRestaurant(w http.ResponseWriter, r *http.Request, restaurantID RestaurantID) {
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
