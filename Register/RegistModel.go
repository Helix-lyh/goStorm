package Register

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"net"
	"strconv"
	"sync"
	"time"
)

type SlaveInfo struct {
	Id string `json:"id"`
	Host string `json:"host"`
	LastHeartbeat int64 `json:"lastHeartbeat"`
}


type RegistInfo struct {
	Host string `json:"host"`
	Mac string `json:"mac"`
	Pid int `json:"pid"`
}

type MasterInfo struct {
	Host string `json:"host"`
	Port int `json:"port"`
	Slaves []SlaveInfo `json:"slaves"`
}

type StandardRes struct {
	Code string `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

func GetMac() {
	// 获取本机的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Error : " + err.Error())
	}
	for _, inter := range interfaces {
		mac := inter.HardwareAddr //获取本机MAC地址
		fmt.Println("MAC = ", mac)
	}
}

func GenSlaveId(info *RegistInfo) string {
	infoStr := info.Host + info.Mac + strconv.Itoa(info.Pid)
	data := []byte(infoStr)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}


//互斥锁
var lock sync.Mutex

type Master struct {
	AllSlaves map[string] *SlaveInfo
}

func (master *Master) GetSlave(slaveId string) (*SlaveInfo, error) {
	if slave, ok := master.AllSlaves[slaveId]; ok{
		return slave, nil
	} else {
		return &SlaveInfo{}, errors.New(fmt.Sprintf("slaveId= %s 不存在!", slaveId))
	}
}

func (master *Master) SetHbTime(salveId string) (int64, error) {
	lock.Lock()
	now := time.Now().UnixNano()
	if slaveInfo, ok := master.AllSlaves[salveId]; ok{
		fmt.Printf("befor master= %v slaveInfo=%v \n", master, slaveInfo)
		slaveInfo.LastHeartbeat = now
		fmt.Printf("after master= %v slaveInfo=%v \n", master, slaveInfo)
		lock.Unlock()
		return now, nil
	} else {
		lock.Unlock()
		return now, errors.New(fmt.Sprintf("slaveId=%s 不存在!", salveId))
	}
}

func (master *Master) AddSlave(salve SlaveInfo) error {
	lock.Lock()
	if _, ok := master.AllSlaves[salve.Id]; ok{
		lock.Unlock()
		return errors.New(fmt.Sprintf("slaveId= %s 已经存在!", salve.Id))
	} else {
		master.AllSlaves[salve.Id] = &salve
		lock.Unlock()
		return nil
	}
}
func (master *Master) GetAllSlave() map[string]*SlaveInfo {
	return master.AllSlaves
}

type Slave struct {
	*RegistInfo `json:"registInfo"`
	*SlaveInfo `json:"slaveInfo"`
}

func (slave *Slave) RegToMaster()  {
	regInfo := RegistInfo{
			Host: "127.0.0.9",
			Mac: "AC-9E-17-91-F3-E4",
			Pid: 123456,
	}
	resp, err := grequests.Post("http://127.0.0.1:9526/regist",
		&grequests.RequestOptions{
			JSON: regInfo,
		})

	if err != nil {
		logger.Info("Unable to make request", resp.Error)
		return
	}

	if resp.Ok != true {
		logger.Info("Request did not return OK")
		return
	}
	res := StandardRes{}
	err = resp.JSON(&res)
	if err != nil{
		return
	}
	logger.Infof("RegToMaster res.Data=%v", res.Data)
	slaInfo := res.Data.(map[string]interface{})
	if slaveId, ok := slaInfo["id"]; ok {
		slave.Id = slaveId.(string)
	}
}

func (slave *Slave) SetSlaveInfo(slaInfo SlaveInfo)  {
	slave.Id = slaInfo.Id
	slave.LastHeartbeat = slaInfo.LastHeartbeat
}

func (slave *Slave) SendHb()  {
	slaveId := slave.Id
	resp, err := grequests.Get("http://127.0.0.1:9526/heartbeat?id=" + slaveId, &grequests.RequestOptions{})

	if err != nil {
		logger.Info("Unable to make request", resp.Error)
		return
	}

	if resp.Ok != true {
		logger.Info("Request did not return OK")
		return
	}
	res := StandardRes{
		Data: "",
	}
	resp.JSON(&res)
	logger.Infof("SendHb res.Data=%v", res.Data)
}















