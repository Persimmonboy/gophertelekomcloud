package security_policy

import (
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/internal/build"
)

type UpdateOpts struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Protocols   []string `json:"protocols,omitempty"`
	Ciphers     []string `json:"ciphers,omitempty"`
}

func Update(client *golangsdk.ServiceClient, opts UpdateOpts, id string) (*SecurityPolicy, error) {
	b, err := build.RequestBody(opts, "security_policy")
	if err != nil {
		return nil, err
	}

	raw, err := client.Put(client.ServiceURL("security-policies", id), b, nil, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return extra(err, raw)
}
