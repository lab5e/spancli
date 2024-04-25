package activity

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/global"
	"github.com/lab5e/spancli/pkg/helpers"
)

// Command holds the subcommands for the inbox command
type Command struct {
	Watch watchActivityCommand `command:"watch" alias:"w" description:"watch activity stream"`
}

type watchActivityCommand struct {
	ID commonopt.CollectionAndOptionalDevice
}

func (w *watchActivityCommand) Execute([]string) error {
	jwtToken := ""
	if global.Options.Token == "" {
		jwtToken = helpers.ReadCredentials()
		if jwtToken == "" {
			fmt.Println("You must either specify an API token or log in to the service")
			return errors.New("not authenticated")
		}
	}

	ws, err := NewActivityEventStream(global.Options.Token, jwtToken, w.ID.CollectionID)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(os.Stdout)
	for {
		defer ws.Close()
		res, err := ws.Recv()
		if res.CollectionId == nil {
			continue
		}
		if w.ID.DeviceID != "" {
			if res.GetDeviceId() != w.ID.DeviceID {
				continue
			}
		}
		if err != nil {
			fmt.Printf("Error receiving activity stream: %v", err)
			return err
		}
		if err := encoder.Encode(res); err != nil {
			fmt.Printf("Error encoding JSON: %v", err)
			return err
		}
	}
}
