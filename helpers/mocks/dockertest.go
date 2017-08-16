package mocks

import (
	"fmt"
	"github.com/ory-am/dockertest"
	"github.com/pkg/errors"
)

type PublishedPort struct {
	Number int

	// currently only "tcp" has been tested
	Protocol string
}

func (port *PublishedPort) String() string {
	return fmt.Sprintf("%s/%s", port.Number, port.Protocol)
}

///////////////////////////////////////////////////////////////////////////////

const (
	_DEFAULT_DOCKER_POOL = ""
)

// Docker related error messages
const (
	ERR_DOCKER_COULD_NOT_CONNECT = "Could not connect to docker: %s"
	ERR_DOCKER_RESOURCE_START = "Could not start resource: %s"
)

// TODO look in to possibly replacing this with dockertest.RunOptions
type DockertestContainer struct {
	// Image that should be pulled down to build the docker container
	Image string

	// This should match what we are currently using in production
	Tag string

	// Hostname or IP address that will be used when connecting to the container
	Host string

	// The name of the published port rather than the port on the container
	Port PublishedPort

	// Connection to the docker API that will be used to create and remove docker images
	Pool *dockertest.Pool

	// Represents the actual container
	Resource *dockertest.Resource
}

func (container *DockertestContainer) Start() (err error) {
	if container.Pool != nil {
		return // container already started
	}

	container.Pool, err = dockertest.NewPool(_DEFAULT_DOCKER_POOL)
	if err != nil {
		return errors.Wrap(err, ERR_DOCKER_COULD_NOT_CONNECT)
	}

	container.Resource, err = container.Pool.Run(container.Image, container.Tag, nil)
	if err != nil {
		return errors.Wrap(err, ERR_DOCKER_RESOURCE_START)
	}

	return err
}

func (container *DockertestContainer) Stop() error {
	return container.Pool.Purge(container.Resource)
}
