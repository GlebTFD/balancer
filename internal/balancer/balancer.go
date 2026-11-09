package balancer

import (
	"net/http"
	"net/http/httputil"
	"sync/atomic"

	"github.com/hashicorp/go-hclog"
)

type LoadBalancer struct {
	log     hclog.Logger
	Servers []*Backend
	Counter uint64
	Proxy   *httputil.ReverseProxy
}

func NewLoadBalacner(log hclog.Logger, servers map[string]string) *LoadBalancer {
	backends := make([]*Backend, 0, len(servers))
	counter := uint64(0)

	for i, v := range servers {
		backends = append(backends, NewBackend(log, i, v))
	}

	lb := &LoadBalancer{
		log:     log,
		Servers: backends,
		Counter: counter,
		Proxy:   nil,
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetXForwarded()
			// Set server
			index := atomic.AddUint64(&lb.Counter, 1) % uint64(len(lb.Servers))
			r.SetURL(lb.Servers[index].TargetURL)
		},
		Transport: &http.Transport{
			MaxIdleConns: 100,
		},
	}

	lb.Proxy = proxy
	return lb
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lb.log.Info("Serve")
	lb.Proxy.ServeHTTP(w, r)
}
