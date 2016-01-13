package serve

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/factionlabs/quorra/api"
	"github.com/factionlabs/quorra/version"
)

var Command = cli.Command{
	Name:   "serve",
	Usage:  "start server",
	Action: serveAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "listen, l",
			Usage: "listen address",
			Value: ":8080",
		},
		cli.StringFlag{
			Name:  "public-dir, s",
			Usage: "path to public media directory",
			Value: "public",
		},
	},
}

func serveAction(c *cli.Context) {
	log.Infof("%s %s", version.FullName(), version.FullVersion())

	listenAddr := c.String("listen")
	publicDir := c.String("public-dir")
	dbAddr := c.GlobalString("db-addr")
	dbName := c.GlobalString("db-name")
	storeKey := c.GlobalString("store-key")

	cfg := &api.APIConfig{
		ListenAddr: listenAddr,
		PublicDir:  publicDir,
		DBName:     dbName,
		DBAddr:     dbAddr,
		StoreKey:   storeKey,
	}

	a, err := api.NewAPI(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
