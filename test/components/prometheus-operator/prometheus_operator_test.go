// Copyright 2020 The Lokomotive Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build packet
// +build e2e

package prometheusoperator

import (
	"fmt"
	"testing"
	"time"

	testutil "github.com/kinvolk/lokomotive/test/components/util"
)

func TestPrometheusOperatorDeployment(t *testing.T) {
	namespace := "monitoring"

	client, err := testutil.CreateKubeClient(t)
	if err != nil {
		t.Errorf("could not create Kubernetes client: %v", err)
	}
	t.Log("got kubernetes client")

	deployments := []string{
		"prometheus-operator-operator",
		"prometheus-operator-kube-state-metrics",
		"prometheus-operator-grafana",
	}

	for _, deployment := range deployments {
		t.Run("deployment", func(t *testing.T) {
			t.Parallel()

			testutil.WaitForDeployment(t, client, namespace, deployment, time.Second*5, time.Minute*5)
			t.Logf("Required replicas")
		})
	}

	statefulSets := []string{
		"alertmanager-prometheus-operator-alertmanager",
		"prometheus-prometheus-operator-prometheus",
	}

	for _, statefulset := range statefulSets {
		t.Run(fmt.Sprintf("statefulset %s", statefulset), func(t *testing.T) {
			t.Parallel()
			replicas := 1

			testutil.WaitForStatefulSet(t, client, namespace, statefulset, replicas, time.Second*5, time.Minute*5)
			t.Logf("Required replicas: %d", replicas)
		})
	}

	testutil.WaitForDaemonSet(t, client, namespace, "prometheus-operator-prometheus-node-exporter", time.Second*5, time.Minute*10)
}
