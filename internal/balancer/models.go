package balancer

import (
	"net/url"

	"github.com/hashicorp/go-hclog"
)

type Backend struct {
	log       hclog.Logger
	Name      string
	TargetURL *url.URL
}

func NewBackend(log hclog.Logger, serverName string, serverPath string) *Backend {
	url, err := url.Parse(serverPath)
	if err != nil {
		log.Error("Error to parse url")
	}

	return &Backend{
		Name:      serverName,
		TargetURL: url,
	}
}
