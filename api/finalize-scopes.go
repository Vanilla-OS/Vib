package api

// / Get the final tagged Image name
var IMAEGNAME = 1

// / Get the final tagged Image ID
var IMAGEID = 2

// / Get the build recipe
var RECIPE = 4

// / Get a read-write filesystem of the Image
var RWFS = 8

// / Get a read-only filesystem of the Image
var ROFS = 16

// / Prepare the filesystem to be chrooted into, requires either RWFilesystem or ROFilesystem
var CHROOTFS = 32
