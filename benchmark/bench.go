package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
)

/*
input

algorithm name / regular expression
number of iterations / time
count
input size
generate cpu profile
generate memory profile
number of threads

*/

type Input struct {
	Algorithm  string `json:"algorithm"`
	Iterations string `json:"iterations"`
	Time       string `json:"time"`
	Size       string `json:"size"`
	Threads    string `json:"threads"`
}

type Output struct {
	Os           string `json:"os"`
	Architecture string `json:"architecture"`
	Cpu          string `json:"cpu"`
	Iterations   string `json:"iterations"`
	Time         string `json:"time"`
	Average      string `json:"average"`
	Allocations  string `json:"allocations"`
	Bytes        string `json:"bytes"`
}

func main() {

	var i Input

	err := json.NewDecoder(os.Stdin).Decode(&i)
	checkError(err)

	args := []string{
		"test",
		"-bench=" + i.Algorithm + "/" + i.Size + "$",
		"-benchmem",
		"-cpu=" + i.Threads,
	}

	var benchtime string
	if i.Time == "" {
		benchtime = i.Iterations + "x"
	} else {
		benchtime = i.Time + "s"
	}

	args = append(args, "-benchtime="+benchtime)

	// execute the benchmark with the given settings
	data, err := exec.Command("go", args...).Output()
	checkError(err)

	raw := string(data)

	var o Output

	o.Os = find(raw, "goos: (.*)")
	o.Architecture = find(raw, "goarch: (.*)")
	o.Cpu = find(raw, "cpu: (.*)")
	o.Time = find(raw, "ok\\s*\\S*\\s*(\\S*)s")
	o.Allocations = find(raw, "(\\S*) allocs/op")
	o.Bytes = find(raw, "(\\S*) B/op")
	o.Average = find(raw, "(\\S*) ns/op")
	o.Iterations = find(raw, "(\\S*)\\s*\\S* ns/op")

	output, err := json.MarshalIndent(o, "", "  ")
	checkError(err)

	fmt.Printf("%s\n", output)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func find(target, pattern string) string {
	re := regexp.MustCompile(pattern)
	return re.FindStringSubmatch(target)[1]
}
