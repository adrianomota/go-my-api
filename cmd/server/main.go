package main

import (
	"net/http"

	"github.com/adrianomota/fullcycle/my-api/configs"
	"github.com/adrianomota/fullcycle/my-api/internal/entity"
	"github.com/adrianomota/fullcycle/my-api/internal/infra/database"
	"github.com/adrianomota/fullcycle/my-api/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("test_dev.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})
	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB, config.TokenAuhth, config.JWTExpiresIn)

	r := chi.NewRouter()

	//middlewares
	r.Use(middleware.Logger)

	//products
	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuhth))
		r.Use(jwtauth.Authenticator)
		r.Get("/", productHandler.All)
		r.Get("/{id}", productHandler.Get)
		r.Post("/", productHandler.Create)
		r.Put("/{id}", productHandler.Update)
		r.Delete("/{id}", productHandler.Delete)
	})

	//users
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.Create)
		r.Get("/{id}", userHandler.Get)
		r.Post("/token", userHandler.GetJWT)
	})

	//server
	http.ListenAndServe(":8000", r)
}
