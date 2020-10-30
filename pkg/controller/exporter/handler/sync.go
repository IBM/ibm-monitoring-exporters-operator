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
	"context"

	secv1client "github.com/openshift/client-go/security/clientset/versioned/typed/security/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	monitoringv1alpha1 "github.com/IBM/ibm-monitoring-exporters-operator/pkg/apis/monitoring/v1alpha1"
	"github.com/IBM/ibm-monitoring-exporters-operator/pkg/controller/exporter/model"
)

var log = logf.Log.WithName("controller_exporter.Sync")

// Handler handles a CR and it should be recreated for each reconsile
type Handler struct {
	Context      context.Context
	Client       client.Client
	CR           *monitoringv1alpha1.Exporter
	CurrentState *ClusterState
	Schema       *runtime.Scheme
	SecClient    secv1client.SecurityV1Interface
}

// Sync is entry point of Handler and makes cluster status as expected
func (h *Handler) Sync() error {
	err := model.CreateOrUpdateSCC(h.SecClient)
	if err != nil {
		log.Error(err, "Fail to reconsile SCC")
		return err
	}
	log.Info("SCC is reconciled")
	if err := h.updateStatus(); err != nil {
		return err
	}

	if err := h.syncCertSecret(); err != nil {
		return err
	}

	if err := h.syncRouterConfigmap(); err != nil {
		return err
	}
	if err := h.syncCollectd(); err != nil {
		return err
	}
	if err := h.syncNodeExporter(); err != nil {
		return err
	}
	if err := h.syncKubeStateMetrics(); err != nil {
		return err
	}
	return nil

}
func (h *Handler) updateStatus() error {
	//cert
	if h.CurrentState.CertSecret != nil {
		h.CR.Status.Cert = model.Ready
	} else {
		h.CR.Status.Cert = model.NotReady
	}
	//router configmap
	if h.CurrentState.RouterConfig != nil {
		h.CR.Status.RouterConfigMap = model.Ready
	} else {
		h.CR.Status.RouterConfigMap = model.NotReady
	}
	//collectd
	if h.CurrentState.CollectdState.Deployment != nil {
		h.CR.Status.Collectd = h.CurrentState.CollectdState.Deployment.Status
	}

	//nodeexporter
	if h.CurrentState.NodeExporterState.DeamonSet != nil {
		h.CR.Status.NodeExporter = h.CurrentState.NodeExporterState.DeamonSet.Status
	}
	//kube-state
	if h.CurrentState.KubeStateMetricsState.Deployment != nil {
		h.CR.Status.KubeState = h.CurrentState.CollectdState.Deployment.Status
	}

	if err := h.Client.Status().Update(h.Context, h.CR); err != nil {
		log.Error(err, "Failed to update status")

	}
	return nil
}
func (h *Handler) syncRouterConfigmap() error {
	if h.CurrentState.RouterConfig == nil && !h.CR.Spec.Collectd.Enable {
		return nil
	}
	if h.CurrentState.RouterConfig != nil && h.CR.Spec.Collectd.Enable {
		log.Info("Update router configmap")
		configmap, err := model.RouterConfigmap(h.CR)
		if err != nil {
			log.Error(err, "Failed to create router configmap data model")
			return err
		}
		if err := h.updateObject(configmap); err != nil {
			log.Error(err, "Failed to update router configmap")
		}
		return nil

	}
	if h.CurrentState.RouterConfig == nil && h.CR.Spec.Collectd.Enable {
		log.Info("Create collectd configmap")
		configmap, err := model.RouterConfigmap(h.CR)
		if err != nil {
			log.Error(err, "Failed to create router configmap data model")
			return err
		}
		if err := h.createObject(configmap); err != nil {
			log.Error(err, "Failed to create router configmap")
			return err
		}
		return nil
	}
	if h.CurrentState.RouterConfig != nil && !h.CR.Spec.Collectd.Enable {
		log.Info("Delete router configmap")
		if err := h.Client.Delete(h.Context, h.CurrentState.RouterConfig); err != nil {
			log.Error(err, "Failed to delete router configmap and ignore it")
		}
		return nil
	}
	return nil
}
func (h *Handler) createObject(obj runtime.Object) error {
	if err := controllerutil.SetControllerReference(h.CR, obj.(v1.Object), h.Schema); err != nil {
		return err
	}
	return h.Client.Create(h.Context, obj)
}

func (h *Handler) updateObject(obj runtime.Object) error {
	if err := controllerutil.SetControllerReference(h.CR, obj.(v1.Object), h.Schema); err != nil {
		return err
	}
	if err := h.Client.Update(h.Context, obj); err != nil {
		if kerrors.IsConflict(err) {
			return model.NewRequeueError("sync.UpdateObject", "Object version conflict when updating and requeue it")
		}
		return err

	}
	return nil

}
