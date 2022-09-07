module github.com/cloudfoundry/bosh-deployment-resource

go 1.17

require (
	github.com/cloudfoundry/bosh-cli v6.2.1+incompatible
	github.com/cloudfoundry/bosh-utils v0.0.0-20200222100218-2c99d7618fe7
	github.com/cloudfoundry/socks5-proxy v0.2.0
	github.com/cppforlife/go-patch v0.2.0
	github.com/cppforlife/go-semi-semantic v0.0.0-20160921010311-576b6af77ae4
	github.com/jessevdk/go-flags v1.0.0
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.9.0
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/api v0.18.0
	gopkg.in/yaml.v2 v2.2.8
)

require (
	cloud.google.com/go v0.53.0 // indirect
	cloud.google.com/go/storage v1.5.0 // indirect
	code.cloudfoundry.org/clock v1.0.0 // indirect
	code.cloudfoundry.org/tlsconfig v0.0.0-20200131000646-bbe0f8da39b3 // indirect
	code.cloudfoundry.org/workpool v0.0.0-20200131000409-2ac56b354115 // indirect
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/aws/aws-sdk-go v1.29.9 // indirect
	github.com/bmatcuk/doublestar v1.2.2 // indirect
	github.com/charlievieth/fs v0.0.0-20170613215519-7dc373669fa1 // indirect
	github.com/cheggaaa/pb v1.0.29 // indirect
	github.com/cloudfoundry/bosh-agent v2.304.0+incompatible // indirect
	github.com/cloudfoundry/bosh-davcli v0.0.44 // indirect
	github.com/cloudfoundry/bosh-gcscli v0.0.18 // indirect
	github.com/cloudfoundry/bosh-s3cli v0.0.95 // indirect
	github.com/cloudfoundry/config-server v0.0.124 // indirect
	github.com/cloudfoundry/go-socks5 v0.0.0-20180221174514-54f73bdb8a8e // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.3.3 // indirect
	github.com/googleapis/gax-go/v2 v2.0.5 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.0.0 // indirect
	github.com/hpcloud/tail v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af // indirect
	github.com/jpillora/backoff v0.0.0-20170918002102-8eab2debe79d // indirect
	github.com/jstemmer/go-junit-report v0.9.1 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mattn/go-runewidth v0.0.8 // indirect
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	github.com/pivotal-cf/paraphernalia v0.0.0-20180203224945-a64ae2051c20 // indirect
	github.com/vito/go-interact v1.0.0 // indirect
	go.opencensus.io v0.22.3 // indirect
	golang.org/x/crypto v0.0.0-20200221231518-2aa609cf4a9d // indirect
	golang.org/x/exp v0.0.0-20200224162631-6cc2880d07d6 // indirect
	golang.org/x/lint v0.0.0-20200130185559-910be7a94367 // indirect
	golang.org/x/mod v0.2.0 // indirect
	golang.org/x/net v0.0.0-20200222125558-5a598a2470a0 // indirect
	golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae // indirect
	golang.org/x/text v0.3.2 // indirect
	golang.org/x/tools v0.0.0-20200225190036-fefc8d187781 // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	google.golang.org/genproto v0.0.0-20200225123651-fc8f55426688 // indirect
	google.golang.org/grpc v1.27.1 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	honnef.co/go/tools v0.0.1-2020.1.3 // indirect
)

replace github.com/jessevdk/go-flags v1.0.0 => github.com/cppforlife/go-flags v0.0.0-20170707010757-351f5f310b26

replace gopkg.in/fsnotify.v1 v1.4.7 => gopkg.in/fsnotify/fsnotify.v1 v1.4.7
