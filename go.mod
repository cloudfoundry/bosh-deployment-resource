module github.com/cloudfoundry/bosh-deployment-resource

go 1.17

require (
	github.com/cloudfoundry/bosh-cli/v7 v7.0.1
	github.com/cloudfoundry/bosh-utils v0.0.331
	github.com/cloudfoundry/socks5-proxy v0.2.71
	github.com/cppforlife/go-patch v0.2.0
	github.com/cppforlife/go-semi-semantic v0.0.0-20160921010311-576b6af77ae4
	github.com/jessevdk/go-flags v1.5.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.20.2
	golang.org/x/oauth2 v0.0.0-20220822191816-0ebed06d0094
	google.golang.org/api v0.95.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	cloud.google.com/go v0.104.0 // indirect
	cloud.google.com/go/compute v1.9.0 // indirect
	cloud.google.com/go/iam v0.4.0 // indirect
	cloud.google.com/go/storage v1.26.0 // indirect
	code.cloudfoundry.org/clock v1.0.0 // indirect
	code.cloudfoundry.org/tlsconfig v0.0.0-20220621140725-0e6fbd869921 // indirect
	code.cloudfoundry.org/workpool v0.0.0-20200131000409-2ac56b354115 // indirect
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/aws/aws-sdk-go v1.44.94 // indirect
	github.com/bmatcuk/doublestar v1.3.4 // indirect
	github.com/charlievieth/fs v0.0.3 // indirect
	github.com/cheggaaa/pb/v3 v3.1.0 // indirect
	github.com/cloudfoundry/bosh-agent v2.367.0+incompatible // indirect
	github.com/cloudfoundry/bosh-davcli v0.0.78 // indirect
	github.com/cloudfoundry/bosh-gcscli v0.0.54 // indirect
	github.com/cloudfoundry/bosh-s3cli v0.0.152 // indirect
	github.com/cloudfoundry/config-server v0.1.71 // indirect
	github.com/cloudfoundry/go-socks5 v0.0.0-20180221174514-54f73bdb8a8e // indirect
	github.com/creack/pty v1.1.9 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.1.0 // indirect
	github.com/googleapis/gax-go/v2 v2.5.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/pivotal-cf/paraphernalia v0.0.0-20180203224945-a64ae2051c20 // indirect
	github.com/rivo/uniseg v0.3.4 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/vito/go-interact v1.0.1 // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/crypto v0.0.0-20220829220503-c86fa9a7ed90 // indirect
	golang.org/x/net v0.0.0-20220907135653-1e95f45603a7 // indirect
	golang.org/x/sys v0.0.0-20220908164124-27713097b956 // indirect
	golang.org/x/term v0.0.0-20220722155259-a9ba230a4035 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220908141613-51c1cc9bc6d0 // indirect
	google.golang.org/grpc v1.49.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/jessevdk/go-flags v1.0.0 => github.com/cppforlife/go-flags v0.0.0-20170707010757-351f5f310b26

replace gopkg.in/fsnotify.v1 v1.4.7 => gopkg.in/fsnotify/fsnotify.v1 v1.4.7
