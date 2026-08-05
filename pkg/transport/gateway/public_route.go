package gateway

import (
	"regexp"
	"strings"
)

// PublicRoute описывает gateway-маршрут, не требующий access token.
type PublicRoute struct {
	Method string
	Path   string
}

type routeMatcher struct {
	routes []compiledPublicRoute
}

type compiledPublicRoute struct {
	method  string
	pattern *regexp.Regexp
}

func newRouteMatcher(routes []PublicRoute) routeMatcher {
	compiled := make([]compiledPublicRoute, 0, len(routes))
	for _, route := range routes {
		compiled = append(compiled, compiledPublicRoute{
			method:  route.Method,
			pattern: regexp.MustCompile(routePattern(route.Path)),
		})
	}

	return routeMatcher{routes: compiled}
}

func (m routeMatcher) Match(method, path string) bool {
	path = strings.TrimRight(path, "/")
	if path == "" {
		path = "/"
	}

	for _, route := range m.routes {
		if strings.EqualFold(route.method, method) && route.pattern.MatchString(path) {
			return true
		}
	}

	return false
}

func routePattern(route string) string {
	parts := strings.Split(strings.Trim(route, "/"), "/")
	for i, part := range parts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			parts[i] = `[^/]+`
			continue
		}

		parts[i] = regexp.QuoteMeta(part)
	}

	return "^/" + strings.Join(parts, "/") + "$"
}
