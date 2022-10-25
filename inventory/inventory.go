package inventory

type Inventory struct {
	Inventory HostList `yaml:"inventory"`
}

type HostList struct {
	Hosts map[string]HostOptions `yaml:"hosts"`
}

type HostOptions struct {
	AnsiblePort int `yaml:"ansible_port,omitempty"`
	AsibleHost string `yaml:"ansible_host,omitempty"`
	AnsibleOsFamily string `yaml:"ansible_os_family,omitempty"`
}
