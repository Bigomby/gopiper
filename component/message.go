// Process data using a pipeline.
// Copyright (C) 2017 Diego Fernández Barrera
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package component

// Message is used to send data through the pipeline
type Message struct {
	data       []byte
	attributes map[string]interface{}

	report *Report
}

// NewMessage creates a new instance of Message
func NewMessage() *Message {
	return &Message{
		data:       []byte{},
		attributes: make(map[string]interface{}),

		report: &Report{},
	}
}

// SetData store data on an LIFO queue so the nexts handlers can use it
func (m *Message) SetData(data []byte) {
	m.data = data
}

// GetData get the data stored by the previous handler
func (m *Message) GetData() []byte {
	return m.data
}

// Status returns a report of the data being processes by a worker
func (m *Message) Status() *Report {
	return m.report
}

// Release free resources asociated to the message
func (m *Message) Release() {}
