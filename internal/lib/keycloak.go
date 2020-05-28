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
	"github.com/Nerzal/gocloak/v5"
)

type KeycloakService struct {
	client       gocloak.GoCloak
	token        *gocloak.JWT
	clientId     string
	clientSecret string
	realm        string
	userName     string
	password     string
}

func NewKeycloakService(url string, clientId string, clientSecret string, realm string, userName string, password string) *KeycloakService {
	client := gocloak.NewClient(url)
	return &KeycloakService{client, nil, clientId, clientSecret, realm, userName, password}
}

func (k *KeycloakService) Login() {
	token, err := k.client.Login(k.clientId, k.clientSecret, k.realm, k.userName, k.password)
	if err != nil {
		fmt.Println("Login failed:" + err.Error())
	}
	k.token = token
}

func (k *KeycloakService) Logout() {
	err := k.client.Logout(k.clientId, k.clientSecret, k.realm, k.token.RefreshToken)
	if err != nil {
		fmt.Println("Logout failed:" + err.Error())
	}
}

func (k *KeycloakService) GetAccessToken() string {
	return k.token.AccessToken
}

func (k *KeycloakService) GetUserInfo() (*gocloak.UserInfo, error) {
	user, err := k.client.GetUserInfo(k.token.AccessToken, k.realm)
	return user, err
}

func (k *KeycloakService) GetUserByID(id string) (user *gocloak.User, err error) {
	user, err = k.client.GetUserByID(k.token.AccessToken, k.realm, id)
	return
}
