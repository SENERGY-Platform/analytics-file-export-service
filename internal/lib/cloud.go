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
	"io"
	"os"

	"github.com/studio-b12/gowebdav"
)

type CloudService struct {
	Client *gowebdav.Client
}

func NewCloudService(host string, username string, password string) *CloudService {
	c := gowebdav.NewClient(host, username, password)
	err := c.Connect()
	if err != nil {
		fmt.Println(err.Error())
	}
	return &CloudService{c}
}

func (cs *CloudService) MkDir(path string, mode os.FileMode) {
	err := cs.Client.MkdirAll(path, mode)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (cs *CloudService) UploadFile(path string, file io.Reader, mode os.FileMode) (err error) {
	err = cs.Client.WriteStream(path, file, mode)
	return
}

func (cs *CloudService) UploadFileFromByteArray(path string, file []byte, mode os.FileMode) (err error) {
	err = cs.Client.Write(path, file, mode)
	return
}
