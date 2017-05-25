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

import (
	"sync"

	"github.com/Bigomby/gopiper/component"
)

// pool contains all the workers for a component and a factory to create
// workers for that component.
//
// pool also contains the channels to send and receive data from/to the pool:
// - input: To inject messages on the worker pool.
// - produce: Allow workers to produce their own messages.
// - output: Messages already processed.
// - terminate: Use to shutdown workers.
//
// The done WaitGroup is used to wait for all workers to shutdown.
type pool struct {
	factory   component.Factory
	workers   chan component.Component
	input     chan component.Message
	produce   chan component.Message
	output    chan component.Message
	terminate chan struct{}
	done      *sync.WaitGroup
}

// Pipeline contains an array of worker pools of components.
type Pipeline struct {
	pools []*pool
}

// NewPipeline creates a new Pipeline. It uses factories to spawn workers and
// push them to the worker pool.
//
// The output channel for a pool is the input channel for the next pool.
// The input channel for the first component is nil. The first component Handle
// function won't be called so it should not be implemented.
func NewPipeline(factories []component.Factory) *Pipeline {
	pipeline := &Pipeline{}

	for i, factory := range factories {
		pool := pool{
			factory:   factory,
			terminate: make(chan struct{}),
			done:      &sync.WaitGroup{},
			workers:   make(chan component.Component, factory.PoolSize()),
		}

		pool.produce = make(chan component.Message, factory.ChannelSize())
		pool.output = make(chan component.Message, factory.ChannelSize())
		if i > 0 {
			pool.input = pipeline.pools[i-1].output
		}

		for j := 0; j < factory.PoolSize(); j++ {
			spawnWorker(factory, pool)
		}

		pipeline.pools = append(pipeline.pools, &pool)
	}

	go func() {
		for msg := range pipeline.pools[len(factories)-1].output {
			msg.Release()
		}
	}()

	return pipeline
}

// Close broadcast a terminate signal to all pools, wait for all gorutines to
// finish and release associated resources
func (p *Pipeline) Close() {
	for _, pool := range p.pools {
		close(pool.terminate)
		pool.done.Wait()
		pool.factory.Destroy()
	}
}

func spawnWorker(factory component.Factory, pool pool) {
	go func() {
		pool.workers <- factory.Create(pool.produce)
		pool.done.Add(1)

	running:
		for {
			select {
			// Handle the message coming from the previous component
			case msg := <-pool.input:
				worker := <-pool.workers
				process(worker.Handle, pool.output, msg)
				pool.workers <- worker

			// Handle messages coming from this component
			case msg := <-pool.produce:
				pool.output <- msg

			// Handle terminate signal
			case <-pool.terminate:
				break running
			}
		}

		pool.done.Done()
	}()
}

func process(
	handle func(component.Message) *component.Report,
	output chan<- component.Message,
	msg component.Message,
) {
	report := handle(msg)

	switch {
	case report.Status == component.Done:
		output <- msg

	case report.Status == component.Continue:
		// TODO

	case report.Status >= component.Fail && report.Status < component.Retry:
		// TODO Retry messages

	case report.Status >= component.Retry && report.Status < component.Drop:
		// TODO

	default:
		// TODO

	}
}
