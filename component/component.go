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

// Component represents any element in the pipeline that performs some work
// on a message.
//
// A component should call the callback function to report the status of the
// action performed. Depending on the report, the following actions will
// occur:
//
// - Done (code 0): The message has been successfully processed and should
// be sent to the next component of the pipeline.
// - Continue (code 1): The message has been successfully processed, but
// the message should not be sent to the next component of the pipeline.
// - Retry (code 100 - 199): The component has failed to process the message
// and should be retried.
// - Fail (code 200-299): The component has failed to process the message and
// the message should be discarded.
// - Discard (300+): Silently drop the message.
type Component interface {
	Handle(Message) *Report
}
