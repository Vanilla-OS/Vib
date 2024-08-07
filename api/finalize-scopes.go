package api

// / Get the final tagged Image name
var IMAGENAME int32 = 1

// / Get the final tagged Image ID
var IMAGEID int32 = 2

// / Get the build recipe
var RECIPE int32 = 4

// Get the used build runtime
var RUNTIME int32 = 8

// / Get a read-only filesystem of the Image
var FS int32 = 16

type ScopeData struct {
	ImageName string
	ImageID   string
	Recipe    Recipe
	Runtime   string
	FS        string
}
