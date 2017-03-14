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

// Component represents any element in the pipeline that performs some work
// on a message.
//
// A component should call the callback function to report the status of the
// action performed. Depending on the report, the following actions will
// occurr:
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
	Handle(*Message, HandledCallback)
}
