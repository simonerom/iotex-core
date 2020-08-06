module github.com/iotexproject/iotex-core

go 1.13

require (
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/ethereum/go-ethereum v1.8.27
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9
	github.com/golang/mock v1.4.0
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/iotexproject/go-fsm v1.0.0
	github.com/iotexproject/go-p2p v0.2.13-0.20200805220926-3607734c3277
	github.com/iotexproject/go-pkgs v0.1.2-0.20200523040337-5f1d9ddaa8ee
	github.com/iotexproject/iotex-address v0.2.2
	github.com/iotexproject/iotex-antenna-go/v2 v2.3.2
	github.com/iotexproject/iotex-election v0.3.2
	github.com/iotexproject/iotex-proto v0.3.2-0.20200729044038-c22fbb206571
	github.com/libp2p/go-libp2p-peerstore v0.2.6
	github.com/mattn/go-sqlite3 v1.11.0
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1
	github.com/multiformats/go-multiaddr v0.2.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.7.1
	github.com/rs/zerolog v1.14.3
	github.com/schollz/progressbar/v2 v2.15.0
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.6.1
	go.etcd.io/bbolt v1.3.5
	go.uber.org/automaxprocs v1.2.0
	go.uber.org/config v1.3.1
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20200423211502-4bdfaf469ed5
	golang.org/x/net v0.0.0-20191204025024-5ee1b9f4859a
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a
	google.golang.org/genproto v0.0.0-20190530194941-fb225487d101
	google.golang.org/grpc v1.21.0
	gopkg.in/yaml.v2 v2.2.5
)

replace github.com/ethereum/go-ethereum => github.com/iotexproject/go-ethereum v0.3.1
