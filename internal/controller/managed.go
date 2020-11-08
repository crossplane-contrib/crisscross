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

package controller

import (
	"context"

	"github.com/hasheddan/crisscross/internal/client"

	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
)

type connector struct {
	endpoint string
}

// Connect creates a client to connect to the remote controller.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	return &external{client: client.NewControllerClient(c.endpoint)}, nil
}

type external struct {
	client client.Client
}

// Observe makes an observe request to the remote controller.
func (e *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	return e.client.Observe(mg)
}

// Create makes a create request to the remote controller.
func (e *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	return e.client.Create(mg)
}

// Update makes an update request to the remote controller.
func (e *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	return e.client.Update(mg)
}

// Delete makes a delete request to the remote controller
func (e *external) Delete(ctx context.Context, mg resource.Managed) error {
	return e.client.Delete(mg)
}
