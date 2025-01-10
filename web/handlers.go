package web

import (
	"fmt"
	"net/http"
	"time"
	"xray-checker/checker"
	"xray-checker/config"
	"xray-checker/metrics"
	"xray-checker/models"
)

var registeredEndpoints []EndpointInfo

type EndpointInfo struct {
	Name      string
	URL       string
	ProxyPort int
}

func IndexHandler(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		data := PageData{
			Version:            version,
			Port:               config.CLIConfig.Metrics.Port,
			CheckInterval:      config.CLIConfig.Proxy.CheckInterval,
			IPCheckUrl:         config.CLIConfig.Proxy.IpCheckUrl,
			CheckMethod:        config.CLIConfig.Proxy.CheckMethod,
			StatusCheckUrl:     config.CLIConfig.Proxy.StatusCheckUrl,
			SimulateLatency:    config.CLIConfig.Proxy.SimulateLatency,
			Timeout:            config.CLIConfig.Proxy.Timeout,
			SubscriptionUpdate: config.CLIConfig.Subscription.Update,
			StartPort:          config.CLIConfig.Xray.StartPort,
			Instance:           config.CLIConfig.Metrics.Instance,
			PushUrl:            metrics.GetPushURL(config.CLIConfig.Metrics.PushURL),
			Endpoints:          registeredEndpoints,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := RenderIndex(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func BasicAuthMiddleware(username, password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok || user != username || pass != password {
				w.Header().Set("WWW-Authenticate", `Basic realm="metrics"`)
				http.Error(w, "Unauthorized.", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func ConfigStatusHandler(proxyChecker *checker.ProxyChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len("/config/"):]
		if path == "" {
			http.Error(w, "Config path is required", http.StatusBadRequest)
			return
		}

		var found *models.ProxyConfig
		for _, proxy := range proxyChecker.GetProxies() {
			proxyPath := fmt.Sprintf("%s-%s-%d", proxy.Protocol, proxy.Server, proxy.Port)
			if proxyPath == path {
				found = proxy
				break
			}
		}

		if found == nil {
			http.Error(w, "Config not found", http.StatusNotFound)
			return
		}

		status, latency, err := proxyChecker.GetProxyStatus(found.Name)
		if err != nil {
			http.Error(w, "Config not found", http.StatusNotFound)
			return
		}

		if config.CLIConfig.Proxy.SimulateLatency {
			time.Sleep(time.Duration(latency))
		}

		if status {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Failed"))
		}
	}
}

func RegisterConfigEndpoints(proxies []*models.ProxyConfig, startPort int) {
	registeredEndpoints = make([]EndpointInfo, 0, len(proxies))

	for i, proxy := range proxies {
		endpoint := fmt.Sprintf("/config/%s-%s-%d", proxy.Protocol, proxy.Server, proxy.Port)

		registeredEndpoints = append(registeredEndpoints, EndpointInfo{
			Name:      fmt.Sprintf("%s (%s:%d)", proxy.Name, proxy.Server, proxy.Port),
			URL:       endpoint,
			ProxyPort: startPort + i,
		})
	}
}
