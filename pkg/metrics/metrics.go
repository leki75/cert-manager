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

// Package metrics contains global structures related to metrics collection
// cert-manager exposes the following metrics:
// certificate_expiration_timestamp_seconds{name, namespace}
package metrics

import (
	"context"
	"crypto/x509"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	"github.com/leki75/cert-manager/pkg/apis/certmanager/v1alpha1"
	cmlisters "github.com/leki75/cert-manager/pkg/client/listers/certmanager/v1alpha1"
	logf "github.com/leki75/cert-manager/pkg/logs"
	"github.com/leki75/cert-manager/pkg/util/errors"
	"github.com/leki75/cert-manager/pkg/util/kube"
)

const (
	// Namespace is the namespace for cert-manager metric names
	namespace                              = "certmanager"
	prometheusMetricsServerAddress         = "0.0.0.0:9402"
	prometheusMetricsServerShutdownTimeout = 5 * time.Second
	prometheusMetricsServerReadTimeout     = 8 * time.Second
	prometheusMetricsServerWriteTimeout    = 8 * time.Second
	prometheusMetricsServerMaxHeaderBytes  = 1 << 20 // 1 MiB
)

// Default set of metrics
var Default = New(logf.NewContext(context.Background(), logf.Log.WithName("metrics")))

var CertificateExpiryTimeSeconds = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "certificate_expiration_timestamp_seconds",
		Help:      "The date after which the certificate expires. Expressed as a Unix Epoch Time.",
	},
	[]string{"name", "namespace"},
)

// ACMEClientRequestCount is a Prometheus summary to collect the number of
// requests made to each endpoint with the ACME client.
var ACMEClientRequestCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "acme_client_request_count",
		Help:      "The number of requests made by the ACME client.",
		Subsystem: "http",
	},
	[]string{"scheme", "host", "path", "method", "status"},
)

// ACMEClientRequestDurationSeconds is a Prometheus summary to collect request
// times for the ACME client.
var ACMEClientRequestDurationSeconds = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  namespace,
		Name:       "acme_client_request_duration_seconds",
		Help:       "The HTTP request latencies in seconds for the ACME client.",
		Subsystem:  "http",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"scheme", "host", "path", "method", "status"},
)

var ControllerSyncCallCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "controller_sync_call_count",
		Help:      "The number of sync() calls made by a controller.",
	},
	[]string{"controller"},
)

// registeredCertificates holds the set of all certificates which are currently
// registered by Prometheus
var registeredCertificates = &struct {
	certificates map[string]struct{}
	mtx          sync.Mutex
}{
	certificates: make(map[string]struct{}),
}

var activeCertificates cmlisters.CertificateLister

type Metrics struct {
	ctx context.Context
	http.Server
	activeCertificates cmlisters.CertificateLister

	// TODO (@dippynark): switch this to use an interface to make it testable
	registry                         *prometheus.Registry
	CertificateExpiryTimeSeconds     *prometheus.GaugeVec
	ACMEClientRequestDurationSeconds *prometheus.SummaryVec
	ACMEClientRequestCount           *prometheus.CounterVec
	ControllerSyncCallCount          *prometheus.CounterVec
}

func New(ctx context.Context) *Metrics {
	router := mux.NewRouter()

	// Create server and register prometheus metrics handler
	s := &Metrics{
		ctx: ctx,
		Server: http.Server{
			Addr:           prometheusMetricsServerAddress,
			ReadTimeout:    prometheusMetricsServerReadTimeout,
			WriteTimeout:   prometheusMetricsServerWriteTimeout,
			MaxHeaderBytes: prometheusMetricsServerMaxHeaderBytes,
			Handler:        router,
		},
		activeCertificates:               nil,
		registry:                         prometheus.NewRegistry(),
		CertificateExpiryTimeSeconds:     CertificateExpiryTimeSeconds,
		ACMEClientRequestDurationSeconds: ACMEClientRequestDurationSeconds,
		ACMEClientRequestCount:           ACMEClientRequestCount,
		ControllerSyncCallCount:          ControllerSyncCallCount,
	}

	router.Handle("/metrics", promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{}))

	return s
}

