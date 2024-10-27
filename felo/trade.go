package felo

import (
	"encoding/json"
	"fmt"
)

type ProductSt struct {
	ID            string           `json:"id"`
	ProductID     string           `json:"product_id"` // huawei
	Name          string           `json:"name,omitempty"`
	Price         int64            `json:"price"`
	OutCode       string           `json:"out_code"`
	Desc          string           `json:"desc,omitempty"`
	Coins         int64            `json:"coins"`
	Dt            int              `json:"dt"`
	Type          string           `json:"type"`
	Currency      string           `json:"currency"`
	ChatTag       string           `json:"chat_tag,omitempty"`
	ActivityBonus *ActivityBonusSt `json:"activity_bonus"`
	RechargeTag   *RechargeTag     `json:"recharge_tag"`
}

type RechargeTag struct {
	TagType string `json:"tag_type"` // tag类型
	TagDesc string `json:"tag_desc"` // tag文案
}

type ActivityBonusSt struct {
	Tag       string `json:"tag"`
	TotalCoin int64  `json:"total_coin"`
	Num       int64  `json:"num"`
	Limit     int64  `json:"limit"`
	Weight    int64  `json:"weight"`
	Desc      string `json:"desc"`
}

type OpenDiscountCfg struct {
	ConfigName string                          `json:"config_name"` //配置名称，需唯一
	OpenArea   []int64                         `json:"open_area"`   //支持的大区
	OpenPrice  []int64                         `json:"open_price"`  //支持的官方金币档位
	St         int64                           `json:"st"`          //配置开始时间
	Et         int64                           `json:"et"`          //配置结束时间
	Score      int64                           `json:"score"`       //生效排序分，分数越大优先级越高
	LimitNum   int64                           `json:"limit_num"`   //返利最多可生效次数
	BonusMap   map[int32]*OpenDiscountBonusCfg `json:"bonus_map"`
	TagMap     map[int64]*Tag                  `json:"tag_map"`
}

type UserDiscountCfg struct {
	OfficialBonusMap map[int64]*DiscountInfoSt `json:"official_bonus_map"` //key为金币数，值为包含返利配置来源配置名及配置值
	OfficialRatioMap map[int64]float64         `json:"official_ratio_map"` //官方折扣返利折合USD值
	ThreeBonusMap    map[int64]*DiscountInfoSt `json:"three_bonus_map"`    //key为金币数，值为包含返利配置来源配置名及配置值
}

type DiscountInfoSt struct {
	ConfigName string `json:"config_name"`
	IntValue   int64  `json:"int_value"`
	LimitNum   int64  `json:"limit_num"`
	TagType    string `json:"tag_type"`
	TagDsc     string `json:"tag_dsc"`
}

type OpenDiscountBonusCfg struct {
	OfficialBonusMap map[int64]int64   `json:"official_bonus_map"`
	OfficialRatioMap map[int64]float64 `json:"official_ratio_map"`
	ThreeBonusMap    map[int64]int64   `json:"three_bonus_map"`
	ThirdPayTag      []int64           `json:"third_pay_tap"`
}

var tagConfigList []*Tag
var tagConfigMap map[int64]*Tag

