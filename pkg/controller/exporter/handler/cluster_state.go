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

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	monitoringv1alpha1 "github.com/IBM/ibm-monitoring-exporters-operator/pkg/apis/monitoring/v1alpha1"
	"github.com/IBM/ibm-monitoring-exporters-operator/pkg/controller/exporter/model"
)

// ClusterState is current status of all exporters  for current CR in cluster
type ClusterState struct {
	CertSecret            *v1.Secret
	RouterConfig          *v1.ConfigMap
	CollectdState         CollectdState
	NodeExporterState     NodeExporterState
	KubeStateMetricsState KubeStateMetricsState
}

func (c *ClusterState) Read(ctx context.Context, cr *monitoringv1alpha1.Exporter, client client.Client) error {
	if err := c.readCertSecret(ctx, cr, client); err != nil {
		return err
	}
	if err := c.readRouterConfigmap(ctx, cr, client); err != nil {
		return err
	}
	if err := c.CollectdState.read(ctx, cr, client); err != nil {
		return err
	}
	if err := c.NodeExporterState.read(ctx, cr, client); err != nil {
		return err
	}
	if err := c.KubeStateMetricsState.read(ctx, cr, client); err != nil {
		return err
	}
	return nil
}
func (c *ClusterState) readCertSecret(ctx context.Context, cr *monitoringv1alpha1.Exporter, cl client.Client) error {
	key := client.ObjectKey{
		Namespace: cr.Namespace,
		Name:      cr.Spec.Certs.ExporterSecret,
	}
	secret := &v1.Secret{}
	if err := cl.Get(ctx, key, secret); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err

	}
	c.CertSecret = secret
	return nil
}

func (c *ClusterState) readRouterConfigmap(ctx context.Context, cr *monitoringv1alpha1.Exporter, cl client.Client) error {
	configmap, err := model.RouterConfigmap(cr)
	if err != nil {
		return err
	}
	key := client.ObjectKey{Name: configmap.Name, Namespace: configmap.Namespace}
	if err := cl.Get(ctx, key, configmap); err != nil {
		if errors.IsNotFound(err) {
			c.RouterConfig = nil
			return nil
		}
		return err
	}
	c.RouterConfig = configmap
	return nil
}

// NodeExporterState is current status of node exporter daemonset in cluster
type NodeExporterState struct {
	Service   *v1.Service
	DeamonSet *appsv1.DaemonSet
}

func (no *NodeExporterState) read(ctx context.Context, cr *monitoringv1alpha1.Exporter, cl client.Client) error {
	//read service
	service := model.NodeExporterService(cr)
	key := client.ObjectKey{Name: service.Name, Namespace: service.Namespace}
	if err := cl.Get(ctx, key, service); err != nil {
		if errors.IsNotFound(err) {
			no.Service = nil
			return nil
		}
		return err
	}
	no.Service = service
	//read deployment
	daemonSet := model.NodeExporterDaemonset(cr)
	key = client.ObjectKey{Name: daemonSet.Name, Namespace: daemonSet.Namespace}
	if err := cl.Get(ctx, key, daemonSet); err != nil {
		if errors.IsNotFound(err) {
			no.DeamonSet = nil
			return nil
		}
		return err
	}
	no.DeamonSet = daemonSet
	return nil
}

// KubeStateMetricsState is current status of kube-state-metrics in cluster
type KubeStateMetricsState struct {
	Service    *v1.Service
	Deployment *appsv1.Deployment
}

func (k *KubeStateMetricsState) read(ctx context.Context, cr *monitoringv1alpha1.Exporter, cl client.Client) error {
	//read service
	service := model.KubeStateService(cr)
	key := client.ObjectKey{Name: service.Name, Namespace: service.Namespace}
	if err := cl.Get(ctx, key, service); err != nil {
		if errors.IsNotFound(err) {
			k.Service = nil
			return nil
		}
		return err
	}
	k.Service = service
	//read deployment
	deployment := model.KubeStateDeployment(cr)
	key = client.ObjectKey{Name: deployment.Name, Namespace: deployment.Namespace}
	if err := cl.Get(ctx, key, deployment); err != nil {
		if errors.IsNotFound(err) {
			k.Deployment = nil
			return nil
		}
		return err
	}
	k.Deployment = deployment
	return nil

}

// CollectdState is current status of collectd service and depoyment in cluster
type CollectdState struct {
	Service    *v1.Service
	Deployment *appsv1.Deployment
}

func (c *CollectdState) read(ctx context.Context, cr *monitoringv1alpha1.Exporter, cl client.Client) error {
	//check service
	service := model.CollectdService(cr)
	key := client.ObjectKey{Namespace: cr.Namespace, Name: model.GetCollectdObjName(cr)}
	if err := cl.Get(ctx, key, service); err != nil {
		if errors.IsNotFound(err) {
			c.Service = nil
			return nil
		}
		return err
	}
	c.Service = service
	//check deployment
	deployment := model.CollectdDeployment(cr)
	if err := cl.Get(ctx, key, deployment); err != nil {
		if errors.IsNotFound(err) {
			c.Deployment = nil
			return nil
		}
		return err

	}
	c.Deployment = deployment
	return nil
}
