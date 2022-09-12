package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/images"
	"github.com/opentelekomcloud/gophertelekomcloud/pagination"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
	fake "github.com/opentelekomcloud/gophertelekomcloud/testhelper/client"
)

func TestListImages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/images/detail", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		_ = r.ParseForm()
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			_, _ = fmt.Fprint(w, `
				{
					"images": [
						{
							"status": "ACTIVE",
							"updated": "2014-09-23T12:54:56Z",
							"id": "f3e4a95d-1f4f-4989-97ce-f3a1fb8c04d7",
							"OS-EXT-IMG-SIZE:size": 476704768,
							"name": "F17-x86_64-cfntools",
							"created": "2014-09-23T12:54:52Z",
							"minDisk": 0,
							"progress": 100,
							"minRam": 0,
							"metadata": {
								"architecture": "x86_64",
								"block_device_mapping": {
									"guest_format": null,
									"boot_index": 0,
									"device_name": "/dev/vda",
									"delete_on_termination": false
								}
              }
						},
						{
							"status": "ACTIVE",
							"updated": "2014-09-23T12:51:43Z",
							"id": "f90f6034-2570-4974-8351-6b49732ef2eb",
							"OS-EXT-IMG-SIZE:size": 13167616,
							"name": "cirros-0.3.2-x86_64-disk",
							"created": "2014-09-23T12:51:42Z",
							"minDisk": 0,
							"progress": 100,
							"minRam": 0
						}
					]
				}
			`)
		case "2":
			_, _ = fmt.Fprint(w, `{ "images": [] }`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})

	pages := 0
	options := images.ListOpts{Limit: 2}
	err := images.ListDetail(fake.ServiceClient(), options).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := images.ExtractImages(page)
		if err != nil {
			return false, err
		}

		expected := []images.Image{
			{
				ID:       "f3e4a95d-1f4f-4989-97ce-f3a1fb8c04d7",
				Name:     "F17-x86_64-cfntools",
				Created:  "2014-09-23T12:54:52Z",
				Updated:  "2014-09-23T12:54:56Z",
				MinDisk:  0,
				MinRAM:   0,
				Progress: 100,
				Status:   "ACTIVE",
				Metadata: map[string]interface{}{
					"architecture": "x86_64",
					"block_device_mapping": map[string]interface{}{
						"guest_format":          interface{}(nil),
						"boot_index":            float64(0),
						"device_name":           "/dev/vda",
						"delete_on_termination": false,
					},
				},
			},
			{
				ID:       "f90f6034-2570-4974-8351-6b49732ef2eb",
				Name:     "cirros-0.3.2-x86_64-disk",
				Created:  "2014-09-23T12:51:42Z",
				Updated:  "2014-09-23T12:51:43Z",
				MinDisk:  0,
				MinRAM:   0,
				Progress: 100,
				Status:   "ACTIVE",
			},
		}

		th.AssertDeepEquals(t, expected, actual)

		return false, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, pages)
}

func TestGetImage(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/images/12345678", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `
			{
				"image": {
					"status": "ACTIVE",
					"updated": "2014-09-23T12:54:56Z",
					"id": "f3e4a95d-1f4f-4989-97ce-f3a1fb8c04d7",
					"OS-EXT-IMG-SIZE:size": 476704768,
					"name": "F17-x86_64-cfntools",
					"created": "2014-09-23T12:54:52Z",
					"minDisk": 0,
					"progress": 100,
					"minRam": 0,
					"metadata": {
						"architecture": "x86_64",
						"block_device_mapping": {
							"guest_format": null,
							"boot_index": 0,
							"device_name": "/dev/vda",
							"delete_on_termination": false
						}
					}
				}
			}
		`)
	})

	actual, err := images.Get(fake.ServiceClient(), "12345678")
	th.AssertNoErr(t, err)

	expected := &images.Image{
		Status:   "ACTIVE",
		Updated:  "2014-09-23T12:54:56Z",
		ID:       "f3e4a95d-1f4f-4989-97ce-f3a1fb8c04d7",
		Name:     "F17-x86_64-cfntools",
		Created:  "2014-09-23T12:54:52Z",
		MinDisk:  0,
		Progress: 100,
		MinRAM:   0,
		Metadata: map[string]interface{}{
			"architecture": "x86_64",
			"block_device_mapping": map[string]interface{}{
				"guest_format":          interface{}(nil),
				"boot_index":            float64(0),
				"device_name":           "/dev/vda",
				"delete_on_termination": false,
			},
		},
	}

	th.AssertDeepEquals(t, expected, actual)
}

func TestNextPageURL(t *testing.T) {
	var page images.ImagePage
	bodyString := []byte(`{"images":{"links":[{"href":"http://192.154.23.87/12345/images/image3","rel":"bookmark"}]}, "images_links":[{"href":"http://192.154.23.87/12345/images/image4","rel":"next"}]}`)
	page.Body = bodyString

	expected := "http://192.154.23.87/12345/images/image4"
	actual, err := page.NextPageURL()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, expected, actual)
}

// Test Image delete
func TestDeleteImage(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/images/12345678", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})

	res := images.Delete(fake.ServiceClient(), "12345678")
	th.AssertNoErr(t, res)
}
