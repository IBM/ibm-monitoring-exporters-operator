package handler

import (
	"errors"

	"github.com/IBM/ibm-monitoring-exporters-operator/pkg/controller/exporter/model"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

func (h *Handler) syncCertSecret() error {
	//TODO: it should handle exporter client secret too.
	//Go back here when cert manager go module conflict issue being fixed
	if h.CurrentState.CertSecret != nil {
		log.Info("Exporter cert secret exists")
		return nil
	}
	if h.CR.Spec.Certs.Provider == "" {
		log.Error(nil, "No provider specified and secret "+h.CR.Spec.Certs.ExporterSecret+" does not exist")
		return errors.New("No provider specified and secret " + h.CR.Spec.Certs.ExporterSecret + " does not exist")
	}
	if h.CR.Spec.Certs.Provider == "certmanager" {
		//TODO:
		//certmanager api version is alpha3 now but it is alpha1 in cs
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
		//return errors.New("Wait for secret for next loop")
		return model.NewRequeueError("syncCertSecret", "wait for cert secret to be created after creating certification object")

	}
	if h.CR.Spec.Certs.Provider == "ocp" {
		certService := model.OCPCertService(h.CR)
		if err := h.createObject(certService); err != nil {
			if kerrors.IsAlreadyExists(err) {
				log.Info("Service for cert creation already exists.")
				return model.NewRequeueError("syncCertSecret", "wait for cert secret to be created after creating service object")

			}
			log.Error(err, "Faild to create service for cert creation")
			return err

		}
		log.Info("Successfully create cert service")
		// We can not verify if secret is created or not for now. so return to next loop
		//return errors.New("Wait for secret for next loop")
		return model.NewRequeueError("syncCertSecret", "wait for cert secret to be created after creating service object")

	}
	log.Error(nil, "Unsupported exporter cert provider")
	return errors.New("Unsupported exporter cert provider: " + h.CR.Spec.Certs.Provider)

}
