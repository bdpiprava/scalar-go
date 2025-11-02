package sanitizer_test

import (
	"testing"

	"github.com/bdpiprava/scalar-go/sanitizer"
	"github.com/stretchr/testify/assert"
)

func TestCSS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Edge Cases
		{
			name:  "should return empty string for empty input",
			input: "",
			want:  "",
		},
		{
			name:  "should preserve clean CSS unchanged",
			input: ".button { color: red; background: blue; }",
			want:  ".button { color: red; background: blue; }",
		},
		{
			name:  "should preserve legitimate url() with safe schemes",
			input: "background: url('https://example.com/image.png');",
			want:  "background: url('https://example.com/image.png');",
		},
		{
			name:  "should preserve http url scheme",
			input: "background: url(http://example.com/bg.jpg);",
			want:  "background: url(http://example.com/bg.jpg);",
		},

		// @import Attack Vectors
		{
			name:  "should remove @import rule with url",
			input: "@import url('http://evil.com/malicious.css');",
			want:  "",
		},
		{
			name:  "should remove @import rule with double quotes",
			input: `@import "http://evil.com/malicious.css";`,
			want:  "",
		},
		{
			name:  "should remove @import rule without quotes",
			input: "@import http://evil.com/malicious.css;",
			want:  "",
		},
		{
			name:  "should remove @import with case variation (uppercase)",
			input: "@IMPORT url('http://evil.com/malicious.css');",
			want:  "",
		},
		{
			name:  "should remove @import with mixed case",
			input: "@ImPoRt url('http://evil.com/malicious.css');",
			want:  "",
		},
		{
			name:  "should remove @import without semicolon",
			input: "@import url('http://evil.com/malicious.css')",
			want:  "",
		},
		{
			name:  "should remove multiple @import rules",
			input: "@import url('a.css'); @import url('b.css');",
			want:  " ",
		},

		// CSS Expression Attacks (IE-specific)
		{
			name:  "should remove expression() with javascript",
			input: "width: expression(alert('XSS'));",
			want:  "width: );",
		},
		{
			name:  "should remove expression with case variation (uppercase)",
			input: "width: EXPRESSION(alert('XSS'));",
			want:  "width: );",
		},
		{
			name:  "should remove expression with mixed case",
			input: "width: ExPrEsSiOn(alert('XSS'));",
			want:  "width: );",
		},
		{
			name:  "should remove expression with extra whitespace",
			input: "width: expression  (  alert('XSS')  );",
			want:  "width:   );",
		},
		{
			name:  "should remove nested expressions",
			input: "width: expression(expression(alert('XSS')));",
			want:  "width: ));",
		},

		// Dangerous URL Scheme Attacks
		{
			name:  "should remove javascript: url scheme",
			input: "background: url('javascript:alert(\"XSS\")');",
			want:  "background: url(alert(\"XSS\")');",
		},
		{
			name:  "should remove javascript: url scheme without quotes",
			input: "background: url(javascript:alert('XSS'));",
			want:  "background: url(alert('XSS'));",
		},
		{
			name:  "should remove javascript: with case variation (uppercase)",
			input: "background: url('JAVASCRIPT:alert(1)');",
			want:  "background: url(alert(1)');",
		},
		{
			name:  "should remove javascript: with mixed case",
			input: "background: url('JaVaScRiPt:alert(1)');",
			want:  "background: url(alert(1)');",
		},
		{
			name:  "should remove javascript: with whitespace before colon",
			input: "background: url('javascript :alert(1)');",
			want:  "background: url('javascript :alert(1)');",
		},
		{
			name:  "should remove data: url scheme",
			input: "background: url('data:text/html,<script>alert(1)</script>');",
			want:  "background: url(text/html,alert(1)');",
		},
		{
			name:  "should remove data: url with base64",
			input: "background: url('data:image/svg+xml;base64,PHN2Zz48c2NyaXB0PmFsZXJ0KDEpPC9zY3JpcHQ+PC9zdmc+');",
			want:  "background: url(image/svg+xml;base64,PHN2Zz48c2NyaXB0PmFsZXJ0KDEpPC9zY3JpcHQ+PC9zdmc+');",
		},
		{
			name:  "should remove vbscript: url scheme",
			input: "background: url('vbscript:msgbox(\"XSS\")');",
			want:  "background: url(msgbox(\"XSS\")');",
		},
		{
			name:  "should remove vbscript: with case variation",
			input: "background: url('VBSCRIPT:msgbox(1)');",
			want:  "background: url(msgbox(1)');",
		},
		{
			name:  "should remove url with extra whitespace around scheme",
			input: "background: url(  '  javascript  :alert(1)');",
			want:  "background: url(  '  javascript  :alert(1)');",
		},

		// HTML Tag Injection Attacks
		{
			name:  "should remove HTML script tag",
			input: ".class { <script>alert('XSS')</script> }",
			want:  ".class { alert('XSS') }",
		},
		{
			name:  "should remove HTML img tag with onerror",
			input: ".class { <img src=x onerror=alert(1)> }",
			want:  ".class {  }",
		},
		{
			name:  "should remove malformed HTML tags with extra spaces",
			input: ".class { <  script  >alert(1)<  /script  > }",
			want:  ".class { alert(1) }",
		},
		{
			name:  "should remove multiple HTML tags",
			input: "<div><script>alert(1)</script></div>",
			want:  "alert(1)",
		},
		{
			name:  "should remove HTML tags with attributes",
			input: `<div class="evil" onclick="alert(1)">CSS</div>`,
			want:  `CSS`,
		},

		// Style Breaking Sequences
		{
			name:  "should remove </style> closing tag",
			input: ".class { content: '</style><script>alert(1)</script>'; }",
			want:  ".class { content: 'alert(1)'; }",
		},
		{
			name:  "should remove </style> with extra spaces",
			input: ".class { content: '</  style  >'; }",
			want:  ".class { content: ''; }",
		},
		{
			name:  "should remove </style> with case variation",
			input: ".class { content: '</STYLE>'; }",
			want:  ".class { content: ''; }",
		},
		{
			name:  "should remove </style> with mixed case",
			input: ".class { content: '</StYlE>'; }",
			want:  ".class { content: ''; }",
		},
		{
			name:  "should remove </style> with attributes",
			input: ".class { content: '</style type=\"text/css\">'; }",
			want:  ".class { content: ''; }",
		},

		// Event Handler Attacks
		{
			name:  "should remove onclick event handler",
			input: ".class { onclick=alert(1); }",
			want:  ".class { alert(1); }",
		},
		{
			name:  "should remove onload event handler",
			input: ".class { onload=alert(1); }",
			want:  ".class { alert(1); }",
		},
		{
			name:  "should remove onerror event handler",
			input: ".class { onerror=alert(1); }",
			want:  ".class { alert(1); }",
		},
		{
			name:  "should remove onmouseover event handler",
			input: ".class { onmouseover=alert(1); }",
			want:  ".class { alert(1); }",
		},
		{
			name:  "should remove event handlers with case variation",
			input: ".class { ONCLICK=alert(1); }",
			want:  ".class { alert(1); }",
		},
		{
			name:  "should remove event handlers with whitespace",
			input: ".class { onclick  =  alert(1); }",
			want:  ".class {   alert(1); }",
		},
		{
			name:  "should remove multiple event handlers",
			input: ".class { onclick=a(); onerror=b(); }",
			want:  ".class { a(); b(); }",
		},

		// HTML Entity Attacks
		{
			name:  "should remove numeric HTML entity",
			input: ".class { content: '&#60;script&#62;'; }",
			want:  ".class { content: 'script'; }",
		},
		{
			name:  "should remove named HTML entity",
			input: ".class { content: '&lt;script&gt;'; }",
			want:  ".class { content: 'script'; }",
		},
		{
			name:  "should remove hex HTML entity",
			input: ".class { content: '&#x3c;script&#x3e;'; }",
			want:  ".class { content: 'script'; }",
		},
		{
			name:  "should remove multiple HTML entities",
			input: ".class { content: '&lt;&gt;&amp;&quot;'; }",
			want:  ".class { content: ''; }",
		},

		// Combined Attack Vectors
		{
			name:  "should handle multiple attack vectors in one input",
			input: `@import url('evil.css'); <script>alert(1)</script> width: expression(alert(2)); background: url('javascript:alert(3)');`,
			want:  ` alert(1) width: ); background: url(alert(3)');`,
		},
		{
			name:  "should handle nested attack patterns",
			input: `<div onclick=alert(1)><style>@import url('javascript:alert(2)');</style></div>`,
			want:  "",
		},
		{
			name:  "should handle complex real-world XSS attempt",
			input: `.evil {
				background: url('javascript:void(document.cookie)');
				width: expression(alert('XSS'));
				content: '</style><script>alert(1)</script>';
			}`,
			want: `.evil {
				background: url(void(document.cookie)');
				width: );
				content: 'alert(1)';
			}`,
		},
		{
			name:  "should sanitize multiple dangerous patterns with entities",
			input: `onclick=&lt;script&gt; background: url('data:text/html,<img src=x onerror=alert(1)>');`,
			want:  `script background: url(text/html,');`,
		},

		// Edge Cases and Bypass Attempts
		{
			name:  "should handle url() without space after url",
			input: "background: url('javascript:alert(1)');",
			want:  "background: url(alert(1)');",
		},
		{
			name:  "should handle url() with multiple spaces",
			input: "background: url   (   'javascript:alert(1)'   );",
			want:  "background: url(alert(1)'   );",
		},
		{
			name:  "should preserve legitimate CSS with similar keywords",
			input: ".important { font-weight: bold; }",
			want:  ".important { font-weight: bold; }",
		},
		{
			name:  "should handle CSS with newlines and tabs",
			input: ".class {\n\tbackground: url('javascript:alert(1)');\n}",
			want:  ".class {\n\tbackground: url(alert(1)');\n}",
		},
		{
			name:  "should handle empty url()",
			input: "background: url();",
			want:  "background: url();",
		},
		{
			name:  "should handle legitimate CSS custom properties",
			input: ":root { --main-color: #06c; --accent-color: #006; }",
			want:  ":root { --main-color: #06c; --accent-color: #006; }",
		},
		{
			name:  "should handle CSS with calc() function",
			input: "width: calc(100% - 20px);",
			want:  "width: calc(100% - 20px);",
		},
		{
			name:  "should handle CSS with rgb() function",
			input: "color: rgb(255, 0, 0);",
			want:  "color: rgb(255, 0, 0);",
		},
		{
			name:  "should handle legitimate @media rules",
			input: "@media screen and (max-width: 600px) { .class { color: red; } }",
			want:  "@media screen and (max-width: 600px) { .class { color: red; } }",
		},
		{
			name:  "should handle CSS comments (not removed by sanitizer)",
			input: "/* This is a comment */ .class { color: red; }",
			want:  "/* This is a comment */ .class { color: red; }",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := sanitizer.CSS(tc.input)

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCSS_MultiPassSanitization(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "should handle obfuscated javascript url with entities",
			input: `background: url('&#106;&#97;&#118;&#97;&#115;&#99;&#114;&#105;&#112;&#116;&#58;alert(1)');`,
			want:  "background: url('alert(1)');",
		},
		{
			name: "should handle expression within imported CSS",
			input: `@import url('data:text/css,*{width:expression(alert(1))}');`,
			want:  "",
		},
		{
			name:  "should handle style tag breaking with event handler",
			input: `content: '</style><img src=x onerror=alert(1)>`,
			want:  "content: '",
		},
		{
			name: "should handle all attack vectors combined in complex CSS",
			input: `
				@import url('javascript:void(0)');
				.class1 { width: expression(alert(1)); }
				.class2 { background: url('data:text/html,<script>alert(2)</script>'); }
				<script>alert(3)</script>
				.class3 { content: '</style><script>alert(4)</script>'; }
				.class4 onclick=alert(5) { color: red; }
				.class5 { content: '&lt;script&gt;alert(6)&lt;/script&gt;'; }
			`,
			want: `

				.class1 { width: ); }
				.class2 { background: url(text/html,alert(2)'); }
				alert(3)
				.class3 { content: 'alert(4)'; }
				.class4 alert(5) { color: red; }
				.class5 { content: 'scriptalert(6)/script'; }
			`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := sanitizer.CSS(tc.input)

			assert.Equal(t, tc.want, got)
		})
	}
}

// BenchmarkCSS measures the performance of CSS sanitization
func BenchmarkCSS(b *testing.B) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "clean CSS",
			input: ".button { color: red; background: blue; padding: 10px; }",
		},
		{
			name:  "single attack vector",
			input: "background: url('javascript:alert(1)');",
		},
		{
			name: "multiple attack vectors",
			input: `@import url('evil.css');
				width: expression(alert(1));
				background: url('javascript:alert(2)');
				<script>alert(3)</script>`,
		},
		{
			name: "complex real-world CSS",
			input: `
				.header { background: url('https://example.com/bg.jpg'); }
				.button { color: #fff; padding: 10px 20px; }
				.card { box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
				@media (max-width: 768px) { .header { font-size: 14px; } }
			`,
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = sanitizer.CSS(tc.input)
			}
		})
	}
}
