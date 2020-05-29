/*
 * Copyright 2020 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	cron "github.com/robfig/cron/v3"

	"github.com/SENERGY-Platform/analytics-file-export-service/v2/internal/lib"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	serving := *lib.NewServingService(
		lib.GetEnv("SERVING_API_ENDPOINT", ""),
	)
	influx := *lib.NewInfluxService(
		lib.GetEnv("INFLUX_API_URL", ""),
	)
	cloud := *lib.NewCloudService(
		lib.GetEnv("NEXTCLOUD_HOST", ""),
		lib.GetEnv("NEXTCLOUD_USER", ""),
		lib.GetEnv("NEXTCLOUD_PW", ""),
	)
	keycloak := *lib.NewKeycloakService(
		lib.GetEnv("KEYCLOAK_ADDRESS", "http://test"),
		lib.GetEnv("KEYCLOAK_CLIENT_ID", "test"),
		lib.GetEnv("KEYCLOAK_CLIENT_SECRET", "test"),
		lib.GetEnv("KEYCLOAK_REALM", "test"),
		lib.GetEnv("KEYCLOAK_USER", "test"),
		lib.GetEnv("KEYCLOAK_PW", "test"),
	)

	if lib.GetEnv("CRON_SCHEDULE", "* * * * *") == "false" {
		es := lib.NewExportService(keycloak, serving, influx, cloud, lib.GetEnv("CLOUD_PATH", ""))
		es.StartExportService()
		os.Exit(0)
	} else {
		c := cron.New()
		_, err = c.AddFunc(lib.GetEnv("CRON_SCHEDULE", "* * * * *"), func() {
			log.Println("Start backup")
			es := lib.NewExportService(keycloak, serving, influx, cloud, lib.GetEnv("CLOUD_PATH", ""))
			es.StartExportService()
		})
		if err != nil {
			log.Fatal("Error starting job: " + err.Error())
		}
		c.Start()
	}
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	sig := <-shutdown
	log.Println("received shutdown signal", sig)
}