func ObjectToString(o interface{}) string {
	data, err := json.Marshal(o)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

//func init() {
//	consul.Init("127.0.0.1:8500", "me-srvconfs")
//	InitProductConfigCfg()
//}

func GetProductTag(list []*ProductSt) {
	for i, item := range list {
		if v, has := tagConfigMap[item.Coins]; has {
			list[i].RechargeTag = &RechargeTag{
				TagType: v.TagType,
				TagDesc: v.Desc,
			}
		}
	}
}

func GetBonusOfficialPayList(param BaseParam, from string, list []*ProductSt) {
	// 开业大酬宾 - 官方列表
	cfg, useDiscountRecordMap, err := GetUserDiscountCfg(param)
	if err != nil {
		fmt.Printf("GetBonusOfficialPayList failed, err;%v", err)
		return
	}
	if cfg == nil {
		return
	}
	for i, item := range list {
		if cfg.OfficialBonusMap[item.Coins] != nil && cfg.OfficialBonusMap[item.Coins].IntValue != 0 {
			bonusCoin := cfg.OfficialBonusMap[item.Coins].IntValue
			cfgName := cfg.OfficialBonusMap[item.Coins].ConfigName
			//ratio := int64(float64(bonusCoin) / float64(item.Coins) * 100)
			//bonusPrice := cfg.OfficialRatioMap[item.Coins]
			num := useDiscountRecordMap[fmt.Sprintf("%v_%v", cfgName, item.Coins)]
			activity := &ActivityBonusSt{
				TotalCoin: item.Coins + bonusCoin,
				Tag:       cfg.OfficialBonusMap[item.Coins].TagType,
				Num:       num,
				Limit:     cfg.OfficialBonusMap[item.Coins].LimitNum,
				Weight:    int64(i + 1),
				Desc:      cfg.OfficialBonusMap[item.Coins].TagDsc, // todo z 多语
			}
			list[i].ActivityBonus = activity
		}
	}
	fmt.Println("..")
}

func GetUserDiscountCfg(param BaseParam) (userDiscountCfg *UserDiscountCfg, discountRecordMap map[string]int64, err error) {
	//userDiscountCfg = &UserDiscountCfg{}
	//
	//
	//userDiscountCfg.ThreeBonusMap = make(map[int64]*DiscountInfoSt, 0)
	//userDiscountCfg.OfficialBonusMap = make(map[int64]*DiscountInfoSt, 0)
	//userDiscountCfg.OfficialRatioMap = make(map[int64]float64, 0)
	//if err != nil {
	//	fmt.Printf("GetUserDiscountTimesMap failed:%v", err)
	//	return
	//}
	//for _, cfg := range tagConfigList {
	//	//if !utils.InInt64Slice(cfg.OpenArea, int64(param.AreaCode)) { //用户区域与配置不符
	//	//	continue
	//	//}
	//	nowUnix := time.Now().Unix()
	//	if nowUnix < cfg.St || nowUnix > cfg.Et { //不在配置时间范围
	//		continue
	//	}
	//	if cfg.BonusMap[param.AreaCode] != nil {
	//		bonusCfg := cfg.BonusMap[param.AreaCode]
	//
	//		//组装此用户仍能生效的三方渠道bonusMap
	//		for coins, bonus := range bonusCfg.ThreeBonusMap {
	//			//if !utils.InInt64Slice(bonusCfg.ThirdPayTag, coins) { //不在开放档位列表跳过
	//			//	continue
	//			//}
	//			//recordKey := fmt.Sprintf("%v_%v", cfg.ConfigName, coins)
	//			//if discountRecordsMap[recordKey] >= cfg.LimitNum { //已达返利生效次数上限
	//			//	continue
	//			//}
	//			if _, ok := userDiscountCfg.ThreeBonusMap[coins]; !ok {
	//				userDiscountCfg.ThreeBonusMap[coins] = &DiscountInfoSt{
	//					ConfigName: cfg.ConfigName,
	//					IntValue:   bonus,
	//					LimitNum:   cfg.LimitNum,
	//				}
	//			}
	//		}
	//		//组装此用户仍能生效的官方渠道bonusMap
	//		for coins, bonus := range bonusCfg.OfficialBonusMap {
	//			//if !utils.InInt64Slice(cfg.OpenPrice, coins) { //不在开放档位列表跳过
	//			//	continue
	//			//}
	//			//recordKey := fmt.Sprintf("%v_%v", cfg.ConfigName, coins)
	//			//if discountRecordsMap[recordKey] >= cfg.LimitNum { //已达返利生效次数上限
	//			//	continue
	//			//}
	//			if _, ok := userDiscountCfg.OfficialBonusMap[coins]; !ok {
	//				userDiscountCfg.OfficialBonusMap[coins] = &DiscountInfoSt{
	//					ConfigName: cfg.ConfigName,
	//					IntValue:   bonus,
	//					LimitNum:   cfg.LimitNum,
	//				}
	//				userDiscountCfg.OfficialRatioMap[coins] = bonusCfg.OfficialRatioMap[coins]
	//
	//
	//			}
	//		}
	//	}
	//}
	//if utils.IsEnvTest() {
	//	logger.Info("mid:%v, open discount cfg:%v", param.MemId, common.Object2Str(userDiscountCfg))
	//}
	return
}
