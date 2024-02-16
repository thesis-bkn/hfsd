package templates

import "embed"

//go:embed "static"
var Files embed.FS

//go:embed "asset"
var Assets embed.FS
