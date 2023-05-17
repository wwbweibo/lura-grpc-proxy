package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	luraconfig "github.com/luraproject/lura/config"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/wwbweibo/lura-grpc-proxy/internal/domain"
	"github.com/wwbweibo/lura-grpc-proxy/internal/utils"
)

type Config struct {
	// Settings specified the dir of lura settings file
	Settings string `json:"settings"`
}

func (cfg *Config) Hook(v *viper.Viper) {
}

// RuntimeConfig is init during the runtime for per defined routes. something like kv
type RuntimeConfig struct {
	// ExecuteOn 使用正则表达式来匹配路径，获取该插件是否应该在该路径上执行
	ExecuteOn []string `json:"execute_on"`
}

func LoadRouteFromConfiguration(gRPCMethodName string) (string, string, error) {
	endpoint, host, err := loadAllInOneConfiguration(gRPCMethodName)
	if err == nil {
		return endpoint, host, nil
	}
	endpoint, host, err = loadTemplateConfiguration(gRPCMethodName)
	if err == nil {
		return endpoint, host, nil
	}
	return "", "", errors.New("could not find matched route")
}

func loadAllInOneConfiguration(gRPCMethodName string) (string, string, error) {
	path := tryParseConfigLocation()
	if path == "" {
		path = "/etc/krakend/krakend.json"
		// will try default path
	}
	config, err := luraconfig.NewParser().Parse(path)
	if err != nil {
		return "", "", errors.Wrap(err, "parse from configuration file "+path+" error")
	}
	for _, endpoint := range config.Endpoints {
		if config, ok := endpoint.Backend[0].ExtraConfig[HttpClientPluginName]; ok {
			cfg := config.(map[string]interface{})
			if proxyConfig, ok := cfg[PluginName]; ok {
				methodName := proxyConfig.(map[string]string)["method_name"]
				if methodName == gRPCMethodName {
					return fmt.Sprintf("%s:%s", endpoint.Method, endpoint.Endpoint), endpoint.Backend[0].Host[0], nil
				}
			}
		}
	}
	return "", "", errors.New("could not find matched route")
}

func loadTemplateConfiguration(gRPCMethodName string) (string, string, error) {
	settingsDir, exist := os.LookupEnv("FC_SETTINGS")
	if !exist || settingsDir == "" {
		return "", "", errors.New("not using template file")
	}
	files, err := ioutil.ReadDir(settingsDir)
	if err != nil {
		return "", "", errors.Wrap(err, "list file in dir "+settingsDir+" error.")
	}
	for _, f := range files {
		fmt.Println(f.Name())
		bts, err := os.ReadFile(path.Join(settingsDir, f.Name()))
		if err != nil {
			continue
		}
		tmpl := domain.Template{}
		err = json.Unmarshal(bts, &tmpl)
		if err != nil {
			continue
		}
		for _, route := range tmpl.Routes {
			if route.Backends[0].GrpcMethod == gRPCMethodName {
				return fmt.Sprintf("%s:%s", route.Method, route.Endpoint), route.Backends[0].Host, nil
			}
		}
	}
	return "", "", errors.New("could not find matched route")
}

func tryParseConfigLocation() string {
	cmdArgs := os.Args[1:]
	for idx, arg := range cmdArgs {
		if utils.StartWith(arg, "--") {
			if arg[2:] == "config" || arg[2:] == "c" {
				return cmdArgs[idx+1]
			}
			if strings.Contains(arg[2:], "=") {
				return strings.Split(arg[2:], "=")[1]
			}
		} else if utils.StartWith(arg, "-") {
			if arg[1:] == "config" || arg[1:] == "c" {
				return cmdArgs[idx+1]
			}
			if strings.Contains(arg[1:], "=") {
				return strings.Split(arg[1:], "=")[1]
			}
		} else {
			continue
		}
	}
	return ""
}
