package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	core "github.com/shoplineapp/captin/core"
	models "github.com/shoplineapp/captin/models"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.Info("* Starting captin (Press ENTER to quit)")

	pwd, _ := os.Getwd()
	path := os.Args[1:][0]
	absPath := filepath.Join(pwd, path)

	configMapper := models.NewConfigurationMapperFromPath(absPath)
	captin := core.NewCaptin(*configMapper)

	reader := bufio.NewReader(os.Stdin)
	for {
		log.Debug("> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		parsedInt, err := strconv.Atoi(text)

		if err == nil {
			for i := 0; i < parsedInt; i++ {
				captin.Execute(models.IncomingEvent{
					Key:        "product.update",
					Source:     "core",
					Payload:    map[string]interface{}{"field1": 1},
					TargetType: "Product",
					TargetId:   "product_id",
				})
			}
		} else {
			os.Exit(0)
		}
	}
}
