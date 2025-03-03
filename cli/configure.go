package spinup

// These are the updated struct definitions for the container.go file

type ContainerServiceWrapperUpdateInput struct {
	ForceRedeploy bool                         `json:"force_redeploy"`
	Service       *ContainerServiceUpdateInput `json:"service"`
	Size          *FlexInt                     `json:"size_id"`
}

type ContainerServiceUpdateInput struct {
	CapacityProviderStrategy []*CapacityProviderStrategyInput `json:"capacity_provider_strategy,omitempty"`
	ContainerDefinitions     []*ContainerDefinition           `json:"container_definitions,omitempty"`
	DesiredCount             int64                            `json:"desired_count,omitempty"`
	PlatformVersion          string                           `json:"platform_version,omitempty"`
}

type CapacityProviderStrategyInput struct {
	Base             int64  `json:"base"`
	CapacityProvider string `json:"capacity_provider"`
	Weight           int64  `json:"weight"`
}