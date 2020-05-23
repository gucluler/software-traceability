package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	auth "traceability/auth"
	data "traceability/data"
	"traceability/database"
	handlers "traceability/handlers"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const address = ":8001"
const (
	dbHost = "localhost"
	dbPort = 27017
	dbName = "pqsltestdb"
)

func main() {

	connectDB()

	sm := mux.NewRouter()
	l := log.New(os.Stdout, "chat-api", log.LstdFlags)
	v := data.NewValidation()
	uh := handlers.NewUsers(l, v)
	getUs := sm.Methods(http.MethodGet).Subrouter()
	getUs.HandleFunc("/users", uh.ListAll)
	postUs := sm.Methods(http.MethodPost).Subrouter()
	postUs.HandleFunc("/user", uh.CreateUser)
	postUs.Use(auth.Middleware)
	postUs.Use(uh.MiddlewareValidateUser)

	loginUser := sm.Methods(http.MethodPost).Subrouter()
	loginUser.HandleFunc("/login", uh.LoginUser)
	loginUser.Use(uh.MiddlewareValidateAuth)

	s := http.Server{
		Addr:         address,           // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}
	go func() {
		fmt.Println("server is starting at", address)
		err := s.ListenAndServe()

		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	if cancel != nil {
		fmt.Println("cancel != nil")
	}
	s.Shutdown(ctx)
}

func connectDB() {
	dbURI := fmt.Sprintf("mongodb://%s:%d", dbHost, dbPort)
	clientOptions := options.Client().ApplyURI(dbURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(dbName)
	database.DBCon = client
	database.DB = db

	fmt.Println("Connected to MongoDB!")
}
