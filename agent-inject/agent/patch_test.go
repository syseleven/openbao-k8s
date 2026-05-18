// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package agent

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestAddVolumes(t *testing.T) {
	t.Run("no existing volumes", func(t *testing.T) {
		target := []corev1.Volume{}
		volumes := []corev1.Volume{
			{Name: "v1", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
		}
		patch := addVolumes(target, volumes, "/spec/volumes")
		assert.Len(t, patch, 1)
		var op, path string
		assert.NoError(t, json.Unmarshal(*patch[0]["op"], &op))
		assert.NoError(t, json.Unmarshal(*patch[0]["path"], &path))
		assert.Equal(t, "add", op)
		assert.Equal(t, "/spec/volumes", path)
	})

	t.Run("existing volumes no conflict", func(t *testing.T) {
		target := []corev1.Volume{
			{Name: "v0", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
		}
		volumes := []corev1.Volume{
			{Name: "v1", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
		}
		patch := addVolumes(target, volumes, "/spec/volumes")
		assert.Len(t, patch, 1)
		var op, path string
		assert.NoError(t, json.Unmarshal(*patch[0]["op"], &op))
		assert.NoError(t, json.Unmarshal(*patch[0]["path"], &path))
		assert.Equal(t, "add", op)
		assert.Equal(t, "/spec/volumes/-", path)
	})

	t.Run("skip identical volume", func(t *testing.T) {
		target := []corev1.Volume{
			{Name: "v1", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{Medium: "Memory"}}},
		}
		volumes := []corev1.Volume{
			{Name: "v1", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{Medium: "Memory"}}},
		}
		patch := addVolumes(target, volumes, "/spec/volumes")
		assert.Len(t, patch, 0)
	})

	t.Run("do not skip volume with different attributes", func(t *testing.T) {
		target := []corev1.Volume{
			{Name: "v1", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{Medium: ""}}},
		}
		volumes := []corev1.Volume{
			{Name: "v1", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{Medium: "Memory"}}},
		}
		patch := addVolumes(target, volumes, "/spec/volumes")
		assert.Len(t, patch, 1)
		var path string
		assert.NoError(t, json.Unmarshal(*patch[0]["path"], &path))
		assert.Equal(t, "/spec/volumes/-", path)
	})
}
