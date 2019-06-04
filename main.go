package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

// Types
type Tutorial struct {
	Title    string
	Author   Author
	Comments []Comment
}

type Author struct {
	Name      string
	Tutorials []int
}

type Comment struct {
	Body string
}

// Function to return an array of tutorials
func populate() []Tutorial {
	author := &Author{Name: "Kevin Terah", Tutorials: []int{1}}
	tutorial := Tutorial{
		ID:     1,
		Title:  "Environmental Hygiene",
		Author: *author,
		Comments: []Comment{
			Comment{Body: "Climate change is ..."},
		},
	}

	var tutorials []Tutorial
	tutorials = append(tutorials, tutorial)

	return tutorials
}

// Objects
var commentType = graphql.NewObject(
	graphql.ObjectConfig{
		// Define the name and fields of the object
		Name: "Comment",
		Fields: graphql.Fields{
			"body": &graphql.Field{Type: graphql.String},
		},
	},
)

var authorType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Author",
		Fields: graphql.Fields{
			"Name": &graphql.Field{Type: graphql.String},
			"Tutorials": &graphql.Field{
				// Use NewList to hold an array of integers
				Type: graphql.NewList(graphql.Int),
			},
		},
	},
)

var tutorialType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Tutorial",
		Fields: graphql.Fields{
			"id":       &graphql.Field{Type: graphql.Int},
			"title":    &graphql.Field{Type: graphql.String},
			"author":   &graphql.Field{Type: authorType},
			"comments": &graphql.Field{Type: graphql.NewList(commentType)},
		},
	},
)

func main() {
	// Schema
	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Failed to create new schema, error: %v", err)
	}

	// Query
	query := `
		{
			hello
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("Failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON) // {"data": {"hello" : "world"}}
}
