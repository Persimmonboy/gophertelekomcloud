package nodes

import "github.com/opentelekomcloud/gophertelekomcloud"

// Delete will permanently delete a particular node based on its unique ID and cluster ID.
func Delete(client *golangsdk.ServiceClient, clusterID, nodeID string) (r DeleteResult) {
	raw, err := client.Delete(client.ServiceURL("clusters", clusterID, "nodes", nodeID), &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}
