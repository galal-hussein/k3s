package proxy

import (
	sysnet "net"
	"net/url"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/rancher/k3s/pkg/agent/loadbalancer"
)

type Proxy interface {
	Update(addresses []string)
	StartAPIServerProxy(port int) error
	SupervisorURL() string
	SupervisorAddresses() []string
	APIServerURL() string
	DisableProxy() error
	IsAPIServerProxyEnabled() bool
}

func NewAPIProxy(enabled bool, dataDir, supervisorURL string, lbServerPort int, etcdNode bool) (Proxy, error) {
	p := proxy{
		lbEnabled:            enabled,
		dataDir:              dataDir,
		initialSupervisorURL: supervisorURL,
		supervisorURL:        supervisorURL,
		apiServerURL:         supervisorURL,
		lbServerPort:         lbServerPort,
	}

	if enabled {
		lb, err := loadbalancer.New(dataDir, loadbalancer.SupervisorServiceName, supervisorURL, p.lbServerPort)
		if err != nil {
			return nil, err
		}
		lb.ETCDNode = etcdNode
		p.supervisorLB = lb
		p.supervisorURL = lb.LoadBalancerServerURL()
		p.apiServerURL = p.supervisorURL
	}

	u, err := url.Parse(p.initialSupervisorURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s", p.initialSupervisorURL)
	}
	p.fallbackSupervisorAddress = u.Host
	p.supervisorPort = u.Port()

	return &p, nil
}

type proxy struct {
	dataDir      string
	lbEnabled    bool
	lbServerPort int

	initialSupervisorURL      string
	supervisorURL             string
	supervisorPort            string
	fallbackSupervisorAddress string
	supervisorAddresses       []string
	supervisorLB              *loadbalancer.LoadBalancer
	updatedSupervisorURL      string

	apiServerURL     string
	apiServerLB      *loadbalancer.LoadBalancer
	apiServerEnabled bool
}

func (p *proxy) Update(addresses []string) {
	apiServerAddresses := addresses
	supervisorAddresses := addresses

	if p.apiServerEnabled {
		supervisorAddresses = p.setSupervisorPort(supervisorAddresses)
	}
	logrus.Info(apiServerAddresses)
	logrus.Info(supervisorAddresses)
	logrus.Infof("apiserverlb %#v supervisorlb %#v", p.apiServerLB, p.supervisorLB)
	if p.apiServerLB != nil {
		p.apiServerLB.Update(apiServerAddresses)
	}
	if p.supervisorLB != nil {
		p.supervisorLB.Update(supervisorAddresses)
	}
	logrus.Infof("supervisor address in update: %v", supervisorAddresses)
	p.supervisorAddresses = supervisorAddresses
	logrus.Infof("supervisor address in proxy: %v", p.supervisorAddresses)
}

func (p *proxy) setSupervisorPort(addresses []string) []string {
	var newAddresses []string
	for _, address := range addresses {
		h, _, err := sysnet.SplitHostPort(address)
		if err != nil {
			logrus.Errorf("Failed to parse address %s, dropping: %v", address, err)
			continue
		}
		newAddresses = append(newAddresses, sysnet.JoinHostPort(h, p.supervisorPort))
	}
	return newAddresses
}

func (p *proxy) StartAPIServerProxy(port int) error {
	u, err := url.Parse(p.initialSupervisorURL)
	if err != nil {
		return errors.Wrapf(err, "failed to parse server URL %s", p.initialSupervisorURL)
	}
	u.Host = sysnet.JoinHostPort(u.Hostname(), strconv.Itoa(port))

	p.apiServerURL = u.String()
	p.apiServerEnabled = true

	if p.lbEnabled {
		lbServerPort := p.lbServerPort
		if lbServerPort != 0 {
			lbServerPort = lbServerPort + 1
		}
		lb, err := loadbalancer.New(p.dataDir, loadbalancer.APIServerServiceName, p.apiServerURL, lbServerPort)
		if err != nil {
			return err
		}
		p.apiServerURL = lb.LoadBalancerServerURL()
		p.apiServerLB = lb
	}

	return nil
}

func (p *proxy) SupervisorURL() string {
	return p.supervisorURL
}

func (p *proxy) SupervisorAddresses() []string {
	logrus.Infof("supervisor addresses: %v", p.supervisorAddresses)
	if len(p.supervisorAddresses) > 0 {
		return p.supervisorAddresses
	}
	return []string{p.fallbackSupervisorAddress}
}

func (p *proxy) APIServerURL() string {
	return p.apiServerURL
}

func (p *proxy) DisableProxy() error {
	if p.apiServerLB != nil {
		if err := p.apiServerLB.Listener.Close(); err != nil {
			return err
		}
	}
	if p.supervisorLB != nil {
		if err := p.supervisorLB.Listener.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (p *proxy) IsAPIServerProxyEnabled() bool {
	return p.apiServerEnabled
}
