package config

type (
	Settings struct {
		Database DatabaseSettings `yaml:"database"`
		Telegram string           `yaml:"telegram"`
	}

	DatabaseSettings struct {
		Arguments string `yaml:"arguments"`
		Type      string `yaml:"type"`
	}

	Admins struct {
		Admins []string `yaml:"admins"`
	}
)
