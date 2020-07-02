module github.com/disaster37/gobot-fat

go 1.13

require (
	github.com/DATA-DOG/go-sqlmock v1.4.1
	github.com/creack/goselect v0.1.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/elastic/go-elasticsearch/v7 v7.6.0
	github.com/gofrs/uuid v3.3.0+incompatible // indirect
	github.com/jinzhu/gorm v1.9.11
	github.com/labstack/echo/v4 v4.1.11
	github.com/mitchellh/mapstructure v1.1.2
	github.com/pkg/errors v0.8.1
	github.com/raff/goble v0.0.0-20200327175727-d63360dcfd80 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.4.0
	github.com/x-cray/logrus-prefixed-formatter v0.5.2
	gobot.io/x/gobot v1.14.0

)

//replace gobot.io/x/gobot => github.com/disaster37/gobot v1.14.1-0.20200214195251-3e9e12582a3d

replace gobot.io/x/gobot => github.com/disaster37/gobot v1.14.1-0.20200702173906-5aed348fad6f
