local gopiper = require("gopiper")

local components = {}
table.insert(
  components,
  gopiper.loadComponent("components/stdin_component.so", {})
)
table.insert(
  components,
  gopiper.loadComponent("components/stdout_component.so", {})
)

gopiper.createPipeline(components)

print("Pipeline started")
