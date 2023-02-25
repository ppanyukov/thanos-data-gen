module github.com/ppanyukov/thanos-data-gen

go 1.13

// Not sure about all these versions!

// require github.com/prometheus/prometheus/tsdb v1.8.2

require (
	github.com/go-kit/kit v0.9.0
	github.com/oklog/run v1.0.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/common v0.7.0
	github.com/prometheus/prometheus v1.8.2-0.20190913102521-8ab628b35467
	github.com/stretchr/testify v1.4.0 // indirect
	go.uber.org/automaxprocs v1.2.0
	golang.org/x/sys v0.1.0 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
