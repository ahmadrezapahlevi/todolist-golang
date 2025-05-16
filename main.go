package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var rnd *renderer.Render
var todoCollection *mongo.Collection
var ctx = context.Background()

const (
	dbName         = "demo_todo"
	collectionName = "todo"
	port           = ":9010"
)

type todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type todoModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Completed bool               `bson:"completed"`
	CreatedAt time.Time          `bson:"createdAt"`
}

func init() {
	rnd = renderer.New()

	uri := "isi dengan mongo uri"
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	checkErr(err)

	todoCollection = client.Database(dbName).Collection(collectionName)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := rnd.Template(w, http.StatusOK, []string{"static/home.tpl"}, nil)
	checkErr(err)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var t todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusBadRequest, err)
		return
	}

	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "Title is required"})
		return
	}

	newTodo := todoModel{
		ID:        primitive.NewObjectID(),
		Title:     t.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	_, err := todoCollection.InsertOne(ctx, newTodo)
	if err != nil {
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{"message": "Insert failed", "error": err})
		return
	}

	rnd.JSON(w, http.StatusCreated, renderer.M{"message": "Todo created", "todo_id": newTodo.ID.Hex()})
}

func fetchTodos(w http.ResponseWriter, r *http.Request) {
	cursor, err := todoCollection.Find(ctx, bson.M{})
	if err != nil {
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{"message": "Fetch failed", "error": err})
		return
	}
	defer cursor.Close(ctx)

	var todos []todo
	for cursor.Next(ctx) {
		var t todoModel
		if err := cursor.Decode(&t); err != nil {
			continue
		}
		todos = append(todos, todo{
			ID:        t.ID.Hex(),
			Title:     t.Title,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt,
		})
	}

	rnd.JSON(w, http.StatusOK, renderer.M{"data": todos})
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimSpace(chi.URLParam(r, "id"))
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "Invalid ID"})
		return
	}

	var t todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusBadRequest, err)
		return
	}

	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "Title is required"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"title":     t.Title,
			"completed": t.Completed,
		},
	}

	_, err = todoCollection.UpdateByID(ctx, id, update)
	if err != nil {
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{"message": "Update failed", "error": err})
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{"message": "Todo updated"})
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimSpace(chi.URLParam(r, "id"))
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "Invalid ID"})
		return
	}

	_, err = todoCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{"message": "Delete failed", "error": err})
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{"message": "Todo deleted"})
}

func main() {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Mount("/todo", todoHandlers())

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Listening on port", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-stopChan
	log.Println("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctxShutdown)
	log.Println("Server gracefully stopped!")
}

func todoHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", fetchTodos)
		r.Post("/", createTodo)
		r.Put("/{id}", updateTodo)
		r.Delete("/{id}", deleteTodo)
	})
	return rg
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}