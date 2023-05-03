package oam

import (
	"encoding/json"

	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
)

// ParseApplicationComponent converts json application component to go struct
func ParseApplicationComponent(jsn string) (v1alpha1.Component, error) {
	var acomp v1alpha1.Component
	err := json.Unmarshal([]byte(jsn), &acomp)
	return acomp, err
}

// ParseApplicationConfiguration converts json application configuration to go struct
func ParseApplicationConfiguration(jsn string) (v1alpha1.Configuration, error) {
	var acomp v1alpha1.Configuration
	err := json.Unmarshal([]byte(jsn), &acomp)
	return acomp, err
}
