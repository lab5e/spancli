package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/lab5e/go-spanapi/v4"
)

type dataCmd struct {
	CollectionID string `long:"collection-id" description:"Span collection ID" required:"yes"`
	DeviceID     string `long:"device-id" description:"device id"`
	Limit        int32  `long:"limit" description:"max number of entries to fetch" default:"30"`
	BatchSize    int32  `long:"batch-size" description:"size of batches we request from Span" default:"500"`
	Start        string `long:"start" description:"start of time range in milliseconds since epoch"`
	End          string `long:"end" description:"end of time range in milliseconds since epoch"`
	Decode       bool   `long:"decode" description:"decode payload"`
	ISODate      bool   `long:"iso-date" description:"format date as ISO-8601 date"`
	JSONOutput   bool   `long:"json" description:"output as JSON"`
	JSONPretty   bool   `long:"pretty" description:"output as indented JSON"`
}

func (r *dataCmd) Execute([]string) error {
	client, ctx, cancel := newSpanAPIClient()
	defer cancel()

	if r.DeviceID == "" {
		return r.listCollectionData(client, ctx)
	}
	return r.listDeviceData(client, ctx)
}

func (r *dataCmd) listCollectionData(client *spanapi.APIClient, ctx context.Context) error {
	rowCount := int32(0)
	lastMessageID := ""

	if r.BatchSize > r.Limit {
		r.BatchSize = r.Limit
	}

	start := int64(0)
	end := int64(time.Now().UnixNano() / time.Hour.Milliseconds())

	for {
		// Don't ask for more than we need
		remainder := r.Limit - rowCount
		if remainder < r.BatchSize {
			r.BatchSize = remainder
			fmt.Printf("\nREMAINDER %d\n", remainder)
		}

		req := client.CollectionsApi.
			ListCollectionData(ctx, r.CollectionID).
			Limit(r.BatchSize).
			Offset(lastMessageID)

		// First time around the loop we use the start and end parameters
		if rowCount == 0 {
			req = req.
				Start(fmt.Sprintf("%d", start)).
				End(fmt.Sprintf("%d", end))
		} else {
			req.Offset(lastMessageID)
		}

		items, res, err := req.Execute()
		if err != nil {
			return apiError(res, err)
		}

		if items.Data == nil || len(*items.Data) == 0 {
			// Zero rows returned
			return nil
		}

		for _, data := range *items.Data {
			// bail out if we have reached the end of the interval
			ms, _ := msSinceEpochToTime(*data.Received)
			if ms < start {
				fmt.Printf("End of interval, %d rows returned\n", rowCount)
				return nil
			}

			fmt.Printf("%6d | %s %s %s %s\n", rowCount, *data.Received, *data.MessageId, *data.Device.DeviceId, *data.Payload)
			lastMessageID = *data.MessageId
			rowCount++

			if rowCount >= r.Limit {
				return nil
			}
		}
	}
}

func (r *dataCmd) listDeviceData(client *spanapi.APIClient, ctx context.Context) error {
	rowCount := int32(0)
	lastMessageID := ""

	if r.BatchSize > r.Limit {
		r.BatchSize = r.Limit
	}

	start := int64(0)
	end := int64(time.Now().UnixNano() / time.Hour.Milliseconds())

	for {
		// Don't ask for more than we need
		remainder := r.Limit - rowCount
		if remainder < r.BatchSize {
			r.BatchSize = remainder
			fmt.Printf("\nREMAINDER %d\n", remainder)
		}

		req := client.DevicesApi.
			ListDeviceData(ctx, r.CollectionID, r.DeviceID).
			Limit(r.BatchSize).
			Offset(lastMessageID)

		// First time around the loop we use the start and end parameters
		if rowCount == 0 {
			req = req.
				Start(fmt.Sprintf("%d", start)).
				End(fmt.Sprintf("%d", end))
		} else {
			req.Offset(lastMessageID)
		}

		items, res, err := req.Execute()
		if err != nil {
			return apiError(res, err)
		}

		if items.Data == nil || len(*items.Data) == 0 {
			// Zero rows returned
			return nil
		}

		for _, data := range *items.Data {
			// bail out if we have reached the end of the interval
			ms, _ := msSinceEpochToTime(*data.Received)
			if ms < start {
				fmt.Printf("End of interval, %d rows returned\n", rowCount)
				return nil
			}

			fmt.Printf("%6d | %s %s %s %s\n", rowCount, *data.Received, *data.MessageId, *data.Device.DeviceId, *data.Payload)
			lastMessageID = *data.MessageId
			rowCount++

			if rowCount >= r.Limit {
				return nil
			}
		}
	}
}

func msSinceEpochToTime(ts string) (int64, time.Time) {
	r, err := strconv.ParseInt(ts, 10, 63)
	if err != nil {
		return time.Now().UnixNano() / int64(time.Millisecond), time.Now()
	}
	return r, time.Unix(0, r*int64(time.Millisecond))
}
