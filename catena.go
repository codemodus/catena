// Package catena aides gRPC interceptor catenation.
package catena

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// UnaryServerCatena ...
type UnaryServerCatena struct {
	is []grpc.UnaryServerInterceptor
}

// NewUnaryServerCatena ...
func NewUnaryServerCatena(is ...grpc.UnaryServerInterceptor) *UnaryServerCatena {
	return &UnaryServerCatena{is: is}
}

func appendInterceptors(is []grpc.UnaryServerInterceptor, ais ...grpc.UnaryServerInterceptor) []grpc.UnaryServerInterceptor {
	lcur := len(is)
	ltot := lcur + len(ais)
	if ltot > cap(is) {
		nis := make([]grpc.UnaryServerInterceptor, ltot)
		copy(nis, is)
		is = nis
	}

	copy(is[lcur:], ais)

	return is
}

// Append ...
func (c *UnaryServerCatena) Append(is ...grpc.UnaryServerInterceptor) *UnaryServerCatena {
	c = NewUnaryServerCatena(appendInterceptors(c.is, is...)...)

	return c
}

// Merge ...
func (c *UnaryServerCatena) Merge(cs ...*UnaryServerCatena) *UnaryServerCatena {
	for k := range cs {
		c = NewUnaryServerCatena(appendInterceptors(c.is, cs[k].is...)...)
	}

	return c
}

// Copy ...
func (c *UnaryServerCatena) Copy(catena *UnaryServerCatena) {
	c.is = make([]grpc.UnaryServerInterceptor, len(catena.is))

	for k := range catena.is {
		c.is[k] = catena.is[k]
	}
}

// Interceptor ...
func (c *UnaryServerCatena) Interceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		concat := func(wrap grpc.UnaryServerInterceptor, h grpc.UnaryHandler) grpc.UnaryHandler {
			return func(outerCtx context.Context, outerReq interface{}) (interface{}, error) {
				return wrap(outerCtx, outerReq, info, h)
			}
		}

		for i := len(c.is) - 1; i >= 0; i-- {
			handler = concat(c.is[i], handler)
		}

		return handler(ctx, req)
	}
}

// ServerOption ...
func (c *UnaryServerCatena) ServerOption() grpc.ServerOption {
	return grpc.UnaryInterceptor(c.Interceptor())
}
