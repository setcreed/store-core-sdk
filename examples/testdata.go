package examples

type User struct {
	UserId       int    `mapstructure:"user_id"`
	UserName     string `mapstructure:"user_name"`
	UserPassword string `mapstructure:"user_password"`
}
