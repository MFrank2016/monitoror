package usecase

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/monitoror/monitoror/models"
	monitorableConfig "github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	"github.com/monitoror/monitoror/pkg/monitoror/utils"

	"github.com/jsdidierlaurent/echo-middleware/cache"
)

// Versions
const (
	CurrentVersion = Version1000
	MinimalVersion = Version1000

	Version1000 = "1.0" // Initial version
)

const (
	EmptyTileType models.TileType = "EMPTY"
	GroupTileType models.TileType = "GROUP"

	MetaTileStoreKeyPrefix = "monitoror.config.metaTile.key"
)

type (
	configUsecase struct {
		repository monitorableConfig.Repository

		tileConfigs     map[models.TileType]map[string]*TileConfig
		metaTileConfigs map[models.TileType]map[string]*MetaTileConfig

		// meta tile cache. used in case of timeout
		metaTileStore   cache.Store
		cacheExpiration time.Duration
	}

	// TileConfig struct is used by GetConfig endpoint to check / hydrate config
	TileConfig struct {
		Validator       utils.Validator
		Path            string
		InitialMaxDelay int
	}

	// MetaTileConfig struct is used by GetConfig endpoint to check / hydrate config
	MetaTileConfig struct {
		Validator utils.Validator
		Builder   builder.MetaTileBuilder
	}
)

func NewConfigUsecase(repository monitorableConfig.Repository, store cache.Store, downstreamStoreExpiration int) monitorableConfig.Usecase {
	tileConfigs := make(map[models.TileType]map[string]*TileConfig)

	// Used for authorized type
	tileConfigs[EmptyTileType] = nil
	tileConfigs[GroupTileType] = nil

	dynamicTileConfigs := make(map[models.TileType]map[string]*MetaTileConfig)

	return &configUsecase{
		repository:      repository,
		tileConfigs:     tileConfigs,
		metaTileConfigs: dynamicTileConfigs,
		metaTileStore:   store,
		cacheExpiration: time.Millisecond * time.Duration(downstreamStoreExpiration),
	}
}

func (cu *configUsecase) RegisterTile(
	tileType models.TileType, variant string, clientConfigValidator utils.Validator, path string, initialMaxDelay int,
) {
	value, exists := cu.tileConfigs[tileType]
	if !exists {
		value = make(map[string]*TileConfig)
		cu.tileConfigs[tileType] = value
	}

	value[variant] = &TileConfig{
		Path:            path,
		Validator:       clientConfigValidator,
		InitialMaxDelay: initialMaxDelay,
	}
}

func (cu *configUsecase) RegisterMetaTile(
	tileType models.TileType, variant string, clientConfigValidator utils.Validator, builder builder.MetaTileBuilder,
) {
	// Used for authorized type
	cu.tileConfigs[tileType] = nil

	value, exists := cu.metaTileConfigs[tileType]
	if !exists {
		value = make(map[string]*MetaTileConfig)
	}

	value[variant] = &MetaTileConfig{
		Validator: clientConfigValidator,
		Builder:   builder,
	}
	cu.metaTileConfigs[tileType] = value
}

func (cu *configUsecase) DisableTile(tileType models.TileType, variant string) {
	// TODO
}

// --- Utility functions ---
func keys(m interface{}) string {
	keys := reflect.ValueOf(m).MapKeys()
	strKeys := make([]string, len(keys))

	for i := 0; i < len(keys); i++ {
		strKeys[i] = fmt.Sprintf(`%v`, keys[i])
	}

	return strings.Join(strKeys, ", ")
}

func stringify(v interface{}) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}
