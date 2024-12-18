package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{

	Route{
		"Status",
		"GET",
		"/status",
		getStatus,
	},

	Route{
		"DBlist",
		"GET",
		"/db_list",
		getDBlist,
	},

	Route{
		"MLlist",
		"GET",
		"/ml_list",
		getMLlist,
	},

	Route{
		"CreateMapping",
		"POST",
		"/create_mapping",
		createMapping,
	},

	Route{
		"PerformRequest",
		"POST",
		"/perform_request",
		performRequest,
	},
}
