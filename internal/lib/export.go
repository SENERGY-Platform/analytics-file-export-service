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
	gocloak "github.com/Nerzal/gocloak/v5"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ExportService struct {
	keycloak  KeycloakService
	serving   ServingService
	influx    InfluxService
	cloud     CloudService
	cloudPath string
	filePath  string
	wg        sync.WaitGroup
}

var NOW = time.Now()
var didNotExport []string

func NewExportService(keycloak KeycloakService, serving ServingService, influx InfluxService, cloud CloudService, cloudPath string) *ExportService {
	if GetEnv("NOW_DATE", "") != "" {
		NOW, _ = time.Parse("2006-01-02", GetEnv("NOW_DATE", ""))
	}
	filePath := GetEnv("FILES_PATH", "files")
	return &ExportService{keycloak: keycloak, serving: serving, influx: influx, cloud: cloud, cloudPath: cloudPath, filePath: filePath}
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
	//es.uploadFiles()
	log.Println("[CORE]", "All files created, waiting for upload to finish...")
	es.wg.Wait()
	if len(didNotExport) > 0 {
		log.Println("Did not upload:")
		for _, export := range didNotExport {
			log.Println(export)
		}
	}
	log.Println("[CORE]", "All done, bye!")
}

func (es *ExportService) createCsvFiles(user *gocloak.UserInfo) {
	servings, err := es.serving.GetServingServices(*user.Sub, es.keycloak.GetAccessToken())
	if err != nil {
		log.Fatal("GetServingServices failed: " + err.Error())
	} else {
		days, err := strconv.Atoi(GetEnv("DAYS_BACK", "1"))
		if err != nil {
			log.Println(err)
		}
		es.getInfluxDataOfExportLastDays(servings, days)
	}
}

func (es *ExportService) getInfluxDataOfExportLastDays(servings []ServingInstance, days int) {
	for day := 0; day > -days; day-- {
		startDate := NOW.AddDate(0, 0, day-1)
		for hour := 0; hour < 24; hour++ {
			start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), hour, 0, 0, 0, time.UTC)
			end := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), hour+1, 0, 0, 0, time.UTC)
			data, err := es.influx.GetData(es.keycloak.GetAccessToken(), servings, start, end)
			if err != nil {
				log.Println(err.Error())
			}
			es.writeCsv(data, start)
		}
		log.Println("... done")
	}
}

func (es *ExportService) uploadFile(path string) {
	es.wg.Add(1)
	log.Println("[CORE]", "Starting upload:", path)
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
		//bytes, _ := ioutil.ReadFile(path)
		//err := es.cloud.UploadFileFromByteArray(strings.Replace(path, LOCAL_PATH, es.cloudPath, -1), bytes, 0755)
		_ = es.cloud.UploadFile(strings.Replace(path, es.filePath, es.cloudPath, -1), f, 0755)
		_ = f.Close()
		err = os.Remove(path)
		if err != nil {
			log.Println("[CORE]", "Could not delete", path)
		}
	}
	es.wg.Done()
}

func (es *ExportService) writeCsv(data [][]interface{}, t time.Time) {
	folder := es.filePath + "/" + t.Format("2006-01-02")
	filePath := folder + "/" + t.Format(time.RFC3339) + ".csv"
	defer func() {
		if r := recover(); r != nil {
			didNotExport = append(didNotExport, filePath)
			_ = os.Remove(filePath)
		}
	}()
	if _, err := os.Stat(es.filePath); os.IsNotExist(err) {
		_ = os.MkdirAll(es.filePath, 0755)
	}
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		_ = os.MkdirAll(folder, 0755)
	}

	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	w := csv.NewWriter(f)
	dataS := getValuesAsString(data)
	_ = w.WriteAll(dataS)

	log.Println("[CORE]", "Writing file", filePath)
	// w.Flush() already called by WriteAll()
	_ = f.Close()
	log.Println("[CORE]", "Launching upload in background...")
	go es.uploadFile(filePath)
}

func getValuesAsString(data [][]interface{}) (stringValues [][]string) {
	stringValues = make([][]string, len(data))
	for i := range data {
		row := make([]string, len(data[i]))
		for j := range data[i] {
			if data[i][j] == nil {
				row[j] = ""
			} else {
				row[j] = fmt.Sprintf("%v", data[i][j])
			}
		}
		stringValues[i] = row
	}
	return
}
