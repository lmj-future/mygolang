package integration

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/h3c/mygolang/config"
	"github.com/h3c/mygolang/integration/mqtt"
)

var integration Integration

// Setup configures the integration.
func Setup(conf config.Config) error {
	fmt.Println("Setup")
	fmt.Println(conf)
	var err error
	integration, err = mqtt.NewBackend(conf)
	if err != nil {
		return errors.Wrap(err, "setup mqtt integration error")
	}

	return nil
}

// GetIntegration returns the integration.
func GetIntegration() Integration {
	return integration
}

// Integration defines the interface that an integration must implement.
type Integration interface {
	// Start starts the integration.
	Start() error

	// Stop stops the integration.
	Stop() error
}
