package controllers

import (
	"encoding/json"
	"enian_blog/lib/cmn"
	"enian_blog/models"
	"errors"
	"fmt"

	"github.com/astaxie/beego"
)

type BaseApiController struct {
	beego.Controller
	BodyParam cmn.Msi
	UserInfo  models.User
}

func (c *BaseApiController) ApiReturn(code int, msg string, data interface{}) {
	returnData := cmn.Msi{
		"code": code,
		"msg":  msg,
	}
	if data != nil {
		returnData["data"] = data
	}
	c.Data["json"] = &returnData
	c.ServeJSON()
	c.StopRun()
}

// 返回成功
func (c *BaseApiController) ApiSuccess(data interface{}) {
	c.ApiReturn(0, "OK", data)
}

// 返回成功，没有data数据
func (c *BaseApiController) ApiOk() {
	c.ApiReturn(0, "OK", nil)
}

// 返回列表数据
func (c *BaseApiController) ApiListData(list interface{}, total int64) {
	data := cmn.Msi{
		"list":  list,
		"total": total,
	}
	c.ApiReturn(0, "OK", data)
}

// 返回错误
func (c *BaseApiController) ApiError(code int, msg string) {
	c.ApiReturn(code, msg, nil)
}

// 简单提示型错误
func (c *BaseApiController) ApiErrorEasy(msg string) {
	c.ApiReturn(-1, msg, nil)
}

// bodypost转map
func (c *BaseApiController) ParseBodyJsonToMss(params interface{}) error {
	// params := cmn.Msi{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	return err
}
func (c *BaseApiController) ParseBodyJsonToInterface(params interface{}) error {
	// params := cmn.Msi{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	return err
}

// // 获取指定key，如果key 不存在，则返回错误
// func (c *BaseApiController) GetJsonPostByKey(key ...string) cmn.Msi {
// 	for _, v := range key {

// 	}
// 	cParseJsonPost()
// 	params := cmn.Msi{}
// 	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

// 	return params
// }

// bodypost转mss 并验证key是否存在
func (c *BaseApiController) ParseBodyJsonToMssAndKeyExistCheck(keys ...string) (cmn.Mss, error) {
	mapParam := cmn.Mss{}
	if err := c.ParseBodyJsonToMss(&mapParam); err != nil {
		c.ApiError(-2, "参数格式不正确")
		return mapParam, err
	}
	defectField, err := c.MssKeyExistCheck(mapParam, keys...)
	if err != nil {
		c.ApiError(-2, "缺少参数"+defectField)
	}
	return mapParam, nil
}

func (c *BaseApiController) MssKeyExistCheck(mss cmn.Mss, keys ...string) (string, error) {
	for _, v := range keys {
		if _, ok := mss[v]; !ok {
			return v, errors.New("缺少参数:" + v)
		}
	}
	return "", nil
}

// bodypost转mss 并验证key是否存在
func (c *BaseApiController) ParseBodyJsonToMsiAndKeyExistCheck(keys ...string) (cmn.Msi, error) {
	mapParam := cmn.Msi{}
	if err := c.ParseBodyJsonToInterface(&mapParam); err != nil {
		c.ApiError(-2, "格式不正确")
		return mapParam, err
	}
	defectField, err := c.MsiKeyExistCheck(mapParam, keys...)
	if err != nil {
		c.ApiError(-2, "缺少参数"+defectField)
	}
	return mapParam, nil
}

func (c *BaseApiController) MsiKeyExistCheck(msi cmn.Msi, keys ...string) (string, error) {
	for _, v := range keys {
		if _, ok := msi[v]; !ok {
			return v, errors.New("缺少参数:" + v)
		}
	}
	return "", nil
}

func (c *BaseApiController) GetValueByMssKey(mss cmn.Mss, key string) (string, error) {
	if v, ok := mss[key]; ok {
		return v, nil
	}

	return "", errors.New("non-existent")
}

func (c *BaseApiController) GetValueByMsiKeyInterface(msi cmn.Msi, key string) interface{} {
	return msi[key]
}

func (c *BaseApiController) GetValueByMsiKeyString(msi cmn.Msi, key string) (string, error) {
	if v, ok := msi[key]; ok {
		if v1, ok1 := v.(string); ok1 {
			return v1, nil
		}
	}

	return "", errors.New("non-existent")
}

func (c *BaseApiController) GetValueByMsiKeyStringDefault(msi cmn.Msi, key, defaultValue string) string {
	if v, err := c.GetValueByMsiKeyString(msi, key); err == nil {
		return v
	} else {
		return defaultValue
	}
}

func (c *BaseApiController) GetValueByMsiKeyFloat64(msi cmn.Msi, key string) (float64, error) {
	// fmt.Println("类型", reflect.TypeOf(msi[key]), msi[key])
	if v, ok := msi[key]; ok {
		if v1, ok1 := v.(float64); ok1 {
			return v1, nil
		}
	}
	return 0, errors.New("non-existent")
}

func (c *BaseApiController) GetValueByMsiKeyInt(msi cmn.Msi, key string) (int, error) {
	if v, err := c.GetValueByMsiKeyFloat64(msi, key); err == nil {
		return int(v), nil
	}
	return 0, errors.New("non-existent")
}

func (c *BaseApiController) GetValueByMsiKeyIntDefault(msi cmn.Msi, key string, defaultValue int) int {
	if v, err := c.GetValueByMsiKeyFloat64(msi, key); err == nil {
		return int(v)
	}
	return defaultValue
}

// 获取分页信息
func (c *BaseApiController) GetPage() (page int, limit int) {
	var err error
	if page, err = c.GetInt("page"); err != nil {
		page = 0
	}
	if limit, err = c.GetInt("limit"); err != nil {
		limit = 0
	}
	return
}

func (c *BaseApiController) UploadImage() {

}

func (c *BaseApiController) Test() {
	mUser := models.User{}
	_, err := mUser.AddOne(models.User{
		Username:   "li",
		Mail:       "123456@qq.com",
		Password:   "123456",
		Name:       "郝",
		Autograph:  "我是签名",
		Head_image: "",
		Status:     2,
		Role:       1,
	})
	if err != nil {
		fmt.Println("有错误", err)
	}
	c.ApiSuccess(nil)
}
