package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/go-spanapi/v4"
)

type collectionCmd struct {
	Add    addCollection    `command:"add" description:"create new collection"`
	Get    getCollection    `command:"get" description:"get collection details"`
	List   listCollection   `command:"list" alias:"ls" description:"list collections"`
	Delete deleteCollection `command:"delete" alias:"del" description:"delete collection"`
	Update updateCollection `command:"update" alias:"up" description:"update collection"`
}

type addCollection struct {
	TeamID       string   `long:"team-id" description:"team the collection belongs to" required:"yes"`
	Tags         []string `long:"tag" description:"Set tag value (name:value)"`
	Name         string   `long:"name" description:"name of the collection"`
	MaskIMSI     bool     `long:"mask-imsi" description:"mask IMSI"`
	MaskIMEI     bool     `long:"mask-imei" description:"mask IMEI"`
	MaskLocation bool     `long:"mask-location" description:"mask location"`
	MaskMSISDN   bool     `long:"mask-msisdn" description:"mask MSISDN (phone number)"`
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
	tags, err := tagsToMap(r.Tags)
	if err != nil {
		return err
	}

	client, ctx, cancel := newSpanAPIClient()
	defer cancel()

	collection := spanapi.Collection{
		TeamId: &r.TeamID,
		FieldMask: &spanapi.FieldMask{
			Imsi:     &r.MaskIMSI,
			Imei:     &r.MaskIMEI,
			Msisdn:   &r.MaskMSISDN,
			Location: &r.MaskLocation,
		},
		Tags: &tags,
	}

	if r.Name != "" {
		(*collection.Tags)["name"] = r.Name
	}

	col, res, err := client.CollectionsApi.CreateCollection(ctx).Body(collection).Execute()
	if err != nil {
		return apiError(res, err)
	}

	fmt.Printf("created collection %s\n", *col.CollectionId)
	return nil
}

func (r *getCollection) Execute([]string) error {
	client, ctx, cancel := newSpanAPIClient()
	defer cancel()

	col, res, err := client.CollectionsApi.RetrieveCollection(ctx, r.CollectionID).Execute()
	if err != nil {
		return apiError(res, err)
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
		if !verifyDeleteIntent() {
			return fmt.Errorf("user aborted delete")
		}
	}

	client, ctx, cancel := newSpanAPIClient()
	defer cancel()

	col, res, err := client.CollectionsApi.DeleteCollection(ctx, r.CollectionID).Execute()
	if err != nil {
		return apiError(res, err)
	}

	fmt.Printf("deleted collection %s", *col.CollectionId)
	return nil
}

func (r *listCollection) Execute([]string) error {
	client, ctx, cancel := newSpanAPIClient()
	defer cancel()

	collections, res, err := client.CollectionsApi.ListCollections(ctx).Execute()
	if err != nil {
		return apiError(res, err)
	}

	if r.Format == "json" {
		json, err := json.MarshalIndent(collections, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(json))
		return nil
	}

	t := newTableOutput(r.Format, r.NoColor, r.PageSize)
	t.SetTitle("Collections")
	t.AppendHeader(table.Row{"ID", "Name", "TeamID"})
	for _, col := range *collections.Collections {
		// only truncate name if we output as 'text'
		name := col.GetTags()["name"]
		if r.Format == "text" {
			name = truncateString(name, 25)
		}

		t.AppendRow(table.Row{
			*col.CollectionId,
			name,
			*col.TeamId,
		})
	}
	renderTable(t, r.Format)
	return nil
}

func (u *updateCollection) Execute([]string) error {
	client, ctx, cancel := newSpanAPIClient()
	defer cancel()

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
