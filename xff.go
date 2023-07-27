package xff

import (
	"context"
	"github.com/czyt/xff/internal/mask"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"net"
)

const (
	xForwardForHeader = "X-Forwarded-For"
)

func Server(opts ...Option) middleware.Middleware {
	xffOpt := &xffOption{allowAll: true}
	for _, opt := range opts {
		opt(xffOpt)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				if hr, ok := tr.(*http.Transport); ok {
					if xffh := hr.RequestHeader().Get(xForwardForHeader); xffh != "" {
						if needUpdateRemoteAddr, address := checkNeedUpdateRemoteAddr(xffOpt, hr.Request().RemoteAddr, xffh); needUpdateRemoteAddr {
							hr.Request().RemoteAddr = address
						}
					}
				}
			}
			return handler(ctx, req)
		}
	}
}

func checkNeedUpdateRemoteAddr(opt *xffOption, currentRemoteAddr, xffAddress string) (needUpdate bool, remoteAddress string) {
	if sip, sport, err := net.SplitHostPort(currentRemoteAddr); err == nil && sip != "" {
		ip := net.ParseIP(sip)
		if (ip != nil && mask.CheckIpInMasks(ip, opt.allowedMasks)) || opt.allowAll {
			if xffIp := mask.GetPublicIpFrom(xffAddress); xffIp != "" {
				return true, net.JoinHostPort(xffIp, sport)
			}
		}
	}
	return false, currentRemoteAddr
}
