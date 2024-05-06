package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	core "github.com/shoplineapp/captin/v2/core"
	models "github.com/shoplineapp/captin/v2/models"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.Info("* Starting captin (Press ctrl+c to quit)")

	pwd, _ := os.Getwd()
	path := os.Args[1:][0]
	absPath := filepath.Join(pwd, path)

	configMapper := models.NewConfigurationMapperFromPath(absPath)
	captin := core.NewCaptin(*configMapper)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	enabled := true

	go func() {
		for enabled && captin.IsRunning() != true {
			captin.Execute(context.Background(), models.IncomingEvent{
				Key:        "product.update",
				Source:     "core",
				Payload:    map[string]interface{}{"field1": 1},
				TargetType: "Product",
				TargetId:   "product_id",
			})
			time.Sleep(1 * time.Second)
		}
	}()

	<-quit
	enabled = false
	log.Println("Gracefully shutting down...")
	for {
		time.Sleep(1 * time.Second)
		if captin.IsRunning() != true {
			log.Println("Tasks released")
			break
		}
	}

	log.Println("Bye!")
	os.Exit(0)
}
