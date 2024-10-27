package felo

type LocSt struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type BaseParam struct {
	App              string      `json:"h_app,omitempty"` //me-live
	AppVer           string      `json:"h_av"`            //app版本号
	DevType          int32       `json:"h_dt"`            //操作系统 0安卓 1ios
	DevId            string      `json:"h_did"`           //设备号
	NetType          int32       `json:"h_nt"`            //网络类型
	MemId            int64       `json:"h_m,omitempty"`   //mid
	Loc              LocSt       `json:"h_loc,omitempty"` //坐标
	Chan             string      `json:"h_ch,omitempty"`  //下载渠道
	Timestamp        int64       `json:"h_ts"`            //当前时间
	Model            string      `json:"h_model"`         //手机型号
	Token            string      `json:"token"`           //token
	Ver              string      `json:"ver"`             //Deprecated
	ClientIP         string      `json:"client_ip"`
	Language         string      `json:"h_language,omitempty"`        //客户端使用语言 "en": 英语 "ar"：阿拉伯语
	DeviceLanguage   string      `json:"h_device_language,omitempty"` //客户端设备语言 安卓："en": 英语 "ar"：ios: ar-US
	Adjust           string      `json:"h_adjust"`                    //Adjust标识
	Sys              string      `json:"h_sys"`                       //Android上报android id, iOS上报idfv
	Gender           int         `json:"h_gender,omitempty"`          //性别
	Age              int         `json:"h_age,omitempty"`             //年龄
	Ip               string      `json:"h_ip,omitempty"`              //客户端ip
	ZoneName         string      `json:"h_zone_name"`                 // 用户时区的名称
	ZoneAbbreviation string      `json:"h_zone_abbreviation"`         // 用户时区的缩写
	ZoneOffset       int64       `json:"h_zone_offset"`               // 用户时区与零时区的间隔秒数
	RequestId        string      `json:"request_uuid"`                //请求唯一id
	AppsFlyerId      string      `json:"h_appsflyerid"`               //AppsFlyer 标识
	ShuMeiId         string      `json:"h_shumei_id"`                 //数美设备标识
	HeaderRegionCode int64       `json:"h_region_code"`               //国家代码
	AreaCode         int32       `json:"h_area_code"`                 //区域代码
	Carrier          string      `json:"h_carrier"`                   //运营商代码 2.5.0添加 示例：中国联通,cn,46001
	Test             interface{} `json:"h_test"`                      //android int 0,1; ios bool
	OS               interface{} `json:"h_os"`                        //系统版本
	Tz               string      `json:"tz"`                          // op操作系统时区
	Cpu              string      `json:"h_cpu"`                       //系统cpu占用，目前是ios有
}

//
//func verifyApplePayReceiptWithoutOrderid(receiptdata, transactionID string) (sandbox bool, ioscode string, err error) {
//	client := appstore.New()
//	req := appstore.IAPRequest{
//		ReceiptData: receiptdata,
//		Password:    "646957271c814eccbcc9c1c5c200f062",
//	}
//	//resp := &appstore.IAPResponse{}
//	resp := &appstore.IAPResponse{}
//	ctx, cancelFunc := context.WithCancel(context.Background())
//	go func() {
//		time.Sleep(15000 * time.Millisecond)
//		cancelFunc()
//	}()
//	err = client.Verify(ctx, req, resp)
//	if err != nil {
//		fmt.Printf("verify apple pay err1: %s", err)
//		return
//	}
//	err = appstore.HandleError(resp.Status)
//	if err != nil {
//		fmt.Printf("verify apple pay err2: %s", err)
//		return
//	}
//	if resp.Environment == appstore.Sandbox {
//		sandbox = true
//	}
//	bundleID := resp.Receipt.BundleID
//	if bundleID != "com.cynet.xingqi" {
//		err = fmt.Errorf("bundle id mismatch")
//		return
//	}
//
//	var valid bool
//
//	for _, v := range resp.Receipt.InApp {
//		if v.TransactionID == transactionID {
//			valid = true
//			ioscode = v.ProductID
//			fmt.Printf("verifyApplePayReceipt OK: %v, inapp:%+v", transactionID, v)
//		}
//	}
//
//	if !valid {
//		err = fmt.Errorf("transaction mismatch")
//		return
//	}
//	return
//}
