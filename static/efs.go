package static

import "embed"

//go:embed "js" "css"
var StaticFS embed.FS
