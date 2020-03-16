package builder

import "github.com/monitoror/monitoror/models"

type (
	MetaTileBuilder func(params interface{}) ([]Result, error)

	Result struct {
		TileType models.TileType
		Label    string
		Params   map[string]interface{}
	}
)
