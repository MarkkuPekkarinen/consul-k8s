package mock

import (
	"net/http"
	"net/http/httptest"
)

func NewFakeServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r == nil || r.Method != "GET" {
			return
		}

		if r.URL.Path == "/namespaces/consul/poddisruptionbudgets/consul-server" {
			response := `{
"apiVersion": "policy/v1",
"kind": "PodDisruptionBudget",
"metadata": {
	"annotations": {
		"meta.helm.sh/release-name": "peppertrout",
		"meta.helm.sh/release-namespace": "default"
	},
	"creationTimestamp": "2022-04-01T18:35:55Z",
	"generation": 1,
	"labels": {
		"app": "consul",
		"app.kubernetes.io/managed-by": "Helm",
		"chart": "consul-helm",
		"component": "server",
		"heritage": "Helm",
		"release": "peppertrout"
	},
	"name": "peppertrout-consul-server",
	"namespace": "default",
	"resourceVersion": "33307489",
	"uid": "4fcd9849-1899-4bce-bbce-345ed67bda4f"
},
"spec": {
	"maxUnavailable": 0,
	"selector": {
		"matchLabels": {
			"app": "consul",
			"component": "server",
			"release": "peppertrout"
		}
	}
},
"status": {
	"conditions": [
		{
			"lastTransitionTime": "2022-04-01T18:36:27Z",
			"message": "",
			"observedGeneration": 1,
			"reason": "InsufficientPods",
			"status": "False",
			"type": "DisruptionAllowed"
		}
	],
	"currentHealthy": 1,
	"desiredHealthy": 1,
	"disruptionsAllowed": 0,
	"expectedPods": 1,
	"observedGeneration": 1
}
}`
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(response))
			return
		}
	}))
}
