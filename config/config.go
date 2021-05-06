package config

import (
	"time"
)

// Config defines the configuration structure.
type Config struct {
	Integration struct {
		MQTT struct {
			EventTopicTemplate      string        `mapstructure:"event_topic_template"`
			CommandTopicTemplate    string        `mapstructure:"command_topic_template"`
			KeepAlive               time.Duration `mapstructure:"keep_alive"`
			MaxReconnectInterval    time.Duration `mapstructure:"max_reconnect_interval"`
			TerminateOnConnectError bool          `mapstructure:"terminate_on_connect_error"`

			Auth struct {
				Type string `mapstructure:"type"`

				Generic struct {
					Server       string   `mapstructure:"server"`
					Servers      []string `mapstructure:"servers"`
					Username     string   `mapstructure:"username"`
					Password     string   `mapstrucure:"password"`
					CACert       string   `mapstructure:"ca_cert"`
					TLSCert      string   `mapstructure:"tls_cert"`
					TLSKey       string   `mapstructure:"tls_key"`
					QOS          uint8    `mapstructure:"qos"`
					CleanSession bool     `mapstructure:"clean_session"`
					ClientID     string   `mapstructure:"client_id"`
				} `mapstructure:"generic"`
			} `mapstructure:"auth"`
		} `mapstructure:"mqtt"`
	} `mapstructure:"integration"`
}

// C holds the global configuration.
var C Config
