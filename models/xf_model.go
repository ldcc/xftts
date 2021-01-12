package models

type Resp struct {
	Code int         `json:"code"` // 状态码
	Msg  string      `json:"msg"`  // 返回消息 500
	Data interface{} `json:"data"` // 数据实体
}

type SpeechReq struct {
	Txt  string   `json:"txt"`  // 语音合成文本
	Lang []string `json:"lang"` // "jyut" -> 粤语；"mandarin" -> 国语
	Hash [16]byte `json:"-"`
}
