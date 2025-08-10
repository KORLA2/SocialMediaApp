package main

import (
	"fmt"
	"log"
	"time"

	env "github.com/KORLA2/SocialMedia/internal"
	"github.com/KORLA2/SocialMedia/internal/db"
	"github.com/KORLA2/SocialMedia/internal/store"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")
	cfg := config{
		addr: env.GetString("ADDR", ":8008"),

		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://user:password@localhost:5433/social?sslmode=disable"),
			maxOpenConns: env.GetInt("MAX_CONNS", 25),
			maxIdleConns: env.GetInt("MAX_IDLE_CONNS", 25),
			maxIdleTime:  env.GetString("MAX_IDLE_TIME", (20 * time.Minute).String()),
		},
	}

	fmt.Println("Listening on port", cfg.addr)

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)

	if err != nil {
		log.Fatal(err)
	}
	storage := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  storage,
	}
	mux := app.mount()
	log.Fatal(app.Run(mux))

}
