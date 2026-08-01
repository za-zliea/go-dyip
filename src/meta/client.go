package meta

type ClientMeta struct {
	Server    string   `yaml:"server"`
	Token     string   `yaml:"token"`
	Domain    string   `yaml:"domain"`
	Protocol  Protocol `yaml:"protocol"`
	Local     bool     `yaml:"local"`
	Interface string   `yaml:"interface,omitempty"`
	Auth      string   `yaml:"auth"`
	Interval  int      `yaml:"interval"`
}

func (c *ClientMeta) Generate() {
	c.Server = "http://127.0.0.1:8080/"
	c.Token = "your-token-abcde12345"
	c.Domain = "your-subdomain.your-doamin"
	c.Auth = "your-doamin-token-abce12345"
	c.Local = false
	c.Interval = 300
	c.Protocol = IPV4
}
func (c *ClientMeta) Empty() bool {
	return c.Server == "" || c.Token == ""
}
