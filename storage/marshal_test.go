package storage

import "testing"

func TestParseEndpoint(t *testing.T) {
	dataAndExpected := [][]string{
		[]string{"foo:bar", "foo", "bar"},
		[]string{"foo:", "foo", ""},
		[]string{":bar", "", "bar"},
		[]string{"hoge.hoge:fuga/fuga", "hoge.hoge", "fuga/fuga"},
		[]string{":", "", ""},
	}
	for _, row := range dataAndExpected {
		info := EndpointInfo(row[0])
		if string(info.Joint()) != row[1] {
			t.Errorf("Joint not match: %s != %s (in %s)", info.Joint(), row[1], row[0])
		}

		if string(info.Port()) != row[2] {
			t.Errorf("Port not match: %s != %s (in %s)", info.Port(), row[2], row[0])
		}
	}
}


