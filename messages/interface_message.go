// Copyright 2017 Diego Fernández Barrera
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package messages

import "github.com/Bigomby/gopiper/component"

// InterfaceMessage is a message implementation which can contains anything
type InterfaceMessage struct {
	data       interface{}
	attributes map[string]interface{}
	report     *component.Report
}

// GetData get the bytes array stored by the previous handler
func (m InterfaceMessage) GetData() interface{} {
	return m.data
}

// SetData receives an interface that should be an array of bytes
func (m *InterfaceMessage) SetData(data interface{}) {
	m.data = data
}

// GetAttribute returns the attribute value for a given key
func (m InterfaceMessage) GetAttribute(attr string) interface{} {
	return m.attributes[attr]
}

// SetAttribute sets the attribute value for a given key
func (m *InterfaceMessage) SetAttribute(attr string, value interface{}) {
	m.attributes[attr] = value
}

// Status returns a report of the data being processes by a worker
func (m InterfaceMessage) Status() *component.Report {
	return m.report
}

// Release free resources associated to the message
func (m *InterfaceMessage) Release() {}
