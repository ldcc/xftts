package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/core/logs"

	beego "github.com/beego/beego/v2/server/web"
	"xftts/models"
	"xftts/service"
)

type XfController struct {
	beego.Controller
}

func (c *XfController) MakeTTS() {
	var (
		req  models.SpeechReq
		resp models.Resp
		buf  []byte
		err  error
	)
	defer func() {
		if err != nil {
			resp.Code = 500
			resp.Msg = err.Error()
			c.Data["json"] = resp
			c.Ctx.ResponseWriter.WriteHeader(500)
			_ = c.ServeJSON()
			logs.Error(err)
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

	buf, err = service.NewXfService().MakeTTS(&req)
	if err != nil {
		err = fmt.Errorf("语音生成失败，%v", err)
		return
	}

	_, err = c.Ctx.ResponseWriter.Write(buf)
	if err != nil {
		err = fmt.Errorf("发送失败，%v", err)
		return
	}
}
