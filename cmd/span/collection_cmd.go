package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/helpers"
)

type collectionCmd struct {
	Add    addCollection    `command:"add" description:"create collection"`
	Get    getCollection    `command:"get" description:"get collection"`
	List   listCollection   `command:"list" alias:"ls" description:"list collections"`
	Delete deleteCollection `command:"delete" alias:"del" description:"delete collection"`
	Update updateCollection `command:"update" alias:"up" description:"update collection"`
}

type addCollection struct {
	TeamID string   `long:"team-id" description:"team the collection belongs to" required:"yes"`
	Tags   []string `long:"tag" short:"t" description:"Set tag value (name:value)"`
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
	tags, err := tagsToMap(r.Tags)
	if err != nil {
		return err
	}
	client := spanapi.NewAPIClient(clientConfig())
	ctx, cancel := spanContext()
	defer cancel()

	helpers.CheckVersion(ctx, client)

	collection, _, err := client.CollectionsApi.CreateCollectionExecute(spanapi.ApiCreateCollectionRequest{}.Body(
		spanapi.Collection{
			TeamId: &r.TeamID,
			Tags:   &tags,
		},
	))
	if err != nil {
		return err
	}

	fmt.Printf("created collection with id '%s'\n", *collection.CollectionId)
	return nil
}

func (r *getCollection) Execute([]string) error {
	client := spanapi.NewAPIClient(clientConfig())
	ctx, cancel := spanContext()
	defer cancel()

	helpers.CheckVersion(ctx, client)

	collection, _, err := client.CollectionsApi.RetrieveCollection(ctx, r.CollectionID).Execute()
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
	client := spanapi.NewAPIClient(clientConfig())
	ctx, cancel := spanContext()
	defer cancel()

	helpers.CheckVersion(ctx, client)

	collections, _, err := client.CollectionsApi.ListCollections(ctx).Execute()
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintf(w, strings.Join([]string{"CollectionID", "TeamID", "FW Mgmt", "Tags"}, "\t")+"\n")

	for _, col := range *collections.Collections {
		fmt.Fprintf(w, strings.Join([]string{
			strPtr(col.CollectionId),
			strPtr(col.TeamId),
			strPtr((*string)(col.Firmware.Management)),
			tagsToString(*col.Tags),
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

	client := spanapi.NewAPIClient(clientConfig())
	ctx, cancel := spanContext()
	defer cancel()

	helpers.CheckVersion(ctx, client)

	_, _, err := client.CollectionsApi.DeleteCollection(ctx, r.CollectionID).Execute()
	if err != nil {
		return err
	}

	fmt.Printf("deleted collection '%s'\n", r.CollectionID)
	return nil
}

func (u *updateCollection) Execute([]string) error {
	client := spanapi.NewAPIClient(clientConfig())
	ctx, cancel := spanContext()
	defer cancel()

	helpers.CheckVersion(ctx, client)

	collection, resp, err := client.CollectionsApi.RetrieveCollection(ctx, u.CollectionID).Execute()
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return errors.New("unknown collection")
		}
		return err
	}
	if collection.Tags == nil {
		tags := make(map[string]string)
		collection.Tags = &tags
	}
	t := *collection.Tags
	for _, val := range u.Tags {
		nameValue := strings.Split(val, ":")
		if len(nameValue) != 2 {
			return errors.New("tag name incorrectly formatted (needs name:value)")
		}
		t[nameValue[0]] = nameValue[1]
	}
	_, _, err = client.CollectionsApi.UpdateCollection(ctx, u.CollectionID).Body(collection).Execute()
	if err != nil {
		return err
	}
	fmt.Printf("updated collection '%s'\n", u.CollectionID)
	return nil
}
