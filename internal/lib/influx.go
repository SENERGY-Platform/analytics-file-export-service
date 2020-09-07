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
)

type InfluxService struct {
	url string
}

func NewInfluxService(url string) *InfluxService {
	return &InfluxService{url: url}
}

func (i *InfluxService) GetData(accessToken string, id string, start string, end string) (influxResponse InfluxResponse, err error) {
	TmpPath := GetEnv("FILES_PATH", "files") + "/tmp/"
	data := InfluxRequest{
		Time:    InfluxTime{Start: start, End: end},
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
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		_ = out.Close()
		fmt.Println(err)
	}
	_ = out.Close()
	jsonFile, err := os.Open(TmpPath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	err = json.NewDecoder(jsonFile).Decode(&influxResponse)
	os.Remove(TmpPath)
	return
}
