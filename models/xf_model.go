package models

type Resp struct {
	Code int         `json:"code"` //状态码
	Msg  string      `json:"msg"`  //返回消息 500
	Data interface{} `json:"data"` //数据实体
}

type Speech struct {
	Hash    [32]byte `json:"hash,omitempty"`
	Txt     string   `json:"txt"`
	Speaker string   `json:"speaker"`
}
