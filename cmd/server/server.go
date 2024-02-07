package main

import (
	"database/sql"
	"github.com/carloseduribeiro/working-with-graphql-go/internal/database"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/carloseduribeiro/working-with-graphql-go/graph"

	_ "github.com/mattn/go-sqlite3"
)

const defaultPort = "8080"

func main() {
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec(`create table if not exists category(id string, name string, description string)`); err != nil {
		panic(err)
	}
	if _, err := db.Exec(`create table if not exists course(id string, name string, description string, category_id string)`); err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	categoryDB := database.NewCategory(db)
	courseDB := database.NewCourse(db)
	resolver := &graph.Resolver{
		CategoryDB: categoryDB,
		CourseDB:   courseDB,
	}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
