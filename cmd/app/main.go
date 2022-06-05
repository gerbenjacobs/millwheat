package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/gerbenjacobs/millwheat/game/data"
	"github.com/gerbenjacobs/millwheat/handler"
	"github.com/gerbenjacobs/millwheat/services"
	"github.com/gerbenjacobs/millwheat/storage"
)

func main() {
	// flags
	cfgPath := flag.String("config-path", ".", "path to config file")
	flag.Parse()

	// handle shutdown signals
	shutdown := make(chan os.Signal, 3)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// set output logging (specifically for windows)
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(colorable.NewColorableStdout())
	log.SetLevel(log.DebugLevel)

	// load configuration
	var c Configuration
	viper.SetConfigName("config")
	viper.AddConfigPath(*cfgPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %s", err)
	}
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode into struct: %v", err)
	}

	// set up and check database
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?parseTime=true", c.DB.User, c.DB.Password, c.DB.Database))
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	db.SetConnMaxLifetime(time.Second)

	// create repositories and services
	auth := services.NewAuth([]byte(c.Svc.SecretToken))
	userSvc, err := services.NewUserSvc(storage.NewUserRepository(db), auth)
	if err != nil {
		log.Fatalf("failed to start user service: %v", err)
	}

	townSvc := services.NewTownSvc(storage.NewTownRepository(db))
	prodSvc := services.NewProductionSvc(storage.NewProductionRepository(db))
	battleSvc := services.NewBattleSvc(storage.NewBattleRepo(db))

	gameSvc := services.NewGameSvc(townSvc, prodSvc, battleSvc, data.Items, data.Buildings)

	// set up the route handler and server
	app := handler.New(handler.Dependencies{
		Auth:    auth,
		UserSvc: userSvc,

		GameSvc:       gameSvc,
		TownSvc:       townSvc,
		ProductionSvc: prodSvc,
		BattleSvc:     battleSvc,

		Items:     data.Items,
		Buildings: data.Buildings,
	})
	srv := &http.Server{
		Addr:         c.Svc.Address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      app,
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	app.Tick(ctx)

	// start running the server
	go func() {
		log.Print("Server started on " + srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to listen: %v", err)
		}
	}()

	// wait for shutdown signals
	<-shutdown
	cancelFunc() // global ctx for services
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Print("Server stopped successfully")
}

type Configuration struct {
	Svc struct {
		Name        string
		Version     string
		Env         string
		Address     string
		SecretToken string
	}
	DB struct {
		User     string
		Password string
		Database string
	}
}
