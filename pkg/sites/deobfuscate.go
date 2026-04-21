package sites

import (
	"html"
	"net/url"
	"strings"
)

func deobfuscateURL(raw string) string {
	if raw == "" {
		return ""
	}

	cleaned := strings.ReplaceAll(raw, `\/`, "/")
	cleaned = strings.ReplaceAll(cleaned, `\\`, `\`)

	decoded := html.UnescapeString(cleaned)

	if result, err := url.PathUnescape(decoded); err == nil {
		decoded = result
	}

	return decoded
}
