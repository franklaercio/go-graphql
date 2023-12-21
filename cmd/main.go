package main

import (
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"log"
)

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

func populate() []Tutorial {
	author := &Author{Name: "Frank Júnior", Tutorials: []int{1}}
	tutorial := Tutorial{
		ID:     1,
		Title:  "Go Graphql",
		Author: *author,
		Comments: []Comment{
			{Body: "Nice!"},
		},
	}

	var tutorials []Tutorial
	tutorials = append(tutorials, tutorial)
	return tutorials
}

func main() {
	fmt.Println("Graphql tutorial")

	tutorials := populate()

	var commentType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Comment",
			Fields: graphql.Fields{
				"Body": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	var authorType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Author",
			Fields: graphql.Fields{
				"Name":      &graphql.Field{Type: graphql.String},
				"Tutorials": &graphql.Field{Type: graphql.NewList(graphql.Int)},
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

	fields := graphql.Fields{
		"tutorial": &graphql.Field{
			Type:        tutorialType,
			Description: "Get Tutorial by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for _, tutorial := range tutorials {
						if tutorial.ID == id {
							return tutorial, nil
						}
					}
				}
				return nil, nil
			},
		},
		"list": &graphql.Field{
			Type:        graphql.NewList(tutorialType),
			Description: "Get Full Tutorial List",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return tutorials, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}

	schema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		log.Fatalf("Failed to create a new Graphql schema, err %v", err)
	}

	query := `
		{
			tutorial(id:1) {
				title
				author {
					Name
					Tutorials
				}
			}
		}
	`

	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)

	if len(r.Errors) > 0 {
		log.Fatalf("Failed to execute Graphql operation, erros %v", r.Errors)
	}

	rJSON, _ := json.Marshal(r)

	fmt.Printf("%s \n", rJSON)
}
