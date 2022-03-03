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
	"reflect"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
)

type ServingInstance struct {
	ID               uuid.UUID              `json:"ID,omitempty"`
	Name             string                 `json:"Name,omitempty"`
	Description      string                 `json:"Description,omitempty"`
	EntityName       string                 `json:"EntityName,omitempty"`
	ServiceName      string                 `json:"ServiceName,omitempty"`
	Topic            string                 `json:"Topic,omitempty"`
	Database         string                 `json:"Database,omitempty"`
	Measurement      string                 `json:"Measurement,omitempty"`
	Filter           string                 `json:"Filter,omitempty"`
	FilterType       string                 `json:"FilterType,omitempty"`
	TimePath         string                 `json:"TimePath,omitempty"`
	UserId           string                 `json:"UserId,omitempty"`
	RancherServiceId string                 `json:"RancherServiceId,omitempty"`
	Offset           string                 `json:"Offset,omitempty"`
	Values           []ServingInstanceValue `json:"Values,omitempty"`
	CreatedAt        time.Time              `json:"CreatedAt,omitempty"`
	UpdatedAt        time.Time              `json:"UpdatedAt,omitempty"`
}

type ServingInstanceValue struct {
	InstanceID uuid.UUID `json:"InstanceID,omitempty"`
	Name       string    `json:"Name,omitempty"`
	Type       string    `json:"Type,omitempty"`
	Path       string    `json:"Path,omitempty"`
}

type InfluxRequest struct {
	Time    InfluxTime    `json:"time,omitempty"`
	Queries []InfluxQuery `json:"queries,omitempty"`
}

type InfluxTime struct {
	Last  string `json:"last,omitempty"`
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

type InfluxQuery struct {
	Id string `json:"id,omitempty"`
}

type InfluxResponse struct {
	Results []InfluxResults `json:"results,omitempty"`
}

type InfluxResults struct {
	Series []InfluxSeries `json:"series,omitempty"`
}

type InfluxSeries struct {
	Columns []string        `json:"columns,omitempty"`
	Name    string          `json:"name,omitempty"`
	Values  [][]interface{} `json:"values,omitempty"`
}

func (i InfluxSeries) GetValuesAsString() (stringValues [][]string) {
	for _, val := range i.Values {
		var a []string
		for _, data := range val {
			if data != nil {
				switch reflect.TypeOf(data).Kind() {
				case reflect.Float64:
					a = append(a, strconv.FormatFloat(data.(float64), 'f', -1, 64))
					break
				case reflect.String:
					a = append(a, data.(string))
					break
				case reflect.Bool:
					a = append(a, strconv.FormatBool(data.(bool)))
					break
				default:
					break
				}
			}
		}
		stringValues = append(stringValues, a)
	}
	return
}
