package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"

	"github.com/gerbenjacobs/millwheat/game"
	"github.com/gerbenjacobs/millwheat/game/data"
	"github.com/gerbenjacobs/millwheat/handler"
	"github.com/gerbenjacobs/millwheat/services"
	"github.com/gerbenjacobs/millwheat/storage"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
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
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %s", err)
	}
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode into struct: %v", err)
	}

	// load game data
	tempTowns, tempItems, tempBuildings := tempGameData()

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

	townSvc := services.NewTownSvc(storage.NewTownRepository(tempTowns))
	prodSvc := services.NewProductionSvc(storage.NewProductionRepository())

	// TEMP
	fakeJobs(prodSvc)

	// set up the route handler and server
	app := handler.New(handler.Dependencies{
		Auth:    auth,
		UserSvc: userSvc,

		TownSvc:       townSvc,
		ProductionSvc: prodSvc,

		Items:     tempItems,
		Buildings: tempBuildings,
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

func tempGameData() (game.Towns, game.Items, game.Buildings) {
	tempTowns := map[uuid.UUID]*game.Town{
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"): {
			ID:    uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			Name:  "Northbrook",
			Owner: uuid.MustParse("273d94bb1cf7408da4c85feda0eeff75"),
			Buildings: []game.TownBuilding{
				{
					ID:           uuid.MustParse("f8b93eab-b11d-44ca-bba6-162c60e4762e"),
					Type:         game.BuildingFarm,
					CurrentLevel: 2,
				},
				{
					ID:           uuid.MustParse("578ad258-e913-4831-998a-2983dd4920ed"),
					Type:         game.BuildingMill,
					CurrentLevel: 5,
				},
				//{
				//	ID:           uuid.MustParse("c5e9a879-b1a8-4670-9aab-73b98eb21f6b"),
				//	Type:         game.BuildingBakery,
				//	CurrentLevel: 1,
				//},
				//{
				//	ID:           uuid.MustParse("578ad257-b1a7-4671-bba5-73b98eb21f6c"),
				//	Type:         game.BuildingButcher,
				//	CurrentLevel: 1,
				//},
				{
					ID:           uuid.MustParse("1afb2cf0-6671-4d54-9824-c15940073ce2"),
					Type:         game.BuildingWeaponSmith,
					CurrentLevel: 1,
				},
			},
			CreatedAt: time.Now().Add(-5 * time.Minute),
			UpdatedAt: time.Now(),
		},
	}

	return tempTowns, data.Items, data.Buildings
}

func fakeJobs(prodSvc services.ProductionService) {
	ctx := context.WithValue(context.Background(), services.CtxKeyTownID, uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"))
	_ = prodSvc.CreateJob(ctx, &game.InputJob{
		Type: game.JobTypeBuilding,
		BuildingJob: &game.BuildingJob{
			ID:    uuid.New(),
			Type:  game.BuildingBakery,
			Level: 1,
		},
		Hours: 5 * time.Hour,
	})
	_ = prodSvc.CreateJob(ctx, &game.InputJob{
		Type: game.JobTypeBuilding,
		BuildingJob: &game.BuildingJob{
			ID:    uuid.MustParse("f8b93eab-b11d-44ca-bba6-162c60e4762e"),
			Type:  game.BuildingFarm,
			Level: 3,
		},
		Hours: 2 * time.Hour,
	})
}
