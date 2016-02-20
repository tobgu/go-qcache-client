package qclient

import "testing"

// Emacs
// -----
// godef describe - C-c C-d
// godef jump - M-.   (own binding)


// Command line
// ------------
// go test
// go test github.com/user/stringutil  - Run relative to GOPATH
// go install  - Build and install result in $GOPATH/bin/
// go build    - Build package but don't produce any output
// go run      - Build and run

func TestQuery(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Hello, world", "hello"},
	}

	var client = New([]string{""})
	for _, c := range cases {
		got := client.Query(c.in, c.in)
		if got != c.want {
			t.Errorf("Query(%q, %q) == %q, want %q", c.in, c.in, got, c.want)
		}
	}
}


func TestGetWithNoMatchingKey(t *testing.T) {
	var client = New([]string{"http://localhost:9401"})
	result := client.Get("foo", Query{Where: []interface{}{"foo", []interface{}{"+", "bar", 1}}})
	if result != "hello" {
		t.Errorf("Failed!")
	}
}
