package e2e

import (
	"testing"

	"github.com/brancz/kube-rbac-proxy/test/kubetest"
	"k8s.io/client-go/kubernetes"
)

func testHTTP2(client kubernetes.Interface) kubetest.TestSuite {
	return func(t *testing.T) {
		command := `HTTP_VERSION=$(curl -sI --http2 --connect-timeout 5 -k --fail -w "%{http_version}\n" -o /dev/null https://kube-rbac-proxy.default.svc.cluster.local:8443/metrics); if [[ "$HTTP_VERSION" == "2" ]]; then echo "Did not expect HTTP/2. Actual protocol: $HTTP_VERSION" > /proc/self/fd/2; exit 1; fi`

		kubetest.Scenario{
			Name: "With failing HTTP2-client",
			Description: `
				Expecting http/2 capable client to fail to connect with http/2.
			`,

			Given: kubetest.Actions(
				kubetest.CreatedManifests(
					client,
					"ignorepaths/clusterRole.yaml",
					"ignorepaths/clusterRoleBinding.yaml",
					"ignorepaths/deployment.yaml",
					"ignorepaths/service.yaml",
					"ignorepaths/serviceAccount.yaml",
					"ignorepaths/clusterRole-client.yaml",
					"ignorepaths/clusterRoleBinding-client.yaml",
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
					nil,
				),
			),
		}.Run(t)
	}
}
