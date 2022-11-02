// Package extensions contains common functions for creating block storage
// resources that are extensions of the block storage API. See the `*_test.go`
// files for example usages.
package extensions

import (
	"fmt"
	"strings"
	"testing"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/images"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/servers"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/evs/extensions/backups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/evs/extensions/volumeactions"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/evs/v2/volumes"
	v3 "github.com/opentelekomcloud/gophertelekomcloud/openstack/evs/v3/volumes"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/evs/v3/volumetypes"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
)

// CreateUploadImage will upload volume it as volume-baked image. An name of new image or err will be
// returned
func CreateUploadImage(t *testing.T, client *golangsdk.ServiceClient, volume *volumes.Volume) (volumeactions.VolumeImage, error) {
	if testing.Short() {
		t.Skip("Skipping test that requires volume-backed image uploading in short mode.")
	}

	imageName := tools.RandomString("ACPTTEST", 16)
	uploadImageOpts := volumeactions.UploadImageOpts{
		ImageName: imageName,
		Force:     true,
	}

	volumeImage, err := volumeactions.UploadImage(client, volume.ID, uploadImageOpts)
	if err != nil {
		return volumeImage, err
	}

	t.Logf("Uploading volume %s as volume-backed image %s", volume.ID, imageName)

	if err := volumes.WaitForStatus(client, volume.ID, "available", 60); err != nil {
		return volumeImage, err
	}

	t.Logf("Uploaded volume %s as volume-backed image %s", volume.ID, imageName)

	return volumeImage, nil

}

// DeleteUploadedImage deletes uploaded image. An error will be returned
// if the deletion request failed.
func DeleteUploadedImage(t *testing.T, client *golangsdk.ServiceClient, imageID string) error {
	if testing.Short() {
		t.Skip("Skipping test that requires volume-backed image removing in short mode.")
	}

	t.Logf("Removing image %s", imageID)

	err := images.Delete(client, imageID)
	if err != nil {
		return err
	}

	return nil
}

// CreateVolumeAttach will attach a volume to an instance. An error will be
// returned if the attachment failed.
func CreateVolumeAttach(t *testing.T, client *golangsdk.ServiceClient, volume *volumes.Volume, server *servers.Server) error {
	if testing.Short() {
		t.Skip("Skipping test that requires volume attachment in short mode.")
	}

	attachOpts := volumeactions.AttachOpts{
		MountPoint:   "/mnt",
		Mode:         "rw",
		InstanceUUID: server.ID,
	}

	t.Logf("Attempting to attach volume %s to server %s", volume.ID, server.ID)

	if err := volumeactions.Attach(client, volume.ID, attachOpts); err != nil {
		return err
	}

	if err := volumes.WaitForStatus(client, volume.ID, "in-use", 60); err != nil {
		return err
	}

	t.Logf("Attached volume %s to server %s", volume.ID, server.ID)

	return nil
}

// CreateVolumeReserve creates a volume reservation. An error will be returned
// if the reservation failed.
func CreateVolumeReserve(t *testing.T, client *golangsdk.ServiceClient, volume *volumes.Volume) error {
	if testing.Short() {
		t.Skip("Skipping test that requires volume reservation in short mode.")
	}

	t.Logf("Attempting to reserve volume %s", volume.ID)

	if err := volumeactions.Reserve(client, volume.ID); err != nil {
		return err
	}

	t.Logf("Reserved volume %s", volume.ID)

	return nil
}

// DeleteVolumeAttach will detach a volume from an instance. A fatal error will
// occur if the snapshot failed to be deleted. This works best when used as a
// deferred function.
func DeleteVolumeAttach(t *testing.T, client *golangsdk.ServiceClient, volume *volumes.Volume) {
	t.Logf("Attepting to detach volume volume: %s", volume.ID)

	detachOpts := volumeactions.DetachOpts{
		AttachmentID: volume.Attachments[0].AttachmentID,
	}

	if err := volumeactions.Detach(client, volume.ID, detachOpts); err != nil {
		t.Fatalf("Unable to detach volume %s: %v", volume.ID, err)
	}

	if err := volumes.WaitForStatus(client, volume.ID, "available", 60); err != nil {
		t.Fatalf("Volume %s failed to become unavailable in 60 seconds: %v", volume.ID, err)
	}

	t.Logf("Detached volume: %s", volume.ID)
}

// DeleteVolumeReserve deletes a volume reservation. A fatal error will occur
// if the deletion request failed. This works best when used as a deferred
// function.
func DeleteVolumeReserve(t *testing.T, client *golangsdk.ServiceClient, volume *volumes.Volume) {
	if testing.Short() {
		t.Skip("Skipping test that requires volume reservation in short mode.")
	}

	t.Logf("Attempting to unreserve volume %s", volume.ID)

	if err := volumeactions.Unreserve(client, volume.ID); err != nil {
		t.Fatalf("Unable to unreserve volume %s: %v", volume.ID, err)
	}

	t.Logf("Unreserved volume %s", volume.ID)
}

// ExtendVolumeSize will extend the size of a volume.
func ExtendVolumeSize(t *testing.T, client *golangsdk.ServiceClient, volume *volumes.Volume) error {
	t.Logf("Attempting to extend the size of volume %s", volume.ID)

	extendOpts := volumeactions.ExtendSizeOpts{
		NewSize: 2,
	}

	err := volumeactions.ExtendSize(client, volume.ID, extendOpts)
	if err != nil {
		return err
	}

	if err := volumes.WaitForStatus(client, volume.ID, "available", 60); err != nil {
		return err
	}

	return nil
}

