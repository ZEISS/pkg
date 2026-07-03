package tplx

import (
	"text/template"
)

var genericMap = map[string]interface{}{
	"hello": func() string { return "Hello!" },
}

// FuncMap returns a 'text/template'.FuncMap.
func FuncMap() template.FuncMap {
	return HtmlFuncMap()
}

// TxtFuncMap returns a 'text/template'.FuncMap.
func TxtFuncMap() template.FuncMap {
	return template.FuncMap(GenericFuncMap())
}

// HtmlFuncMap returns an 'html/template'.Funcmap.
func HtmlFuncMap() template.FuncMap {
	return template.FuncMap(GenericFuncMap())
}

// GenericFuncMap returns a copy of the basic function map as a map[string]interface{}.
func GenericFuncMap() map[string]interface{} {
	gfm := make(map[string]interface{}, len(genericMap))
	for k, v := range genericMap {
		gfm[k] = v
	}

	return gfm
}
