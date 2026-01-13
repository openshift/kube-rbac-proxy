/*
Copyright 2026 the kube-rbac-proxy maintainers. All rights reserved.

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

package e2e

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/brancz/kube-rbac-proxy/test/kubetest"
	"github.com/brancz/kube-rbac-proxy/test/kubetest/testtemplates"
)

const (
	hardcodedAuthorizerNamespace = "openshift-monitoring"
	hardcodedAuthorizerClientSA  = "prometheus-k8s"
)

func testHardcodedAuthorizer(client kubernetes.Interface) kubetest.TestSuite {
	return func(t *testing.T) {
		command := `curl --connect-timeout 5 -v -s -k --fail -H "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" https://kube-rbac-proxy.openshift-monitoring.svc.cluster.local:8443/metrics`

		runScenarioInNamespace(t, hardcodedAuthorizerNamespace, kubetest.Scenario{
			Name: "OpenShift Hardcoded Authorizer",
			Description: `
				Verify that the ServiceAccount prometheus-k8s can access the metrics endpoint
				of the kube-rbac-proxy via the OpenShift hardcoded authorizer
			`,

			Given: kubetest.Actions(
				setupHardcodedAuthorizerNamespace(client),
				kubetest.NewBasicKubeRBACProxyTestConfig().
					AddSAClusterRoleBinding("kube-rbac-proxy", testtemplates.GetKRPAuthDelegatorRole()).
					Launch(client),
			),
			When: kubetest.Actions(
				kubetest.PodsAreReady(
					client,
					1,
					"app=kube-rbac-proxy",
				),
				kubetest.ServiceIsReady(
					client,
					"kube-rbac-proxy",
				),
			),
			Then: kubetest.Actions(
				kubetest.ClientSucceeds(
					client,
					command,
					&kubetest.RunOptions{
						ServiceAccount: hardcodedAuthorizerClientSA,
					},
				),
			),
		})
	}
}

// runScenarioInNamespace runs a scenario with a custom namespace context.
// This is a local helper to avoid modifying upstream kubetest package.
func runScenarioInNamespace(t *testing.T, namespace string, s kubetest.Scenario) bool {
	ctx := &kubetest.ScenarioContext{
		Namespace: namespace,
	}

	defer func(ctx *kubetest.ScenarioContext) {
		for _, f := range ctx.CleanUp {
			if err := f(); err != nil {
				t.Logf("cleanup error: %v", err)
			}
		}
	}(ctx)

	return t.Run(s.Name, func(t *testing.T) {
		if s.Given != nil {
			if err := s.Given(ctx); err != nil {
				t.Fatalf("failed to create given setup: %v", err)
			}
		}

		if s.When != nil {
			if err := s.When(ctx); err != nil {
				t.Errorf("failed to evaluate state: %v", err)
			}
		}

		if s.Then != nil {
			if err := s.Then(ctx); err != nil {
				t.Errorf("checks failed: %v", err)
			}
		}
	})
}

// setupHardcodedAuthorizerNamespace creates the openshift-monitoring namespace
// and the prometheus-k8s service account for the hardcoded authorizer test.
func setupHardcodedAuthorizerNamespace(client kubernetes.Interface) kubetest.Action {
	return func(ctx *kubetest.ScenarioContext) error {
		// Create namespace
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: hardcodedAuthorizerNamespace,
			},
		}
		if _, err := client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{}); err != nil {
			return err
		}
		ctx.AddCleanUp(func() error {
			return client.CoreV1().Namespaces().Delete(context.TODO(), hardcodedAuthorizerNamespace, metav1.DeleteOptions{})
		})

		// Create client service account (prometheus-k8s)
		sa := &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      hardcodedAuthorizerClientSA,
				Namespace: hardcodedAuthorizerNamespace,
			},
		}
		if _, err := client.CoreV1().ServiceAccounts(hardcodedAuthorizerNamespace).Create(context.TODO(), sa, metav1.CreateOptions{}); err != nil {
			return err
		}
		ctx.AddCleanUp(func() error {
			return client.CoreV1().ServiceAccounts(hardcodedAuthorizerNamespace).Delete(context.TODO(), hardcodedAuthorizerClientSA, metav1.DeleteOptions{})
		})

		return nil
	}
}
