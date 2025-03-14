package webserver

import (
	"labs-two-service-b/config"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Route struct {
	Method  string
	Handler http.HandlerFunc
}

type WebServer struct {
	Router   *chi.Mux
	conf     *config.AppSettings
	Handlers map[string][]Route
}

func NewWebServer(conf *config.AppSettings) *WebServer {
	return &WebServer{
		Router:   chi.NewRouter(),
		Handlers: make(map[string][]Route),
		conf:     conf,
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
