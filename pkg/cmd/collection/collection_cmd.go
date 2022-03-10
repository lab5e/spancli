package collection

import (
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/helpers"
)

type Command struct {
	Add    addCollection    `command:"add" description:"create new collection"`
	Get    getCollection    `command:"get" description:"get collection details"`
	List   listCollection   `command:"list" alias:"ls" description:"list collections"`
	Delete deleteCollection `command:"delete" alias:"del" description:"delete collection"`
	Update updateCollection `command:"update" alias:"up" description:"update collection"`
}

func (*Command) Execute([]string) error {
	return nil
}

type addCollection struct {
	TeamID string   `long:"team-id" description:"team the collection belongs to" required:"yes"`
	Tags   []string `long:"tag" description:"Set tag value (name:value)"`
	Name   string   `long:"name" description:"name of the collection"`
}

type getCollection struct {
	CollectionID string `long:"collection-id" description:"Span collection ID" required:"yes"`
}

type listCollection struct {
	Format   string `long:"format" default:"text" description:"which output format to use" choice:"csv" choice:"html" choice:"markdown" choice:"text" choice:"json"`
	NoColor  bool   `long:"no-color" env:"SPAN_NO_COLOR" description:"turn off coloring"`
	PageSize int    `long:"page-size" description:"if set, chop output into pages of page-size length"`
}

type deleteCollection struct {
	CollectionID string `long:"collection-id" description:"Span collection ID" required:"yes"`
	YesIAmSure   bool   `long:"yes-i-am-sure" description:"disable prompt for 'are you sure'"`
}

type updateCollection struct {
	CollectionID string   `long:"collection-id" env:"SPAN_COLLECTION_ID" description:"Span collection ID" required:"yes"`
	Tags         []string `long:"tag" description:"Set tag value (name:value)"`
}

func (r *addCollection) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	collection := spanapi.CreateCollectionRequest{
		TeamId: &r.TeamID,
		Tags:   helpers.TagMerge(nil, r.Tags),
	}

	if r.Name != "" {
		(*collection.Tags)["name"] = r.Name
	}

	col, res, err := client.CollectionsApi.CreateCollection(ctx).Body(collection).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("created collection %s\n", *col.CollectionId)
	return nil
}

func (r *getCollection) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	col, res, err := client.CollectionsApi.RetrieveCollection(ctx, r.CollectionID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	jsonData, err := json.MarshalIndent(col, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonData))
	return nil
}

func (r *deleteCollection) Execute([]string) error {
	if !r.YesIAmSure {
		if !helpers.VerifyDeleteIntent() {
			return fmt.Errorf("user aborted delete")
		}
	}

	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	col, res, err := client.CollectionsApi.DeleteCollection(ctx, r.CollectionID).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("deleted collection %s\n", *col.CollectionId)
	return nil
}

func (r *listCollection) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	collections, res, err := client.CollectionsApi.ListCollections(ctx).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	if collections.Collections == nil {
		fmt.Printf("no collections\n")
		return nil
	}

	if r.Format == "json" {
		json, err := json.MarshalIndent(collections, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(json))
		return nil
	}

	t := helpers.NewTableOutput(r.Format, r.NoColor, r.PageSize)
	t.SetTitle("Collections")
	t.AppendHeader(table.Row{"ID", "Name", "TeamID"})
	for _, col := range collections.Collections {
		// only truncate name if we output as 'text'
		name := col.GetTags()["name"]
		if r.Format == "text" {
			name = helpers.EllipsisString(name, 25)
		}

		t.AppendRow(table.Row{
			*col.CollectionId,
			name,
			*col.TeamId,
		})
	}
	helpers.RenderTable(t, r.Format)
	return nil
}

func (u *updateCollection) Execute([]string) error {
	client, ctx, cancel := helpers.NewSpanAPIClient()
	defer cancel()

	update := spanapi.UpdateCollectionRequest{
		Tags: helpers.TagMerge(nil, u.Tags),
	}

	collectionUpdated, res, err := client.CollectionsApi.UpdateCollection(ctx, u.CollectionID).Body(update).Execute()
	if err != nil {
		return helpers.ApiError(res, err)
	}

	fmt.Printf("updated collection %s\n", *collectionUpdated.CollectionId)
	return nil
}
