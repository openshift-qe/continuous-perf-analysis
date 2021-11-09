// module github.com/kedark3/cpa

// go 1.16

module github.com/kedark3/cpa

go 1.16

require (
	github.com/alexflint/go-arg v1.4.2
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.16.0
	github.com/openshift/openshift-tests v0.0.0-20210916082130-4fca21c38ee6
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/common v0.6.0
	github.com/slack-go/slack v0.9.5 // indirect
	k8s.io/api v0.17.1
	k8s.io/apimachinery v0.17.1
	k8s.io/kubernetes v1.21.0
	sigs.k8s.io/yaml v1.2.0
)

replace (
	bitbucket.org/ww/goautoneg => github.com/munnerz/goautoneg v0.0.0-20120707110453-a547fc61f48d
	github.com/golang/glog => github.com/openshift/golang-glog v0.0.0-20190322123450-3c92600d7533
	github.com/google/cadvisor => github.com/openshift/google-cadvisor v0.33.2-0.20190902063809-4db825a8ad0d
	github.com/jteeuwen/go-bindata => github.com/jteeuwen/go-bindata v3.0.8-0.20151023091102-a0ff2567cfb7+incompatible
	github.com/onsi/ginkgo => github.com/openshift/onsi-ginkgo v1.4.1-0.20190902091932-d0603c19fe78
	github.com/opencontainers/runc => github.com/openshift/opencontainers-runc v1.0.0-rc4.0.20190926164333-b942ff4cc6f8
	github.com/openshift/api => github.com/openshift/api v0.0.0-20200117162508-e7ccdda6ba67
	k8s.io/api => k8s.io/api v0.17.1
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.1
	k8s.io/apimachinery => github.com/openshift/kubernetes-apimachinery v0.0.0-20191121175448-79c2a76c473a
	k8s.io/apiserver => github.com/openshift/kubernetes-apiserver v0.0.0-20200109101329-ed563d1b80a1
	k8s.io/cli-runtime => github.com/openshift/kubernetes-cli-runtime v0.0.0-20200115000600-01f2488fd0b7
	k8s.io/client-go => github.com/openshift/kubernetes-client-go v0.0.0-20200106170045-1fda2942f64d
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.17.1
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.17.1
	k8s.io/code-generator => k8s.io/code-generator v0.17.1
	k8s.io/component-base => k8s.io/component-base v0.17.1
	k8s.io/cri-api => k8s.io/cri-api v0.17.1
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.17.1
	k8s.io/kube-aggregator => github.com/openshift/kubernetes-kube-aggregator v0.0.0-20191209133208-1e3c0eec4d61
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.17.1
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.17.1
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.17.1
	k8s.io/kubectl => k8s.io/kubectl v0.17.1
	k8s.io/kubelet => k8s.io/kubelet v0.17.1

	k8s.io/kubernetes => github.com/openshift/kubernetes v1.17.0-alpha.0.0.20200120180958-5945c3b07163
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.0.0-20191121182806-cdbd52110e91
	k8s.io/metrics => k8s.io/metrics v0.0.0-20191121181631-c7d4ee0ffc0e
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.0.0-20191121181040-36c9528858d2
	sigs.k8s.io/yaml => sigs.k8s.io/yaml v1.2.0
)
