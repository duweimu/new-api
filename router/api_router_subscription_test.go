package router

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionBalancePayRoute_IsRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()
	SetApiRouter(engine)

	routes := engine.Routes()
	for _, route := range routes {
		if route.Method == "POST" && route.Path == "/api/subscription/balance/pay" {
			return
		}
	}

	require.Fail(t, "subscription balance pay route is not registered")
}
