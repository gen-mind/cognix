package main

import (
	"cognix.ch/api/v2/core/proto"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
)

// nsSecret is a constant representing the file path to the Kubernetes Service Account Namespace Secret.
const nsSecret = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

// K8SServer represents a Kubernetes server with additional methods for ConfigMap operations.
//
// It embeds the proto.UnsafeConfigMapServer interface and includes the following fields:
// - namespace: the namespace in which the ConfigMaps are managed
// - client: a Kubernetes clientset used to interact with the Kubernetes cluster
//
// The K8SServer type provides the following methods:
// - GetList: retrieves the list of ConfigMaps in the namespace
// - Save: saves a ConfigMap to the namespace
// - Delete: deletes a ConfigMap from the namespace
//
// Example Usage:
//
// server, err := NewK8SServer()
//
//	if err != nil {
//	  // handle error
//	}
//
// ctx := context.Background()
// listReq := &proto.ConfigMapList{Name: "example"}
// listResp, err := server.GetList(ctx, listReq)
//
//	if err != nil {
//	  // handle error
//	}
//
//	saveReq := &proto.ConfigMapSave{
//	  Name: "example",
//	  Value: &proto.ConfigMapValue{
//	    Key:   "key",
//	    Value: "value",
//	  },
//	}
//
// _, err = server.Save(ctx, saveReq)
//
//	if err != nil {
//	  // handle error
//	}
//
// deleteReq := &proto.ConfigMapDelete{Name: "example", Key: "key"}
// _, err = server.Delete(ctx, deleteReq)
//
//	if err != nil {
//	  // handle error
//	}
type K8SServer struct {
	proto.UnsafeConfigMapServer
	namespace string
	client    *kubernetes.Clientset
}

// GetList retrieves a list of config maps based on the provided ConfigMapList request.
// The method first fetches the config map identified by the name in the request.
// It then constructs a ConfigMapListResponse containing the key-value pairs of the config map.
// Finally, it returns the resulting ConfigMapListResponse or an error if any occurred.
func (k *K8SServer) GetList(ctx context.Context, r *proto.ConfigMapList) (*proto.ConfigMapListResponse, error) {
	configMap, err := k.client.CoreV1().ConfigMaps(k.namespace).Get(ctx, r.GetName(), v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	result := proto.ConfigMapListResponse{
		Values: make([]*proto.ConfigMapRecord, 0),
	}
	for key, value := range configMap.Data {
		result.Values = append(result.Values, &proto.ConfigMapRecord{
			Key:   key,
			Value: value,
		})
	}
	return &result, nil
}

// Save saves the value of a specific key in a ConfigMap.
func (k *K8SServer) Save(ctx context.Context, r *proto.ConfigMapSave) (*empty.Empty, error) {
	configMap, err := k.client.CoreV1().ConfigMaps(k.namespace).Get(ctx, r.GetName(), v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	configMap.Data[r.GetValue().GetKey()] = r.GetValue().GetValue()
	if _, err = k.client.CoreV1().ConfigMaps(k.namespace).Update(ctx, configMap, v1.UpdateOptions{}); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

// Delete deletes a data entry from a ConfigMap.
//
// The Delete method retrieves the ConfigMap with the given name using the specified
// namespace and gets the data entry corresponding to the specified key. It deletes
// the entry from the ConfigMap by removing it from the map. Finally, it updates the
// ConfigMap in the Kubernetes cluster with the updated data and returns an empty
// response if the operation is successful, or an error if any error occurs while
// performing the deletion or update.
func (k *K8SServer) Delete(ctx context.Context, r *proto.ConfigMapDelete) (*empty.Empty, error) {
	configMap, err := k.client.CoreV1().ConfigMaps(k.namespace).Get(ctx, r.GetName(), v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	delete(configMap.Data, r.GetKey())

	if _, err = k.client.CoreV1().ConfigMaps(k.namespace).Update(ctx, configMap, v1.UpdateOptions{}); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

// NewK8SServer returns a new instance of K8SServer that implements the ConfigMapServer interface and
// communicates with Kubernetes API to perform CRUD operations on ConfigMap resources.
//
// The function first retrieves the in-cluster configuration using rest.InClusterConfig().
// Then it creates a new Kubernetes clientset using the obtained configuration.
// It then calls the getCurrentNamespace() function to get the current namespace that the server will operate on.
// Finally, it creates a new instance of K8SServer with the client and namespace, and returns it along with any error encountered.
// If any error occurs during the process, nil is returned as the server instance and the error is propagated.
func NewK8SServer() (proto.ConfigMapServer, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	ns, err := getCurrentNamespace()
	if err != nil {
		return nil, err
	}
	return &K8SServer{
		client:    client,
		namespace: ns,
	}, nil
}

// getCurrentNamespace retrieves the current Kubernetes namespace by reading the contents
// of the nsSecret file. If an error occurs while reading the file, an empty string and
// the error are returned.
func getCurrentNamespace() (string, error) {
	ns, err := os.ReadFile(nsSecret)
	if err != nil {
		return "", err
	}
	return string(ns), nil
}
