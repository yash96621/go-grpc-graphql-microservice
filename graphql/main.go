package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/handler"
	"github.com/kelseyhightower/envconfig"
)

type AppConifg struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL string `envconfig:"CATALOG_ACCOUNT_SERVICE_URL_URL"`
	OrderUrl   string `envconfig:"ORDER_ACCOUNT_SERVICE_URL"`
}

func Main() {
	var cfg AppConifg
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	s, err := NewGraphQlServer(cfg.AccountURL, cfg.CatalogURL, cfg.OrderUrl)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/graphql", handler.New(s.ToExecutableSchema()))
	http.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
