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
