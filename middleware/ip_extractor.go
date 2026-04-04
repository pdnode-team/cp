package middleware

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

// ConfigureIPExtractor 设置 RealIP / 限流所用的客户端 IP 提取逻辑。
//
// 若环境变量 TRUSTED_PROXY_CIDRS 为空：使用直连 IP（不读取 X-Forwarded-For），适合本机直连或不确定反代配置时，避免客户端伪造 XFF。
//
// 若已设置（逗号分隔的 CIDR，例如 "10.0.0.0/8,172.16.0.0/12,192.168.0.0/16"）：仅当 RemoteAddr 落在这些网段时才按 X-Forwarded-For 解析客户端 IP；并关闭对任意私网/回环的默认信任，避免误信。
func ConfigureIPExtractor(e *echo.Echo) {
	raw := strings.TrimSpace(strings.ReplaceAll(os.Getenv("TRUSTED_PROXY_CIDRS"), " ", ""))
	if raw == "" {
		e.IPExtractor = echo.ExtractIPDirect()
		return
	}

	var opts []echo.TrustOption
	opts = append(opts,
		echo.TrustLoopback(false),
		echo.TrustLinkLocal(false),
		echo.TrustPrivateNet(false),
	)

	for _, part := range strings.Split(raw, ",") {
		if part == "" {
			continue
		}
		_, ipNet, err := net.ParseCIDR(part)
		if err != nil {
			log.Fatalf("TRUSTED_PROXY_CIDRS 含无效 CIDR %q: %v", part, err)
		}
		opts = append(opts, echo.TrustIPRange(ipNet))
	}

	if len(opts) == 3 {
		log.Fatal("TRUSTED_PROXY_CIDRS 已设置但未解析出任何 CIDR")
	}

	e.IPExtractor = echo.ExtractIPFromXFFHeader(opts...)
}
