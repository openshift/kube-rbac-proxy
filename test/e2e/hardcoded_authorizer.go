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

		kubetest.Scenario{
			Name: "OpenShift Hardcoded Authorizer",
			Description: `
				Verify that the ServiceAccount prometheus-k8s can access the metrics endpoint
				of the kube-rbac-proxy via the OpenShift hardcoded authorizer
			`,

			Given: kubetest.Actions(
				WithNamespace(client, hardcodedAuthorizerNamespace),
				WithServiceAccount(client, hardcodedAuthorizerClientSA),
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
		}.Run(t)
	}
}

func WithNamespace(client kubernetes.Interface, namespace string) kubetest.Action {
	return func(ctx *kubetest.ScenarioContext) error {
		ctx.Namespace = namespace

		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: ctx.Namespace,
			},
		}

		if _, err := client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{}); err != nil {
			return err
		}

		ctx.AddCleanUp(func() error {
			return client.CoreV1().Namespaces().Delete(context.TODO(), ctx.Namespace, metav1.DeleteOptions{})
		})

		return nil
	}
}

func WithServiceAccount(client kubernetes.Interface, name string) kubetest.Action {
	return func(ctx *kubetest.ScenarioContext) error {
		// Create client service account (prometheus-k8s)
		sa := &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: ctx.Namespace,
			},
		}

		if _, err := client.CoreV1().ServiceAccounts(ctx.Namespace).Create(context.TODO(), sa, metav1.CreateOptions{}); err != nil {
			return err
		}

		ctx.AddCleanUp(func() error {
			return client.CoreV1().ServiceAccounts(ctx.Namespace).Delete(context.TODO(), sa.Name, metav1.DeleteOptions{})
		})

		return nil
	}
}
