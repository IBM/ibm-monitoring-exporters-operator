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

	"github.com/IBM/ibm-monitoring-exporters-operator/pkg/controller/exporter/model"
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
	// We can not verify if secret is created or not for now so return to next loop
	return model.NewRequeueError("syncCertSecret", "wait for cert secret to be created after creating certification object")

}
