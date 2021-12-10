package functions

import "net/url"

func init() {
	registerFunction("urlescape", url.QueryEscape)
}
