package conf

import (
	"errors"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"net/url"
)

type Config struct {
	Settings YamlConfig
	FilePath string
}

type YamlConfig struct {
	Destination   string
	LogLevel      string `bool:"log-level"`
	Elasticsearch struct {
		Uri          string
		Username     string
		Password     string
		UseZincAwsS3 bool `yaml:"use-zinc-aws-s3"`
		NoDuplicates bool `yaml:"no-duplicates"`
	}

	Github struct {
		Token string
	}

	Script struct {
		Path       string
		MaxWorkers int    `yaml:"max-workers"`
		LogFile    string `yaml:"log-file"`
	}

	Slack struct {
		Token     string
		ChannelID string `yaml:"channel-id"`
	}

	Truffle struct {
		Path        string
		MaxWorkers  int    `yaml:"max-workers"`
		LogFile     string `yaml:"log-file"`
		GitHubToken string `yaml:"github-token"`
		Ignore      []string
	}
	Database struct {
		Engine string
		Dsn    string
		Path   string
	}
}

func (c *Config) validate() error {
	s := c.Settings
	if s.Destination == "elastic" {
		esUri := s.Elasticsearch.Uri
		if _, err := url.ParseRequestURI(esUri); err != nil {
			return errors.New("Invalid elastic search Uri: " + esUri)
		}
	}
	return nil
}

func (c *Config) Load() error {
	if c.FilePath == "" {
		c.FilePath = "./config.yaml"
	}
	f, err := os.Open(c.FilePath)
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("Unable to load config file: " + c.FilePath + ", using runtime settings.")
	}
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&c.Settings)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
