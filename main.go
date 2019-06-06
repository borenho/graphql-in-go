package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

// Types
type Tutorial struct {
	ID       int
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

var tutorials []Tutorial

// Function to return an array of tutorials
func populate() []Tutorial {
	author := &Author{Name: "Kevin Terah", Tutorials: []int{1, 2}}
	tutorial := Tutorial{
		ID:     1,
		Title:  "Environmental Hygiene",
		Author: *author,
		Comments: []Comment{
			Comment{Body: "Climate change is ..."},
		},
	}

	tutorial2 := Tutorial{
		ID:     2,
		Title:  "Environmental Cost",
		Author: *author,
		Comments: []Comment{
			Comment{Body: "Weather is good ..."},
		},
	}

	//	var tutorials []Tutorial

	tutorials = append(tutorials, tutorial)
	tutorials = append(tutorials, tutorial2)

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

var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"create": &graphql.Field{
			Type:        tutorialType,
			Description: "Create a new tutorial",
			Args: graphql.FieldConfigArgument{
				"title": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				tutorial := Tutorial{Title: params.Args["title"].(string)}
				tutorials = append(tutorials, tutorial)

				return tutorial, nil
			},
		},
	},
})

func main() {
	tutorials = populate()

	// Schema
	fields := graphql.Fields{
		"tutorial": &graphql.Field{
			Type:        tutorialType,
			Description: "Get a tutorial by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Take in the ID argument
				id, ok := p.Args["id"].(int)
				if ok {
					// Parse tutorial array for matching ID
					for _, tutorial := range tutorials {
						if int(tutorial.ID) == id {
							return tutorial, nil
						}
					}
				}

				return nil, nil
			},
		},

		// list endpoint to returm all tutorials
		"list": &graphql.Field{
			Type:        graphql.NewList(tutorialType),
			Description: "Get Tutorial List",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return tutorials, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(rootQuery),
		Mutation: mutationType,
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Failed to create new schema, error: %v", err)
	}

	// Query
	query := `
		mutation {
			 create(title: "Kele Ne Oo") {
				title
			}
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("Failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)
}
