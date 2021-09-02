package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"regexp"
)

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
	// use flags to read input
	algorithm := flag.String("algorithm", "", "")
	iterations := flag.String("iterations", "", "")
	time := flag.String("time", "", "")
	size := flag.String("size", "", "")
	threads := flag.String("threads", "", "")

	flag.Parse()

	args := []string{
		"test",
		"-bench=" + *algorithm + "/" + *size + "$",
		"-benchmem",
		"-cpu=" + *threads,
	}

	var benchtime string
	if *time == "" {
		benchtime = *iterations + "x"
	} else {
		benchtime = *time + "s"
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
