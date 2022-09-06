package defsecrules

import (
	"encoding/json"

	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/internal/extract"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/secgroups"
	"github.com/opentelekomcloud/gophertelekomcloud/pagination"
)

// DefaultRule represents a rule belonging to the "default" security group.
// It is identical to an openstack/compute/v2/extensions/secgroups.Rule.
type DefaultRule secgroups.Rule

func (r *DefaultRule) UnmarshalJSON(b []byte) error {
	var s secgroups.Rule
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = DefaultRule(s)
	return nil
}

// DefaultRulePage is a single page of a DefaultRule collection.
type DefaultRulePage struct {
	pagination.SinglePageBase
}

// IsEmpty determines whether or not a page of default rules contains any results.
func (page DefaultRulePage) IsEmpty() (bool, error) {
	users, err := ExtractDefaultRules(page)
	return len(users) == 0, err
}

// ExtractDefaultRules returns a slice of DefaultRules contained in a single
// page of results.
func ExtractDefaultRules(r pagination.Page) ([]DefaultRule, error) {
	var res struct {
		DefaultRules []DefaultRule `json:"security_group_default_rules"`
	}
	err := (r.(DefaultRulePage)).ExtractInto(&res)
	return res, err
}

type commonResult struct {
	golangsdk.Result
}

// CreateResult represents the result of a create operation.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

// Extract will extract a DefaultRule struct from a Create or Get response.
func (raw commonResult) Extract() (*DefaultRule, error) {
	var res struct {
		DefaultRule DefaultRule `json:"security_group_default_rule"`
	}
	err = extract.Into(raw, &res)
	return &s.DefaultRule, err
}

// DeleteResult is the response from a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
