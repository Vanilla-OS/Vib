package api

// / Get the final tagged Image name
var IMAGENAME int32 = 1

// / Get the final tagged Image ID
var IMAGEID int32 = 2

// / Get the build recipe
var RECIPE int32 = 4

// Get the used build runtime
var RUNTIME int32 = 8

// / Get a read-write filesystem of the Image
var RWFS int32 = 16

// / Get a read-only filesystem of the Image
var ROFS int32 = 32

// / Prepare the filesystem to be chrooted into, requires either RWFilesystem or ROFilesystem
var CHROOTFS int32 = 64

type ScopeData struct {
	ImageName string
	ImageID   string
	Recipe    Recipe
	Runtime   string
	RWFS      string
	ROFS      string
}
