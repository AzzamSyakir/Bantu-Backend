package route

import (
	"bantu-backend/src/gateway-service/delivery/http"

	"github.com/gorilla/mux"
)

type RootRoute struct {
	Router       *mux.Router
	GatewayRoute *GatewayRoute
}

func NewRootRoute(
	router *mux.Router,
	authRoute *GatewayRoute,

) *RootRoute {
	rootRoute := &RootRoute{
		Router:       router,
		GatewayRoute: authRoute,
	}
	return rootRoute
}

func (rootRoute *RootRoute) Register() {
	rootRoute.GatewayRoute.Register()
}

type GatewayRoute struct {
	Router            *mux.Router
	GatewayController *http.GatewayController
}

func NewGatewayRoute(router *mux.Router, GatewayController *http.GatewayController) *GatewayRoute {
	GatewayRoute := &GatewayRoute{
		Router:            router.PathPrefix("/auths").Subrouter(),
		GatewayController: GatewayController,
	}
	return GatewayRoute
}

func (authRoute *GatewayRoute) Register() {
	authRoute.Router.HandleFunc("/register", authRoute.GatewayController.Register).Methods("POST")

	authRoute.Router.HandleFunc("/login", authRoute.GatewayController.Login).Methods("POST")
	authRoute.Router.HandleFunc("/access-token", authRoute.GatewayController.GetNewAccessToken).Methods("POST")
	authRoute.Router.HandleFunc("/logout", authRoute.GatewayController.Logout).Methods("POST")
}
