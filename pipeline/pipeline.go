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

package pipeline

import "github.com/Bigomby/gopiper/component"

// workerPool contains the worker pool and the channels to send and receive data
// from/to the pool.
type workerPool struct {
	workers   chan component.Component
	input     chan component.Message
	postponed chan component.Message
	output    chan component.Message
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

		if i > 0 {
			pool.input = pipeline.pools[i-1].output
		}
		pool.postponed = make(chan component.Message, 100)
		pool.output = make(chan component.Message, 100)

		for j := 0; j < factories[i].Amount(); j++ {
			spawnWorker(factories[i], pool.workers, pool.input, pool.postponed,
				pool.output)
		}

		pipeline.pools[i] = &pool
	}

	return pipeline
}

// Close closes all resources
func (p *Pipeline) Close() { /* TODO */ }

func spawnWorker(
	factory component.Factory,
	workerPool chan component.Component,
	input <-chan component.Message,
	postponed chan component.Message,
	output chan<- component.Message,
) {
	go func() {
		workerPool <- factory.Create(postponed)

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
	output chan<- component.Message, msg component.Message,
) {
	worker.Handle(msg, func(report *component.Report) {
		switch {
		case report.Status == component.Done:
			output <- msg

		case report.Status == component.Continue:
		default:
		}

		msg.Release()
	})
}
