package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"github.com/yash96621/go-grpc-graphql-microservice/order"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL" `
	AccountURL  string `envconfig:"ACCOUNT_SERVICE_URL" `
	CatalogURL  string `envconfig:"CATALOG_SERVICE_URL" `
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	var r order.Repository

	retry.ForeverSleep(2*time.Second, func(_ int) error {
		var err error
		r, err = order.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println("could not connect to db:", err)
			return err
		}
		return nil
	})
	defer r.Close()

	log.Println("listening on port 8080......")

	s := order.NewService(r)
	log.Fatal(order.ListenGRPC(s, cfg.AccountURL, cfg.CatalogURL, 8080))

}
