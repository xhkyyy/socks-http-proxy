package main

import (
	"reflect"
	"testing"
)

func TestLoadConf(t *testing.T) {
	type args struct {
		filePath string
	}

	var confInfo = Conf{DomainFile: "rule file", Addr: "addr string", Socks5Addr: "socks5 addr"}

	tests := []struct {
		name string
		args args
		want *Conf
	}{
		{"t1", args{"./cnf.json.test"}, &confInfo},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadConf(tt.args.filePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConf() = %v, want %v", got, tt.want)
			}
		})
	}
}
