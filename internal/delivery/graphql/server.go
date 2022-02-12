package graphql

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rasulov-emirlan/pukbot/internal/delivery/graphql/graph"
	"github.com/rasulov-emirlan/pukbot/internal/delivery/graphql/graph/generated"
	"github.com/rasulov-emirlan/pukbot/internal/puk"
)

func Start(port string, p puk.Service) {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		PukService: p,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func NewHandler(p puk.Service) (*handler.Server, http.HandlerFunc) {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		PukService: p,
	}}))
	play := playground.Handler("GraphQL playground", "/query")
	return srv, play
}
