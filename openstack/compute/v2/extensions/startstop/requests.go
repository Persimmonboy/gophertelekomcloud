package startstop

import "github.com/opentelekomcloud/gophertelekomcloud"

func actionURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("servers", id, "action")
}

// Start is the operation responsible for starting a Compute server.
func Start(client *golangsdk.ServiceClient, id string) (err error) {
	_, err = client.Post(actionURL(client, id), map[string]interface{}{"os-start": nil}, nil, nil)
	return
}

// Stop is the operation responsible for stopping a Compute server.
func Stop(client *golangsdk.ServiceClient, id string) (err error) {
	_, err = client.Post(actionURL(client, id), map[string]interface{}{"os-stop": nil}, nil, nil)
	return
}
