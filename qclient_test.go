package qclient

import ("testing"
	"encoding/csv"
	"encoding/json"
	"bytes"
	"log"
	"os/exec"
	"time"
)

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

	var client = NewClient([]string{""})
	for _, c := range cases {
		got := client.Query(c.in, c.in)
		if got != c.want {
			t.Errorf("Query(%q, %q) == %q, want %q", c.in, c.in, got, c.want)
		}
	}
}


func init() {
	cmd := exec.Command("docker", "kill", "qcache-go-test")
	cmd.Run()
	
     	cmd = exec.Command(
		"docker", "run", "--rm", "-p", "9401:9401", "--name", "qcache-go-test", "tobgu/qcache")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	// Give QCache some time to start
	time.Sleep(1 * time.Second)
}

func TestGetPostGet(t *testing.T) {
	client := NewClient([]string{"http://localhost:9401"})
	key := "baz"
	query := Query{Where: And(Eq("bar", 10))}
	result, _ := client.Get(key, query)
	if result != nil {
		t.Errorf("Did not expect any result before inserting data!")
	}

	// Prepare and post CSV	
	records := [][]string{
		{"foo", "bar"},
		{"x", "1"},
		{"y", "10"},
		{"z", "100"},
	}

	var buf bytes.Buffer;
	w := csv.NewWriter(&buf)
	w.WriteAll(records) // calls Flush internally

	if err := w.Error(); err != nil {
		log.Fatalln("Error writing csv:", err)
	}

	err := client.Post(key, "text/csv", &buf)
	if err != nil {
		t.Errorf("Did not expect any result before inserting data!")
	}

	result, _ = client.Get(key, query)
	if result == nil {
		t.Errorf("Expected result after inserting data!")
	}

	var f []interface{}
	json.Unmarshal(result, &f)

	r := f[0].(map[string]interface{})
	if !(r["foo"].(string) == "y" && r["bar"].(float64) == 10) {
		t.Errorf("Unexpected result", r)
	}
}
