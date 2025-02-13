/*
Copyright 2019 The Jetstack cert-manager contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ca

import (
	apiutil "github.com/leki75/cert-manager/pkg/api/util"
	controllerpkg "github.com/leki75/cert-manager/pkg/controller"
	"github.com/leki75/cert-manager/pkg/controller/certificaterequests"
)

const (
	CRControllerName = "certificaterequests-issuer-ca"
)

func init() {
	// create certificate request controller for ca issuer
	controllerpkg.Register(CRControllerName, func(ctx *controllerpkg.Context) (controllerpkg.Interface, error) {
		controller := certificaterequests.New(apiutil.IssuerCA)

		c, err := controllerpkg.New(ctx, CRControllerName, controller)
		if err != nil {
			return nil, err
		}
		return c.Run, nil
	})
}
