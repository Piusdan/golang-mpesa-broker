package paymentsgateway
import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range gatewayRoutes {
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

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

var gatewayRoutes = routes{
	route{
		"Index",
		"GET",
		"/",
		IndexEndpoint,
	}, route{
		"Disburse",
		"POST",
		"/disburse",
		DisburseEndpoint,
	},
}
