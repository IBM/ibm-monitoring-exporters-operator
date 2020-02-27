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

package model

import (
	v1 "k8s.io/api/core/v1"

	monitoringv1alpha1 "github.com/IBM/ibm-monitoring-exporters-operator/pkg/apis/monitoring/v1alpha1"
)

//getVolumes create volume objects for all kind of exporters
func getVolumes(cr *monitoringv1alpha1.Exporter, exporter ExporterKind) []v1.Volume {
	var volumes []v1.Volume
	routerConfigVol := v1.Volume{
		Name: routerConfigVolName,
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{Name: getRouterConfigmapName(cr)},
			},
		},
	}
	var routerEntryItems []v1.KeyToPath
	routerEntryItems = append(routerEntryItems, v1.KeyToPath{Key: routerEntryMapKey, Path: routerEntryMapKey})
	var defaultMode int32 = 0744
	routerEntryVol := v1.Volume{
		Name: routerEntryVolName,
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{Name: getRouterConfigmapName(cr)},
				Items:                routerEntryItems,
				DefaultMode:          &defaultMode,
			},
		},
	}
	caCertsVol := v1.Volume{
		Name: caCertsVolName,
		VolumeSource: v1.VolumeSource{
			Secret: &v1.SecretVolumeSource{
				SecretName: cr.Spec.Certs.CASecret,
			},
		},
	}
	tlsCertsVol := v1.Volume{
		Name: tlsCertsVolName,
		VolumeSource: v1.VolumeSource{
			Secret: &v1.SecretVolumeSource{
				SecretName: cr.Spec.Certs.ExporterSecret,
			},
		},
	}
	switch exporter {
	case COLLECTD:
		var items []v1.KeyToPath
		items = append(items, v1.KeyToPath{Key: collectdNginxMapKey, Path: "nginx.conf"})
		routerConfigVol.VolumeSource.ConfigMap.Items = items

	case KUBE:
		var items []v1.KeyToPath
		items = append(items, v1.KeyToPath{Key: kubeNginxMapKey, Path: "nginx.conf"})
		routerConfigVol.VolumeSource.ConfigMap.Items = items
	case NODE:
		var items []v1.KeyToPath
		items = append(items, v1.KeyToPath{Key: nodeNginxMapKey, Path: "nginx.conf"})
		routerConfigVol.VolumeSource.ConfigMap.Items = items
	default:
		panic("Impossible exporter type when creating volume objects")

	}

	volumes = append(volumes, routerConfigVol, routerEntryVol, caCertsVol, tlsCertsVol)
	if exporter == NODE {
		procVolume := v1.Volume{
			Name: procVolName,
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/proc",
				},
			},
		}
		sysVolume := v1.Volume{
			Name: sysVolName,
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/sys",
				},
			},
		}
		volumes = append(volumes, procVolume, sysVolume)

	}
	return volumes

}
