package output

type oid struct {
	OutputID string `long:"output-id" description:"output ID" required:"yes"`
}
type otp struct {
	//lint:ignore SA5008 Multiple choices makes the linter unhappy
	Type string `long:"type" description:"type of output" default:"webhook" required:"yes" choice:"webhook" choice:"udpout" choice:"mqttclient" choice:"mqttbroker"`
}
