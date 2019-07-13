package storm

import (
	"github.com/levigross/grequests"
	"goStorm/Logger"
	"time"
)

var logger = Logger.GetLogger("goStorm")

type blackHoleReq struct {
	ReqCode      string `json:"reqCode" binding:"required"`
	ReqStartTime int64  `json:"reqStartTime" binding:"required"`
}

type blackHoleData struct {
	ReqCode       string `json:"reqCode" binding:"required"`
	ReqStartTime  int64  `json:"reqStartTime" binding:"required"`
	ReqArriveTime int64  `json:"reqArriveTime" binding:"required"`
	ResStartTime  int64  `json:"resStartTime" binding:"required"`
}

type blackHoleRes struct {
	Code string        `json:"code" binding:"required"`
	Msg  string        `json:"msg" binding:"required"`
	Data blackHoleData `json:"data" binding:"required"`
}

func (b *blackHoleReq) Do() {
	//log.Println("num:", b.ReqCode)
	blackHoleTest(b)
}

func blackHoleTest(reqParams *blackHoleReq) {
	resp, err := grequests.Post("http://129.204.212.59:9527/blackhole",
		&grequests.RequestOptions{
			JSON: reqParams,
		})

	if err != nil {
		logger.Info("Unable to make request", resp.Error)
		return
	}

	if resp.Ok != true {
		logger.Info("Request did not return OK")
		return
	}
	res := blackHoleRes{
		Data: blackHoleData{},
	}
	resp.JSON(&res)
	logger.Infof("res.Data=%v", res.Data)
	logger.Infof("请求编号 %s - 请求发起时间 %s - 请求到达时间 %s - 请求返回时间 %s - 请求结束时间 %s",
		res.Data.ReqCode,
		time.Unix(res.Data.ReqStartTime/1e9, res.Data.ReqStartTime % 1e9).Format("2006-01-02 15:04:05.999999"),
		time.Unix(res.Data.ReqArriveTime/1e9, res.Data.ReqArriveTime % 1e9).Format("2006-01-02 15:04:05.999999"),
		time.Unix(res.Data.ResStartTime/1e9, res.Data.ReqArriveTime % 1e9).Format("2006-01-02 15:04:05.999999"),
		time.Now().Format("2006-01-02 15:04:05.999999"),
	)

}
