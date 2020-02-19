package handler

import (
	"github.com/IBM/ibm-monitoring-exporters-operator/pkg/controller/exporter/model"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

func (h *Handler) syncCertSecret() error {
	if h.CurrentState.CertSecret != nil {
		log.Info("Exporter cert secret exists")
		return nil
	}
	cert := model.CertManagerCert(h.CR)
	if err := h.createObject(cert); err != nil {
		if kerrors.IsAlreadyExists(err) {
			log.Info("certificate object already exists.")
			return model.NewRequeueError("syncCertSecret", "wait for cert secret to be created after creating certificate object")
		}
		log.Error(err, "Failed to create certificate")
		return err
	}
	// We can not verify if secret is created or not for now. so return to next loop
	return model.NewRequeueError("syncCertSecret", "wait for cert secret to be created after creating certification object")

}
