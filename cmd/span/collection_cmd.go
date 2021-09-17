package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/lab5e/spanclient-go/v4"
)

type collectionCmd struct {
	Add    addCollection    `command:"add" description:"create collection"`
	Get    getCollection    `command:"get" description:"get collection"`
	List   listCollection   `command:"list" alias:"ls" description:"list collections"`
	Delete deleteCollection `command:"delete" alias:"del" description:"delete collection"`
	Update updateCollection `command:"update" alias:"up" description:"update collection"`
}

type addCollection struct {
	TeamID string `long:"team-id" description:"team the collection belongs to" required:"yes"`
	Name   string `long:"name" description:"name of the collection" required:"yes"`
}

type getCollection struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
}

type listCollection struct{}

type deleteCollection struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	YesIAmSure   bool   `long:"yes-i-am-sure" description:"disable prompt for 'are you sure'"`
}

type updateCollection struct {
	CollectionID string   `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	Tags         []string `long:"tag" short:"t" description:"Set tag value (name:value)"`
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
	fmt.Fprintf(w, strings.Join([]string{"CollectionID", "TeamID", "FW Mgmt", "Tags"}, "\t")+"\n")

	for _, col := range collections.Collections {
		fmt.Fprintf(w, strings.Join([]string{
			col.CollectionId,
			col.TeamId,
			string(col.Firmware.Management),
			tagsToString(col.Tags),
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

func (u *updateCollection) Execute([]string) error {
	client := spanclient.NewAPIClient(clientConfig())
	ctx, cancel := spanContext()
	defer cancel()

	collection, resp, err := client.CollectionsApi.RetrieveCollection(ctx, u.CollectionID)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return errors.New("unknown collection")
		}
		return err
	}
	if collection.Tags == nil {
		collection.Tags = make(map[string]string)
	}
	for _, val := range u.Tags {
		nameValue := strings.Split(val, ":")
		if len(nameValue) != 2 {
			return errors.New("tag name incorrectly formatted (needs name:value)")
		}
		collection.Tags[nameValue[0]] = nameValue[1]
	}
	_, _, err = client.CollectionsApi.UpdateCollection(ctx, u.CollectionID, collection)
	if err != nil {
		return err
	}
	fmt.Printf("updated collection '%s'\n", u.CollectionID)
	return nil
}
