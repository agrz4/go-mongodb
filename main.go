package main

import (
	"go-mongodb/usecase"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	// userservice instance
	userService := usecase.UserService{}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello from server"))
		})

		r.Post("/users", userService.CreateUser)
		r.Get("/users/{id}", userService.GetUserByID)
		r.Get("/users", userService.GetAllUsers)
		r.Put("/users/{id}", userService.UpdateUserAgeByID)
		r.Delete("/users/{id}", userService.DeleteUserByID)
		r.Delete("/users", userService.DeleteAllUsers)
	})

	slog.Info("server is starting at :8080")
	http.ListenAndServe(":8080", r)
}