// SetImageMetadata will apply the metadata to a volume.
func SetImageMetadata(t *testing.T, client *golangsdk.ServiceClient, volume *volumes.Volume) error {
	t.Logf("Attempting to apply image metadata to volume %s", volume.ID)

	imageMetadataOpts := volumeactions.ImageMetadataOpts{
		Metadata: map[string]string{
			"image_name": "testimage",
		},
	}

	err := volumeactions.SetImageMetadata(client, volume.ID, imageMetadataOpts)
	if err != nil {
		return err
	}

	return nil
}

// CreateBackup will create a backup based on a volume. An error will be
// will be returned if the backup could not be created.
func CreateBackup(t *testing.T, client *golangsdk.ServiceClient, volumeID string) (*backups.Backup, error) {
	t.Logf("Attempting to create a backup of volume %s", volumeID)

	backupName := tools.RandomString("ACPTTEST", 16)
	createOpts := backups.CreateOpts{
		VolumeID: volumeID,
		Name:     backupName,
	}

	backup, err := backups.Create(client, createOpts)
	if err != nil {
		return nil, err
	}

	err = WaitForBackupStatus(client, backup.ID, "available")
	if err != nil {
		return nil, err
	}

	backup, err = backups.Get(client, backup.ID)
	if err != nil {
		return nil, err
	}

	t.Logf("Successfully created backup %s", backup.ID)
	tools.PrintResource(t, backup)

	th.AssertEquals(t, backup.Name, backupName)

	return backup, nil
}

// DeleteBackup will delete a backup. A fatal error will occur if the backup
// could not be deleted. This works best when used as a deferred function.
func DeleteBackup(t *testing.T, client *golangsdk.ServiceClient, backupID string) {
	if err := backups.Delete(client, backupID); err != nil {
		t.Fatalf("Unable to delete backup %s: %s", backupID, err)
	}

	t.Logf("Deleted backup %s", backupID)
}

// WaitForBackupStatus will continually poll a backup, checking for a particular
// status. It will do this for the amount of seconds defined.
func WaitForBackupStatus(client *golangsdk.ServiceClient, id, status string) error {
	return tools.WaitFor(func() (bool, error) {
		current, err := backups.Get(client, id)
		if err != nil {
			return false, err
		}

		if current.Status == status {
			return true, nil
		}

		return false, nil
	})
}

// SetBootable will set a bootable status to a volume.
func SetBootable(t *testing.T, client *golangsdk.ServiceClient, volume *volumes.Volume) error {
	t.Logf("Attempting to apply bootable status to volume %s", volume.ID)

	bootableOpts := volumeactions.BootableOpts{
		Bootable: true,
	}

	err := volumeactions.SetBootable(client, volume.ID, bootableOpts)
	if err != nil {
		return err
	}

	vol, err := v3.Get(client, volume.ID)
	if err != nil {
		return err
	}

	if strings.ToLower(vol.Bootable) != "true" {
		return fmt.Errorf("Volume bootable status is %q, expected 'true'", vol.Bootable)
	}

	bootableOpts = volumeactions.BootableOpts{
		Bootable: false,
	}

	err = volumeactions.SetBootable(client, volume.ID, bootableOpts)
	if err != nil {
		return err
	}

	vol, err = v3.Get(client, volume.ID)
	if err != nil {
		return err
	}

	if strings.ToLower(vol.Bootable) == "true" {
		return fmt.Errorf("Volume bootable status is %q, expected 'false'", vol.Bootable)
	}

	return nil
}

// ChangeVolumeType will extend the size of a volume.
func ChangeVolumeType(t *testing.T, client *golangsdk.ServiceClient, volume *v3.Volume, vt *volumetypes.VolumeType) error {
	t.Logf("Attempting to change the type of volume %s from %s to %s", volume.ID, volume.VolumeType, vt.Name)

	changeOpts := volumeactions.ChangeTypeOpts{
		NewType:         vt.Name,
		MigrationPolicy: volumeactions.MigrationPolicyOnDemand,
	}

	err := volumeactions.ChangeType(client, volume.ID, changeOpts)
	if err != nil {
		return err
	}

	if err := volumes.WaitForStatus(client, volume.ID, "available", 60); err != nil {
		return err
	}

	return nil
}

// ReImage will re-image a volume
func ReImage(t *testing.T, client *golangsdk.ServiceClient, volume *volumes.Volume, imageID string) error {
	t.Logf("Attempting to re-image volume %s", volume.ID)

	reimageOpts := volumeactions.ReImageOpts{
		ImageID:         imageID,
		ReImageReserved: false,
	}

	err := volumeactions.ReImage(client, volume.ID, reimageOpts)
	if err != nil {
		return err
	}

	err = volumes.WaitForStatus(client, volume.ID, "available", 60)
	if err != nil {
		return err
	}

	vol, err := v3.Get(client, volume.ID)
	if err != nil {
		return err
	}

	if vol.VolumeImageMetadata == nil {
		return fmt.Errorf("volume does not have VolumeImageMetadata map")
	}

	if strings.ToLower(vol.VolumeImageMetadata["image_id"]) != imageID {
		return fmt.Errorf("volume image id '%s', expected '%s'", vol.VolumeImageMetadata["image_id"], imageID)
	}

	return nil
}
