package config

type Queue struct {
	Uri      string `default:"amqp://guest:guest@localhost:5672"`
	Username string `default:"guest"`
	Password string `default:"guest"`
	Exchange string `default:"amq.topic"`
	Vhost    string `default:"%2F"`
}
