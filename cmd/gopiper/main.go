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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"

	"github.com/Bigomby/gopiper/pipeline"
	lua "github.com/yuin/gopher-lua"
)

var (
	version string
	pipe    string

	terminate = make(chan struct{})

	// Pipeline is the global pipeline object
	Pipeline *pipeline.Pipeline
)

func printVersion() {
	fmt.Printf("Gopiper :: %s\n", version)
	fmt.Printf("Go      :: %s\n", runtime.Version())
}

func init() {
	versionFlag := flag.Bool("version", false, "Show version info")
	pipeFlag := flag.String("pipe", "pipeline.lua",
		"Lua file containing the description of the pipeline")
	flag.Parse()

	if *versionFlag {
		printVersion()
		os.Exit(0)
	}

	if *pipeFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	pipe = *pipeFlag
}

func main() {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("gopiper", LuaLoader)
	if err := L.DoFile(pipe); err != nil {
		log.Fatal(err)
	}

	sigint := make(chan os.Signal)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	Pipeline.Close()
}
