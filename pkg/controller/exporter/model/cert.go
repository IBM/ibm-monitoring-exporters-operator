package model

import (
	monitoringv1alpha1 "github.com/IBM/ibm-monitoring-exporters-operator/pkg/apis/monitoring/v1alpha1"
	cmv1alpha3 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha3"
	cmmetav1 "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CertManagerCert creates Certificate object
func CertManagerCert(cr *monitoringv1alpha1.Exporter) *cmv1alpha3.Certificate {
	return &cmv1alpha3.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.Certs.ExporterSecret,
			Namespace: cr.Namespace,
			Labels:    getCertLabels(cr),
		},
		Spec: cmv1alpha3.CertificateSpec{
			SecretName: cr.Spec.Certs.ExporterSecret,
			IssuerRef: cmmetav1.ObjectReference{
				Name: cr.Spec.Certs.ClusterIssuer,
				Kind: "ClusterIssuer",
			},
		},
	}
}

// OCPCertService creates service which is only used for cert creation
func OCPCertService(cr *monitoringv1alpha1.Exporter) *v1.Service {
	return &v1.Service{
		ObjectMeta: v12.ObjectMeta{
			Name:        cr.Spec.Certs.ExporterSecret,
			Namespace:   cr.Namespace,
			Labels:      getCertLabels(cr),
			Annotations: getCertSvcAnnotations(cr),
		},
		Spec: v1.ServiceSpec{
			Type:  "",
			Ports: getOCPCertServicePorts(cr),
		},
	}
}
func getOCPCertServicePorts(cr *monitoringv1alpha1.Exporter) []v1.ServicePort {
	return []v1.ServicePort{
		v1.ServicePort{
			Name: "nouse",
			Port: 4499,
		},
	}
}
func getCertLabels(cr *monitoringv1alpha1.Exporter) map[string]string {
	lables := make(map[string]string)
	lables["app"] = "ibm-monitoring"
	for key, v := range cr.Labels {
		lables[key] = v
	}
	return lables
}
func getCertSvcAnnotations(cr *monitoringv1alpha1.Exporter) map[string]string {
	annotations := make(map[string]string)
	annotations["service.beta.openshift.io/serving-cert-secret-name"] = cr.Spec.Certs.ExporterSecret
	return annotations
}
