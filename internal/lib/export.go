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
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	gocloak "github.com/Nerzal/gocloak/v5"
	humanize "github.com/dustin/go-humanize"
)

type ExportService struct {
	keycloak  KeycloakService
	serving   ServingService
	influx    InfluxService
	cloud     CloudService
	cloudPath string
	filePath  string
}

func NewExportService(keycloak KeycloakService, serving ServingService, influx InfluxService, cloud CloudService, cloudPath string) *ExportService {
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
}

func (es *ExportService) createCsvFiles(user *gocloak.UserInfo) {
	servings, err := es.serving.GetServingServices(*user.Sub, es.keycloak.GetAccessToken())
	if err != nil {
		log.Fatal("GetServingServices failed: " + err.Error())
	} else {
		var wg sync.WaitGroup
		if len(GetEnv("SINGLE_MEASUREMENT", "")) > 0 {
			for _, serving := range servings {
				if serving.Measurement == GetEnv("SINGLE_MEASUREMENT", "") {
					servings = []ServingInstance{serving}
				}
			}
		}
		servingsTotal := strconv.Itoa(len(servings))
		for no, serving := range servings {
			wg.Add(1)
			func() {
				fmt.Println("Get (" + strconv.Itoa(no+1) + "/" + servingsTotal + "):" + serving.Measurement + " - " + serving.Name)
				days, err := strconv.Atoi(GetEnv("DAYS_BACK", "1"))
				if err != nil {
					fmt.Println(err)
				}
				es.getInfluxDataOfExportLastDays(serving, days)
				defer wg.Done()
			}()
		}
		wg.Wait()
	}
}

func (es *ExportService) getInfluxDataOfExportLastDays(serving ServingInstance, days int) {
	now := time.Now()
	for day := 0; day > -days; day-- {
		startDate := now.AddDate(0, 0, day-1)
		endDate := startDate.AddDate(0, 0, 1)
		fmt.Println(startDate.Format("2006-01-02"))
		start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
		end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, time.UTC)
		data, _ := es.influx.GetData(es.keycloak.GetAccessToken(), serving.Measurement, start.Format(time.RFC3339), end.Format(time.RFC3339))
		for _, i := range data.Results {
			es.writeCsv(i, serving, start.Format("2006-01-02"))
		}
	}
}

func (es *ExportService) uploadFiles() {
	files, err := walkMatch(es.filePath, "*.csv")
	if err != nil {
		panic(err)
	}
	filesTotal := strconv.Itoa(len(files))
	for no, path := range files {
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
			fmt.Println("Uploading (" + strconv.Itoa(no+1) + "/" + filesTotal + "): " + path)
			//bytes, _ := ioutil.ReadFile(path)
			//err := es.cloud.UploadFileFromByteArray(strings.Replace(path, LOCAL_PATH, es.cloudPath, -1), bytes, 0755)
			err := es.cloud.UploadFile(strings.Replace(path, es.filePath, es.cloudPath, -1), f, 0755)
			if err != nil {
				fmt.Println("Could not upload " + path + " " + err.Error())
			} else {
				_ = os.Remove(path)
				fmt.Println("... done")
			}
			err = f.Close()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func walkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (es *ExportService) writeCsv(i InfluxResults, serving ServingInstance, fileName string) {
	PATH := "./" + es.filePath + "/" + serving.Measurement + "_" + strings.Replace(serving.Name, " ", "_", -1) + "/"
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
	PrintMemUsage()
	w.Flush()
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", humanize.Bytes(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", humanize.Bytes(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", humanize.Bytes(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
