package constant

// Order-experiment Redis keys (depth-scenario and element-count-scenario
// investigation, replacing the retired Student dataset).
const (
	RedisKeyOrderDepthZero = "order:depth-zero"
	RedisKeyOrderDepthTwo  = "order:depth-two"
	RedisKeyOrderDepthFour = "order:depth-four"
	RedisKeyOrderOne       = "order:one"
	RedisKeyOrderHundred   = "order:hundred"
	RedisKeyOrderThousand  = "order:thousand"
)

// Shape-experiment Redis keys (structural-depth investigation, separate
// from the Student dataset).
const (
	RedisKeyShapeDepth0Compact       = "shape:depth0:compact"
	RedisKeyShapeDepth0Large         = "shape:depth0:large"
	RedisKeyShapeDepth1WideCompact   = "shape:depth1-wide:compact"
	RedisKeyShapeDepth1WideLarge     = "shape:depth1-wide:large"
	RedisKeyShapeDepth3NarrowCompact = "shape:depth3-narrow:compact"
	RedisKeyShapeDepth3NarrowLarge   = "shape:depth3-narrow:large"
	RedisKeyShapeDepth4WideCompact   = "shape:depth4-wide:compact"
	RedisKeyShapeDepth4WideLarge     = "shape:depth4-wide:large"
)
