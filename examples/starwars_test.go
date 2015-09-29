package examples

import (
	"bytes"
	"github.com/theplant/graphql/examples/starwars"
	"github.com/theplant/graphql/graphql"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStarwars(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	starwarCasesPath := filepath.Join(gopath, "/src/github.com/theplant/graphql/examples/starwars/cases/*.gql")
	gqls, err := filepath.Glob(starwarCasesPath)
	if err != nil {
		panic(err)
	}

	on := starwars.DefaultQuery
	for _, gqlFilePath := range gqls {
		json := bytes.NewBuffer(nil)

		gqlText := text(gqlFilePath)
		err = graphql.ExecuteQL(on, gqlText, json)

		if err != nil {
			panic(err)
		}

		expectedJsonText := text(strings.Replace(gqlFilePath, ".gql", ".json", 1))
		segs := strings.Split(gqlFilePath, "/")

		if json.String() != expectedJsonText {
			t.Errorf("File: %s, Query:\n%s\nResult is: %s\nResult should: %s", segs[len(segs)-1], gqlText, json.String(), string(expectedJsonText))
		}
	}
}

func text(path string) string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	textBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return string(textBytes)
}
