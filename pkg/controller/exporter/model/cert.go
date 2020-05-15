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
	cmv1alpha1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"

	monitoringv1alpha1 "github.com/IBM/ibm-monitoring-exporters-operator/pkg/apis/monitoring/v1alpha1"

	//cmmetav1 "github.com/jetstack/cert-manager/pkg/apis/meta/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CertManagerCert creates Certificate object
func CertManagerCert(cr *monitoringv1alpha1.Exporter) *cmv1alpha1.Certificate {
	return &cmv1alpha1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.Certs.ExporterSecret,
			Namespace: cr.Namespace,
			Labels:    getCertLabels(cr),
		},
		Spec: cmv1alpha1.CertificateSpec{
			SecretName: cr.Spec.Certs.ExporterSecret,
			IssuerRef: cmv1alpha1.ObjectReference{
				Name: cr.Spec.Certs.Issuer,
				Kind: cmv1alpha1.ClusterIssuerKind,
			},
			CommonName: AppLabekValue,
			DNSNames:   []string{"*." + cr.Namespace + ".pod"},
		},
	}
}

func getCertLabels(cr *monitoringv1alpha1.Exporter) map[string]string {
	labels := make(map[string]string)
	labels[AppLabelKey] = AppLabekValue
	labels = appendCommonLabels(labels)
	for key, v := range cr.Labels {
		labels[key] = v
	}
	return labels
}
