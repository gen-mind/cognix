package main

import (
	"cognix.ch/api/v2/core/proto"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestK8SServer_GetList(t *testing.T) {
	service, err := NewDockerServer("mock")
	assert.NoError(t, err)
	resp, err := service.GetList(context.Background(), &proto.ConfigMapList{
		Name: "api-srv",
	})
	for _, row := range resp.GetValues() {
		if row.GetKey() == "PORT" {
			t.Logf(" ------ %s", row.GetValue())
		}
		t.Logf("%s -- %s ", row.GetKey(), row.GetValue())
	}

	_, err = service.Save(context.Background(), &proto.ConfigMapSave{
		Name: "api-srv",
		Value: &proto.ConfigMapRecord{
			Key:   "NEW-PORT",
			Value: "99499",
		},
	})

	assert.NoError(t, err)
	_, err = service.Delete(context.Background(), &proto.ConfigMapDelete{
		Name: "api-srv",
		Key:  "PORT",
	})
	assert.NoError(t, err)
}
