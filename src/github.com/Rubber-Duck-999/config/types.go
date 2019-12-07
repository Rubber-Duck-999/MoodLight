package config

type ConfigTypes struct {
	EmailSettings struct {
		Email    string `yaml:"email"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		To_email string `yaml:"to_email"`
	} `yaml:"email_settings"`
}
