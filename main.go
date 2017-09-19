package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/skuid/dewey/registry"
	"github.com/skuid/spec"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var pretty = flag.Bool("pretty", false, "pretty print output")
var interval = flag.String("interval", "30s", "sync interval")
var outputPath = flag.String("dir", "/opt/dewey/catalogs", "catalog file output directory")
var configFile = flag.String("config", "/opt/dewey/config.yaml", "config file location")

func init() {
	logger, err := spec.NewStandardLogger()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	viper.SetConfigFile("./config.yaml")
	viper.BindPFlags(flag.CommandLine)
}

func catalogRegistries() {
	err := viper.ReadInConfig()
	if err != nil {
		zap.L().Error(fmt.Sprintf("unable to read configuration: %s", err.Error()))
		return
	}
	r := []registry.RepoConfig{}
	err = viper.UnmarshalKey("registries", &r)
	if err != nil {
		zap.L().Error(err.Error())
		return
	}

	for _, reg := range r {

		catalog, err := registry.ConvertToCatalogableRegistry(reg)
		if err != nil {
			zap.L().Error(err.Error())
			continue
		}
		if catalog == nil {
			continue
		}

		c, err := catalog.GetCatalog()
		if err != nil {
			zap.L().Error(err.Error())
			continue
		}

		content, err := c.FileContent(viper.GetBool("pretty"))
		if err != nil {
			zap.L().Error(err.Error())
			continue
		}
		zap.L().Info("Updating catalog file", zap.String("registry", reg.Name))
		if err := ioutil.WriteFile(reg.Filename(viper.GetString("dir")), content, 0644); err != nil {
			zap.L().Error("error writing catalog file", zap.Error(err))
		}
	}
}

func main() {
	flag.Parse()
	window, err := time.ParseDuration(viper.GetString("interval"))
	if err != nil {
		zap.L().Fatal("unable to parse duration", zap.Error(err))
	}

	zap.L().Info("Bootstrapping catalog files")
	catalogRegistries()

	tickChannel := time.Tick(window)
	for range tickChannel {
		catalogRegistries()
	}
}