func (m *Metrics) waitShutdown(stopCh <-chan struct{}) {
	log := logf.FromContext(m.ctx)
	<-stopCh
	log.Info("stopping Prometheus metrics server...")

	ctx, cancel := context.WithTimeout(context.Background(), prometheusMetricsServerShutdownTimeout)
	defer cancel()

	if err := m.Shutdown(ctx); err != nil {
		log.Error(err, "prometheus metrics server shutdown failed", err)
		return
	}

	log.Info("prometheus metrics server gracefully stopped")
}

func (m *Metrics) Start(stopCh <-chan struct{}) {
	log := logf.FromContext(m.ctx)

	m.registry.MustRegister(m.CertificateExpiryTimeSeconds)
	m.registry.MustRegister(m.ACMEClientRequestDurationSeconds)
	m.registry.MustRegister(m.ACMEClientRequestCount)
	m.registry.MustRegister(m.ControllerSyncCallCount)

	go func() {
		log := log.WithValues("address", m.Addr)
		log.Info("listening for connections on")
		if err := m.ListenAndServe(); err != nil {
			log.Error(err, "error running prometheus metrics server")
			return
		}

		log.Info("prometheus metrics server exited")

	}()

	go wait.Until(func() { m.cleanUp() }, time.Minute, stopCh)

	m.waitShutdown(stopCh)
}

// UpdateCertificateExpiry updates the expiry time of a certificate
func (m *Metrics) UpdateCertificateExpiry(crt *v1alpha1.Certificate, secretLister corelisters.SecretLister) {
	log := logf.FromContext(m.ctx)
	log = logf.WithResource(log, crt)
	log = logf.WithRelatedResourceName(log, crt.Spec.SecretName, crt.Namespace, "Secret")

	log.V(logf.DebugLevel).Info("attempting to retrieve secret for certificate")
	// grab existing certificate
	cert, err := kube.SecretTLSCert(m.ctx, secretLister, crt.Namespace, crt.Spec.SecretName)
	if err != nil {
		if !apierrors.IsNotFound(err) && !errors.IsInvalidData(err) {
			log.Error(err, "error reading secret for certificate")
		}
		return
	}

	updateX509Expiry(crt, cert)
}

func updateX509Expiry(crt *v1alpha1.Certificate, cert *x509.Certificate) {
	expiryTime := cert.NotAfter
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(crt)
	if err != nil {
		return
	}

	registeredCertificates.mtx.Lock()
	defer registeredCertificates.mtx.Unlock()
	// set certificate expiry time
	CertificateExpiryTimeSeconds.With(prometheus.Labels{
		"name":      crt.Name,
		"namespace": crt.Namespace}).Set(float64(expiryTime.Unix()))
	registeredCertificates.certificates[key] = struct{}{}
}

func (m *Metrics) SetActiveCertificates(cl cmlisters.CertificateLister) {
	m.activeCertificates = cl
}

func (m *Metrics) cleanUp() {
	log := logf.FromContext(m.ctx)
	log.V(logf.DebugLevel).Info("attempting to clean up metrics for recently deleted certificates")

	if activeCertificates == nil {
		log.V(logf.DebugLevel).Info("active certificates is still uninitialized")
		return
	}

	activeCrts, err := activeCertificates.List(labels.Everything())
	if err != nil {
		log.Error(err, "error retrieving active certificates")
		return
	}

	cleanUpCertificates(activeCrts)
}

func cleanUpCertificates(activeCrts []*v1alpha1.Certificate) {
	activeMap := make(map[string]struct{}, len(activeCrts))
	for _, crt := range activeCrts {
		key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(crt)
		if err != nil {
			continue
		}

		activeMap[key] = struct{}{}
	}

	registeredCertificates.mtx.Lock()
	defer registeredCertificates.mtx.Unlock()
	var toCleanUp []string
	for key := range registeredCertificates.certificates {
		if _, found := activeMap[key]; !found {
			toCleanUp = append(toCleanUp, key)
		}
	}

	for _, key := range toCleanUp {
		namespace, name, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			continue
		}

		CertificateExpiryTimeSeconds.Delete(prometheus.Labels{
			"name":      name,
			"namespace": namespace,
		})
		delete(registeredCertificates.certificates, key)
	}
}

func (m *Metrics) IncrementSyncCallCount(controllerName string) {
	log := logf.FromContext(m.ctx)
	log.V(logf.DebugLevel).Info("incrementing controller sync call count", "controllerName", controllerName)
	ControllerSyncCallCount.WithLabelValues(controllerName).Inc()
}
