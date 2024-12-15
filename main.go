package main

import (
	"context"
	"fmt"
	"go-mongodb/repository/mongodb"
	"go-mongodb/usecase"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error while reading .env")
	}
	slog.Info("env loaded successfully.")
}

func mongoConnection() *mongo.Client {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGO_URI")).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	// Send a ping to confirm a successful connection
	if err := client.Database(os.Getenv("MONGO_DBNAME")).RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client
}

func main() {
	mongoClient := mongoConnection()
	defer mongoClient.Disconnect(context.Background())

	// mongo connection to coleection
	collection := mongoClient.Database(os.Getenv("MONGO_DBNAME")).Collection(os.Getenv("MONGO_COLLECTION_NAME"))

	// userservice instance
	userService := usecase.UserService{
		DBClient: mongodb.MongoClient{
			Client: *collection,
		},
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.SetHeader("CONTENT-TYPE", "application/json"))

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
