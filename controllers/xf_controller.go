package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type XfController struct {
	beego.Controller
}

func (c *XfController) Once() {

}

//func (o *ObjectController) Get() {
//	objectId := o.Ctx.Input.Param(":objectId")
//	if objectId != "" {
//		ob, err := models.GetOne(objectId)
//		if err != nil {
//			o.Data["json"] = err.Error()
//		} else {
//			o.Data["json"] = ob
//		}
//	}
//	o.ServeJSON()
//}
//
//func (o *ObjectController) Put() {
//	objectId := o.Ctx.Input.Param(":objectId")
//	var ob models.Object
//	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
//
//	err := models.Update(objectId, ob.Score)
//	if err != nil {
//		o.Data["json"] = err.Error()
//	} else {
//		o.Data["json"] = "update success!"
//	}
//	o.ServeJSON()
//}
