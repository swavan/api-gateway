package user

type Config struct {
	Create struct {
		Scripts []string `mapstructure:"script"`
	} `mapstructure:"create"`
}
