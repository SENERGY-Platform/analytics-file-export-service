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
	keycloak := *lib.NewKeycloakService(
		lib.GetEnv("KEYCLOAK_ADDRESS", "http://test"),
		lib.GetEnv("KEYCLOAK_CLIENT_ID", "test"),
		lib.GetEnv("KEYCLOAK_CLIENT_SECRET", "test"),
		lib.GetEnv("KEYCLOAK_REALM", "test"),
		lib.GetEnv("KEYCLOAK_USER", "test"),
		lib.GetEnv("KEYCLOAK_PW", "test"),
	)
	es := lib.NewExportService(keycloak, serving, influx)
	es.StartExportService()
}
