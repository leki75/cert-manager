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

package controller

import (
	"context"
	"crypto/x509"
	"time"

	cmapi "github.com/leki75/cert-manager/pkg/apis/certmanager/v1alpha1"
	"github.com/leki75/cert-manager/pkg/logs"
)

func (o IssuerOptions) ResourceNamespace(iss cmapi.GenericIssuer) string {
	ns := iss.GetObjectMeta().Namespace
	if ns == "" {
		ns = o.ClusterResourceNamespace
	}
	return ns
}

func (o IssuerOptions) CanUseAmbientCredentials(iss cmapi.GenericIssuer) bool {
	switch iss.(type) {
	case *cmapi.ClusterIssuer:
		return o.ClusterIssuerAmbientCredentials
	case *cmapi.Issuer:
		return o.IssuerAmbientCredentials
	}
	return false
}

func (o IssuerOptions) CertificateNeedsRenew(ctx context.Context, cert *x509.Certificate, crt *cmapi.Certificate) bool {
	return o.CalculateDurationUntilRenew(ctx, cert, crt) <= 0
}

// to help testing
var now = time.Now

// CalculateDurationUntilRenew calculates how long cert-manager should wait to
// until attempting to renew this certificate resource.
func (o IssuerOptions) CalculateDurationUntilRenew(ctx context.Context, cert *x509.Certificate, crt *cmapi.Certificate) time.Duration {
	log := logs.FromContext(ctx, "CalculateDurationUntilRenew")

	// validate if the certificate received was with the issuer configured
	// duration. If not we generate an event to warn the user of that fact.
	certDuration := cert.NotAfter.Sub(cert.NotBefore)
	if crt.Spec.Duration != nil && certDuration < crt.Spec.Duration.Duration {
		log.Info("requested certificate validity period differs from period given on returned certificate", "requested_duration", crt.Spec.Duration.Duration, "actual_duration", certDuration)
		// TODO Use the message as the reason in a 'renewal status' condition
	}

	// renew is the duration before the certificate expiration that cert-manager
	// will start to try renewing the certificate.
	renewBefore := o.RenewBeforeExpiryDuration
	if crt.Spec.RenewBefore != nil {
		renewBefore = crt.Spec.RenewBefore.Duration
	}

	// Verify that the renewBefore duration is inside the certificate validity duration.
	// If not we notify with an event that we will renew the certificate
	// before (certificate duration / 3) of its expiration duration.
	if renewBefore > certDuration {
		log.Info("certificate renewal duration was changed to fit inside the received certificate validity duration from issuer.")
		// TODO Use the message as the reason in a 'renewal status' condition
		// We will renew 1/3 before the expiration date.
		renewBefore = certDuration / 3
	}

	// calculate the amount of time until expiry
	durationUntilExpiry := cert.NotAfter.Sub(now())
	// calculate how long until we should start attempting to renew the certificate
	renewIn := durationUntilExpiry - renewBefore

	return renewIn
}
