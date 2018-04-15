package query

import (
	"github.com/valyala/fasttemplate"
)

var (
	startTag = "{{"
	endTag   = "}}"
	// AllTemplates allows template reflection and generation by external clients eg. queryjs
	AllTemplates = map[string]*fasttemplate.Template{}
)
