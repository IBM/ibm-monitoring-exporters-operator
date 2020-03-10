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

package handler

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/IBM/ibm-monitoring-exporters-operator/pkg/controller/exporter/model"
)

func (h *Handler) syncCertSecret() error {
	cert := model.CertManagerCert(h.CR)
	if h.CurrentState.CertSecret != nil {
		if h.CR.Spec.Certs.AutoClean {
			key := client.ObjectKey{Name: cert.Name, Namespace: cert.Namespace}
			if err := h.Client.Get(h.Context, key, cert); err != nil {
				if kerrors.IsNotFound(err) {
					//secret exists but no certicate
					//delete the secret and create new one
					log.Info("Deleting tls secret which is old and out of control")
					if err = h.Client.Delete(h.Context, h.CurrentState.CertSecret); err != nil {
						log.Error(err, "Failed to delete old tls secret")
						return err
					}

				} else {
					//failed to get certificate because other errors
					log.Error(err, "Failed to get certification object")
					return err

				}
			} else {
				return nil
			}

		} else {
			// when it is not autoclean keep secret no matter who created it
			log.Info("Exporter cert secret exists")
			return nil
		}
	}

	if err := h.createObject(cert); err != nil {
		if kerrors.IsAlreadyExists(err) {
			log.Info("certificate object already exists.")
			return model.NewRequeueError("syncCertSecret", "wait for cert secret to be created after creating certificate object")
		}
		log.Error(err, "Failed to create certificate")
		return err
	}
	// We can not verify if secret is created or not for now so return to next loop
	return model.NewRequeueError("syncCertSecret", "wait for cert secret to be created after creating certification object")

}
