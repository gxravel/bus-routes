package assets

import "embed"

// SwaggerFiles swagger data
//go:embed swagger/*
var SwaggerFiles embed.FS
