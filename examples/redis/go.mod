module redis_example

go 1.12

require (
	github.com/gin-gonic/gin v1.4.0
	github.com/go-redis/redis v6.15.2+incompatible
	github.com/shoplineapp/captin v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.3.0
)

replace github.com/shoplineapp/captin => ../../
