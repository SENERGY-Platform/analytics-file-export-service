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
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	gocloak "github.com/Nerzal/gocloak/v5"
)

const LOCAL_PATH = "files"

type ExportService struct {
	keycloak  KeycloakService
	serving   ServingService
	influx    InfluxService
	cloud     CloudService
	cloudPath string
}

func NewExportService(keycloak KeycloakService, serving ServingService, influx InfluxService, cloud CloudService, cloudPath string) *ExportService {
	return &ExportService{keycloak: keycloak, serving: serving, influx: influx, cloud: cloud, cloudPath: cloudPath}
}

func (es *ExportService) StartExportService() {
	es.keycloak.Login()
	defer es.keycloak.Logout()
	user, err := es.keycloak.GetUserInfo()
	if err != nil {
		log.Fatal("GetUserInfo failed:" + err.Error())
	}
	if user != nil {
		es.createCsvFiles(user)
	}
	es.uploadFiles()
}

func (es *ExportService) createCsvFiles(user *gocloak.UserInfo) {
	servings, err := es.serving.GetServingServices(*user.Sub, es.keycloak.GetAccessToken())
	if err != nil {
		log.Fatal("GetServingServices failed: " + err.Error())
	} else {
		var wg sync.WaitGroup
		for _, serving := range servings {
			now := time.Now()
			yesterday := now.AddDate(0, 0, -1)
			start := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)
			end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
			wg.Add(1)
			func() {
				defer wg.Done()
				data, _ := es.influx.GetData(es.keycloak.GetAccessToken(), serving.Measurement, start.Format(time.RFC3339), end.Format(time.RFC3339))
				for _, i := range data.Results {
					es.writeCsv(i, serving, start.Format("2006-01-02"))
				}
			}()
			time.Sleep(2 * time.Second)
		}
		wg.Wait()
	}
}

func (es *ExportService) uploadFiles() {
	var files []string
	err := filepath.Walk("files", func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, path := range files {
		fi, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
			return
		}
		switch mode := fi.Mode(); {
		case mode.IsDir():
			// do directory stuff
		case mode.IsRegular():
			f, _ := os.Open(path)
			defer f.Close()
			log.Println("Uploading: " + path)
			//bytes, _ := ioutil.ReadFile(path)
			//err := es.cloud.UploadFileFromByteArray(strings.Replace(path, LOCAL_PATH, es.cloudPath, -1), bytes, 0755)
			err := es.cloud.UploadFile(strings.Replace(path, LOCAL_PATH, es.cloudPath, -1), f, 0755)
			if err != nil {
				log.Fatalln("Could not upload " + path + " " + err.Error())
			} else {
				_ = os.Remove(path)
				log.Println("Uploading to " + strings.Replace(path, LOCAL_PATH, es.cloudPath, -1) + "... done")
			}
		}
	}
}

func (es *ExportService) writeCsv(i InfluxResults, serving ServingInstance, fileName string) {
	PATH := "./" + LOCAL_PATH + "/" + serving.Measurement + "_" + strings.Replace(serving.Name, " ", "_", -1) + "/"
	if _, err := os.Stat(PATH); os.IsNotExist(err) {
		_ = os.MkdirAll(PATH, 0755)
	}
	filePath := PATH + fileName + ".csv"
	// Create csv file
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// Write Unmarshaled json data to CSV file
	w := csv.NewWriter(f)
	//Columns
	_ = w.Write(i.Series[0].Columns[:])
	//Data
	for _, d := range i.Series[0].GetValuesAsString() {
		_ = w.Write(d)
	}
	w.Flush()
}
