package main

import (
	"testing"
)

func init() {
	rule = map[string]bool{
		"google.com":  true,
		"twitter.com": true,
		"abc.com":     true,
	}
}

func TestUrlMatcheAny(t *testing.T) {
	type args struct {
		urlStr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"ruleTest", args{urlStr: "www.google.com"}, true},
		{"ruleTest", args{urlStr: "google.com"}, true},
		{"ruleTest", args{urlStr: "google.cn"}, false},
		{"ruleTest", args{urlStr: "twitter.com"}, true},
		{"ruleTest", args{urlStr: "www.abc.com"}, true},
		{"ruleTest", args{urlStr: "www.xxx.abc.com"}, true},
		{"ruleTest", args{urlStr: "xxx.www.xxx.abc.com"}, true},
		{"ruleTest", args{urlStr: "abc.com"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UrlMatcheHost(tt.args.urlStr); got != tt.want {
				t.Errorf("UrlMatcheHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkUrlMatcheAny(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UrlMatcheHost("www.abc.com:443")
	}
}

func TestGetSubDomains(t *testing.T) {
	want := []string{"xx.com"}
	if got := getSubDomains("www.xx.com"); !eqSlice(want, got) {
		t.Errorf("getSubDomains(www.xx.com) = %q, want %q", got, want)
	}

	want = []string{"cd.123.xx.com", "123.xx.com", "xx.com"}

	if got := getSubDomains("ab.cd.123.xx.com"); !eqSlice(want, got) {
		t.Errorf("getSubDomains(ab.cd.123.xx.com) = %q, want %q", got, want)
	}

	want = nil

	if got := getSubDomains("xx.com"); !eqSlice(want, got) {
		t.Errorf("getSubDomains(xx.com) = %q, want %q", got, want)
	}

	want = nil

	if got := getSubDomains(""); !eqSlice(want, got) {
		t.Errorf("getSubDomains('') = %q, want %q", got, want)
	}

	want = nil
	var param string
	if got := getSubDomains(param); !eqSlice(want, got) {
		t.Errorf("getSubDomains('') = %q, want %q", got, want)
	}
}

func eqSlice(a, b []string) bool {

	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
