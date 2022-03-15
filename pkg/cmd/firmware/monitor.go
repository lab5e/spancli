package firmware

import (
	"context"
	"fmt"
	"time"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/go-spanapi/v4/apitools"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/global"
	"github.com/lab5e/spancli/pkg/helpers"
)

// This monitors firmware states on devices by polling the device list every N seconds and printing the
// changes. Nothing fancy.
type monitorFirmware struct {
	ID      commonopt.Collection
	ImageID string `long:"image-id" description:"firmware image id" required:"no"`
}

type deviceState struct {
	State             string
	CurrentFirmwareId string
	Version           string
	SerialNumber      string
	Manufacturer      string
	ModelNumber       string
}

func fromFirmware(fw *spanapi.FirmwareMetadata) deviceState {
	return deviceState{
		State:             fw.GetState(),
		CurrentFirmwareId: fw.GetCurrentFirmwareId(),
		Version:           fw.GetFirmwareVersion(),
		SerialNumber:      fw.GetSerialNumber(),
		Manufacturer:      fw.GetManufacturer(),
		ModelNumber:       fw.GetModelNumber(),
	}
}

func (s *deviceState) Changed(other deviceState) bool {
	if s.State != other.State {
		return true
	}
	if s.CurrentFirmwareId != other.CurrentFirmwareId {
		return true
	}
	// This shouldn't really change unless the current firmware id changes but it might be out of sync if
	// the firmware images are out of sync with what's deployed.
	if s.Version != other.Version || s.SerialNumber != other.SerialNumber ||
		s.Manufacturer != other.Manufacturer || s.ModelNumber != other.ModelNumber {
		return true
	}
	return false
}

func (c *monitorFirmware) Execute([]string) error {
	client, _, authDone := helpers.NewSpanAPIClient()
	authDone()

	ctx, done := context.WithCancel(apitools.ContextWithAuth(global.Options.Token))
	defer done()

	if c.ImageID != "" {
		fmt.Printf("Filtering on firmware image %s\n", c.ImageID)
	}
	fmt.Printf("Monitoring collection %s for firmware updates...\n", c.ID.CollectionID)

	states := make(map[string]deviceState)
	for {
		list, res, err := client.DevicesApi.ListDevices(ctx, c.ID.CollectionID).Execute()
		if err != nil {
			return helpers.ApiError(res, err)
		}

		for _, d := range list.Devices {
			if d.Firmware != nil {
				if c.ImageID != "" && c.ImageID != d.Firmware.GetCurrentFirmwareId() && c.ImageID != d.Firmware.GetTargetFirmwareId() {
					continue
				}
				newState := fromFirmware(d.Firmware)
				state, ok := states[d.GetDeviceId()]
				if ok {
					if state.Changed(newState) {
						c.reportChange(d.GetDeviceId(), state, newState)
						states[d.GetDeviceId()] = newState
					}
				} else {
					// just add
					states[d.GetDeviceId()] = newState
				}
			}
		}
		time.Sleep(time.Second * 5)
	}
}

func (c *monitorFirmware) reportChange(deviceID string, oldState deviceState, newState deviceState) {
	if oldState.State != newState.State {
		fmt.Printf("Device %s changed state from %s to %s\n", deviceID, oldState.State, newState.State)
	}
	if oldState.CurrentFirmwareId != newState.CurrentFirmwareId {
		fmt.Printf("Device %s changed current image from %s to %s\n", deviceID, oldState.CurrentFirmwareId, newState.CurrentFirmwareId)
	}
	if oldState.Version != newState.Version {
		fmt.Printf("Device %s changed version from %s to %s\n", deviceID, oldState.Version, newState.Version)
	}
	if oldState.Manufacturer != newState.Manufacturer {
		fmt.Printf("Device %s changed manufacturer from %s to %s\n", deviceID, oldState.Manufacturer, newState.Manufacturer)
	}
	if oldState.SerialNumber != newState.SerialNumber {
		fmt.Printf("Device %s changed serial number from %s to %s\n", deviceID, oldState.SerialNumber, newState.SerialNumber)
	}
	if oldState.ModelNumber != newState.ModelNumber {
		fmt.Printf("Device %s changed model number from %s to %s\n", deviceID, oldState.ModelNumber, newState.ModelNumber)
	}
}
