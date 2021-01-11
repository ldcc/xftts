package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"io/ioutil"
	"xftts/models"
	"xftts/xf"
)

const (
	wavsuffix = ".wav"
)

type XfController struct {
	beego.Controller
}

func (c *XfController) Once() {
	var (
		req  models.Speech
		resp models.Resp
		err  error
	)
	defer func() {
		if err != nil {
			resp.Code = 500
			resp.Msg = err.Error()
			c.Data["json"] = resp
			_ = c.ServeJSON()
		}
	}()

	if len(c.Ctx.Input.RequestBody) == 0 {
		err = errors.New("请求数据为空")
		return
	}

	err = json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		err = fmt.Errorf("参数解析错误，%v", err)
		return
	}

	var (
		desPath string
		hash    [16]byte
		buf     []byte
	)
	hash = md5.Sum([]byte(req.Txt))
	req.Hash = hex.EncodeToString(hash[:])
	desPath = "out/" + req.Hash + wavsuffix

	err = xf.TTSSrv.Once(req.Txt, desPath)
	if err != nil {
		return
	}

	buf, err = ioutil.ReadFile(desPath)
	if err != nil {
		err = fmt.Errorf("获取文件失败，%v", err)
		return
	}

	if len(buf) == 0 || err != nil {
		err = fmt.Errorf("资源长度为空，%v", err)
		return
	}

	_, err = c.Ctx.ResponseWriter.Write(buf)
	if err != nil {
		err = fmt.Errorf("发送失败，%v", err)
		return
	}
}
