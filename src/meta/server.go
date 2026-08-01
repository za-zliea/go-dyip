package meta

import (
	"crypto/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	adminDefaultUserName  = "admin"
	adminPasswordLen      = 16
	adminPasswordAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	adminFallbackPassword = "CHANGE_ME_IMMEDIATELY"
)

type ServerMeta struct {
	Address string             `yaml:"address"`
	Port    int                `yaml:"port"`
	Token   string             `yaml:"token"`
	RealIp  *string            `yaml:"realip,omitempty"`
	Admin   *Admin             `yaml:"admin"`
	Metas   []*IpMeta          `yaml:"ips"`
	MetaMap map[string]*IpMeta `yaml:"-"`
}

type Admin struct {
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
}

// HistoryEntry is one element of IpMeta.History: the IP that was synced and when.
type HistoryEntry struct {
	Ip   string    `yaml:"ip" json:"ip"`
	Time time.Time `yaml:"time" json:"time"`
}

type IpMeta struct {
	Provider        string         `yaml:"provider"`
	Accesskey       string         `yaml:"ak,omitempty"`
	AccessKeySecret string         `yaml:"sk,omitempty"`
	Domain          string         `yaml:"domain"`
	Subdomain       string         `yaml:"subdomain"`
	Auth            string         `yaml:"auth"`
	Protocol        Protocol       `yaml:"protocol"`
	SyncType        string         `yaml:"synctype"`
	Ip              *string        `yaml:"ip,omitempty"`
	History         []HistoryEntry `yaml:"history,omitempty"`
}

func (s *ServerMeta) Generate() {
	metas := make([]*IpMeta, 1)
	ipMeta := IpMeta{
		Provider:        "your-provider (TENCENT/ALIYUN/GODADDY/GOOGLE/CLOUDFLARE)",
		Accesskey:       "abcde12345",
		AccessKeySecret: "abcde12345",
		Domain:          "your-doamin",
		Subdomain:       "your-subdomain",
		Auth:            "your-doamin-token-abce12345",
		Protocol:        "IPV4",
		SyncType:        "00",
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
		metaMap[s.Metas[i].Subdomain+"."+s.Metas[i].Domain+"."+string(s.Metas[i].Protocol)] = s.Metas[i]
	}

	s.MetaMap = metaMap
}

// generateRandomPassword generates a cryptographically secure alphanumeric password.
func generateRandomPassword(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	out := make([]byte, n)
	for i, b := range buf {
		out[i] = adminPasswordAlphabet[int(b)%len(adminPasswordAlphabet)]
	}
	return string(out), nil
}

// EnsureAdmin populates s.Admin with a freshly generated bcrypt-hashed credential if no
// usable admin account is present. It returns the generated plaintext password and true
// when a new account was created, or "" and false when the existing admin is retained.
func (s *ServerMeta) EnsureAdmin() (string, bool) {
	if s.Admin != nil && s.Admin.UserName != "" && s.Admin.Password != "" {
		return "", false
	}

	plaintext, err := generateRandomPassword(adminPasswordLen)
	if err != nil {
		// Extremely unlikely; surface a static placeholder so operators notice and replace it.
		plaintext = adminFallbackPassword
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		// Should never happen with bcrypt; fall back to storing the plaintext verbatim.
		hashed = []byte(plaintext)
	}

	s.Admin = &Admin{
		UserName: adminDefaultUserName,
		Password: string(hashed),
	}
	return plaintext, true
}
