package inbox

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/go-spanapi/v4/apitools"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/global"
	"github.com/lab5e/spancli/pkg/helpers"
)

// watchInboxCmd connects to the MQTT service to monitor the inbox data
type watchInboxCmd struct {
	ID commonopt.CollectionAndDeviceOrGateway

	//lint:ignore SA5008 Multiple choices makes the linter unhappy
	Format string `long:"format" default:"text" description:"which output format to use" choice:"text" choice:"json"`
}

func (w *watchInboxCmd) Execute([]string) error {
	jwtToken := ""
	var stream apitools.DataStream
	var err error
	if global.Options.Token == "" {
		jwtToken = helpers.ReadCredentials()
		if jwtToken == "" {
			fmt.Println("You must either specify an API token or log in to the service")
			return errors.New("not authenticated")
		}
		fmt.Println("Logged in; using web socket to read the inbox activity")
		stream, err = NewDataStream(global.Options.Token, jwtToken, w.ID.CollectionID, w.ID.DeviceID, w.ID.GatewayID)
		if err != nil {
			return err
		}
	}

	if global.Options.Token != "" {
		fmt.Println("Reading activity stream from MQTT broker")

		buf := make([]byte, 8)
		rand.Read(buf)
		clientID := fmt.Sprintf("spancli_%d", binary.BigEndian.Uint64(buf))
		opts := make([]apitools.MQTTOption, 0)
		opts = append(opts, apitools.WithAPIToken(global.Options.Token))
		opts = append(opts, apitools.WithCollectionID(w.ID.CollectionID))
		opts = append(opts, apitools.WithClientID(clientID))
		if global.Options.OverrideEndpoint != "" {
			opts = append(opts, apitools.WithBrokerOverride(global.Options.MQTTOverrideEndpoint))
		}
		stream, err = apitools.NewMQTTStream(opts...)
		if err != nil {
			return err
		}

	}

	defer stream.Close()

	for {
		msg, err := stream.Recv()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return err
		}
		w.printMessage(msg)
	}
}

func (w *watchInboxCmd) printMessage(m spanapi.OutputDataMessage) {
	switch w.Format {
	case "json":
		buf, err := json.Marshal(&m)
		if err != nil {
			fmt.Printf("\nerror: %v\n", err)
			return
		}

		fmt.Printf("%s\n", string(buf))
	case "text":
		fmt.Printf("%s\t%s\t%s\t%s\n", m.GetMessageId(), helpers.DateFormat(m.GetReceived(), false), m.GetTransport(), m.GetPayload())
	default:
		panic(fmt.Sprintf("Unknown format: %s", w.Format))
	}
}
