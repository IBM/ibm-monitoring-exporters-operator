//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package handler

import (
	"github.com/IBM/ibm-monitoring-exporters-operator/pkg/controller/exporter/model"
)

func (h *Handler) syncCollectd() error {
	if err := h.syncCollectdService(); err != nil {
		return err
	}

	if err := h.syncCollectdDeployment(); err != nil {
		return err
	}
	return nil
}
func (h *Handler) syncCollectdService() error {
	//create service
	if h.CurrentState.CollectdState.Service == nil && h.CR.Spec.Collectd.Enable {
		service := model.CollectdService(h.CR)
		if err := h.createObject(service); err != nil {
			log.Error(err, "Failed to create collectd service")
			return err
		}
		log.Info("Successfully create collectd service")
	}
	//delete service
	if h.CurrentState.CollectdState.Service != nil && !h.CR.Spec.Collectd.Enable {
		if err := h.Client.Delete(h.Context, h.CurrentState.CollectdState.Service); err != nil {
			log.Error(err, "Failed to delete collectd service")
			return err
		}
		log.Info("Successfully delete collectd service")
	}
	// update service
	if h.CurrentState.CollectdState.Service != nil && h.CR.Spec.Collectd.Enable {
		newService := model.UpdatedCollectdService(h.CR, h.CurrentState.CollectdState.Service)
		if err := h.updateObject(newService); err != nil {
			log.Error(err, "Failed to update collectd Service")
			return err
		}
		log.Info("Successfully update collectd service")
	}

	return nil
}
func (h *Handler) syncCollectdDeployment() error {
	//create
	if h.CR.Spec.Collectd.Enable && h.CurrentState.CollectdState.Deployment == nil {
		deployment := model.CollectdDeployment(h.CR)
		if err := h.createObject(deployment); err != nil {
			log.Error(err, "Fail to create collectd deployment")
			return err
		}
		log.Info("Create collectd deployment successfully")

	}
	//delete
	if !h.CR.Spec.Collectd.Enable && h.CurrentState.CollectdState.Deployment != nil {
		if err := h.Client.Delete(h.Context, h.CurrentState.CollectdState.Deployment); err != nil {
			log.Error(err, "Failed to delete collectd deployment")
			return err
		}
		log.Info("Successfully delete collectd deployment")
	}
	//update
	if h.CR.Spec.Collectd.Enable && h.CurrentState.CollectdState.Deployment != nil {
		newDeployment := model.UpdatedCollectdDeployment(h.CR, h.CurrentState.CollectdState.Deployment)
		if err := h.updateObject(newDeployment); err != nil {
			log.Error(err, "Failed to update collectd deployment")
			return err
		}
		log.Info("Successfully update collectd deployment")
	}
	return nil
}
