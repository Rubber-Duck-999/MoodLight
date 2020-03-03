package config

type ConfigTypes struct {
	EmailSettings struct {
		Email    string `yaml:"email"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		To_email string `yaml:"to_email"`
	} `yaml:"email_settings"`
	MessageSettings struct {
		Sid      string `yaml:"sid"`
		Token    string `yaml:"token"`
		From_num string `yaml:"from_num"`
		To_num   string `yaml:"to_num"`
	} `yaml:"message_settings"`
}
