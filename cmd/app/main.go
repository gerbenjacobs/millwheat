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

	// set up the route handler and server
	app := handler.New(handler.Dependencies{
		Auth:    auth,
		UserSvc: userSvc,

		TownSvc: townSvc,

		Items:     tempItems,
		Buildings: tempBuildings,
	})
	srv := &http.Server{
		Addr:         c.Svc.Address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      app,
	}

	// start running the server
	go func() {
		log.Print("Server started on " + srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to listen: %v", err)
		}
	}()

	// wait for shutdown signals
	<-shutdown
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
					Type:         game.BuildingFarm,
					CurrentLevel: 2,
				},
				{
					Type:         game.BuildingMill,
					CurrentLevel: 1,
				},
				{
					Type:         game.BuildingBakery,
					CurrentLevel: 1,
				},
			},
			CreatedAt: time.Now().Add(-5 * time.Minute),
			UpdatedAt: time.Now(),
		},
	}

	tempBuildings := game.Buildings{
		game.BuildingFarm: {
			Name:        "Farm",
			Description: "Grows wheat in the fields.",
			Image:       "https://www.knightsandmerchants.net/application/files/7515/6823/6441/farm.png",
			Consumes:    nil,
			Produces:    []game.ItemID{"wheat"},
			Mechanics: []game.BuildingMechanic{
				{
					Name: "Wheat per hour",
					Levels: map[int]int{
						1: 1,
						2: 2,
					},
				},
			},
			BuildCosts: map[int]game.BuildCost{
				1: {1, 3},
				2: {2, 6},
				3: {3, 15},
				4: {5, 50},
				5: {7, 75},
			},
		},
		game.BuildingMill: {
			Name:        "Mill",
			Description: "Mills wheat into bags of flour.",
			Image:       "https://www.knightsandmerchants.net/application/files/9415/6823/6446/mill.png",
			Consumes:    []game.ItemID{"wheat"},
			Produces:    []game.ItemID{"flour"},
			Mechanics: []game.BuildingMechanic{
				{
					Name: "Flour per wheat",
					Levels: map[int]int{
						1: 1,
					},
				},
				{
					Name: "Flour per hour",
					Levels: map[int]int{
						1: 1,
					},
				},
			},
		},
		game.BuildingBakery: {
			Name:        "Bakery",
			Description: "Bakes bread for the soldiers using flour from the mill.",
			Image:       "https://www.knightsandmerchants.net/application/files/1215/6823/6439/bakery.png",
			Consumes:    []game.ItemID{"flour"},
			Produces:    []game.ItemID{"bread"},
			Mechanics: []game.BuildingMechanic{
				{
					Name: "Bread per flour",
					Levels: map[int]int{
						1: 1,
					},
				},
				{
					Name: "Bread per hour",
					Levels: map[int]int{
						1: 1,
					},
				},
			},
		},
	}

	return tempTowns, data.Items, tempBuildings
}
