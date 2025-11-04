package sanitizer

import "regexp"

var (
	// Patterns for comprehensive CSS sanitization to prevent XSS attacks

	// Remove HTML tags (including malformed ones with extra spaces)
	htmlTagPattern = regexp.MustCompile(`(?i)<[^>]*>`)

	// Remove CSS expressions (IE-specific, but still dangerous)
	cssExpressionPattern = regexp.MustCompile(`(?i)expression\s*\([^)]*\)`)

	// Remove dangerous URL schemes (javascript:, data:, vbscript:)
	// This handles url(), @import, and other CSS contexts where URLs can appear
	dangerousURLPattern = regexp.MustCompile(`(?i)url\s*\(\s*['"]?\s*(javascript|data|vbscript):`)

	// Remove @import rules (can load external malicious CSS)
	importPattern = regexp.MustCompile(`(?i)@import\s+[^;]+;?`)

	// Remove style tag breaking sequences
	styleBreakPattern = regexp.MustCompile(`(?i)</\s*style[^>]*>`)

	// Remove potential event handlers embedded in CSS
	eventHandlerPattern = regexp.MustCompile(`(?i)on\w+\s*=`)

	// Remove HTML entities that could be decoded to dangerous characters
	htmlEntityPattern = regexp.MustCompile(`&[#\w]+;`)
)

// CSS aggressively removes dangerous constructs from CSS to prevent XSS
// This provides defense-in-depth alongside template.CSS() escaping
func CSS(css string) string {
	if css == "" {
		return ""
	}

	// Apply multiple sanitization passes to handle various attack vectors
	sanitized := css

	// 1. Remove @import rules (external CSS loading)
	sanitized = importPattern.ReplaceAllString(sanitized, "")

	// 2. Remove CSS expressions (IE-specific but dangerous)
	sanitized = cssExpressionPattern.ReplaceAllString(sanitized, "")

	// 3. Remove dangerous URL schemes (javascript:, data:, vbscript:)
	sanitized = dangerousURLPattern.ReplaceAllString(sanitized, "url(")

	// 4. Remove any HTML tags (including malformed ones)
	sanitized = htmlTagPattern.ReplaceAllString(sanitized, "")

	// 5. Remove style-breaking sequences
	sanitized = styleBreakPattern.ReplaceAllString(sanitized, "")

	// 6. Remove potential event handlers
	sanitized = eventHandlerPattern.ReplaceAllString(sanitized, "")

	// 7. Remove HTML entities (defense in depth)
	sanitized = htmlEntityPattern.ReplaceAllString(sanitized, "")

	return sanitized
}
