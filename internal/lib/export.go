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

package lib

import (
	"fmt"
	"log"
)

type ExportService struct {
	keycloak KeycloakService
	serving  ServingService
}

func NewExportService(keycloak KeycloakService, serving ServingService) *ExportService {
	return &ExportService{keycloak: keycloak, serving: serving}
}

func (es *ExportService) StartExportService() {
	es.keycloak.Login()
	defer es.keycloak.Logout()
	user, err := es.keycloak.GetUserInfo()
	if err != nil {
		log.Fatal("GetUserInfo failed:" + err.Error())
	}
	if user != nil {
		servings, err := es.serving.GetServingServices(*user.Sub, es.keycloak.GetAccessToken())
		if err != nil {
			log.Fatal("GetServingServices failed: " + err.Error())
		}
		for _, serving := range servings {
			fmt.Println(serving)
		}
	}
}
