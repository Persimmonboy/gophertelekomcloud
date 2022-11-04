package snapshots

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/internal/extract"
)

// Snapshot contains all the information associated with a Cinder Snapshot.
type Snapshot struct {
	// Unique identifier.
	ID string `json:"id"`
	// Date created.
	CreatedAt time.Time `json:"-"`
	// Date updated.
	UpdatedAt time.Time `json:"-"`
	// Display name.
	Name string `json:"name"`
	// Display description.
	Description string `json:"description"`
	// ID of the Volume from which this Snapshot was created.
	VolumeID string `json:"volume_id"`
	// Currect status of the Snapshot.
	Status string `json:"status"`
	// Size of the Snapshot, in GB.
	Size int `json:"size"`
	// User-defined key-value pairs.
	Metadata map[string]string `json:"metadata"`
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

// GetResult contains the response body and error from a Get request.
type GetResult struct {
	commonResult
}

// DeleteResult contains the response body and error from a Delete request.
type DeleteResult struct {
	golangsdk.ErrResult
}

// UpdateResult contains the response body and error from an Update request.
type UpdateResult struct {
	commonResult
}

// UnmarshalJSON converts our JSON API response into our snapshot struct
func (r *Snapshot) UnmarshalJSON(b []byte) error {
	type tmp Snapshot
	var s struct {
		tmp
		CreatedAt golangsdk.JSONRFC3339MilliNoZ `json:"created_at"`
		UpdatedAt golangsdk.JSONRFC3339MilliNoZ `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Snapshot(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)

	return err
}

// UpdateMetadataResult contains the response body and error from an UpdateMetadata request.
type UpdateMetadataResult struct {
	commonResult
}

// ExtractMetadata returns the metadata from a response from snapshots.UpdateMetadata.
func (r UpdateMetadataResult) ExtractMetadata() (map[string]interface{}, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	m := r.Body.(map[string]interface{})["metadata"]
	return m.(map[string]interface{}), nil
}

type commonResult struct {
	golangsdk.Result
}

func extra(err error, raw *http.Response) (*Snapshot, error) {
	if err != nil {
		return nil, err
	}

	var res Snapshot
	err = extract.IntoStructPtr(raw.Body, &res, "snapshot")
	return &res, err
}
