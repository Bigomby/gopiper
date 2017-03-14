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

package component

// Message wraps the data sent through the pipeline
type Message interface {
	SetData(data []byte)
	GetData() []byte
	Status() *Report
	Release()
}

// GopiperMessage is used to send data through the pipeline
type GopiperMessage struct {
	data       []byte
	attributes map[string]interface{}

	report *Report
}

// NewMessage creates a new instance of Message
func NewMessage() Message {
	return &GopiperMessage{
		data:       []byte{},
		attributes: make(map[string]interface{}),

		report: &Report{},
	}
}

// SetData store data on an LIFO queue so the nexts handlers can use it
func (m *GopiperMessage) SetData(data []byte) {
	m.data = data
}

// GetData get the data stored by the previous handler
func (m *GopiperMessage) GetData() []byte {
	return m.data
}

// Status returns a report of the data being processes by a worker
func (m *GopiperMessage) Status() *Report {
	return m.report
}

// Release free resources associated to the message
func (m *GopiperMessage) Release() {}
