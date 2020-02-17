package model

import (
	monitoringv1alpha1 "github.com/IBM/ibm-monitoring-exporters-operator/pkg/apis/monitoring/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

//KubeStateService creates Service object for kube-state-metrics
func KubeStateService(cr *monitoringv1alpha1.Exporter) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        GetKubeStateObjName(cr),
			Namespace:   cr.Namespace,
			Labels:      getKubeStateLabels(cr),
			Annotations: getKubeStateServiceAnnotations(cr),
		},
		Spec: v1.ServiceSpec{
			Ports:    getKubeStateServicePorts(cr),
			Selector: getKubeStateLabels(cr),
			Type:     "ClusterIP",
		},
	}
}

//UpdatedKubeStateService creates updated Service object for kube-state-metrics
func UpdatedKubeStateService(cr *monitoringv1alpha1.Exporter, currService *v1.Service) *v1.Service {
	newService := currService.DeepCopy()
	newService.ObjectMeta.Labels = getKubeStateLabels(cr)
	newService.Spec.Ports = getKubeStateServicePorts(cr)
	newService.Spec.Selector = getKubeStateLabels(cr)
	return newService
}

//KubeStateDeployment creates deployment object for kube-state-metrics
func KubeStateDeployment(cr *monitoringv1alpha1.Exporter) *appsv1.Deployment {
	containers := []v1.Container{*getKubeStateContainer(cr), *getRouterContainer(cr, KUBE)}
	replicas := int32(1)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetKubeStateObjName(cr),
			Namespace: cr.Namespace,
			Labels:    getKubeStateLabels(cr),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: getKubeStateLabels(cr),
			},
			Replicas: &replicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   GetKubeStateObjName(cr),
					Labels: getKubeStateLabels(cr),
					//TODO: it requires special privelege
					//Annotations: map[string]string{"scheduler.alpha.kubernetes.io/critical-pod": ""},
				},
				Spec: v1.PodSpec{
					//TODO: it requires special privelege
					//PriorityClassName: "system-cluster-critical",
					HostPID:     false,
					HostIPC:     false,
					HostNetwork: false,
					Containers:  containers,
					Volumes:     getVolumes(cr, KUBE),
				},
			},
		},
	}
	if cr.Spec.ImagePullSecrets != nil && len(cr.Spec.ImagePullSecrets) != 0 {
		var secrets []v1.LocalObjectReference
		for _, secret := range cr.Spec.ImagePullSecrets {
			secrets = append(secrets, v1.LocalObjectReference{Name: secret})
		}
		deployment.Spec.Template.Spec.ImagePullSecrets = secrets

	}
	if len(cr.Spec.KubeStateMetrics.ServiceAccount) != 0 {
		deployment.Spec.Template.Spec.ServiceAccountName = cr.Spec.KubeStateMetrics.ServiceAccount
	}
	return deployment
}

//UpdatedKubeStateDeployment created updated deployment for kube-state-metrices
func UpdatedKubeStateDeployment(cr *monitoringv1alpha1.Exporter, currDeployment *appsv1.Deployment) *appsv1.Deployment {
	newDeployment := currDeployment.DeepCopy()
	containers := []v1.Container{*getKubeStateContainer(cr), *getRouterContainer(cr, KUBE)}
	newDeployment.ObjectMeta.Labels = getKubeStateLabels(cr)
	newDeployment.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: getKubeStateLabels(cr),
	}
	newDeployment.Spec.Template.ObjectMeta.Labels = getKubeStateLabels(cr)
	newDeployment.Spec.Template.Spec.Containers = containers
	newDeployment.Spec.Template.Spec.Volumes = getVolumes(cr, KUBE)
	if cr.Spec.ImagePullSecrets != nil && len(cr.Spec.ImagePullSecrets) != 0 {
		var secrets []v1.LocalObjectReference
		for _, secret := range cr.Spec.ImagePullSecrets {
			secrets = append(secrets, v1.LocalObjectReference{Name: secret})
		}
		newDeployment.Spec.Template.Spec.ImagePullSecrets = secrets

	}
	if len(cr.Spec.KubeStateMetrics.ServiceAccount) != 0 {
		newDeployment.Spec.Template.Spec.ServiceAccountName = cr.Spec.KubeStateMetrics.ServiceAccount
	}
	return newDeployment

}

func getKubeStateContainer(cr *monitoringv1alpha1.Exporter) *v1.Container {
	drops := []v1.Capability{"ALL"}
	pe := false
	p := false
	rofs := true
	userID := int64(65534)
	noRoot := true
	probePort := intstr.IntOrString{Type: intstr.Int, IntVal: 8080}
	probe := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Path: "/metrics",
				Port: probePort,
			},
		},
		InitialDelaySeconds: 30,
		TimeoutSeconds:      30,
		PeriodSeconds:       10}

	container := &v1.Container{
		Name:            "kubestatemetrics",
		Image:           cr.Spec.KubeStateMetrics.Image,
		ImagePullPolicy: cr.Spec.ImagePolicy,
		Resources:       cr.Spec.KubeStateMetrics.Resource,
		SecurityContext: &v1.SecurityContext{
			RunAsUser:                &userID,
			RunAsNonRoot:             &noRoot,
			AllowPrivilegeEscalation: &pe,
			Privileged:               &p,
			ReadOnlyRootFilesystem:   &rofs,
			Capabilities: &v1.Capabilities{
				Drop: drops,
			},
		},
		ReadinessProbe: &probe,
		LivenessProbe:  &probe,
	}
	return container

}

func getKubeStateServicePorts(cr *monitoringv1alpha1.Exporter) []v1.ServicePort {
	return []v1.ServicePort{
		v1.ServicePort{
			Name:       "metrics",
			Port:       cr.Spec.KubeStateMetrics.Port,
			TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: cr.Spec.KubeStateMetrics.Port},
			Protocol:   "TCP",
		},
	}

}

//GetKubeStateObjName return name of kube-state-metrics service and deployment
func GetKubeStateObjName(cr *monitoringv1alpha1.Exporter) string {
	return cr.Name + "-kube-state"
}

func getKubeStateLabels(cr *monitoringv1alpha1.Exporter) map[string]string {
	lables := make(map[string]string)
	lables["app"] = "ibm-monitoring"
	lables["component"] = "kube-state-metrics"
	for key, v := range cr.Labels {
		lables[key] = v
	}
	return lables
}

func getKubeStateServiceAnnotations(cr *monitoringv1alpha1.Exporter) map[string]string {
	annotations := make(map[string]string)
	annotations["prometheus.io/scrape"] = "true"
	annotations["prometheus.io/scheme"] = "https"
	return annotations
}
