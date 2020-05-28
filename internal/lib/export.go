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
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

type ExportService struct {
	keycloak KeycloakService
	serving  ServingService
	influx   InfluxService
}

func NewExportService(keycloak KeycloakService, serving ServingService, influx InfluxService) *ExportService {
	return &ExportService{keycloak: keycloak, serving: serving, influx: influx}
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
			data, _ := es.influx.GetData(es.keycloak.GetAccessToken(), serving.Measurement)
			for _, i := range data.Results {
				fmt.Println(i.Series[0].Columns)
				path := "./files/"
				if _, err := os.Stat(path); os.IsNotExist(err) {
					os.Mkdir(path, 0755)
				}
				// Create a csv file
				f, err := os.Create(path + serving.Measurement + ".csv")
				if err != nil {
					fmt.Println(err)
				}
				defer f.Close()
				// Write Unmarshaled json data to CSV file
				w := csv.NewWriter(f)
				//Columns
				w.Write(i.Series[0].Columns[:])
				for _, d := range i.Series[0].GetValuesAsString() {
					w.Write(d)
				}
				w.Flush()
			}
		}
	}
}
