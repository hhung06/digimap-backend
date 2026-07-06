package cdn

import "testing"

func TestEncodeInvalidationPath(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "safe ascii path is unchanged",
			in:   "/production/base/public/019f26c0-0b6e-7ed4-a30b-c125e977325d.digimap",
			want: "/production/base/public/019f26c0-0b6e-7ed4-a30b-c125e977325d.digimap",
		},
		{
			name: "theme name with spaces is encoded",
			in:   "/production/venue_themes/custom/019f26c0-0b6e-7ed4-a30b-c125e977325d/Summer Theme.json",
			want: "/production/venue_themes/custom/019f26c0-0b6e-7ed4-a30b-c125e977325d/Summer%20Theme.json",
		},
		{
			name: "theme name with unicode is encoded",
			in:   "/production/venue_themes/custom/019f26c0-0b6e-7ed4-a30b-c125e977325d/夏祭りテーマ.json",
			want: "/production/venue_themes/custom/019f26c0-0b6e-7ed4-a30b-c125e977325d/%E5%A4%8F%E7%A5%AD%E3%82%8A%E3%83%86%E3%83%BC%E3%83%9E.json",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := encodeInvalidationPath(tc.in)
			if got != tc.want {
				t.Errorf("encodeInvalidationPath(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
