package ev

type Config struct {
	QemuPath string
}

func NewConfig() *Config {
	c := &Config{}
	c.WithDefaultValues()
	return c
}

func (c *Config) WithDefaultValues() {
	if c.QemuPath == "" {
		c.QemuPath = "/usr/bin/qemu-system-x86_64"
	}
}
