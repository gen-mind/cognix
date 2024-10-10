package main

import (
	"bytes"
	"cognix.ch/api/v2/core/proto"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"os"
	"strings"
)

// DockerServer represents a server for Docker configuration maps.
type DockerServer struct {
	proto.UnsafeConfigMapServer
	path string
}

// readConfig reads the contents of the specified file and returns them as a slice of strings.
// The function removes any carriage return characters (\r) from the file contents before splitting it into rows.
// It returns the rows as a slice of strings and an error if the file cannot be read.
func (k *DockerServer) readConfig(name string) ([]string, error) {
	buf, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	buf = bytes.ReplaceAll(buf, []byte("\r"), []byte(""))

	rows := strings.Split(string(buf), "\n")
	return rows, nil

}

// GetList retrieves a list of config map records from the DockerServer.
// It takes a context and a ConfigMapList as input and returns a ConfigMapListResponse and an error.
// The method reads the config file for the specified ConfigMapList name and populates the response
// with the values found in the file. Each line in the file is treated as a key-value pair separated by "=".
// If the line does not contain exactly two elements, it is skipped. The response is then returned.
// If an error occurs while reading the file, it is propagated back to the caller.
func (k *DockerServer) GetList(ctx context.Context, r *proto.ConfigMapList) (*proto.ConfigMapListResponse, error) {
	filename := fmt.Sprintf("%s/%s.env", k.path, r.Name)

	rows, err := k.readConfig(filename)
	if err != nil {
		return nil, err
	}

	result := proto.ConfigMapListResponse{
		Values: make([]*proto.ConfigMapRecord, 0),
	}
	for _, row := range rows {
		value := strings.Split(row, "=")
		if len(value) != 2 {
			continue
		}
		result.Values = append(result.Values, &proto.ConfigMapRecord{
			Key:   value[0],
			Value: value[1],
		})
	}
	return &result, nil
}

// Save saves the configuration map to a file.
// It takes a context and a ConfigMapSave request as input.
// It retrieves the filename based on the path and name from the request.
// It reads the existing configuration from the file.
// If the key already exists in the configuration, it updates the value.
// If the key doesn't exist, it adds a new key-value pair.
// Finally, it writes the updated configuration to the file.
// It returns an empty message and any error that occurred.
func (k *DockerServer) Save(ctx context.Context, r *proto.ConfigMapSave) (*empty.Empty, error) {
	filename := fmt.Sprintf("%s/%s.env", k.path, r.Name)
	rows, err := k.readConfig(filename)
	if err != nil {
		return nil, err
	}
	isExists := false
	for i, row := range rows {
		value := strings.Split(row, "=")
		if len(value) != 2 {
			continue
		}
		if value[0] == r.Value.Key {
			rows[i] = fmt.Sprintf("%s=%s", r.Value.Key, r.Value.Value)
			isExists = true
			break
		}
	}
	if !isExists {
		rows = append(rows, fmt.Sprintf("%s=%s", r.Value.Key, r.Value.Value))
	}
	return &empty.Empty{}, os.WriteFile(filename, []byte(strings.Join(rows, "\n")), 0644)
}

// Delete removes a key-value pair from a config map file.
func (k *DockerServer) Delete(ctx context.Context, r *proto.ConfigMapDelete) (*empty.Empty, error) {
	filename := fmt.Sprintf("%s/%s.env", k.path, r.Name)
	rows, err := k.readConfig(filename)
	if err != nil {
		return nil, err
	}
	newRows := make([]string, 0)

	for _, row := range rows {
		value := strings.Split(row, "=")
		if len(value) == 2 && value[0] == r.Key {
			continue
		}
		newRows = append(newRows, row)
	}
	return &empty.Empty{}, os.WriteFile(filename, []byte(strings.Join(newRows, "\n")), 0644)
}

// NewDockerServer creates a new instance of DockerServer with the specified root path.
// It implements the proto.ConfigMapServer interface and returns the server instance and an error.
// If the root path does not exist, it returns an error indicating that the docker config maps
// cannot be found.
//
// NewDockerServer(root string) (proto.ConfigMapServer, error)
func NewDockerServer(root string) (proto.ConfigMapServer, error) {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return nil, fmt.Errorf("can not find docker config maps")
	}
	return &DockerServer{
		path: root,
	}, nil
}
