package componentcheck

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	"github.com/bweston92/healthz/healthz"
)

// NewGRPCHealthCheck takes a gRPC connection and checks the
// status of the connection everytime healthz calls it.
func NewGRPCHealthCheck(con *grpc.ClientConn) healthz.ComponentHealthCheck {
	return func() *healthz.Error {
		switch s := con.GetState(); s {
		case connectivity.Ready, connectivity.Idle:
			return nil
		default:
			return healthz.NewError("Connection is not READY.", healthz.Meta{
				"state": s.String(),
			})
		}
	}
}
