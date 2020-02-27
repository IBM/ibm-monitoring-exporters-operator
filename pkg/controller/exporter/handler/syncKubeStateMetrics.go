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

func (h *Handler) syncKubeStateMetrics() error {
	if err := h.syncKubeStateDeployment(); err != nil {
		return err
	}
	if err := h.syncKubeStateService(); err != nil {
		return err
	}
	return nil
}
func (h *Handler) syncKubeStateDeployment() error {
	//create
	if h.CR.Spec.KubeStateMetrics.Enable && h.CurrentState.KubeStateMetricsState.Deployment == nil {
		deployment := model.KubeStateDeployment(h.CR)
		if err := h.createObject(deployment); err != nil {
			log.Error(err, "Fail to create kube-state-metrics deployment")
			return err
		}
		log.Info("Create kube-state-metrics deployment successfully")
	}
	//delete
	if !h.CR.Spec.KubeStateMetrics.Enable && h.CurrentState.KubeStateMetricsState.Deployment != nil {
		if err := h.Client.Delete(h.Context, h.CurrentState.KubeStateMetricsState.Deployment); err != nil {
			log.Error(err, "Failed to delete kube-state-metrics deployment")
			return err
		}
		log.Info("Successfully delete kube-state-metrics deployment")
	}
	//update
	if h.CR.Spec.KubeStateMetrics.Enable && h.CurrentState.KubeStateMetricsState.Deployment != nil {
		deployment := model.UpdatedKubeStateDeployment(h.CR, h.CurrentState.KubeStateMetricsState.Deployment)
		if err := h.updateObject(deployment); err != nil {
			log.Error(err, "Fail to update kube-state-metrics deployment")
			return err
		}
		log.Info("Successfully update kube-state-metrics deployment")
	}
	return nil
}
func (h *Handler) syncKubeStateService() error {
	//create
	if h.CR.Spec.KubeStateMetrics.Enable && h.CurrentState.KubeStateMetricsState.Service == nil {
		service := model.KubeStateService(h.CR)
		if err := h.createObject(service); err != nil {
			log.Error(err, "Fail to create kube-state-metrics service")
			return err
		}
		log.Info("Create kube-state-metrics service successfully")
	}
	//delete
	if !h.CR.Spec.KubeStateMetrics.Enable && h.CurrentState.KubeStateMetricsState.Service != nil {
		if err := h.Client.Delete(h.Context, h.CurrentState.KubeStateMetricsState.Service); err != nil {
			log.Error(err, "Failed to delete kube-state-metrics service")
			return err
		}
		log.Info("Successfully delete kube-state-metrics service")
	}
	//update
	if h.CR.Spec.KubeStateMetrics.Enable && h.CurrentState.KubeStateMetricsState.Service != nil {
		service := model.UpdatedKubeStateService(h.CR, h.CurrentState.KubeStateMetricsState.Service)
		if err := h.updateObject(service); err != nil {
			log.Error(err, "Fail to update kube-state-metrics service")
			return err
		}
		log.Info("Successfully update kube-state-metrics service")
	}
	return nil
}
