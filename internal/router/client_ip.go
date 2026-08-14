package router

import (
	"fmt"
	"net"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

func buildClientIPFunc(value string) (app.ClientIP, error) {
	values := strings.Split(value, ",")
	trustedCIDRs := make([]*net.IPNet, 0, len(values))

	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}

		_, cidr, err := net.ParseCIDR(value)
		if err != nil {
			return nil, fmt.Errorf("parse trusted proxy CIDR %q: %w", value, err)
		}
		trustedCIDRs = append(trustedCIDRs, cidr)
	}

	return app.ClientIPWithOption(app.ClientIPOptions{
		RemoteIPHeaders: []string{"X-Forwarded-For", "X-Real-IP"},
		TrustedCIDRs:    trustedCIDRs,
	}), nil
}
