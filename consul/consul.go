package consul

import (
	"fmt"
	"strconv"

	"github.com/albertogviana/docker-swarm-registrator/service"
	consulapi "github.com/hashicorp/consul/api"
)

// Consul defines the based structure
type Consul struct {
	ConsulClient *consulapi.Client
}

// Services defines interface with mandatory methods
type Services interface {
	Register(service service.SwarmTask) error
	Deregister(serviceID string) error
}

// NewConsulClient returns a new instance of the Consul structure
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

	if len(service.Checks) > 0 {
		for _, c := range service.Checks {
			agentServiceCheck := consulapi.AgentServiceCheck{}
			agentServiceCheck.CheckID = fmt.Sprintf("%s-%d", c.ID, service.Port)
			agentServiceCheck.Name = c.Name
			agentServiceCheck.Interval = c.Interval
			agentServiceCheck.Timeout = c.Timeout

			http, err := strconv.ParseBool(c.HTTP)
			if err != nil {
				return err
			}

			if http {
				agentServiceCheck.HTTP = fmt.Sprintf("http://%s:%d%s", service.Address, service.Port, c.Path)
			}

			agentServiceCheck.DeregisterCriticalServiceAfter = c.RemoveFailedServiceAfter
			register.Checks = append(register.Checks, &agentServiceCheck)
		}
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
