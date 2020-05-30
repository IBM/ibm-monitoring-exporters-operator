# ibm-monitoring-exporters-operator

> **Important:** Do not install this operator directly. Only install this operator using the IBM Common Services Operator. For more information about installing this operator and other Common Services operators, see [Installer documentation](http://ibm.biz/cpcs_opinstall). If you are using this operator as part of an IBM Cloud Pak, see the documentation for that IBM Cloud Pak to learn more about how to install and use the operator service. For more information about IBM Cloud Paks, see [IBM Cloud Paks that use Common Services](http://ibm.biz/cpcs_cloudpaks).

- ibm-monitoring-exporters-operator is used for installing Prometheus Exporters.

The operator is part of IBM Monitoring Operator stack. It installs node-exporter, kube-state-metrics, and collectd exporters. The IBM Monitoring Operator stack is integrated in IBM Cloud Paks and installed as part of IBM Cloud Platform Common Services. You can also use the Operand Deployment Lifecycle Manager Operator to install IBM Cloud Platform Common Services in your cluster. See the following documentation for prerequisites required to run the stand-alone Operator.

## Supported platforms

Red Hat OpenShift Container Platform 4.x or newer installed on one of the following platforms:
- Linux x86_64
- Linux on Power (ppc64le)
- Linux on IBM Z and LinuxONE

## Operator versions

- 1.8.0
- 1.9.0

## Prerequisites

Before you install this operator, you need to first install the operator dependencies and prerequisites:

- For the list of operator dependencies, see the IBM Knowledge Center [Common Services dependencies documentation](http://ibm.biz/cpcs_opdependencies).

- For the list of prerequisites for installing the operator, see the IBM Knowledge Center [Preparing to install services documentation](http://ibm.biz/cpcs_opinstprereq).

## Documentation

To install the operator with the IBM Common Services Operator follow the installation and configuration instructions within the IBM Knowledge Center.

- If you are using the operator as part of an IBM Cloud Pak, see the documentation for that IBM Cloud Pak. For a list of IBM Cloud Paks, see [IBM Cloud Paks that use Common Services](http://ibm.biz/cpcs_cloudpaks).
- If you are using the operator with an IBM Containerized Software, see the IBM Cloud Platform Common Services Knowledge Center [Installer documentation](http://ibm.biz/cpcs_opinstall).

## SecurityContextConstraints Requirements

The ibm-monitoring-exporters-operator supports running under the OpenShift Container Platform default restricted security context constraints. The kube-state-metrics, node-exporter and collectd-exporter run under privileged security constraints.
For more information about the OpenShift Container Platform Security Context Constraints, see [Managing Security Context Constraints](https://docs.openshift.com/container-platform/4.3/authentication/managing-security-context-constraints.html).

## Developer guide

As a developer, if you want to build and test this operator to try out and learn more about the operator and its capabilities, you can use the following developer guide. The guide provides commands for a quick installation and initial validation for running the operator.

> **Important:** The following developer guide is provided as-is and only for trial and education purposes. IBM and IBM Support does not provide any support for the usage of the operator with this developer guide. For the official supported installation and usage guide for the operator, see the IBM Knowledge Center documentation for your IBM Cloud Pak or for IBM Cloud Platform Common Services.

### End-to-End testing

For more instructions on how to run end-to-end testing with the Operand Deployment Lifecycle Manager, see [ODLM guide](https://github.com/IBM/operand-deployment-lifecycle-manager/blob/master/docs/install/common-service-integration.md#end-to-end-test).
