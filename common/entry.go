/* ©INFINI, All Rights Reserved.
 * mail: contact#infini.ltd */

package common

import (
	"infini.sh/framework/core/config"
)

type EntryConfig struct {
	Enabled          bool                 `config:"enabled"`
	Name             string               `config:"name"`
	MaxConcurrency   int                  `config:"max_concurrency"`
	TLSConfig        config.TLSConfig     `config:"tls"`
	NetworkConfig    config.NetworkConfig `config:"network"`
	RouterConfigName string               `config:"router"`
}

type RuleConfig struct {
	ID          string   `config:"id"`
	Description string   `config:"desc"`
	Method      []string `config:"method"`
	PathPattern []string `config:"pattern"`
	Flow        []string `config:"flow"`
}

type FilterConfig struct {
	Name       string                 `config:"name"`
	Type       string                 `config:"type"`
	Parameters map[string]interface{} `config:"parameters"`
}

type RouterConfig struct {
	Name        string       `config:"name"`
	DefaultFlow string       `config:"default_flow"`
	Rules       []RuleConfig `config:"rules"`
	TracingFlow string       `config:"tracing_flow"`
}

type FlowConfig struct {
	Name    string         `config:"name"`
	Filters []FilterConfig `config:"filter"`
}
