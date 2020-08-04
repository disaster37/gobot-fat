module github.com/disaster37/gobot-fat

go 1.13

require (
	github.com/DATA-DOG/go-sqlmock v1.4.1
	github.com/creack/goselect v0.1.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/disaster37/go-arest v0.0.2-0.20200705070542-e6bcd85bae54

	github.com/elastic/go-elasticsearch/v7 v7.6.0
	github.com/gofrs/uuid v3.3.0+incompatible // indirect
	github.com/jinzhu/gorm v1.9.11
	github.com/labstack/echo/v4 v4.1.11
	github.com/mitchellh/mapstructure v1.1.2
	github.com/pkg/errors v0.9.1
	github.com/raff/goble v0.0.0-20200327175727-d63360dcfd80 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.4.0
	github.com/x-cray/logrus-prefixed-formatter v0.5.2
	gobot.io/x/gobot v1.14.0
	golang.org/x/net v0.0.0-20200625001655-4c5254603344 // indirect
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	google.golang.org/appengine v1.6.1

)

replace gobot.io/x/gobot => github.com/disaster37/gobot v1.14.1-0.20200214195251-3e9e12582a3d

//replace gobot.io/x/gobot => github.com/disaster37/gobot v1.14.1-0.20200702175719-117333e3b1bf
