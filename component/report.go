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

const (
	// Done is the status code for indicating that the message has been
	// successfuly handled and should be sent to the next component
	Done = 0

	// Continue is the status code for indicating that the message has been
	// successfuly handled, but should not be sent to the next component
	Continue = 1

	// Retry is the status code for indicating that the worker has failed to
	// handle the message and should be retried again.
	Retry = 100

	// Fail is the status code for indicating that the worker has failed to
	// handle the message and should not be retried again.
	Fail = 200

	// Drop is the status code for indicating that the message has been discarded
	// by the worker and also shuold be
	Drop = 300
)

// Report contains information about the result of processing a message
type Report struct {
	Status      int
	Description string
}
