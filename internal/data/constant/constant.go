package constant

// RedisKeySmallDataset, RedisKeyMediumDataset, and RedisKeyLargeDataset are the Redis keys used to persist the three seeded datasets.
const (
	RedisKeySmallDataset  = "student:small"
	RedisKeyMediumDataset = "student:medium"
	RedisKeyLargeDataset  = "student:large"
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
