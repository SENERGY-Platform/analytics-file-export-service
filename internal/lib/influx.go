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
	"encoding/json"
	"github.com/parnurzeal/gorequest"
)

type InfluxService struct {
	url string
}

func NewInfluxService(url string) *InfluxService {
	return &InfluxService{url: url}
}

func (i *InfluxService) GetData(accessToken string, id string) (influxResponse InfluxResponse, err error) {
	request := gorequest.New()
	data := InfluxRequest{
		Time:    InfluxTime{Last: "1d"},
		Queries: []InfluxQuery{{Id: id}},
	}
	_, body, _ := request.Post(i.url+"/queries").Set("Authorization", "Bearer "+accessToken).Send(data).End()
	err = json.Unmarshal([]byte(body), &influxResponse)
	return
}
