package model

import (
	monitoringv1alpha1 "github.com/IBM/ibm-monitoring-exporters-operator/pkg/apis/monitoring/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func getCollectdLabels(cr *monitoringv1alpha1.Exporter) map[string]string {
	lables := make(map[string]string)
	lables["app"] = "ibm-monitoring"
	lables["component"] = "collectdexporter"
	for key, v := range cr.Labels {
		lables[key] = v
	}
	return lables
}

func getCollectdServiceAnnotations(cr *monitoringv1alpha1.Exporter) map[string]string {
	annotations := make(map[string]string)
	annotations["prometheus.io/scrape"] = "true"
	annotations["filter.by.port.name"] = "true"
	annotations["prometheus.io/scheme"] = "https"
	return annotations
}
func getCollectdServicePorts(cr *monitoringv1alpha1.Exporter) []v1.ServicePort {
	return []v1.ServicePort{
		v1.ServicePort{
			Name:       "metrics",
			Port:       cr.Spec.Collectd.MetricsPort,
			TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: cr.Spec.Collectd.MetricsPort},
			Protocol:   "TCP",
		},
		v1.ServicePort{
			Name:       "collector",
			Port:       cr.Spec.Collectd.CollectorPort,
			TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 25826},
			Protocol:   "UDP",
		},
	}

}

//GetCollectdObjName return name of collectd service and deployment
func GetCollectdObjName(cr *monitoringv1alpha1.Exporter) string {
	return cr.Name + "-collectd"
}

// CollectdDeployment creates brand new collectd deployment object
func CollectdDeployment(cr *monitoringv1alpha1.Exporter) *appsv1.Deployment {
	containers := []v1.Container{*getCollectdContainer(cr), *getRouterContainer(cr, COLLECTD)}
	replicas := int32(1)
	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetCollectdObjName(cr),
			Namespace: cr.Namespace,
			Labels:    getCollectdLabels(cr),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: getCollectdLabels(cr),
			},
			Replicas: &replicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   GetCollectdObjName(cr),
					Labels: getCollectdLabels(cr),
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
					Volumes:     getVolumes(cr, COLLECTD),
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
	if len(cr.Spec.Collectd.ServiceAccount) != 0 {
		deployment.Spec.Template.Spec.ServiceAccountName = cr.Spec.Collectd.ServiceAccount
	}

	return &deployment

}
func getCollectdContainer(cr *monitoringv1alpha1.Exporter) *v1.Container {
	drops := []v1.Capability{"ALL"}
	pe := false
	p := false
	rofs := true
	userID := int64(65534)
	noRoot := true
	probePort := intstr.IntOrString{Type: intstr.Int, IntVal: 9103}
	cmdArgs := []string{"--collectd.listen-address=:25826"}
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
	container := v1.Container{
		Name:            "collectd-exporter",
		Image:           cr.Spec.Collectd.Image,
		ImagePullPolicy: cr.Spec.ImagePolicy,
		Resources:       cr.Spec.Collectd.Resource,
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
		Args:           cmdArgs,
		ReadinessProbe: &probe,
		LivenessProbe:  &probe,
	}
	return &container
}

//UpdatedCollectdDeployment update existing collectd deployment object
func UpdatedCollectdDeployment(cr *monitoringv1alpha1.Exporter, currDeployment *appsv1.Deployment) *appsv1.Deployment {
	newDeployment := currDeployment.DeepCopy()
	containers := []v1.Container{*getCollectdContainer(cr), *getRouterContainer(cr, COLLECTD)}

	newDeployment.ObjectMeta.Labels = getCollectdLabels(cr)
	newDeployment.Spec.Selector = &metav1.LabelSelector{MatchLabels: getCollectdLabels(cr)}
	newDeployment.Spec.Template.ObjectMeta.Labels = getCollectdLabels(cr)
	newDeployment.Spec.Template.Spec.Containers = containers
	newDeployment.Spec.Template.Spec.Volumes = getVolumes(cr, COLLECTD)

	if cr.Spec.ImagePullSecrets != nil && len(cr.Spec.ImagePullSecrets) != 0 {
		var secrets []v1.LocalObjectReference
		for _, secret := range cr.Spec.ImagePullSecrets {
			secrets = append(secrets, v1.LocalObjectReference{Name: secret})
		}
		newDeployment.Spec.Template.Spec.ImagePullSecrets = secrets

	}
	if len(cr.Spec.Collectd.ServiceAccount) != 0 {
		newDeployment.Spec.Template.Spec.ServiceAccountName = cr.Spec.Collectd.ServiceAccount
	}

	return newDeployment
}

// CollectdService creates brand new Service object for collectd exporter basing on cr
func CollectdService(cr *monitoringv1alpha1.Exporter) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        GetCollectdObjName(cr),
			Namespace:   cr.Namespace,
			Labels:      getCollectdLabels(cr),
			Annotations: getCollectdServiceAnnotations(cr),
		},
		Spec: v1.ServiceSpec{
			Ports:    getCollectdServicePorts(cr),
			Selector: getCollectdLabels(cr),
			Type:     "ClusterIP",
		},
	}
}

// UpdatedCollectdService generated updated collected service
func UpdatedCollectdService(cr *monitoringv1alpha1.Exporter, currService *v1.Service) *v1.Service {
	newService := currService.DeepCopy()
	newService.Labels = getCollectdLabels(cr)
	newService.Annotations = getCollectdServiceAnnotations(cr)
	newService.Spec.Ports = getCollectdServicePorts(cr)
	newService.Spec.Selector = getCollectdLabels(cr)
	return newService
}
