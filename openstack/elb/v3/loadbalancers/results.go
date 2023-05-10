package loadbalancers

import (
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/pagination"
)

type EipInfo struct {
	// Eip ID
	EipID string `json:"eip_id"`
	// Eip Address
	EipAddress string `json:"eip_address"`
	// Eip Address
	IpVersion int `json:"ip_version"`
}

type PublicIpInfo struct {
	// Public IP ID
	PublicIpID string `json:"publicip_id"`
	// Public IP Address
	PublicIpAddress string `json:"publicip_address"`
	// IP Version
	IpVersion int `json:"ip_version"`
}

// StatusTree represents the status of a loadbalancer.
type StatusTree struct {
	Loadbalancer *LoadBalancer `json:"loadbalancer"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a loadbalancer.
func (r commonResult) Extract() (*LoadBalancer, error) {
	s := new(LoadBalancer)
	err := r.ExtractIntoStructPtr(s, "loadbalancer")
	if err != nil {
		return nil, err
	}
	return s, nil
}

// GetStatusesResult represents the result of a GetStatuses operation.
// Call its Extract method to interpret it as a StatusTree.
type GetStatusesResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts the status of
// a Loadbalancer.
func (r GetStatusesResult) Extract() (*StatusTree, error) {
	s := new(StatusTree)
	err := r.ExtractIntoStructPtr(s, "statuses")
	if err != nil {
		return nil, err
	}
	return s, nil
}

// LoadbalancerPage is the page returned by a pager when traversing over a
// collection of loadbalancer.
type LoadbalancerPage struct {
	pagination.PageWithInfo
}

// IsEmpty checks whether a FlavorsPage struct is empty.
func (r LoadbalancerPage) IsEmpty() (bool, error) {
	is, err := ExtractLoadbalancers(r)
	return len(is) == 0, err
}

// ExtractLoadbalancers accepts a Page struct, specifically a LoadbalancerPage struct,
// and extracts the elements into a slice of loadbalancer structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractLoadbalancers(r pagination.Page) ([]LoadBalancer, error) {
	var s []LoadBalancer
	err := (r.(LoadbalancerPage)).ExtractIntoSlicePtr(&s, "loadbalancers")
	if err != nil {
		return nil, err
	}
	return s, nil
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a LoadBalancer.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a LoadBalancer.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a LoadBalancer.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
