/*
Copyright 2020 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/hasheddan/crisscross/pkg/models"
)

// Client is a client for interacting with a remote controller.
type Client interface {
	Observe(resource.Managed) (managed.ExternalObservation, error)
	Create(resource.Managed) (managed.ExternalCreation, error)
	Update(resource.Managed) (managed.ExternalUpdate, error)
	Delete(resource.Managed) error
}

// ControllerClient communicates with a remote controller.
type ControllerClient struct {
	client   *http.Client
	endpoint string
}

// NewControllerClient creates a new ControllerClient.
func NewControllerClient(endpoint string) *ControllerClient {
	return &ControllerClient{
		client:   &http.Client{},
		endpoint: endpoint,
	}
}

// Observe makes an observe request to the remote controller.
func (c *ControllerClient) Observe(mg resource.Managed) (managed.ExternalObservation, error) {
	o := &models.ObservationResponse{}
	if err := c.post(models.ObservationRequest{
		Managed: mg,
	}, o, "/observe"); err != nil {
		return managed.ExternalObservation{}, err
	}
	return managed.ExternalObservation{
		ResourceExists:          o.External.ResourceExists,
		ResourceUpToDate:        o.External.ResourceUpToDate,
		ResourceLateInitialized: o.External.ResourceLateInitialized,
		ConnectionDetails:       managed.ConnectionDetails(o.External.ConnectionDetails),
	}, nil
}

// Create makes an create request to the remote controller.
func (c *ControllerClient) Create(mg resource.Managed) (managed.ExternalCreation, error) {
	o := &models.CreationResponse{}
	if err := c.post(models.CreationRequest{
		Managed: mg,
	}, o, "/create"); err != nil {
		return managed.ExternalCreation{}, err
	}
	return managed.ExternalCreation{
		ConnectionDetails: managed.ConnectionDetails(o.External.ConnectionDetails),
	}, nil
}

// Update makes a update request to the remote controller.
func (c *ControllerClient) Update(mg resource.Managed) (managed.ExternalUpdate, error) {
	o := &models.UpdateResponse{}
	if err := c.post(models.UpdateRequest{
		Managed: mg,
	}, o, "/update"); err != nil {
		return managed.ExternalUpdate{}, err
	}
	return managed.ExternalUpdate{
		ConnectionDetails: managed.ConnectionDetails(o.External.ConnectionDetails),
	}, nil
}

// Delete makes an delete request to the remote controller.
func (c *ControllerClient) Delete(mg resource.Managed) error {
	o := &models.DeletionResponse{}
	if err := c.post(models.DeletionRequest{
		Managed: mg,
	}, o, "/delete"); err != nil {
		return err
	}
	return nil
}

func (c *ControllerClient) post(mg interface{}, res interface{}, path string) error {
	payload, err := json.Marshal(mg)
	if err != nil {
		return err
	}
	resp, err := c.client.Post(c.endpoint+path, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(&res)
}
