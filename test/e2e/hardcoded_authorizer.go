package e2e

import (
	"testing"

	"k8s.io/client-go/kubernetes"

	"github.com/brancz/kube-rbac-proxy/test/kubetest"
)

func testHardcodedAuthorizer(client kubernetes.Interface) kubetest.TestSuite {
	return func(t *testing.T) {
		command := `curl --connect-timeout 5 -v -s -k --fail -H "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" https://kube-rbac-proxy.openshift-monitoring.svc.cluster.local:8443/metrics`
		ctx := &kubetest.ScenarioContext{
			Namespace: "openshift-monitoring",
		}

		s := kubetest.Scenario{
			Name: "OpenShift Hardcoded Authorizer",
			Description: `
				Verify that the ServiceAccount prometheus-k8s can access the metrics endpoint of the kube-rbac-proxy
			`,

			Given: kubetest.Actions(
				kubetest.CreatedManifests(
					client,
					"hardcoded_authorizer/namespace.yaml",
					"hardcoded_authorizer/deployment.yaml",
					"hardcoded_authorizer/clusterRole.yaml",
					"hardcoded_authorizer/clusterRoleBinding.yaml",
					"hardcoded_authorizer/service.yaml",
					"hardcoded_authorizer/serviceAccount.yaml",
					"hardcoded_authorizer/serviceAccount-client.yaml",
				),
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
						ServiceAccount: "prometheus-k8s",
					},
				),
			),
		}

		defer func(ctx *kubetest.ScenarioContext) {
			for _, f := range ctx.CleanUp {
				if err := f(); err != nil {
					panic(err)
				}
			}
		}(ctx)

		t.Run(s.Name, func(t *testing.T) {
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

}
