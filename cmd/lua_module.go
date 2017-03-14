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

package main

import (
	"plugin"

	"github.com/Bigomby/gopiper/component"
	"github.com/Bigomby/gopiper/pipeline"
	lua "github.com/yuin/gopher-lua"
)

// exports contains the exported functions to the lua environment
var exports = map[string]lua.LGFunction{
	"loadComponent":  loadComponent,
	"createPipeline": createPipeline,
}

// LuaLoader is the function to load the Go functions into the lua environment
func LuaLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	mt := L.NewTypeMetatable("factory")
	L.SetGlobal("factory", mt)
	L.Push(mod)

	return 1
}

//////////////////////
// Module functions //
//////////////////////

// loadComponent loads a Go plugin for a component. A component should export a
// NewFactory symbol for a function that returns a Factory interface and set
// the given attributes of the Factory.
//
// The first argument is the path to the plugin (.so file)
// The second argument is a table with the attributes
func loadComponent(L *lua.LState) int {
	pluginPath := L.CheckString(1)
	attributes := L.CheckTable(2)

	p, err := plugin.Open(pluginPath)
	if err != nil {
		panic(err)
	}

	fs, err := p.Lookup("NewFactory")
	if err != nil {
		panic(err)
	}

	factory := fs.(func() component.Factory)()
	attributes.ForEach(func(key lua.LValue, val lua.LValue) {
		factory.SetAttribute(key.String(), val.String())
	})

	ud := L.NewUserData()
	ud.Value = factory
	L.SetMetatable(ud, L.GetTypeMetatable("factory"))

	L.Push(ud)

	return 1
}

// createPipeline receives the path to the plugin to load and a table with
// attributes and creates a pipeline with those plugins.
func createPipeline(L *lua.LState) int {
	var factories []component.Factory

	L.CheckTable(1).ForEach(func(i lua.LValue, val lua.LValue) {
		var (
			ok      bool
			ud      *lua.LUserData
			factory component.Factory
		)

		if ud, ok = val.(*lua.LUserData); !ok {
			panic("Element in table is not user data")
		}

		if factory, ok = ud.Value.(component.Factory); !ok {
			panic("Element in table is not factory")
		}

		factories = append(factories, factory)
	})

	pipeline.NewPipeline(factories)

	return 0
}
