package main

import (
	"fmt"

	"github.com/alecthomas/jsonschema"
	"github.com/google/go-github/v35/github"
)

// create json schema from gitbhu.Issue struct
// used as return arguments in github/issuesearch endpoint

func main() {
	schema := jsonschema.Reflect(&github.Issue{})
	b, _ := schema.MarshalJSON()
	fmt.Println(string(b))
}
