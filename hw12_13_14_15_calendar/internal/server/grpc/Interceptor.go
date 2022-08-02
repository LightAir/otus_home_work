package internalgrpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func (srv *GRPCServer) RequestLogInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	i, err := handler(ctx, req)

	duration := time.Since(start)

	srv.logger.LogGRPCRequest(ctx, info, duration, status.Code(err).String())

	return i, err
}
