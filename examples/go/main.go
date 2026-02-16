package main

import (
	"context"
	"fmt"
	"log"

	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
	"github.com/openfga/go-sdk/configuration"
)

func main() {
	ctx := context.Background()

	configuration, err := configuration.NewConfiguration(configuration.Configuration{
		ApiScheme: "http",
		ApiHost:   "localhost:8080",
	})
	if err != nil {
		log.Fatalf("Failed to create configuration: %v", err)
	}

	fgaClient, err := client.NewSdkClient(configuration)
	if err != nil {
		log.Fatalf("Failed to create OpenFGA client: %v", err)
	}

	store := createStore(ctx, fgaClient)
	createAuthorizationModel(ctx, fgaClient, store)
	createRelationships(ctx, fgaClient, store)
	checkAccess(ctx, fgaClient, store)
	listPermissions(ctx, fgaClient, store)
}

func createStore(ctx context.Context, fgaClient *client.SdkClient) string {
	response, err := fgaClient.CreateStore(ctx).Body(client.CreateStoreRequest{
		Name: "authorization-store",
	}).Execute()
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	fmt.Printf("Created store: %s\n", response.Id)
	return response.Id
}

func createAuthorizationModel(ctx context.Context, fgaClient *client.SdkClient, storeID string) {
	model := `
	model
	  schema 1.1

	type user

	type organization
	  relations
	    define admin: [user]
	    define member: [user]

	type project
	  relations
	    define organization: [organization]
	    define owner: [user]
	    define editor: [user]
	    define viewer: [user]
	`

	response, err := fgaClient.WriteAuthorizationModel(ctx).
		StoreId(storeID).
		Body(client.WriteAuthorizationModelRequest{
			SchemaVersion: "1.1",
			TypeDefinitions: []client.TypeDefinition{
				// Model definitions here
			},
		}).Execute()

	if err != nil {
		log.Fatalf("Failed to write authorization model: %v", err)
	}

	fmt.Printf("Authorization model ID: %s\n", response.AuthorizationModelId)
}

func createRelationships(ctx context.Context, fgaClient *client.SdkClient, storeID string) {
	writes := []client.TupleKey{
		{
			User:     "user:alice",
			Relation: "admin",
			Object:   "organization:acme",
		},
		{
			User:     "user:bob",
			Relation: "member",
			Object:   "organization:acme",
		},
		{
			User:     "organization:acme",
			Relation: "organization",
			Object:   "project:api",
		},
	}

	err := fgaClient.WriteTuples(ctx).
		StoreId(storeID).
		Body(client.WriteRequest{
			Writes: writes,
		}).Execute()

	if err != nil {
		log.Fatalf("Failed to write relationships: %v", err)
	}

	fmt.Println("Relationships created successfully")
}

func checkAccess(ctx context.Context, fgaClient *client.SdkClient, storeID string) {
	response, err := fgaClient.Check(ctx).
		StoreId(storeID).
		Body(client.CheckRequest{
			TupleKey: client.CheckRequestTupleKey{
				User:     "user:alice",
				Relation: "admin",
				Object:   "organization:acme",
			},
		}).Execute()

	if err != nil {
		log.Fatalf("Failed to check access: %v", err)
	}

	fmt.Printf("Alice is admin of acme: %v\n", response.Allowed)
}

func listPermissions(ctx context.Context, fgaClient *client.SdkClient, storeID string) {
	response, err := fgaClient.ListObjects(ctx).
		StoreId(storeID).
		Body(client.ListObjectsRequest{
			AuthorizationModelId: "",
			User:                 "user:alice",
			Relation:             "admin",
			Type:                 "organization",
		}).Execute()

	if err != nil {
		log.Fatalf("Failed to list objects: %v", err)
	}

	fmt.Printf("Alice can admin: %v\n", response.Objects)
}
