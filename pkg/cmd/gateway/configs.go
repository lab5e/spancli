package gateway

import "fmt"

type sampleConfigs struct {
}

func (c *sampleConfigs) Execute([]string) error {
	fmt.Println(`
Sample configuration for a LoRaWAN gateway:
	
	{ "appEUI": "00-11-22-33-44-55-66-77" }
	
Configure with
	--config appEUI:00-11-22-33-44-55-66-77


	`)
	return nil
}
