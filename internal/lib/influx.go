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
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type InfluxService struct {
	url string
}

func NewInfluxService(url string) *InfluxService {
	return &InfluxService{url: url}
}

func (i *InfluxService) GetData(accessToken string, id string, start time.Time) (influxResponse InfluxResponse, err error) {
	for step := 23; step >= 0; step-- {
		var tmpInfluxResponse InfluxResponse
		TmpPath := GetEnv("TMP_PATH", "/tmp/")
		start = time.Date(start.Year(), start.Month(), start.Day(), step, 0, 0, 0, time.UTC)
		end := start.Add(time.Hour * time.Duration(1))
		data := InfluxRequest{
			Time:    InfluxTime{Start: start.Format(time.RFC3339), End: end.Format(time.RFC3339)},
			Queries: []InfluxQuery{{Id: id}},
		}
		jsonData, _ := json.Marshal(data)
		client := http.Client{}
		request, err := http.NewRequest("POST", i.url+"/queries", bytes.NewBuffer(jsonData))
		if request != nil {
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", "Bearer "+accessToken)
		}
		resp, err := client.Do(request)

		if _, err := os.Stat(TmpPath); os.IsNotExist(err) {
			_ = os.MkdirAll(TmpPath, 0755)
		}
		if resp != nil {
			defer resp.Body.Close()
		}
		TmpPath += id + ".tmp"
		out, err := os.Create(TmpPath)
		if err != nil {
			fmt.Println(err)
		}
		if GetEnv("DEBUG", "false") == "true" {
			counter := &WriteCounter{}
			if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
				_ = out.Close()
				fmt.Println(err)
			}
		} else {
			if _, err = io.Copy(out, resp.Body); err != nil {
				_ = out.Close()
				fmt.Println(err)
			}
		}
		_ = out.Close()
		jsonFile, err := os.Open(TmpPath)
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		err = json.NewDecoder(jsonFile).Decode(&tmpInfluxResponse)
		if len(influxResponse.Results) < 1 {
			influxResponse = tmpInfluxResponse
		} else {
			if len(tmpInfluxResponse.Results) > 0 && len(tmpInfluxResponse.Results[0].Series) > 0 {
				if len(influxResponse.Results[0].Series[0].Columns) < 1 {
					influxResponse.Results[0].Series[0].Columns = tmpInfluxResponse.Results[0].Series[0].Columns
				}
				influxResponse.Results[0].Series[0].Values = append(influxResponse.Results[0].Series[0].Values, tmpInfluxResponse.Results[0].Series[0].Values...)
			}
		}
	}
	return
}
