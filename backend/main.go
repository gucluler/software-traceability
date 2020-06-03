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
	"traceability/data"
	"traceability/database"
	archViewHandlers "traceability/handlers/archview"

	componentHandlers "traceability/handlers/archviewcomponents"
	projectHandlers "traceability/handlers/project"
	userHandlers "traceability/handlers/user"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const address = ":8080"
const (
	dbHost = "localhost"
	dbPort = 27017
	dbName = "traceability"
)

func main() {

	connectDB()

	sm := mux.NewRouter()
	l := log.New(os.Stdout, "traceability-api", log.LstdFlags)
	v := data.NewValidation()
	uh := userHandlers.NewUsers(l, v)
	ph := projectHandlers.NewProjects(l, v)
	ah := archViewHandlers.NewArchViews(l, v)
	ch := componentHandlers.NewArchViewComponents(l, v)
	setUserEndpoints(sm, uh)
	setProjectEndpoints(sm, ph)
	setArchViewEndpoints(sm, ah)
	setArchViewComponentEndpoints(sm, ch)

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

func setUserEndpoints(sm *mux.Router, uh *userHandlers.Users) {
	getUserList := sm.Methods(http.MethodGet).Subrouter()
	getUserList.HandleFunc("/users", uh.ListAll)

	getUser := sm.Methods(http.MethodGet).Subrouter()
	getUser.HandleFunc("/users/{id}", uh.GetUser)
	getUser.Use(auth.Middleware)

	postUs := sm.Methods(http.MethodPost).Subrouter()
	postUs.HandleFunc("/users", uh.CreateUser)
	postUs.Use(uh.MiddlewareValidateUser)

	loginUser := sm.Methods(http.MethodPost).Subrouter()
	loginUser.HandleFunc("/login", uh.LoginUser)
	loginUser.Use(uh.MiddlewareValidateAuth)
}

func setProjectEndpoints(sm *mux.Router, ph *projectHandlers.Projects) {
	getProj := sm.Methods(http.MethodGet).Subrouter()
	getProj.HandleFunc("/projects/{id}", ph.GetProject)
	getProj.Use(auth.Middleware)
	getProj.Use(ph.MiddlewareValidatePermission)

	getProjList := sm.Methods(http.MethodGet).Subrouter()
	getProjList.HandleFunc("/projects", ph.ListAll)
	getProjList.Use(auth.Middleware)

	postProj := sm.Methods(http.MethodPost).Subrouter()
	postProj.HandleFunc("/projects", ph.CreateProject)
	postProj.Use(auth.Middleware)
	postProj.Use(ph.MiddlewareValidateProject)
}

func setArchViewEndpoints(sm *mux.Router, ah *archViewHandlers.ArchViews) {
	getArchView := sm.Methods(http.MethodGet).Subrouter()
	getArchView.HandleFunc("/projects/{projectID}/views/{id}", ah.GetArchView)
	getArchView.Use(auth.Middleware)

	postArchView := sm.Methods(http.MethodPost).Subrouter()
	postArchView.HandleFunc("/projects/{projectID}/views", ah.CreateArchView)
	postArchView.Use(auth.Middleware)
	postArchView.Use(ah.MiddlewareValidateArchView)
}

func setArchViewComponentEndpoints(sm *mux.Router, ch *componentHandlers.ArchViewComponents) {
	getComp := sm.Methods(http.MethodGet).Subrouter()
	getComp.HandleFunc("/projects/{projectID}/views/{viewID}/components/{id}", ch.GetArchViewComponent)
	getComp.Use(auth.Middleware)

	postComponent := sm.Methods(http.MethodPost).Subrouter()
	postComponent.HandleFunc("/projects/{projectID}/views/{viewID}/components", ch.AddArchViewComponent)
	postComponent.Use(auth.Middleware)
	postComponent.Use(ch.MiddlewareValidateArchViewComponent)
}
