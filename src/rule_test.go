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
