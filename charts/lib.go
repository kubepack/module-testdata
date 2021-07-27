package charts

import (
	"embed"
)

//go:embed ** **/**/_helpers.tpl **/.helmignore
var FS embed.FS
