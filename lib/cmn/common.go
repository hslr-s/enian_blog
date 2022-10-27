package cmn

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"
)

// =================
// 类型
// =================
type Msi map[string]interface{}
type Mss map[string]string

// =================
// 公共
// =================

// type RunCodeCallBack func(params interface{})

// 变量

// 常量
const (
	TIMEMODE_1 = "2006-01-02 15:04:05"
)

// 运行模式，1.团队版 2.个人版
var RUN_MODE int = 1

func init() {
	// 查询运行模式(暂时取消)
}

// 函数

// 根据运行模式执行代码
// 参数1 团队版需要执行的内容
// 参数2 个人版需要执行的内容

func RunCodeExec(callback ...func()) {
	if RUN_MODE == 1 {
		if len(callback) == 1 {
			callback[0]()
		}
	} else {
		if len(callback) == 2 {
			callback[1]()
		}
	}
}

// 团队版需要执行的代码
func RunCodeExecByTeam(callback func()) {
	if RUN_MODE == 1 {
		callback()
	}
}

// 个人版需要执行的代码
func RunCodeExecByPerson(callback func()) {
	if RUN_MODE == 2 {
		callback()
	}
}

func JsonEncode(c interface{}) string {
	res, _ := json.Marshal(c)
	return string(res)
}

func JsonDecode(jsonStr string, res interface{}) {
	json.Unmarshal([]byte(jsonStr), res)
}

func PasswordEncryption(password string) string {
	return Md5(Md5(Md5(password)))
}

func Md5(str string) string {
	md5Byte := md5.Sum([]byte(str))
	return hex.EncodeToString(md5Byte[:])
}

func InterfaceToString(value interface{}) string {
	val, ok := value.(string)
	if ok {
		return val
	} else {
		return ""
	}
}

// 生成随机字符串
func CreateRandomString(len int) string {
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	return CreateRandomStringByKeyString(str, len)
}

// 生成随机字符串根据key
func CreateRandomStringByKeyString(keystring string, len int) string {
	var container string
	b := bytes.NewBufferString(keystring)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(keystring[randomInt.Int64()])
	}
	return container
}

// 验证邮箱配置是否完整
func CheckMailConfigComplete(mailInfo Msi) error {
	if _, ok := mailInfo["address"]; !ok {
		return errors.New("邮件信息不完整")
	}
	if _, ok := mailInfo["password"]; !ok {
		return errors.New("邮件密码不完整")
	}
	if _, ok := mailInfo["host"]; !ok {
		return errors.New("邮件服务器不完整")
	}
	if _, ok := mailInfo["port"]; !ok {
		return errors.New("邮件端口不完整")
	}
	if _, ok := mailInfo["email_address"]; ok {
		return errors.New("邮件地址不完整")
	}
	return nil
}

func FatalError(msg string) {
	fmt.Printf("\n %c[1;40;31m%s%c[0m\n\n", 0x1B, "ERROR："+msg, 0x1B)
	os.Exit(1)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 时间转相对时间(不包含秒)
// 刚刚 相对时间: 刚刚 (定义为5分钟内)
// 今天 相对时间: 今天 02:02
// 昨天 相对时间: 昨天 10:02
// 今年 相对时间: 02-23 10:02
// 非今年 相对时间: 2021-02-23 10:02
func TimeToRelativeTime(calcTimeTime time.Time) string {
	relativeTime := ""
	nowTime := time.Now()
	sourceTime := calcTimeTime.Format("2006-01-02 15:04:05")
	nowTimeStr := nowTime.Format("20060102")
	calcTime := calcTimeTime.Format("20060102")
	// 最久，今年，前天，昨天，今天 ，刚刚（5分钟内）
	switch {
	case nowTime.Unix()-calcTimeTime.Unix() < 300:
		relativeTime = "刚刚"
	case nowTimeStr[0:8] == calcTime[0:8]:
		relativeTime = "今天 " + sourceTime[11:16]
	case nowTime.AddDate(0, 0, -1).Format("20060102") == calcTime:
		relativeTime = "昨天 " + sourceTime[11:16]
	case nowTimeStr[0:4] != calcTime[0:4]:
		relativeTime = sourceTime[:16]
	case nowTimeStr[0:4] == calcTime[0:4]:
		relativeTime = sourceTime[5:16]
	}
	return relativeTime
}

func TimeStrToRelativeTime(calcTime string) string {
	local, _ := time.LoadLocation("Asia/Shanghai")
	showTime, _ := time.ParseInLocation("2006-01-02 15:04:05", calcTime, local)
	return TimeToRelativeTime(showTime)
}
