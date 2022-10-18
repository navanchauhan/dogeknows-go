package mytemplates

import "embed"

// myTemplates represent the templates used by the application.
//
//go:embed *.html *.gtpl components/*.html
var Templates embed.FS
