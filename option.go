package xff

import (
	"github.com/czyt/xff/internal/mask"
	"net"
)

type Option func(*xffOption)
type xffOption struct {
	// Set to true if all IPs or Subnets are allowed.
	allowAll bool
	// List of IP subnets that are allowed.
	allowedMasks []net.IPNet
}

// AllowAll Set to true if all IPs or Subnets are allowed.
func AllowAll() Option {
	return func(option *xffOption) {
		option.allowAll = true
	}
}

// AllowedSubnets is a list of Subnets from which we will accept the
// X-Forwarded-For header.
// If this list is empty we will accept every Subnets (default).
func AllowedSubnets(subnets []string) Option {
	return func(option *xffOption) {
		masks, _ := mask.ParseFrom(subnets)
		option.allowedMasks = masks
		option.allowAll = false
	}
}
