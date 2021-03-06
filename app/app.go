package app

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/danilomarques1/findmypetapi/handler"
	"github.com/danilomarques1/findmypetapi/lib"
	"github.com/danilomarques1/findmypetapi/repository"
	"github.com/danilomarques1/findmypetapi/service"
	"github.com/danilomarques1/findmypetapi/util"
	validator "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/streadway/amqp"
)

type App struct {
	Router    *mux.Router
	validator *validator.Validate
	DB        *sql.DB
}

func (app *App) Init(sqlFileName, dbstring string) {
	var err error

	app.Router = mux.NewRouter()
	app.Router.Use(contentTypeMiddleware)

	app.validator = validator.New()
	app.DB, err = sql.Open("postgres", dbstring)
	if err != nil {
		log.Fatalf("Error opening database connection %v\n", err)
	}

	sqlFile, err := os.ReadFile(sqlFileName)
	if err != nil {
		log.Fatalf("Error loading databases %v\n", err)
	}
	if _, err := app.DB.Exec(string(sqlFile)); err != nil {
		log.Fatalf("Error creating tables %v\n", err)
	}

	// adding handlers
	userRepo := repository.NewUserRepositorySql(app.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, app.validator)

	connection, err := amqp.Dial(os.Getenv("RABBIT_URL"))
	if err != nil {
		log.Fatalf("Error connecting to rabbit mq %v\n", err)
	}
	producer, err := lib.NewAmqpProducer(connection)
	if err != nil {
		log.Fatalf("Error setting up producer %v\n", err)
	}

	postRepo := repository.NewPostRepositorySql(app.DB)
	postService := service.NewPostService(postRepo, producer)
	postHandler := handler.NewPostHandler(postService, app.validator)

	commentRepo := repository.NewCommentRepositorySql(app.DB)
	commentService := service.NewCommentService(commentRepo, producer)
	commentHandler := handler.NewCommentHandler(commentService, app.validator)

	app.Router.HandleFunc("/user",
		userHandler.Save).Methods(http.MethodPost)
	app.Router.HandleFunc("/session",
		userHandler.CreateSession).Methods(http.MethodPost)
	app.Router.HandleFunc("/session/refresh",
		userHandler.RefreshSession).Methods(http.MethodPut)
	app.Router.Handle("/user",
		util.AuthorizationMiddleware(http.HandlerFunc(
			userHandler.UpdateUser))).Methods(http.MethodPut)

	// Post handelrs
	app.Router.Handle("/post",
		util.AuthorizationMiddleware(http.HandlerFunc(postHandler.CreatePost))).Methods(http.MethodPost)
	app.Router.Handle("/post/user",
		util.AuthorizationMiddleware(http.HandlerFunc(
			postHandler.FindPostsByAuthor))).Methods(http.MethodGet)
	app.Router.Handle("/post",
		util.AuthorizationMiddleware(http.HandlerFunc(
			postHandler.GetAll))).Methods(http.MethodGet)
	app.Router.Handle("/post/{post_id}",
		util.AuthorizationMiddleware(http.HandlerFunc(
			postHandler.GetOne))).Methods(http.MethodGet)
	app.Router.Handle("/post/{post_id}",
		util.AuthorizationMiddleware(http.HandlerFunc(
			postHandler.Update))).Methods(http.MethodPut)

	// comment handlers
	app.Router.Handle("/comment/{post_id}",
		util.AuthorizationMiddleware(http.HandlerFunc(
			commentHandler.CreateComment))).Methods(http.MethodPost)
	app.Router.Handle("/comment/{post_id}",
		util.AuthorizationMiddleware(http.HandlerFunc(
			commentHandler.FindAll))).Methods(http.MethodGet)
}

func (app *App) Listen() {
	port := os.Getenv("PORT")
	handler := cors.AllowAll().Handler(app.Router)
	server := http.Server{
		Handler: handler,
		Addr:    ":" + port,
	}

	log.Printf("Server starting...")
	log.Fatal(server.ListenAndServe())
}

func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
