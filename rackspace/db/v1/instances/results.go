package instances

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack/db/v1/instances"
)

type Datastore struct {
	Type    string
	Version string
}

type Instance struct {
	Created   string //time.Time
	Updated   string //time.Time
	Datastore Datastore
	Flavor    os.Flavor
	Hostname  string
	ID        string
	Links     []gophercloud.Link
	Name      string
	Status    string
	Volume    os.Volume
}

// CreateResult represents the result of a Create operation.
type CreateResult struct {
	os.CreateResult
}

func (r CreateResult) Extract() (*Instance, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response struct {
		Instance Instance `mapstructure:"instance"`
	}

	err := mapstructure.Decode(r.Body, &response)

	return &response.Instance, err
}
