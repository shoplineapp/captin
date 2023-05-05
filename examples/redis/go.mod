module redis_example

go 1.15

require (
	github.com/bsm/redis-lock v8.0.0+incompatible
	github.com/gin-gonic/gin v1.9.0
	github.com/go-redis/redis v6.15.2+incompatible
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/shoplineapp/captin v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.8.1
)

replace github.com/shoplineapp/captin => ../../
