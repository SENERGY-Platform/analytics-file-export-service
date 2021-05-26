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
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type InfluxService struct {
	url string
}

func NewInfluxService(url string) *InfluxService {
	return &InfluxService{url: url}
}

func (i *InfluxService) GetData(accessToken string, servings []ServingInstance, startInit time.Time, endInit time.Time) (influxResponse [][]interface{}, err error) {
	body := make([]QueriesRequestElement, len(servings))
	header := []interface{}{"time"}
	for i := range servings {
		var e QueriesRequestElement
		e.Measurement = servings[i].Measurement
		cols := make([]QueriesRequestElementColumn, len(servings[i].Values))
		for j := range servings[i].Values {
			header = append(header, servings[i].ID.String()+"."+servings[i].Name+"."+servings[i].Values[j].Name)
			cols[j] = QueriesRequestElementColumn{
				Name: servings[i].Values[j].Name,
			}
		}
		e.Columns = cols
		body[i] = e
	}
	influxResponse = append(influxResponse, header)
	interval := time.Minute * 15
	start := startInit
	var end time.Time
	for start.Before(endInit) {
		end = start.Add(interval)
		s := start.Format(time.RFC3339)
		en := end.Format(time.RFC3339)
		log.Println("[INFLUX]", "Getting data between", s, "and", en)
		for i := range body {
			body[i].Time = &QueriesRequestElementTime{
				Start: &s,
				End:   &en,
			}
		}
		jsonData, _ := json.Marshal(body)
		client := http.Client{}
		request, err := http.NewRequest("POST", i.url+"/queries?format=table", bytes.NewBuffer(jsonData))
		if request != nil {
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", "Bearer "+accessToken)
		}
		resp, err := client.Do(request)
		if err != nil {
			return nil, err
		}
		var tmp [][]interface{}
		err = json.NewDecoder(resp.Body).Decode(&tmp)
		if err != nil {
			return nil, err
		}
		influxResponse = append(influxResponse, tmp...)
		log.Println("[INFLUX]", "Collected", len(tmp), "rows, total:", len(influxResponse))
		start = end
	}
	return
}
