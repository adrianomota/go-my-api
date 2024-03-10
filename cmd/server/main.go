package main

import (
	"log"
	"net/http"

	"github.com/adrianomota/fullcycle/my-api/configs"
	_ "github.com/adrianomota/fullcycle/my-api/docs"
	"github.com/adrianomota/fullcycle/my-api/internal/entity"
	"github.com/adrianomota/fullcycle/my-api/internal/infra/database"
	"github.com/adrianomota/fullcycle/my-api/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title My Golang API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2

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
	r.Use(MyLogRequest)

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

	//swagger
	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))

	//server
	http.ListenAndServe(":8000", r)
}

func MyLogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("My Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
