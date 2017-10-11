package cmd

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/andrexus/cloud-initer/api"
	"github.com/andrexus/cloud-initer/conf"
	"github.com/spf13/cobra"
	"github.com/xlab/closer"
)

var serveCmd = cobra.Command{
	Use:   "serve",
	Short: "Start API server",
	Long:  "Start API server on specified host and port",
	Run: func(cmd *cobra.Command, args []string) {
		execWithConfig(cmd, serve)
	},
}

func serve(config *conf.Configuration) {
	db, err := conf.BoltConnect(config)
	if err != nil {
		logrus.Fatalf("Error opening database: %+v", err)
	}

	apiServer := api.NewAPI(config, db)

	l := fmt.Sprintf("%v:%v", config.API.Host, config.API.Port)
	logrus.Infof("API started on: %s", l)

	closer.Bind(func() {
		err := apiServer.Stop()
		if err != nil {
			logrus.Errorf("Error stopping API server: %s", err.Error())
		}
	})

	apiServer.Start()
}
