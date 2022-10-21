module github.com/btcsuite/btcd

require (
	github.com/btcsuite/btcd/btcec/v2 v2.1.3
	github.com/btcsuite/btcd/btcutil v1.1.0
	github.com/btcsuite/btcd/chaincfg/chainhash v1.0.1
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f
	github.com/btcsuite/go-socks v0.0.0-20170105172521-4720035b7bfd
	github.com/btcsuite/websocket v0.0.0-20150119174127-31079b680792
	github.com/btcsuite/winsvc v1.0.0
	github.com/davecgh/go-spew v1.1.1
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1
	github.com/decred/dcrd/lru v1.0.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/jrick/logrotate v1.0.0
	github.com/netsec-ethz/scion-apps v0.5.0
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007
	inet.af/netaddr v0.0.0-20210903134321-85fa6c94624e
)

require (
	github.com/aead/siphash v1.0.1 // indirect
	github.com/antlr/antlr4 v0.0.0-20181218183524-be58ebffde8e // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/britram/borat v0.0.0-20181011130314-f891bcfcfb9b // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/cheekybits/genny v1.0.0 // indirect
	github.com/decred/dcrd/crypto/blake256 v1.0.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-kit/kit v0.10.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/gopacket v1.1.16-0.20190123011826-102d5ca2098c // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645 // indirect
	github.com/inconshreveable/log15 v0.0.0-20180818164646-67afb5ed74ec // indirect
	github.com/kkdai/bstream v0.0.0-20161212061736-f391b8402d23 // indirect
	github.com/lucas-clemente/quic-go v0.23.0 // indirect
	github.com/marten-seemann/qtls-go1-16 v0.1.4 // indirect
	github.com/marten-seemann/qtls-go1-17 v0.1.0 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/netsec-ethz/rains v0.3.0 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pelletier/go-toml v1.9.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.7.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.10.0 // indirect
	github.com/prometheus/procfs v0.1.3 // indirect
	github.com/scionproto/scion v0.6.1-0.20210929154253-764d6e2afe47 // indirect
	github.com/uber/jaeger-client-go v2.29.1+incompatible // indirect
	github.com/uber/jaeger-lib v2.0.0+incompatible // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.17.0 // indirect
	go4.org/intern v0.0.0-20210108033219-3eb7198706b2 // indirect
	go4.org/unsafe/assume-no-moving-gc v0.0.0-20201222180813-1025295fd063 // indirect
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d // indirect
	golang.org/x/text v0.3.6 // indirect
	golang.org/x/tools v0.1.2 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c // indirect
	google.golang.org/grpc v1.38.1 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/btcsuite/btcd/btcutil => ./btcutil

// The retract statements below fixes an accidental push of the tags of a btcd
// fork.
retract (
	v0.18.1
	v0.18.0
	v0.17.1
	v0.17.0
	v0.16.5
	v0.16.4
	v0.16.3
	v0.16.2
	v0.16.1
	v0.16.0

	v0.15.2
	v0.15.1
	v0.15.0

	v0.14.7
	v0.14.6
	v0.14.6
	v0.14.5
	v0.14.4
	v0.14.3
	v0.14.2
	v0.14.1

	v0.14.0
	v0.13.0-beta2
	v0.13.0-beta
)

go 1.17
