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

package gopiper

import "github.com/bigomby/gopiper/component"

const (
	// Done informs that a worker has finshed processing the message
	Done = iota

	// Continue informs that a worker needs more message
	Continue
)

// workerPool contains the worker pool and the channels to send and receive data
// from/to the pool.
type workerPool struct {
	workers   chan component.Component
	input     chan *component.Message
	postponed chan *component.Message
	output    chan *component.Message
}

// Pipeline contains an array of worker pools of components.
type Pipeline struct {
	pools []*workerPool
}

// NewPipeline creates a new Pipeline. It uses factories to spawn workers and
// push them to the worker pool.
//
// Every worker pool has a set of workers, a channel for ingesting messages,
// a channel to receive postponed messages from workers, and channel for send
// processed messages.
//
// The output channel for a pool is the input channel for the next pool.
func NewPipeline(factories []component.Factory) *Pipeline {
	pipeline := &Pipeline{
		pools: make([]*workerPool, len(factories)),
	}

	for i := range factories {
		pool := workerPool{
			workers: make(chan component.Component, factories[i].Amount()),
		}

		if i == 0 {
			pool.input = make(chan *component.Message, 100)
		} else {
			pool.input = pipeline.pools[i-1].output
		}
		pool.output = make(chan *component.Message, 100)
		pool.postponed = make(chan *component.Message, 100)

		for j := 0; j < factories[i].Amount(); j++ {
			component := factories[i].Create(j, pool.postponed)
			pool.workers <- component
		}

		spawnWorker(pool.workers, pool.input, pool.postponed, pool.output)
		pipeline.pools[i] = &pool
	}

	return pipeline
}

// Close closes all resources
func (p *Pipeline) Close() { /* TODO */ }

func spawnWorker(
	workerPool chan component.Component,
	input <-chan *component.Message,
	postponed <-chan *component.Message,
	output chan<- *component.Message,
) {
	go func() {
		for {
			select {
			// Handle the message coming from the previous component
			case msg := <-input:
				worker := <-workerPool
				handleMessage(worker, output, msg)
				workerPool <- worker

				// Handle messages coming from this component
			case msg := <-postponed:
				output <- msg
			}
		}
	}()
}

func handleMessage(
	worker component.Component,
	output chan<- *component.Message,
	msg *component.Message,
) {
	worker.Handle(msg, func(report *component.Report) {
		switch {
		case report.Code == Done:
			output <- msg

		case report.Code == Continue:
		default:
		}

		msg.Release()
	})
}
