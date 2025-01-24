package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
	_ "wallet/docs"
)

type Handler struct {
	DB        *gorm.DB
	tokenAuth *jwtauth.JWTAuth
}

func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()
	h.initToken()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	router.Use(middleware.Logger)
	router.Use(middleware.Heartbeat("/"))
	router.Get("/api/v1/swagger/*", httpSwagger.Handler())
	router.Route("/api", func(api chi.Router) {
		api.Route("/v1", func(v1 chi.Router) {
			v1.Route("/user", func(user chi.Router) {
				user.Group(func(user chi.Router) {
					user.Route("/register", func(register chi.Router) {
						register.Post("/", h.Register)
					})
					user.Route("/login", func(login chi.Router) {
						login.Post("/", h.Login)
					})

				})
				user.Group(func(user chi.Router) {
					user.Use(jwtauth.Verifier(h.tokenAuth))
					user.Use(jwtauth.Authenticator(h.tokenAuth))
					user.Route("/balance", func(balance chi.Router) {
						balance.Get("/", h.GetBalance)
						balance.Post("/deposit", h.Deposit)
						balance.Post("/withdraw", h.Withdraw)
					})
				})
			})
			v1.Route("/exchange", func(exchange chi.Router) {
				exchange.Use(jwtauth.Verifier(h.tokenAuth))
				exchange.Use(jwtauth.Authenticator(h.tokenAuth))
				exchange.Group(func(exchange chi.Router) {
					exchange.Post("/", h.ExchangeCurrencies)
					exchange.Route("/rates", func(rates chi.Router) {
						rates.Use(jwtauth.Verifier(h.tokenAuth))
						rates.Use(jwtauth.Authenticator(h.tokenAuth))
						rates.Get("/", h.GetExchangeRates)
					})
				})
			})
		})
	})
	return router
}
