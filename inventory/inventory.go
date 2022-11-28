package inventory

type Inventory struct {
	Inventory HostList `yaml:"inventory"`
}

type HostList struct {
	Hosts map[string]HostOptions `yaml:"hosts"`
}

type HostOptions struct {
	AnisbleUser string `yaml:"ansible_user,omitempty"`
	AnsiblePassword string `yaml:"ansible_ssh_pass,omitempty"`
	AnsibleSudoPassword string `yaml:"ansible_sudo_pass,omitempty"`
	AnsiblePort int `yaml:"ansible_port,omitempty"`
	AsibleHost string `yaml:"ansible_host,omitempty"`
	AnsibleOsFamily string `yaml:"ansible_os_family,omitempty"`
}
