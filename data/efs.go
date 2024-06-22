package data

import "embed"

const SamplePath = "/assets/samples"

//go:embed assets/samples
var Samples embed.FS
