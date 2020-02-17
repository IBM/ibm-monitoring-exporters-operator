package controller

import (
	"github.com/IBM/ibm-monitoring-exporters-operator/pkg/controller/exporter"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, exporter.Add)
}
