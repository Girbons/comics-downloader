package sites

import "testing"

func TestDeobfuscateURL(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{`https:\/\/foo.com\/img%3Fid%3D1`, "https://foo.com/img?id=1"},
		{"https://foo.com/img&#61;2&#38;token%3Dabc%252F123", "https://foo.com/img=2&token=abc%2F123"},
		{"", ""},
	}

	for _, tc := range cases {
		actual := deobfuscateURL(tc.input)
		if actual != tc.expected {
			t.Fatalf("deobfuscateURL(%q) = %q, expected %q", tc.input, actual, tc.expected)
		}
	}
}
