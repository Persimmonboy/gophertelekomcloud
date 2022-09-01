package clusters

import "github.com/opentelekomcloud/gophertelekomcloud"

func UpdateMasterIp(client *golangsdk.ServiceClient, id string, opts UpdateIpOpts) (r UpdateIpResult) {
	b, err := golangsdk.BuildRequestBody(opts, "spec")
	if err != nil {
		return nil, err
	}

	raw, err := client.Put(client.ServiceURL("clusters", id, "mastereip"), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
