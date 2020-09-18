package server

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"strconv"
	"xftts/xf"
)

type Server struct {
	opts *Options
}

type Options struct {
	TTSParams
	LoginParams

	OutDir    string // 音频输出目录
	BackupDir string // 备份目录，文件输出成功之后，再将文件复制到备份目录
	Level     int    // 音频生成速度级别，越快越耗CPU，级别1~10,数字越小速度越快
	RedisAddr string
	RedisPass string
}

type TTSParams struct {
	Params        string // TTS 合成参数，如果该值被指定，则忽略所有其它字段
	EngineType    string // 引擎类型，cloud：在线引擎，local：离线引擎
	VoiceName     string // 发音人，可在`控制台 -> 发音人授权管理`确认
	TTSResPath    string // 离线资源所在路径，fo|[path]|[offset]|[length]
	Speed         int    // 语速，取值范围：[0, 100]，默认 50
	Volume        int    // 音量，取值范围：[0, 100]，默认 50
	Pitch         int    // 音调，取值范围：[0, 100]，默认 50
	Rdn           int    // 合成音频数字发音
	SampleRate    int    // 合成音频采样率
	TextEncoding  string // 文本编码格式
	Aue           string // 音频编码格式和压缩等级
	TTP           string // 文本类型
	SpeedIncrease string // 语速增强
}

func (p *TTSParams) Format() string {
	if p.Params != "" {
		return p.Params
	}

	var params = &p.Params
	appendParam("engine_type", p.EngineType, params)
	appendParam("voice_name", p.VoiceName, params)
	appendParam("tts_res_path", p.TTSResPath, params)
	appendParam("speed", strconv.Itoa(p.Speed), params)
	appendParam("volume", strconv.Itoa(p.Volume), params)
	appendParam("pitch", strconv.Itoa(p.Pitch), params)
	appendParam("rdn", strconv.Itoa(p.Rdn), params)
	appendParam("sample_rate", strconv.Itoa(p.SampleRate), params)
	appendParam("text_encoding", p.TextEncoding, params)
	appendParam("aue", p.Aue, params)
	appendParam("ttp", p.TTP, params)
	appendParam("speed_increase", p.SpeedIncrease, params)

	return *params
}

type LoginParams struct {
	Params     string // 登录参数，如果该值被指定，则忽略所有其它字段
	Appid      string // XF 提供的 SDK-Appid
	EngineMode string // 离线引擎启动模式
	XXXResPath string //   fo|[path]|[offset]|[length]
}

func (p *LoginParams) Format() string {
	if p.Params != "" {
		return p.Params
	}

	var params = &p.Params
	appendParam("appid", p.Appid, params)
	appendParam("engine_start", p.EngineMode, params)
	appendParam(p.EngineMode+"_res_path", p.XXXResPath, params)
	return *params
}

type Speech struct {
	Id  string `json:"id"`
	Txt string `json:"txt"`
}

func New(opts *Options) *Server {
	return &Server{
		opts: opts,
	}
}

func (s *Server) Once(txt string, desPath string) error {
	logs.Info(fmt.Sprintf("tts:%s", s.opts.TTSParams.Format()))
	logs.Info(fmt.Sprintf("login:%s", s.opts.LoginParams.Format()))
	xf.SetTTSParams(s.opts.TTSParams.Format())
	err := xf.Login(s.opts.LoginParams.Format())

	if err != nil {
		return err
	}
	logs.Info(fmt.Sprintf("txt:%s,des_path:%s", txt, desPath))
	err = xf.TextToSpeech(txt, desPath)
	if err != nil {
		return err
	}

	err = xf.Logout()
	if err != nil {
		return err
	}
	return nil
}

func appendParam(field, value string, src *string) {
	if value == "" {
		return
	}

	if *src == "" {
		*src = fmt.Sprintf("%s = %s", field, value)
	} else {
		*src = fmt.Sprintf("%s = %s, %s", field, value, *src)
	}
}
