package metadata

import (
	"net/http"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/internal/extract"
)

func GetOne(client *golangsdk.ServiceClient, snapshotId string, key string) (map[string]string, error) {
	// GET /v3/{project_id}/snapshots/{snapshot_id}/metadata/{key}
	raw, err := client.Get(client.ServiceURL("snapshots", snapshotId, "metadata", key), nil, nil)
	return extraMeta(err, raw)
}

func extraMeta(err error, raw *http.Response) (map[string]string, error) {
	if err != nil {
		return nil, err
	}

	var res struct {
		Metadata map[string]string `json:"meta"`
	}
	err = extract.Into(raw.Body, &res)
	return res.Metadata, err
}
