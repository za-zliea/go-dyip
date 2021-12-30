package meta

type ServerMeta struct {
	Address string             `yaml:"address"`
	Port    int                `yaml:"port"`
	Token   string             `yaml:"token"`
	RealIp  *string            `yaml:"realip"`
	Metas   []*IpMeta          `yaml:"ips"`
	MetaMap map[string]*IpMeta `yaml:"-"`
}

type IpMeta struct {
	Provider        string   `yaml:"provider"`
	Accesskey       string   `yaml:"ak"`
	AccessKeySecret string   `yaml:"sk"`
	Domain          string   `yaml:"domain"`
	Subdomain       string   `yaml:"subdomain"`
	Auth            string   `yaml:"auth"`
	Ip              *string  `yaml:"ip,omitempty"`
	History         []string `yaml:"history,omitempty"`
}

func (s *ServerMeta) Generate() {
	metas := make([]*IpMeta, 1)
	ipMeta := IpMeta{
		Provider:        "your-provider (TENCENT/ALIYUN/GODADDY/GOOGLE)",
		Accesskey:       "abcde12345",
		AccessKeySecret: "abcde12345",
		Domain:          "your-doamin",
		Subdomain:       "your-subdomain",
		Auth:            "your-doamin-token-abce12345",
		Ip:              nil,
		History:         nil,
	}
	metas[0] = &ipMeta

	s.Address = "127.0.0.1"
	s.Port = 8080
	realIpNote := "x-real-ip"
	s.RealIp = &realIpNote
	s.Token = "your-token-abcde12345"
	s.Metas = metas
}
func (s *ServerMeta) Empty() bool {
	return s.Address == "" || s.Token == "" || s.Metas == nil || len(s.Metas) == 0
}
func (s *ServerMeta) GenerateIpm() {
	metaMap := make(map[string]*IpMeta)

	for i := 0; i < len(s.Metas); i += 1 {
		metaMap[s.Metas[i].Subdomain+"."+s.Metas[i].Domain] = s.Metas[i]
	}

	s.MetaMap = metaMap
}
