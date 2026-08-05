package network

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/proxy"
	"golang.org/x/time/rate"
)

// ProxyConfig содержит настройки SOCKS5-прокси.
type ProxyConfig struct {
	Network string
	Address string
	Auth    *proxy.Auth
}

// RateLimitConfig содержит ограничение частоты запросов и максимальный размер всплеска.
type RateLimitConfig struct {
	RPS   float64
	Burst int
}

// HTTPClientOption настраивает HTTP-клиент, создаваемый через NewHTTPClient.
type HTTPClientOption func(*httpClientConfig)

type httpClientConfig struct {
	transport http.RoundTripper
	timeout   time.Duration
	proxy     *ProxyConfig
	rateLimit *RateLimitConfig
}

// WithTransport задаёт базовый транспорт клиента.
func WithTransport(transport http.RoundTripper) HTTPClientOption {
	return func(config *httpClientConfig) {
		config.transport = transport
	}
}

// WithTimeout задаёт общий таймаут запроса.
func WithTimeout(timeout time.Duration) HTTPClientOption {
	return func(config *httpClientConfig) {
		config.timeout = timeout
	}
}

// WithProxy направляет запросы через SOCKS5-прокси.
func WithProxy(proxyConfig ProxyConfig) HTTPClientOption {
	return func(config *httpClientConfig) {
		config.proxy = &proxyConfig
	}
}

// WithRateLimit ограничивает частоту передачи запросов базовому транспорту.
func WithRateLimit(rateLimitConfig RateLimitConfig) HTTPClientOption {
	return func(config *httpClientConfig) {
		config.rateLimit = &rateLimitConfig
	}
}

// NewHTTPClient создаёт HTTP-клиент с переданными опциями.
func NewHTTPClient(options ...HTTPClientOption) (*http.Client, error) {
	config := httpClientConfig{
		transport: http.DefaultTransport,
	}

	for _, option := range options {
		if option != nil {
			option(&config)
		}
	}

	if config.proxy != nil {
		proxiedTransport, err := newProxyTransport(config.transport, *config.proxy)
		if err != nil {
			return nil, err
		}

		config.transport = proxiedTransport
	}

	if config.rateLimit != nil {
		limitedTransport, err := newRateLimitedTransport(config.transport, *config.rateLimit)
		if err != nil {
			return nil, err
		}

		config.transport = limitedTransport
	}

	return &http.Client{
		Transport: config.transport,
		Timeout:   config.timeout,
	}, nil
}

func newProxyTransport(next http.RoundTripper, config ProxyConfig) (http.RoundTripper, error) {
	if config.Network == "" {
		return nil, fmt.Errorf("proxy network is empty")
	}
	if config.Address == "" {
		return nil, fmt.Errorf("proxy address is empty")
	}

	transport, ok := next.(*http.Transport)
	if !ok {
		return nil, fmt.Errorf("proxy requires *http.Transport, got %T", next)
	}

	dialer, err := proxy.SOCKS5(config.Network, config.Address, config.Auth, &net.Dialer{
		KeepAlive: 30 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("create SOCKS5 dialer: %w", err)
	}

	contextDialer, ok := dialer.(proxy.ContextDialer)
	if !ok {
		return nil, fmt.Errorf("SOCKS5 dialer does not implement proxy.ContextDialer")
	}

	transport = transport.Clone()
	transport.DialContext = contextDialer.DialContext

	return transport, nil
}

type rateLimitedTransport struct {
	limiter *rate.Limiter
	next    http.RoundTripper
}

func newRateLimitedTransport(next http.RoundTripper, config RateLimitConfig) (http.RoundTripper, error) {
	if config.RPS <= 0 {
		return nil, fmt.Errorf("rate limit RPS must be greater than zero")
	}
	if config.Burst <= 0 {
		return nil, fmt.Errorf("rate limit burst must be greater than zero")
	}

	return &rateLimitedTransport{
		limiter: rate.NewLimiter(rate.Limit(config.RPS), config.Burst),
		next:    next,
	}, nil
}

func (t *rateLimitedTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if err := t.limiter.Wait(request.Context()); err != nil {
		return nil, err
	}

	return t.next.RoundTrip(request)
}
