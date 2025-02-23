package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"github.com/yash96621/go-grpc-graphql-microservice/account"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	var r account.Repository

	retry.ForeverSleep(2*time.Second, func(_ int) error {
		var err error
		r, err = account.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println("could not connect to db:", err)
			return err
		}
		return nil
	})
	defer r.Close()

	log.Println("listening on port 8080......")

	s := account.NewService(r)
	log.Fatal(account.ListenGRPC(s, "8080"))

}
