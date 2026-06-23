package actions

import (
	"errors"
	"fmt"
	"strings"
)

type ProviderAction struct {
	PluginID        string           `json:"plugin_id"`
	StepType        string           `json:"step_type"`
	Contract        string           `json:"contract"`
	ContractVersion string           `json:"contract_version,omitempty"`
	RequiredConfig  []string         `json:"required_config,omitempty"`
	Config          []ProviderConfig `json:"config,omitempty"`
	OutputHandling  []string         `json:"output_handling,omitempty"`
	SubmitEndpoint  string           `json:"submit_endpoint,omitempty"`
	Prerequisites   []ProviderAction `json:"prerequisites,omitempty"`
	Fallbacks       []ProviderAction `json:"fallbacks,omitempty"`
}

type ProviderConfig struct {
	Name   string   `json:"name"`
	Value  string   `json:"value,omitempty"`
	Values []string `json:"values,omitempty"`
}

func (a ProviderAction) Validate() error {
	return a.validate(0)
}

func (a ProviderAction) validate(depth int) error {
	var errs []error
	if depth > 1 {
		return errors.New("provider actions must not be nested deeper than one level")
	}
	for _, field := range []struct {
		name  string
		value string
	}{
		{name: "plugin_id", value: a.PluginID},
		{name: "step_type", value: a.StepType},
		{name: "contract", value: a.Contract},
	} {
		if strings.TrimSpace(field.value) == "" {
			errs = append(errs, fmt.Errorf("%s is required", field.name))
			continue
		}
		if hasControlWhitespace(field.value) {
			errs = append(errs, fmt.Errorf("%s must not contain control whitespace", field.name))
		}
	}
	for _, field := range []struct {
		name  string
		value string
	}{
		{name: "contract_version", value: a.ContractVersion},
		{name: "submit_endpoint", value: a.SubmitEndpoint},
	} {
		if field.value != "" && hasControlWhitespace(field.value) {
			errs = append(errs, fmt.Errorf("%s must not contain control whitespace", field.name))
		}
	}
	for idx, field := range a.RequiredConfig {
		if strings.TrimSpace(field) == "" {
			errs = append(errs, fmt.Errorf("required_config[%d] is required", idx))
			continue
		}
		if hasControlWhitespace(field) {
			errs = append(errs, fmt.Errorf("required_config[%d] must not contain control whitespace", idx))
		}
	}
	for idx, field := range a.OutputHandling {
		if strings.TrimSpace(field) == "" {
			errs = append(errs, fmt.Errorf("output_handling[%d] is required", idx))
			continue
		}
		if hasControlWhitespace(field) {
			errs = append(errs, fmt.Errorf("output_handling[%d] must not contain control whitespace", idx))
		}
	}
	for idx, field := range a.Config {
		if err := field.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("config[%d]: %w", idx, err))
		}
	}
	for idx, prereq := range a.Prerequisites {
		if err := prereq.validate(depth + 1); err != nil {
			errs = append(errs, fmt.Errorf("prerequisites[%d]: %w", idx, err))
		}
	}
	for idx, fallback := range a.Fallbacks {
		if err := fallback.validate(depth + 1); err != nil {
			errs = append(errs, fmt.Errorf("fallbacks[%d]: %w", idx, err))
		}
	}
	return errors.Join(errs...)
}

func (c ProviderConfig) Validate() error {
	var errs []error
	if strings.TrimSpace(c.Name) == "" {
		errs = append(errs, errors.New("name is required"))
	} else if hasControlWhitespace(c.Name) {
		errs = append(errs, errors.New("name must not contain control whitespace"))
	}
	if c.Value != "" && hasControlWhitespace(c.Value) {
		errs = append(errs, errors.New("value must not contain control whitespace"))
	}
	for idx, value := range c.Values {
		if strings.TrimSpace(value) == "" {
			errs = append(errs, fmt.Errorf("values[%d] is required", idx))
			continue
		}
		if hasControlWhitespace(value) {
			errs = append(errs, fmt.Errorf("values[%d] must not contain control whitespace", idx))
		}
	}
	if c.Value == "" && len(c.Values) == 0 {
		errs = append(errs, errors.New("value or values is required"))
	}
	return errors.Join(errs...)
}

func hasControlWhitespace(value string) bool {
	return strings.ContainsAny(value, "\r\n\t")
}
