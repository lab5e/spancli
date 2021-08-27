package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/go-spanapi/v4/apitools"
)

type collectionCmd struct {
	Add    addCollection    `command:"add" description:"create collection"`
	Get    getCollection    `command:"get" description:"get collection"`
	List   listCollection   `command:"list" alias:"ls" description:"list collections"`
	Delete deleteCollection `command:"delete" alias:"del" description:"delete collection"`
}

type addCollection struct {
	TeamID string `long:"team-id" description:"team the collection belongs to" required:"yes"`
	Name   string `long:"name" description:"name of the collection" required:"yes"`
}

type getCollection struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
}

type listCollection struct {
	Format   string `long:"format" default:"text" description:"which output format to use" choice:"csv" choice:"html" choice:"markdown" choice:"text" choice:"json"`
	NoColor  bool   `long:"no-color" env:"SPAN_NO_COLOR" description:"turn off coloring"`
	PageSize int    `long:"page-size" description:"if set, chop output into pages of page-size length"`
}

type deleteCollection struct {
	CollectionID string `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	YesIAmSure   bool   `long:"yes-i-am-sure" description:"disable prompt for 'are you sure'"`
}

func (r *listCollection) Execute([]string) error {
	client, ctx, done := newClient()
	defer done()

	collections, _, err := client.CollectionsApi.ListCollections(ctx).Execute()
	if err != nil {
		return err
	}

	// treat JSON formatting as special case that dumps all data
	if r.Format == "json" {
		json, err := json.MarshalIndent(collections, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", json)
		return nil
	}

	t := newTableOutput(r.Format, r.NoColor, r.PageSize)
	t.SetTitle("Collections")
	t.AppendHeader(table.Row{"ID", "Name", "TeamID"})
	for _, col := range *collections.Collections {
		t.AppendRow([]interface{}{*col.CollectionId, col.GetTags()["name"], *col.TeamId})
	}
	renderTable(t, r.Format)
	return nil
}

func newClient() (*spanapi.APIClient, context.Context, context.CancelFunc) {
	config := spanapi.NewConfiguration()
	config.Debug = opt.Debug

	ctx, done := apitools.ContextWithAuthAndTimeout(opt.Token, opt.Timeout)

	return spanapi.NewAPIClient(config), ctx, done
}
