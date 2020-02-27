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

func (h *Handler) syncNodeExporter() error {
	if err := h.syncNodeExporterDaemonset(); err != nil {
		return err
	}
	if err := h.syncNodeExporterService(); err != nil {
		return err
	}
	return nil
}
func (h *Handler) syncNodeExporterService() error {
	//create
	if h.CR.Spec.NodeExporter.Enable && h.CurrentState.NodeExporterState.Service == nil {
		service := model.NodeExporterService(h.CR)
		if err := h.createObject(service); err != nil {
			log.Error(err, "Fail to create node exporter service")
			return err
		}
		log.Info("Create node exporter service successfully")
	}
	//delete
	if !h.CR.Spec.NodeExporter.Enable && h.CurrentState.NodeExporterState.Service != nil {
		if err := h.Client.Delete(h.Context, h.CurrentState.NodeExporterState.Service); err != nil {
			log.Error(err, "Failed to delete node exporter service")
			return err
		}
		log.Info("Successfully delete node exporter service")
	}
	//update
	if h.CR.Spec.NodeExporter.Enable && h.CurrentState.NodeExporterState.Service != nil {
		service := model.UpdatedNodeExporterService(h.CR, h.CurrentState.NodeExporterState.Service)
		if err := h.updateObject(service); err != nil {
			log.Error(err, "Fail to update node exporter service")
			return err
		}
		log.Info("Successfully update node exporter service")
	}
	return nil
}

func (h *Handler) syncNodeExporterDaemonset() error {
	//create
	if h.CR.Spec.NodeExporter.Enable && h.CurrentState.NodeExporterState.DeamonSet == nil {
		daemonSet := model.NodeExporterDaemonset(h.CR)
		if err := h.createObject(daemonSet); err != nil {
			log.Error(err, "Fail to create node exporter daemonset")
			return err
		}
		log.Info("Create node exporter daemonset successfully")
	}
	//delete
	if !h.CR.Spec.NodeExporter.Enable && h.CurrentState.NodeExporterState.DeamonSet != nil {
		if err := h.Client.Delete(h.Context, h.CurrentState.NodeExporterState.DeamonSet); err != nil {
			log.Error(err, "Failed to delete node exporter daemonset")
			return err
		}
		log.Info("Successfully delete node exporter daemonset")
	}
	//update
	if h.CR.Spec.NodeExporter.Enable && h.CurrentState.NodeExporterState.DeamonSet != nil {
		daemonSet := model.UpdatedNodeExporterDeamonset(h.CR, h.CurrentState.NodeExporterState.DeamonSet)
		if err := h.updateObject(daemonSet); err != nil {
			log.Error(err, "Fail to update node exporter daemonset")
			return err
		}
		log.Info("Successfully update node exporter daemonset")
	}
	return nil
}
