[![Build Status](https://travis-ci.org/Bigomby/gopiper.svg?branch=master)](https://travis-ci.org/Bigomby/gopiper)
[![Coverage Status](https://coveralls.io/repos/github/Bigomby/gopiper/badge.svg?branch=master)](https://coveralls.io/github/Bigomby/gopiper?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Bigomby/gopiper)](https://goreportcard.com/report/github.com/Bigomby/gopiper)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/d16082f693d247759084d54ba2f1db3d)](https://www.codacy.com/app/Bigomby/gopiper?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=Bigomby/gopiper&amp;utm_campaign=Badge_Grade)

# gopiper

## Overview

> **NOTE: THIS APP IS ON A PRE-ALPHA STATE. EVERYTHING IS GOING TO CHANGE**

`gopiper` is a tool written in Go for creating high perfomance **message
processing pipelines**. This approach makes heavy use of Go channels.

The purpose of this tool is provide a minimal framework to create
**components** (elements that process messages in a pipeline) and an
abstraction over the orchestration of these components using **lua** as
configuration language.

## Creating components

Components are compiled as Go plugins and are loaded at runtime. To create a
component you just need to implement the `Component` interface and a `Factory`
interface.

```go
type Factory interface {
  Create(postponed chan Message) Component
  Destroy()
  SetAttribute(key string, value interface{}) error
  PoolSize() int
  ChannelSize() int
}

type Component interface {
  Handle(Message, HandledCallback)
}
```

Since the approach of tis app uses a pool of worker, it's recommended to process
messages in a synchronous way. The pipeline's internal logic takes care of
managing workers so you should not return from the `Handle()` function until you
are completely done with the message.

A component receives a `Message` interface, gets data from the message and
performs some work with it. When the job is done, the component should call
the `HandledCallback` to inform the pipeline that the message has been
processed.

```Go
type Message interface {
  GetData() interface{}
  SetData(interface{})
  GetAttribute(string) interface{}
  SetAttribute(string, interface{})
  Status() *Report
  Release()
}
```

### Error handling

After processing a message you should use the `HandledCallback` callback to
inform the result of the processing.

```go
type Report struct {
  Status      int
  Description string
}
```

Status codes are the following:

| Status    | Action                                                           |
|-----------|------------------------------------------------------------------|
| 0         | Done, the next component should process the messag               |
| 1         | Continue, the component expects more messages before send a message to the next worker |
| 2 - 99    | Reserved                                                         |
| 100 - 199 | Retry, the message has been failed to process. Should be retried.|
| 200 - 299 | Fail, the message has been failed to process. Should not be retried |
| 300+      | Drop, silently drop the message.                                 |

You can also include a description for your status for debugging purposes.

## Building a pipeline

### Using lua

The easiest way of build a pipeline is using the lua API. You can load
components and start the pipeline as you can see in the following example:

```lua
-- Import the module
local gopiper = require('gopiper')

-- Insert desired components on compontents table
local components = {
  gopiper.loadComponent('stdin_component.so', {}),
  gopiper.loadComponent('stdout_component.so', {}),
}

-- Create the pipeline
gopiper.createPipeline(components)
```

### Using from Go

`// TODO`

### Using from C

`// TODO`

## Instrumentation

`// TODO`
