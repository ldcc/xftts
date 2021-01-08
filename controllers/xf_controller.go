package controllers

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"golang.org/x/crypto/sha3"
	"io/ioutil"
	"os"
	"xftts/xf"

	beego "github.com/beego/beego/v2/server/web"
	"xftts/models"
)

const (
	wavsuffix = ".wav"
)

type XfController struct {
	beego.Controller
}

func (c *XfController) Once() {
	var (
		req     models.Speech
		resp    models.Resp
		desPath string
		buf     []byte
		err     error
	)
	defer func() {
		if err != nil {
			resp.Code = 500
			resp.Msg = err.Error()
			c.Data["json"] = resp
			_ = c.ServeJSON()
		} else {
			if err = os.Remove(desPath); err != nil {
				logs.Error("删除文件失败:", err)
			}
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
	req.Hash = sha3.Sum256([]byte(req.Txt))
	desPath = "out/" + hex.EncodeToString(req.Hash[:]) + wavsuffix

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
