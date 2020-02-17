package model

import (
	"bytes"
	"html/template"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	monitoringv1alpha1 "github.com/IBM/ibm-monitoring-exporters-operator/pkg/apis/monitoring/v1alpha1"
)

var (
	routerNginxTempl      *template.Template
	routerEntrypointTempl *template.Template
)

//RouterConfigmap create configmap object for colllectd router including nginx config and entrypoint script
func RouterConfigmap(cr *monitoringv1alpha1.Exporter) (*v1.ConfigMap, error) {
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getRouterConfigmapName(cr),
			Namespace: cr.Namespace,
			Labels:    getRouterConfigmapLables(cr),
		},
	}
	datas, err := getConfigmapData(cr)
	if err != nil {
		return cm, err
	}
	cm.Data = datas
	return cm, nil
}
func getRouterConfigmapName(cr *monitoringv1alpha1.Exporter) string {
	return cr.Name + "-router"
}
func getRouterConfigmapLables(cr *monitoringv1alpha1.Exporter) map[string]string {
	labels := make(map[string]string)
	labels["app"] = "ibm-monitoring"
	for key, v := range cr.Labels {
		labels[key] = v
	}
	return labels
}

func getConfigmapData(cr *monitoringv1alpha1.Exporter) (map[string]string, error) {
	datas := make(map[string]string)
	collectdNginx, err := getCollectdNginxConfig(cr)
	if err != nil {
		return datas, err
	}
	nodeNginx, err := getNodeExporterNginxConfig(cr)
	if err != nil {
		return datas, err
	}
	kubestateNginx, err := getKubeStateMetricsNginxConfig(cr)
	if err != nil {
		return datas, err
	}

	datas[collectdNginxMapKey] = collectdNginx
	datas[nodeNginxMapKey] = nodeNginx
	datas[kubeNginxMapKey] = kubestateNginx

	datas[routerEntryMapKey] = routerEntrypoint
	return datas, nil
}
func getCollectdNginxConfig(cr *monitoringv1alpha1.Exporter) (string, error) {
	paras := routerNginxParas{
		UpstreamPort: 9103,
		ListenPort:   cr.Spec.Collectd.MetricsPort,
		SSLCipers:    sslCiphers,
	}
	if cr.Spec.Certs.CACertFileName == "" {
		paras.CACertFileName = "tls.crt"
	} else {
		paras.CACertFileName = cr.Spec.Certs.CACertFileName
	}
	var tplBuffer bytes.Buffer
	if err := routerNginxTempl.Execute(&tplBuffer, paras); err != nil {
		return "", err
	}
	return tplBuffer.String(), nil
}

//getRouterContainer creates router container object whose values are common to all 3 exporters
func getRouterContainer(cr *monitoringv1alpha1.Exporter, exporter ExporterKind) *v1.Container {
	drops := []v1.Capability{"ALL"}
	adds := []v1.Capability{"CHOWN", "NET_ADMIN", "NET_RAW", "LEASE", "SETGID", "SETUID"}
	pe := false
	p := false
	rofs := false
	mounts := []v1.VolumeMount{
		{
			Name:      routerConfigVolName,
			MountPath: "/opt/ibm/router/conf",
		},
		{
			Name:      routerEntryVolName,
			MountPath: "/opt/ibm/router/entry",
		},
		{
			Name:      caCertsVolName,
			MountPath: "/opt/ibm/router/caCerts",
		},
		{
			Name:      tlsCertsVolName,
			MountPath: "/opt/ibm/router/certs",
		},
	}
	container := v1.Container{
		Name:            "router",
		Image:           cr.Spec.RouterImage,
		ImagePullPolicy: cr.Spec.ImagePolicy,
		Command:         []string{"/opt/ibm/router/entry/entrypoint.sh"},
		//Command: []string{"sleep", "3600"},
		SecurityContext: &v1.SecurityContext{
			AllowPrivilegeEscalation: &pe,
			Privileged:               &p,
			ReadOnlyRootFilesystem:   &rofs,
			Capabilities: &v1.Capabilities{
				Drop: drops,
				Add:  adds,
			},
		},
		VolumeMounts: mounts,
	}
	switch exporter {
	case COLLECTD:
		container.Resources = *cr.Spec.Collectd.RouterResource.DeepCopy()
	case KUBE:
		container.Resources = *cr.Spec.KubeStateMetrics.RouterResource.DeepCopy()
	case NODE:
		container.Resources = *cr.Spec.NodeExporter.RouterResource.DeepCopy()
	default:
		panic("Impossible exporter kind when creating router container object")
	}
	return &container
}

func getNodeExporterNginxConfig(cr *monitoringv1alpha1.Exporter) (string, error) {
	paras := routerNginxParas{
		UpstreamPort: cr.Spec.NodeExporter.HostPort,
		ListenPort:   cr.Spec.NodeExporter.ServicePort,
		SSLCipers:    sslCiphers,
		HealthyPort:  cr.Spec.NodeExporter.HealthyPort,
	}
	if cr.Spec.Certs.CACertFileName == "" {
		paras.CACertFileName = "tls.crt"
	} else {
		paras.CACertFileName = cr.Spec.Certs.CACertFileName
	}
	var tplBuffer bytes.Buffer
	if err := routerNginxTempl.Execute(&tplBuffer, paras); err != nil {
		return "", err
	}
	return tplBuffer.String(), nil
}
func getKubeStateMetricsNginxConfig(cr *monitoringv1alpha1.Exporter) (string, error) {
	paras := routerNginxParas{
		UpstreamPort: 8080,
		ListenPort:   cr.Spec.KubeStateMetrics.Port,
		SSLCipers:    sslCiphers,
	}
	if cr.Spec.Certs.CACertFileName == "" {
		paras.CACertFileName = "tls.crt"
	} else {
		paras.CACertFileName = cr.Spec.Certs.CACertFileName
	}
	var tplBuffer bytes.Buffer
	if err := routerNginxTempl.Execute(&tplBuffer, paras); err != nil {
		return "", err
	}
	return tplBuffer.String(), nil
}

type routerNginxParas struct {
	UpstreamPort   int32
	ListenPort     int32
	CACertFileName string
	SSLCipers      string
	HealthyPort    int32
}

func init() {
	routerNginxTempl = template.Must(template.New("router-nginx").Parse(routerNginxConfig))
}
