package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/lab5e/spanclient-go/v4"
)

type collectionCmd struct {
	Add    addCollection    `command:"add" description:"create collection"`
	Get    getCollection    `command:"get" description:"get collection"`
	List   listCollection   `command:"list" description:"list collections"`
	Delete deleteCollection `command:"delete" description:"delete collection"`
}

type addCollection struct {
	TeamID string `long:"team-id" description:"team the collection belongs to" required:"yes"`
	Name   string `long:"name" description:"name of the collection" required:"yes"`
}

type getCollection struct {
	CollectionID string `long:"collection-id" description:"collection id" required:"yes"`
}

type listCollection struct{}

type deleteCollection struct {
	CollectionID string `long:"collection-id" description:"collection id" required:"yes"`
	YesIAmSure   bool   `long:"yes-i-am-sure" description:"disable prompt for 'are you sure'"`
}

func init() {
	parser.AddCommand("collection", "collection management commands", "collection management commands", &collectionCmd{})
}

func (r *addCollection) Execute([]string) error {
	client := spanclient.NewAPIClient(clientConfig())
	ctx, _ := spanContext()
	collection, _, err := client.CollectionsApi.CreateCollection(ctx, spanclient.Collection{
		TeamId: r.TeamID,
		FieldMask: spanclient.FieldMask{
			Imsi:     false,
			Imei:     false,
			Location: false,
		},
		Firmware: spanclient.CollectionFirmware{
			Management: spanclient.DISABLED,
		},
		Tags: map[string]string{"name": r.Name},
	})
	if err != nil {
		return err
	}

	fmt.Printf("created collection with id '%s'\n", collection.CollectionId)
	return nil
}

func (r *getCollection) Execute([]string) error {
	client := spanclient.NewAPIClient(clientConfig())
	ctx, _ := spanContext()
	collection, _, err := client.CollectionsApi.RetrieveCollection(ctx, r.CollectionID)
	if err != nil {
		return err
	}
	json, err := json.MarshalIndent(collection, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to marshal '%v' to JSON: %v", collection, err)
	}
	fmt.Printf("%s\n", json)
	return nil
}

func (r *listCollection) Execute([]string) error {
	client := spanclient.NewAPIClient(clientConfig())
	ctx, _ := spanContext()

	collections, _, err := client.CollectionsApi.ListCollections(ctx)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintf(w, strings.Join([]string{"CollectionID", "TeamID", "Name"}, "\t")+"\n")

	for _, col := range collections.Collections {
		fmt.Fprintf(w, strings.Join([]string{
			col.CollectionId,
			col.TeamId,
			col.Tags["name"],
		}, "\t")+"\n")
	}
	return w.Flush()
}

func (r *deleteCollection) Execute([]string) error {
	if !r.YesIAmSure {
		if !verifyDeleteIntent() {
			return fmt.Errorf("user aborted delete")
		}
	}

	client := spanclient.NewAPIClient(clientConfig())
	ctx, _ := spanContext()

	_, _, err := client.CollectionsApi.DeleteCollection(ctx, r.CollectionID)
	if err != nil {
		return err
	}

	fmt.Printf("deleted collection '%s'\n", r.CollectionID)
	return nil
}
