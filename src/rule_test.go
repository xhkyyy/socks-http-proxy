package main

import (
	"testing"
)

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
		{"ruleTest", args{urlStr: "hwww.sogou.com"}, false},
		{"ruleTest", args{urlStr: "twitter.com"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UrlMatcheHost(tt.args.urlStr); got != tt.want {
				t.Errorf("UrlMatcheHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecondLastIndex(t *testing.T) {
	type args struct {
		s      string
		substr byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"ts", args{"www.abc.com", '.'}, 3},
		{"ts", args{"www.xyz.abc.com", '.'}, 7},
		{"ts", args{"abc.com", '.'}, -1},
		{"ts", args{"abc", '.'}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SecondLastIndex(tt.args.s, tt.args.substr); got != tt.want {
				t.Errorf("SecondLastIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkSecondLastIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SecondLastIndex("www.abc.com", '.')
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
