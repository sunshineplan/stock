package capitalflows

import (
	"encoding/json"

	"github.com/sunshineplan/gohttp"
)

const api = "http://push2.eastmoney.com/api/qt/clist/get?pn=1&pz=500&fields=f14%2Cf62&fs=m%3A90%2Bt%3A2"

// CapitalFlows represents capital flows of all stock sectors.
type CapitalFlows struct {
	AFSB int64 `json:"安防设备" bson:"安防设备"`
	BLTC int64 `json:"玻璃陶瓷" bson:"玻璃陶瓷"`
	BX   int64 `json:"保险" bson:"保险"`
	BZCL int64 `json:"包装材料" bson:"包装材料"`
	CBZZ int64 `json:"船舶制造" bson:"船舶制造"`
	CLHY int64 `json:"材料行业" bson:"材料行业"`
	DLHY int64 `json:"电力行业" bson:"电力行业"`
	DXYY int64 `json:"电信运营" bson:"电信运营"`
	DYJR int64 `json:"多元金融" bson:"多元金融"`
	DZXX int64 `json:"电子信息" bson:"电子信息"`
	DZYJ int64 `json:"电子元件" bson:"电子元件"`
	FDC  int64 `json:"房地产" bson:"房地产"`
	FZFZ int64 `json:"纺织服装" bson:"纺织服装"`
	GCJS int64 `json:"工程建设" bson:"工程建设"`
	GJMY int64 `json:"国际贸易" bson:"国际贸易"`
	GJS  int64 `json:"贵金属" bson:"贵金属"`
	GKSY int64 `json:"港口水运" bson:"港口水运"`
	GSGL int64 `json:"高速公路" bson:"高速公路"`
	GTHY int64 `json:"钢铁行业" bson:"钢铁行业"`
	GYSP int64 `json:"工艺商品" bson:"工艺商品"`
	GYSY int64 `json:"公用事业" bson:"公用事业"`
	HBGC int64 `json:"环保工程" bson:"环保工程"`
	HFHY int64 `json:"化肥行业" bson:"化肥行业"`
	HGHY int64 `json:"化工行业" bson:"化工行业"`
	HQHY int64 `json:"化纤行业" bson:"化纤行业"`
	HTHK int64 `json:"航天航空" bson:"航天航空"`
	JDHY int64 `json:"家电行业" bson:"家电行业"`
	JSZP int64 `json:"金属制品" bson:"金属制品"`
	JXHY int64 `json:"机械行业" bson:"机械行业"`
	JYSB int64 `json:"交运设备" bson:"交运设备"`
	JYWL int64 `json:"交运物流" bson:"交运物流"`
	LYJD int64 `json:"旅游酒店" bson:"旅游酒店"`
	MHJC int64 `json:"民航机场" bson:"民航机场"`
	MTCX int64 `json:"煤炭采选" bson:"煤炭采选"`
	MYJJ int64 `json:"木业家具" bson:"木业家具"`
	NJHY int64 `json:"酿酒行业" bson:"酿酒行业"`
	NMSY int64 `json:"农牧饲渔" bson:"农牧饲渔"`
	NYSY int64 `json:"农药兽药" bson:"农药兽药"`
	QCHY int64 `json:"汽车行业" bson:"汽车行业"`
	QSXT int64 `json:"券商信托" bson:"券商信托"`
	RJFW int64 `json:"软件服务" bson:"软件服务"`
	SJZP int64 `json:"塑胶制品" bson:"塑胶制品"`
	SNJC int64 `json:"水泥建材" bson:"水泥建材"`
	SPDQ int64 `json:"输配电气" bson:"输配电气"`
	SPYL int64 `json:"食品饮料" bson:"食品饮料"`
	SYBH int64 `json:"商业百货" bson:"商业百货"`
	SYHY int64 `json:"石油行业" bson:"石油行业"`
	TXHY int64 `json:"通讯行业" bson:"通讯行业"`
	WHCM int64 `json:"文化传媒" bson:"文化传媒"`
	WJXX int64 `json:"文教休闲" bson:"文教休闲"`
	YH   int64 `json:"银行" bson:"银行"`
	YLGC int64 `json:"园林工程" bson:"园林工程"`
	YLHY int64 `json:"医疗行业" bson:"医疗行业"`
	YQYB int64 `json:"仪器仪表" bson:"仪器仪表"`
	YSJS int64 `json:"有色金属" bson:"有色金属"`
	YYZZ int64 `json:"医药制造" bson:"医药制造"`
	ZBSS int64 `json:"珠宝首饰" bson:"珠宝首饰"`
	ZHHY int64 `json:"综合行业" bson:"综合行业"`
	ZXZS int64 `json:"装修装饰" bson:"装修装饰"`
	ZYSB int64 `json:"专用设备" bson:"专用设备"`
	ZZYS int64 `json:"造纸印刷" bson:"造纸印刷"`
}

// Fetch fetchs capital flows.
func Fetch() (cf CapitalFlows, err error) {
	var res struct {
		Data struct {
			Diff map[string]struct {
				F14 string
				F62 float64
			}
			Total int
		}
	}
	if err = gohttp.Get(api, nil).JSON(&res); err != nil {
		return
	}

	data := make(map[string]int64)
	for _, v := range res.Data.Diff {
		data[v.F14] = int64(v.F62)
	}

	var jsonbody []byte
	jsonbody, err = json.Marshal(data)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonbody, &cf)

	return
}
