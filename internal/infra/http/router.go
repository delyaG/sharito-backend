package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/rs/cors"
	"net/http"
)

func (a *adapter) newRouter() (http.Handler, error) {
	r := chi.NewRouter()

	// Set default middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Logger)

	c := cors.New(cors.Options{
		AllowedOrigins:   a.config.AllowedOrigins,
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})
	r.Use(c.Handler)

	r.Route("/api", func(r chi.Router) {
		r.Route("/test", func(r chi.Router) {
			r.Get("/hello", a.wrap(a.sayHello))
		})

		r.Route("/v1", func(r chi.Router) {
			r.Route("/auth", func(r chi.Router) {
				r.Post("/register", a.wrap(a.register))
				r.Post("/login", a.wrap(a.login))
			})

			r.Group(func(r chi.Router) {
				r.Route("/user", func(r chi.Router) {
					r.Use(jwtauth.Verifier(a.jwtAuth))
					r.Use(a.JWTAuthMiddleware())
					r.Get("/", a.wrap(a.getUser))
				})

				r.Route("/product", func(r chi.Router) {
					r.Get("/", a.wrap(a.getProducts))
					r.Group(func(r chi.Router) {
						r.Use(jwtauth.Verifier(a.jwtAuth))
						r.Use(a.JWTAuthMiddleware())
						r.Post("/", a.wrap(a.addProduct))
						r.Get("/{product_id}", a.wrap(a.getProduct))
					})
				})

				r.Route("/order", func(r chi.Router) {
					r.Use(jwtauth.Verifier(a.jwtAuth))
					r.Use(a.JWTAuthMiddleware())
					r.Post("/{product_id}", a.wrap(a.rentProduct))
					r.Get("/", a.wrap(a.getOrders))
				})
			})
		})
	})

	return r, nil
}
