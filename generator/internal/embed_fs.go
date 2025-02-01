package internal

import "embed"

//go:embed template/*
var EmbedFs embed.FS
