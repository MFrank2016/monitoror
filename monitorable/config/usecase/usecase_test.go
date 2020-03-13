package usecase

import (
	"io/ioutil"
	"strings"

	. "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/config/models"
	"github.com/monitoror/monitoror/monitorable/config/repository"
	"github.com/monitoror/monitoror/monitorable/jenkins"
	jenkinsModels "github.com/monitoror/monitoror/monitorable/jenkins/models"
	"github.com/monitoror/monitoror/monitorable/ping"
	pingModels "github.com/monitoror/monitoror/monitorable/ping/models"
	"github.com/monitoror/monitoror/monitorable/pingdom"
	pindomModels "github.com/monitoror/monitoror/monitorable/pingdom/models"
	"github.com/monitoror/monitoror/monitorable/port"
	portModels "github.com/monitoror/monitoror/monitorable/port/models"

	"github.com/jsdidierlaurent/echo-middleware/cache"
)

func initConfigUsecase(repository config.Repository, store cache.Store) *configUsecase {
	usecase := NewConfigUsecase(repository, store, 5000)

	usecase.RegisterTile(ping.PingTileType, DefaultVariant, &pingModels.PingParams{}, "/ping", 1000)
	usecase.RegisterTile(port.PortTileType, DefaultVariant, &portModels.PortParams{}, "/port", 1000)
	usecase.RegisterTile(jenkins.JenkinsBuildTileType, DefaultVariant, &jenkinsModels.BuildParams{}, "/jenkins/default", 1000)
	usecase.RegisterTile(pingdom.PingdomCheckTileType, DefaultVariant, &pindomModels.CheckParams{}, "/pingdom/default", 1000)

	return usecase.(*configUsecase)
}

func readConfig(input string) (configBag *models.ConfigBag, err error) {
	configBag = &models.ConfigBag{}
	reader := ioutil.NopCloser(strings.NewReader(input))
	configBag.Config, err = repository.ReadConfig(reader)
	return
}
