package config

import (
	"flag"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	serverUrl string
	username  string
	password  string
	secure    string
}

func Get() *Config {
	conf := &Config{}
	flag.StringVar(&conf.secure, "secure", os.Getenv("SECURE"), "Is Secure")
	flag.StringVar(&conf.username, "username", os.Getenv("USERNAME"), "Basic Auth username")
	flag.StringVar(&conf.password, "password", os.Getenv("PASSWORD"), "Basic Auth password")

	flag.StringVar(&conf.serverUrl, "serverUrl", os.Getenv("SERVER_URL"), "Server Url")

	flag.Parse()

	return conf
}

func (c *Config) IsSecure() bool {
	if c.secure == "" {
		return true
	} else if strings.ToLower(c.secure) == "false" {
		return false
	}
	return true
}

func (c *Config) GetUsernameAndPassword() (string, string) {
	return c.username, c.password
}

func (c *Config) GetServerURL() string {
	c.serverUrl = strings.TrimSuffix(c.serverUrl, "/")
	u, _ := url.Parse(c.serverUrl)
	if u.Scheme == "" {
		return "https://" + c.serverUrl
	} else {
		return c.serverUrl
	}
}

func (c *Config) GetServerHost() string {
	c.serverUrl = strings.TrimSuffix(c.serverUrl, "/")
	u, _ := url.Parse(c.serverUrl)
	if u.Scheme == "" {
		return u.Host
	} else {
		return u.Host
	}
}
