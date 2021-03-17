// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package features

import (
	"fmt"
	"net/http"

	"github.com/apache/apisix-ingress-controller/test/e2e/scaffold"
	"github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
)

var _ = ginkgo.Describe("nginx vars", func() {
	opts := &scaffold.Options{
		Name:                    "default",
		Kubeconfig:              scaffold.GetKubeconfig(),
		APISIXConfigPath:        "testdata/apisix-gw-config.yaml",
		APISIXDefaultConfigPath: "testdata/apisix-gw-config-default.yaml",
		IngressAPISIXReplicas:   1,
		HTTPBinServicePort:      80,
		APISIXRouteVersion:      "apisix.apache.org/v2alpha1",
	}
	s := scaffold.NewScaffold(opts)
	ginkgo.It("operator is equal", func() {
		backendSvc, backendPorts := s.DefaultHTTPBackend()

		ar := fmt.Sprintf(`
apiVersion: apisix.apache.org/v2alpha1
kind: ApisixRoute
metadata:
 name: httpbin-route
spec:
 http:
 - name: rule1
   match:
     hosts:
     - httpbin.org
     paths:
       - /ip
     nginxVars:
     - subject: http_x_foo
       op: Equal
       value: bar
   backend:
     serviceName: %s
     servicePort: %d
`, backendSvc, backendPorts[0])

		assert.Nil(ginkgo.GinkgoT(), s.CreateResourceFromString(ar))

		err := s.EnsureNumApisixRoutesCreated(1)
		assert.Nil(ginkgo.GinkgoT(), err, "Checking number of routes")
		err = s.EnsureNumApisixUpstreamsCreated(1)
		assert.Nil(ginkgo.GinkgoT(), err, "Checking number of upstreams")

		_ = s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			WithHeader("X-Foo", "bar").
			Expect().
			Status(http.StatusOK)

		msg := s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			WithHeader("X-Foo", "baz").
			Expect().
			Status(http.StatusNotFound).
			Body().
			Raw()
		assert.Contains(ginkgo.GinkgoT(), msg, "404 Route Not Found")
	})

	ginkgo.It("operator is not_equal", func() {
		backendSvc, backendPorts := s.DefaultHTTPBackend()

		ar := fmt.Sprintf(`
apiVersion: apisix.apache.org/v2alpha1
kind: ApisixRoute
metadata:
 name: httpbin-route
spec:
 http:
 - name: rule1
   match:
     hosts:
     - httpbin.org
     paths:
       - /ip
     nginxVars:
     - subject: http_x_foo
       op: NotEqual
       value: bar
   backend:
     serviceName: %s
     servicePort: %d
`, backendSvc, backendPorts[0])

		assert.Nil(ginkgo.GinkgoT(), s.CreateResourceFromString(ar))

		err := s.EnsureNumApisixRoutesCreated(1)
		assert.Nil(ginkgo.GinkgoT(), err, "Checking number of routes")
		err = s.EnsureNumApisixUpstreamsCreated(1)
		assert.Nil(ginkgo.GinkgoT(), err, "Checking number of upstreams")

		_ = s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			Expect().
			Status(http.StatusOK)

		msg := s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			WithHeader("X-Foo", "bar").
			Expect().
			Status(http.StatusNotFound).
			Body().
			Raw()
		assert.Contains(ginkgo.GinkgoT(), msg, "404 Route Not Found")
	})

	ginkgo.It("operator is greater_than", func() {
		backendSvc, backendPorts := s.DefaultHTTPBackend()

		ar := fmt.Sprintf(`
apiVersion: apisix.apache.org/v2alpha1
kind: ApisixRoute
metadata:
 name: httpbin-route
spec:
 http:
 - name: rule1
   match:
     hosts:
     - httpbin.org
     paths:
       - /ip
     nginxVars:
     - subject: arg_id
       op: GreaterThan
       value: "13"
   backend:
     serviceName: %s
     servicePort: %d
`, backendSvc, backendPorts[0])

		assert.Nil(ginkgo.GinkgoT(), s.CreateResourceFromString(ar))

		err := s.EnsureNumApisixRoutesCreated(1)
		assert.Nil(ginkgo.GinkgoT(), err, "Checking number of routes")
		err = s.EnsureNumApisixUpstreamsCreated(1)
		assert.Nil(ginkgo.GinkgoT(), err, "Checking number of upstreams")

		_ = s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			WithQuery("id", 100).
			Expect().
			Status(http.StatusOK)

		msg := s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			WithQuery("id", 3).
			Expect().
			Status(http.StatusNotFound).
			Body().
			Raw()
		assert.Contains(ginkgo.GinkgoT(), msg, "404 Route Not Found")

		msg = s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			Expect().
			Status(http.StatusNotFound).
			Body().
			Raw()
		assert.Contains(ginkgo.GinkgoT(), msg, "404 Route Not Found")
	})

	ginkgo.It("operator is less_than", func() {
		backendSvc, backendPorts := s.DefaultHTTPBackend()

		ar := fmt.Sprintf(`
apiVersion: apisix.apache.org/v2alpha1
kind: ApisixRoute
metadata:
 name: httpbin-route
spec:
 http:
 - name: rule1
   match:
     hosts:
     - httpbin.org
     paths:
       - /ip
     nginxVars:
     - subject: arg_id
       op: LessThan
       value: "13"
   backend:
     serviceName: %s
     servicePort: %d
`, backendSvc, backendPorts[0])

		assert.Nil(ginkgo.GinkgoT(), s.CreateResourceFromString(ar))

		err := s.EnsureNumApisixRoutesCreated(1)
		assert.Nil(ginkgo.GinkgoT(), err, "Checking number of routes")
		err = s.EnsureNumApisixUpstreamsCreated(1)
		assert.Nil(ginkgo.GinkgoT(), err, "Checking number of upstreams")

		_ = s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			WithQuery("id", 12).
			Expect().
			Status(http.StatusOK)

		msg := s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			WithQuery("id", 13).
			Expect().
			Status(http.StatusNotFound).
			Body().
			Raw()
		assert.Contains(ginkgo.GinkgoT(), msg, "404 Route Not Found")

		msg = s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			Expect().
			Status(http.StatusNotFound).
			Body().
			Raw()
		assert.Contains(ginkgo.GinkgoT(), msg, "404 Route Not Found")
	})

	ginkgo.It("operator is in", func() {
		backendSvc, backendPorts := s.DefaultHTTPBackend()

		ar := fmt.Sprintf(`
apiVersion: apisix.apache.org/v2alpha1
kind: ApisixRoute
metadata:
 name: httpbin-route
spec:
 http:
 - name: rule1
   match:
     hosts:
     - httpbin.org
     paths:
       - /ip
     nginxVars:
     - subject: http_content_type
       op: In
       set: ["text/plain", "text/html", "image/jpeg"]
   backend:
     serviceName: %s
     servicePort: %d
`, backendSvc, backendPorts[0])

		assert.Nil(ginkgo.GinkgoT(), s.CreateResourceFromString(ar))

		err := s.EnsureNumApisixRoutesCreated(1)
		assert.Nil(ginkgo.GinkgoT(), err, "Checking number of routes")
		err = s.EnsureNumApisixUpstreamsCreated(1)
		assert.Nil(ginkgo.GinkgoT(), err, "Checking number of upstreams")

		_ = s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			WithHeader("Content-Type", "text/html").
			Expect().
			Status(http.StatusOK)

		msg := s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			WithHeader("Content-Type", "image/png").
			Expect().
			Status(http.StatusNotFound).
			Body().
			Raw()
		assert.Contains(ginkgo.GinkgoT(), msg, "404 Route Not Found")

		msg = s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			Expect().
			Status(http.StatusNotFound).
			Body().
			Raw()
		assert.Contains(ginkgo.GinkgoT(), msg, "404 Route Not Found")
	})

	ginkgo.It("operator is not_in", func() {
		backendSvc, backendPorts := s.DefaultHTTPBackend()

		ar := fmt.Sprintf(`
apiVersion: apisix.apache.org/v2alpha1
kind: ApisixRoute
metadata:
 name: httpbin-route
spec:
 http:
 - name: rule1
   match:
     hosts:
     - httpbin.org
     paths:
       - /ip
     nginxVars:
     - subject: http_content_type
       op: NotIn
       set: ["text/plain", "text/html", "image/jpeg"]
   backend:
     serviceName: %s
     servicePort: %d
`, backendSvc, backendPorts[0])

		assert.Nil(ginkgo.GinkgoT(), s.CreateResourceFromString(ar))

		err := s.EnsureNumApisixRoutesCreated(1)
		assert.Nil(ginkgo.GinkgoT(), err, "Checking number of routes")
		err = s.EnsureNumApisixUpstreamsCreated(1)
		assert.Nil(ginkgo.GinkgoT(), err, "Checking number of upstreams")

		_ = s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			WithHeader("Content-Type", "text/png").
			Expect().
			Status(http.StatusOK)

		msg := s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			WithHeader("Content-Type", "image/jpeg").
			Expect().
			Status(http.StatusNotFound).
			Body().
			Raw()
		assert.Contains(ginkgo.GinkgoT(), msg, "404 Route Not Found")

		_ = s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			Expect().
			Status(http.StatusOK).
			Body().
			Raw()
	})

	ginkgo.It("operator is regex match", func() {
		backendSvc, backendPorts := s.DefaultHTTPBackend()

		ar := fmt.Sprintf(`
apiVersion: apisix.apache.org/v2alpha1
kind: ApisixRoute
metadata:
 name: httpbin-route
spec:
 http:
 - name: rule1
   match:
     hosts:
     - httpbin.org
     paths:
       - /ip
     nginxVars:
     - subject: http_x_real_uri
       op: RegexMatch
       value: "^/ip/0\\d{2}/.*$"
   backend:
     serviceName: %s
     servicePort: %d
`, backendSvc, backendPorts[0])

		assert.Nil(ginkgo.GinkgoT(), s.CreateResourceFromString(ar))

		err := s.EnsureNumApisixRoutesCreated(1)
		assert.Nil(ginkgo.GinkgoT(), err, "Checking number of routes")
		err = s.EnsureNumApisixUpstreamsCreated(1)
		assert.Nil(ginkgo.GinkgoT(), err, "Checking number of upstreams")

		_ = s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			WithHeader("X-Real-Uri", "/ip/098/v4").
			Expect().
			Status(http.StatusOK)

		msg := s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			WithHeader("X-Real-Uri", "/ip/0983/v4").
			Expect().
			Status(http.StatusNotFound).
			Body().
			Raw()
		assert.Contains(ginkgo.GinkgoT(), msg, "404 Route Not Found")

		msg = s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			Expect().
			Status(http.StatusNotFound).
			Body().
			Raw()
		assert.Contains(ginkgo.GinkgoT(), msg, "404 Route Not Found")
	})

	ginkgo.It("operator is regex not match", func() {
		backendSvc, backendPorts := s.DefaultHTTPBackend()

		ar := fmt.Sprintf(`
apiVersion: apisix.apache.org/v2alpha1
kind: ApisixRoute
metadata:
 name: httpbin-route
spec:
 http:
 - name: rule1
   match:
     hosts:
     - httpbin.org
     paths:
       - /ip
     nginxVars:
     - subject: http_x_real_uri
       op: RegexNotMatch
       value: "^/ip/0\\d{2}/.*$"
   backend:
     serviceName: %s
     servicePort: %d
`, backendSvc, backendPorts[0])

		assert.Nil(ginkgo.GinkgoT(), s.CreateResourceFromString(ar))

		err := s.EnsureNumApisixRoutesCreated(1)
		assert.Nil(ginkgo.GinkgoT(), err, "Checking number of routes")
		err = s.EnsureNumApisixUpstreamsCreated(1)
		assert.Nil(ginkgo.GinkgoT(), err, "Checking number of upstreams")

		_ = s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			WithHeader("X-Real-Uri", "/ip/0983/v4").
			Expect().
			Status(http.StatusOK)

		msg := s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			WithHeader("X-Real-Uri", "/ip/098/v4").
			Expect().
			Status(http.StatusNotFound).
			Body().
			Raw()
		assert.Contains(ginkgo.GinkgoT(), msg, "404 Route Not Found")

		_ = s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			Expect().
			Status(http.StatusOK).
			Body().
			Raw()
	})

	ginkgo.It("operator is regex match in case insensitive mode", func() {
		backendSvc, backendPorts := s.DefaultHTTPBackend()

		ar := fmt.Sprintf(`
apiVersion: apisix.apache.org/v2alpha1
kind: ApisixRoute
metadata:
 name: httpbin-route
spec:
 http:
 - name: rule1
   match:
     hosts:
     - httpbin.org
     paths:
       - /ip
     nginxVars:
     - subject: http_x_real_uri
       op: RegexMatchCaseInsensitive
       value: "^/ip/0\\d{2}/.*$"
   backend:
     serviceName: %s
     servicePort: %d
`, backendSvc, backendPorts[0])

		assert.Nil(ginkgo.GinkgoT(), s.CreateResourceFromString(ar))

		err := s.EnsureNumApisixRoutesCreated(1)
		assert.Nil(ginkgo.GinkgoT(), err, "Checking number of routes")
		err = s.EnsureNumApisixUpstreamsCreated(1)
		assert.Nil(ginkgo.GinkgoT(), err, "Checking number of upstreams")

		_ = s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			WithHeader("X-Real-Uri", "/IP/098/v4").
			Expect().
			Status(http.StatusOK)

		msg := s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			WithHeader("X-Real-Uri", "/ip/0983/v4").
			Expect().
			Status(http.StatusNotFound).
			Body().
			Raw()
		assert.Contains(ginkgo.GinkgoT(), msg, "404 Route Not Found")

		msg = s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			Expect().
			Status(http.StatusNotFound).
			Body().
			Raw()
		assert.Contains(ginkgo.GinkgoT(), msg, "404 Route Not Found")
	})

	ginkgo.It("operator is regex not match in case insensitive mode", func() {
		backendSvc, backendPorts := s.DefaultHTTPBackend()

		ar := fmt.Sprintf(`
apiVersion: apisix.apache.org/v2alpha1
kind: ApisixRoute
metadata:
 name: httpbin-route
spec:
 http:
 - name: rule1
   match:
     hosts:
     - httpbin.org
     paths:
       - /ip
     nginxVars:
     - subject: http_x_real_uri
       op: RegexNotMatchCaseInsensitive
       value: "^/ip/0\\d{2}/.*$"
   backend:
     serviceName: %s
     servicePort: %d
`, backendSvc, backendPorts[0])

		assert.Nil(ginkgo.GinkgoT(), s.CreateResourceFromString(ar))

		err := s.EnsureNumApisixRoutesCreated(1)
		assert.Nil(ginkgo.GinkgoT(), err, "Checking number of routes")
		err = s.EnsureNumApisixUpstreamsCreated(1)
		assert.Nil(ginkgo.GinkgoT(), err, "Checking number of upstreams")

		_ = s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			WithHeader("X-Real-Uri", "/IP/0983/v4").
			Expect().
			Status(http.StatusOK)

		msg := s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			WithHeader("X-Real-Uri", "/IP/098/v4").
			Expect().
			Status(http.StatusNotFound).
			Body().
			Raw()
		assert.Contains(ginkgo.GinkgoT(), msg, "404 Route Not Found")

		_ = s.NewAPISIXClient().GET("/ip").
			WithHeader("Host", "httpbin.org").
			Expect().
			Status(http.StatusOK).
			Body().
			Raw()
	})
})