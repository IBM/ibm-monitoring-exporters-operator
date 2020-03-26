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

//ExporterKind means which kind of exporter
type ExporterKind string

//IReqeueError defines interface for requeueError
type IReqeueError interface {
	Reason() string
}
type requeueError struct {
	component string
	reason    string
}

// NewRequeueError creats requeueError error
func NewRequeueError(component string, reason string) error {
	return &requeueError{component, reason}

}
func (r *requeueError) Error() string {
	return "Component " + r.component + "requires to be requeued: " + r.reason
}
func (r *requeueError) Reason() string {
	return r.reason
}

//IsRequeueErr tells if error type is requeueError
func IsRequeueErr(e error) bool {
	switch e.(type) {
	case IReqeueError:
		return true
	}
	return false
}
func commonAnnotationns() map[string]string {
	return map[string]string{
		"clusterhealth.ibm.com/dependencies": "cert-manager",
		"productName":                        "IBM Cloud Platform Common Services",
		"productID":                          "068a62892a1e4db39641342e592daa25",
		"productVersion":                     "3.3.0",
		"productMetric":                      "FREE",
	}
}

const (
	//Ready is a const string
	Ready = "Ready"
	//NotReady is a const string
	NotReady = "NotReady"
	//AppLabelKey is key of label
	AppLabelKey = "cs/app"
	//AppLabekValue is value of label
	AppLabekValue = "ibm-monitoring"
	//TrueStr String of true value
	TrueStr = "true"
	//HTTPSStr string of https
	HTTPSStr = "https"

	//HealthCheckLabelKey lable key for metering check
	HealthCheckLabelKey = "app.kubernetes.io/instance"
	//HealthCheckLabelValue label value for metering check
	HealthCheckLabelValue = "common-monitoring"

	//DefaultNodeExporterSA is default sa for node exporter
	DefaultNodeExporterSA = "ibm-monitoring-exporter"
	//DefaultExporterSA is default sa for other exporters
	DefaultExporterSA = "ibm-monitoring-exporter"

	routerConfigVolName = "router-config"
	routerEntryVolName  = "router-entry"
	caCertsVolName      = "monitoring-ca-certs"
	tlsCertsVolName     = "monitoring-certs"
	procVolName         = "proc"
	sysVolName          = "sys"

	collectdNginxMapKey = "collectd.nginx.conf"
	kubeNginxMapKey     = "kube.nginx.conf"
	nodeNginxMapKey     = "node.nginx.conf"
	routerEntryMapKey   = "entrypoint.sh"

	//COLLECTD means it is collectd exporter
	COLLECTD = ExporterKind("collectd")
	//NODE means it is node-exporter
	NODE = ExporterKind("node")
	//KUBE means it is kube-state-metrics
	KUBE = ExporterKind("kube")

	sslCiphers        = `ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256`
	routerNginxConfig = `
error_log stderr notice;

events {
	worker_connections 1024;
}

http {
	access_log off;

	default_type application/octet-stream;
	sendfile on;
	keepalive_timeout 65;
	server_tokens off;
	more_set_headers "Server: ";

	# Without this, cosocket-based code in worker
	# initialization cannot resolve localhost.

	upstream metrics {
		server 127.0.0.1:{{.UpstreamPort}};
	}
	proxy_cache_path /tmp/nginx-mesos-cache levels=1:2 keys_zone=mesos:1m inactive=10m;

	server {
		listen {{.ListenPort}} ssl default_server;
		ssl_certificate server.crt;
		ssl_certificate_key server.key;
		ssl_client_certificate /opt/ibm/router/caCerts/ca.crt;
		ssl_verify_client on;
		ssl_protocols TLSv1.2;
		# Ref: https://github.com/cloudflare/sslconfig/blob/master/conf
		# Modulo ChaCha20 cipher.
		ssl_ciphers {{.SSLCipers}};
		ssl_prefer_server_ciphers on;

		server_name dcos.*;
		root /opt/ibm/router/nginx/html;

		location / {
			proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
			proxy_set_header Host $http_host;
			proxy_pass http://metrics/;
		}

		location /index.html {
			return 404;
		}                
	}

	{{if ne .HealthyPort 0}}
	server {
		listen {{.HealthyPort}} ssl default_server;
		ssl_certificate server.crt;
		ssl_certificate_key server.key;
		ssl_client_certificate /opt/ibm/router/caCerts/ca.crt;
		ssl_verify_client off;
		ssl_protocols TLSv1.2;
		# Ref: https://github.com/cloudflare/sslconfig/blob/master/conf
		# Modulo ChaCha20 cipher.
		ssl_ciphers {{.SSLCipers}};
		ssl_prefer_server_ciphers on;

		server_name dcos.*;
		root /opt/ibm/router/nginx/html;

		location /healthy {
			proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
			proxy_set_header Host $http_host;
			proxy_pass http://metrics/healthy;
		}
	}
	{{end}}
}
`
	//Be careful. "#!/bin/sh" can NOT start from new line
	routerEntrypoint = `#!/bin/sh
cp -f /opt/ibm/router/certs/tls.crt /opt/ibm/router/nginx/conf/server.crt
cp -f /opt/ibm/router/certs/tls.key /opt/ibm/router/nginx/conf/server.key
cp -f /opt/ibm/router/conf/nginx.conf /opt/ibm/router/nginx/conf/nginx.conf.monitoring
exec nginx -c /opt/ibm/router/nginx/conf/nginx.conf.monitoring -g 'daemon off;'
`
)
