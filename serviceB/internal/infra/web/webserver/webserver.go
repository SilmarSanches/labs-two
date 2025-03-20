package webserver

import (
	"labs-two-service-b/config"
	"labs-two-service-b/internal/infra/tracing"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Route struct {
	Method  string
	Handler http.HandlerFunc
}

type WebServer struct {
	Router          *chi.Mux
	conf            *config.AppSettings
	Handlers        map[string][]Route
	TracingProvider *tracing.TracingProvider
}

func NewWebServer(conf *config.AppSettings, tracingProvider *tracing.TracingProvider) *WebServer {
	return &WebServer{
		Router:          chi.NewRouter(),
		Handlers:        make(map[string][]Route),
		conf:            conf,
		TracingProvider: tracingProvider,
	}
}

func (w *WebServer) AddHandler(method, path string, handler http.HandlerFunc) {
	w.Handlers[path] = append(w.Handlers[path], Route{Method: method, Handler: handler})
}

func (w *WebServer) Start() {
	w.Router.Use(middleware.Logger)

	for path, routes := range w.Handlers {
		for _, route := range routes {
			w.Router.MethodFunc(route.Method, path, route.Handler)
		}
	}

	http.ListenAndServe(":"+w.conf.Port, w.Router)
}
