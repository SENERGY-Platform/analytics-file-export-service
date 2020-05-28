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
