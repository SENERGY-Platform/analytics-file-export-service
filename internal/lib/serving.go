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
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"strconv"
)

type ServingService struct {
	url string
}

func NewServingService(url string) *ServingService {
	return &ServingService{url: url}
}

func (s *ServingService) GetServingServices(userId string, accessToken string) (servings []ServingInstance, err error) {
	request := gorequest.New()
	resp, body, _ := request.Get(s.url+"/instance").Set("X-UserId", userId).Set("Authorization", "Bearer "+accessToken).End()
	if resp.StatusCode != 200 {
		fmt.Println("could not access serving service: "+strconv.Itoa(resp.StatusCode), resp.Body)
		return servings, errors.New("could not access serving service")
	}
	err = json.Unmarshal([]byte(body), &servings)
	return
}
