package balancer

import (
	"context"

	"github.com/ydb-platform/ydb-go-sdk/v3/internal/conn"
)

// Element is an empty interface that holds some Balancer specific data.
type Element interface{}

// Balancer is an interface that implements particular load-balancing
// algorithm.
//
// Balancer methods called synchronized. That is, implementations must not
// provide additional goroutine safety.
type Balancer interface {
	// Next returns next connection for request.
	// return Err
	Next(ctx context.Context, allowBanned bool) conn.Conn

	// Create makes empty balancer with same implementation
	Create(conns []conn.Conn) Balancer

	// NeedRefresh sync call, which return in one of cases
	// first - if balancer can known about need refresh or not (may be right after call, with pause or never)
	// second - ctx cancelled, must be cancelled be caller for prevent goroutines leak
	// return true if the balancer need refresh
	NeedRefresh(ctx context.Context) bool
}

func IsOkConnection(c conn.Conn, bannedIsOk bool) bool {
	switch c.GetState() {
	case conn.Online, conn.Created, conn.Offline:
		return true
	case conn.Banned:
		return bannedIsOk
	default:
		return false
	}
}
