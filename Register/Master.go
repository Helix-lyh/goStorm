package Register

import (
	"github.com/gin-gonic/gin"
	"goStorm/Logger"
	"net/http"
	"time"
)

var logger = Logger.GetLogger("gostorm")


func SetupRouter(master *Master) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "gostorm master get a request!")
	})


	r.GET("/heartbeat", func(c *gin.Context) {
		slaveId := c.Query("id")
		now, err := master.SetHbTime(slaveId)
		if err != nil{
			logger.Infof("heartbeat error = %v", err)
			c.JSON(200, gin.H{
				"code": "400001",
				"msg": "操作失败!",
				"data": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"code": "000000",
			"msg": "操作成功!",
			"data": now,
		})
	})

	r.GET("/slaves", func(c *gin.Context) {
		slaveId := c.Query("id")
		logger.Debugf("slaves get id=%s", slaveId)
		if slaveId != ""{
			slave, err := master.GetSlave(slaveId)
			if err!= nil{
				c.JSON(200, gin.H{
					"code": "400001",
					"msg": "操作失败!",
					"data": err.Error(),
				})
				return
			}else {
				c.JSON(200, gin.H{
					"code": "000000",
					"msg": "操作成功",
					"data": slave,
				})
			}
		} else {
			slaves := master.GetAllSlave()
			c.JSON(200, gin.H{
				"code": "000000",
				"msg": "操作成功",
				"data": slaves,
			})
		}
	})

	r.POST("/regist", func(c *gin.Context) {
		var registInfo RegistInfo

		err := c.ShouldBind(&registInfo)
		if err != nil {
			logger.Infof("gostorm_Master_regist解析参数失败! err=%v", err)
			c.JSON(500, gin.H{"code": 999999, "msg": "gostorm_Master_regist解析参数失败!", "data": err})
			return
		}

		slave := SlaveInfo{
			Id:GenSlaveId(&registInfo),
			Host:registInfo.Host,
			LastHeartbeat:time.Now().UnixNano(),
		}
		err = master.AddSlave(slave)
		if err != nil {
			c.JSON(200, gin.H{
				"code": "500001",
				"message": "操作失败!",
				"data": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"code": "000000",
			"message": "操作成功!",
			"data": slave,
		})
	})
	return r
}


