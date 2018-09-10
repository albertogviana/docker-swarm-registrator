package consul

import (
	"github.com/albertogviana/docker-swarm-registrator/service"
	consulapi "github.com/hashicorp/consul/api"
)

type Consul struct {
	ConsulClient *consulapi.Client
}

type ConsulServices interface {
	Register(service service.SwarmTask) error
	Deregister(serviceID string) error
}

func NewConsulClient(consulClient *consulapi.Client) *Consul {
	return &Consul{
		consulClient,
	}
}

// Register a service with consul
func (c *Consul) Register(service service.SwarmTask) error {
	register := consulapi.AgentServiceRegistration{
		ID:      service.ID,
		Name:    service.Name,
		Address: service.Address,
		Port:    service.Port,
	}

	err := c.ConsulClient.Agent().ServiceRegister(&register)
	if err != nil {
		return err
	}

	return nil
}

// Deregister a service with consul
func (c *Consul) Deregister(serviceID string) error {
	return c.ConsulClient.Agent().ServiceDeregister(serviceID)
}
