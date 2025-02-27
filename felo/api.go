package felo

import (
	"zlutils/consul"
)

func checkBalance() {

}

type ProductConfig struct {
	ProductId string `json:"product_id"` // 商品id
	Bonus     int64  `json:"bonus"`      // 奖励金币
	Tag       Tag    `json:"tag"`        // 标签
}

type Tag struct {
	TagType string `json:"type"` // tag类型
	Desc    string `json:"desc"` // 文案
}

var productConfigList []*ProductConfig
var productConfigMap map[string]*ProductConfig // productId - productConfig

func InitProductConfigCfg() {
	productConfigList = make([]*ProductConfig, 0)
	productConfigMap = make(map[string]*ProductConfig)

	consul.WatchJson("me-live-buy/client_product_config.json", &productConfigList, func() {
		for _, cfg := range productConfigList {
			productConfigMap[cfg.ProductId] = cfg
		}
	})
}

func getProductConfig(productId string) (pro ProductConfig) {
	if v, has := productConfigMap[productId]; has {
		pro = *v
	}
	return
}
