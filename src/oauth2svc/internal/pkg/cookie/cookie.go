package cookie

import (
	"fmt"
	"strconv"

	"github.com/smallbiznis/oauth2-server/internal/pkg/strings"
)

func SetCookie(tenant, value string) (name, v string, maxAge int, path, domain string, secure, httpOnly bool) {
	name = strings.MustEnv("SESSION_NAME", "_SID")
	domain = fmt.Sprintf("%s.%s", tenant, strings.MustEnv("SESSION_DOMAIN", ".smallbiznis.test"))
	path = strings.MustEnv("SESSION_DOMAIN", "/")
	v = value
	exp, _ := strconv.Atoi(strings.MustEnv("SESSION_MAXAGE", "86400"))
	maxAge = exp
	secure = strings.MustEnv("ENV", "production") == "production"
	httpOnly = true
	return
}
