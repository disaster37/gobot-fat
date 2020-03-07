module github.com/disaster37/gobot-fat

go 1.13

require (
	github.com/elastic/go-elasticsearch/v7 v7.6.0
	github.com/labstack/echo/v4 v4.1.11
	github.com/mitchellh/mapstructure v1.1.2
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.6.1
	github.com/stretchr/testify v1.4.0
	github.com/x-cray/logrus-prefixed-formatter v0.5.2
	gobot.io/x/gobot v1.14.0

)

replace gobot.io/x/gobot => github.com/disaster37/gobot v1.14.1-0.20200214195251-3e9e12582a3d
