module blogX_server

go 1.24

replace github.com/siddontang/go-mysql v1.9.1 => github.com/go-mysql-org/go-mysql v1.9.1

require (
	github.com/BurntSushi/toml v1.3.2
	github.com/PuerkitoBio/goquery v1.10.3
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/elastic/go-elasticsearch/v9 v9.0.0
	github.com/gin-gonic/gin v1.10.0
	github.com/go-mysql-org/go-mysql v1.9.1
	github.com/go-playground/locales v0.14.1
	github.com/go-playground/universal-translator v0.18.1
	github.com/go-playground/validator/v10 v10.20.0
	github.com/jordan-wright/email v4.0.1-0.20210109023952-943e75fe5223+incompatible
	github.com/lionsoul2014/ip2region/binding/golang v0.0.0-20241220152942-06eb5c6e8230
	github.com/mojocn/base64Captcha v1.3.8
	github.com/olivere/elastic/v7 v7.0.32
	github.com/pingcap/check v0.0.0-20211026125417-57bd13f7b5f0
	github.com/pingcap/errors v0.11.5-0.20240311024730-e056997136bb
	github.com/pkg/errors v0.9.1
	github.com/qiniu/go-sdk/v7 v7.25.3
	github.com/redis/go-redis/v9 v9.8.0
	github.com/siddontang/go v0.0.0-20180604090527-bdc77568d726 // indirect
	github.com/siddontang/go-log v0.0.0-20180807004314-8d05993dda07 // indirect
	github.com/siddontang/go-mysql v0.0.0-20190524062908-de6c3a84bcbe
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/crypto v0.38.0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.26.1
	gorm.io/plugin/dbresolver v1.6.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/alex-ant/gomath v0.0.0-20160516115720-89013a210a82 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/bytedance/sonic v1.11.6 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/gomarkdown/markdown v0.0.0-20250311123330-531bef5e742b // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pingcap/failpoint v0.0.0-20240528011301-b51a646c7c86 // indirect
	github.com/pingcap/log v1.1.1-0.20230317032135-a0d097d16e22 // indirect
	github.com/pingcap/parser v0.0.0-20190506092653-e336082eb825 // indirect
	github.com/pingcap/tidb/pkg/parser v0.0.0-20241118164214-4f047be191be // indirect
	github.com/pingcap/tipb v0.0.0-20190428032612-535e1abaa330 // indirect
	github.com/robfig/cron/v3 v3.0.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/tklauser/go-sysconf v0.3.15 // indirect
	github.com/tklauser/numcpus v0.10.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/image v0.23.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/term v0.32.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/fileutil v1.0.0 // indirect

)
