module github.com/ppanyukov/thanos-data-gen

go 1.13

// Not sure about all these versions!

// require github.com/prometheus/prometheus/tsdb v1.8.2

require (
	github.com/go-kit/kit v0.9.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/prometheus v1.8.2-0.20190913102521-8ab628b35467
)

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
