package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/lukaszbudnik/migrator/common"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
)

// Config represents Migrator's yaml configuration file
type Config struct {
	BaseDir           string   `yaml:"baseDir,omitempty"`
	BaseLocation      string   `yaml:"baseLocation" validate:"required"`
	Driver            string   `yaml:"driver" validate:"required"`
	DataSource        string   `yaml:"dataSource" validate:"required"`
	TenantSelectSQL   string   `yaml:"tenantSelectSQL,omitempty"`
	TenantInsertSQL   string   `yaml:"tenantInsertSQL,omitempty"`
	SchemaPlaceHolder string   `yaml:"schemaPlaceHolder,omitempty"`
	SingleMigrations  []string `yaml:"singleMigrations" validate:"min=1"`
	TenantMigrations  []string `yaml:"tenantMigrations,omitempty"`
	SingleScripts     []string `yaml:"singleScripts,omitempty"`
	TenantScripts     []string `yaml:"tenantScripts,omitempty"`
	Port              string   `yaml:"port,omitempty"`
	PathPrefix        string   `yaml:"pathPrefix,omitempty"`
	WebHookURL        string   `yaml:"webHookURL,omitempty"`
	WebHookHeaders    []string `yaml:"webHookHeaders,omitempty"`
}

func (config Config) String() string {
	c, _ := yaml.Marshal(config)
	return strings.TrimSpace(string(c))
}

// FromFile reads config from file which name is passed as an argument
func FromFile(configFileName string) (*Config, error) {
	contents, err := ioutil.ReadFile(configFileName)

	if err != nil {
		return nil, err
	}

	return FromBytes(contents)
}

// FromBytes reads config from raw bytes passed as an argument
func FromBytes(contents []byte) (*Config, error) {
	var config Config

	if err := yaml.Unmarshal(contents, &config); err != nil {
		return nil, err
	}

	if len(config.BaseDir) > 0 && len(config.BaseLocation) == 0 {
		common.Log("WARN", "Deprecated: config property `baseDir` will be removed in migrator v2021.1.0, please rename it to `baseLocation`")
		config.BaseLocation = config.BaseDir
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, err
	}

	substituteEnvVariables(&config)

	return &config, nil
}

func substituteEnvVariables(config *Config) {
	val := reflect.ValueOf(config).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		if val.CanAddr() && val.CanSet() {
			switch typeField.Type.Kind() {
			case reflect.String:
				s := valueField.Interface().(string)
				s = substituteEnvVariable(s)
				valueField.SetString(s)
			case reflect.Slice:
				ss := valueField.Interface().([]string)
				for i := range ss {
					ss[i] = substituteEnvVariable(ss[i])
				}
				valueField.Set(reflect.ValueOf(ss))
			}
		}
	}
}

func substituteEnvVariable(s string) string {
	start := strings.Index(s, "${")
	end := strings.Index(s, "}")
	if start > -1 && end > start && len(s) > end {
		after := s[0:start] + os.Getenv(s[start+2:end]) + s[end+1:]
		return substituteEnvVariable(after)
	}
	return s
}
