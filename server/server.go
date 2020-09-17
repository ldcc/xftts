package server

import (
	"github.com/imroc/log"
	"xftts/xf"
)

type Server struct {
	opts *Options
}

type Options struct {
	OutDir      string //音频输出目录
	BackupDir   string // 备份目录，文件输出成功之后，再将文件复制到备份目录
	Level       int    //音频生成速度级别，越快越耗CPU，级别1~10,数字越小速度越快
	TTSParams   string
	LoginParams string
	RedisAddr   string
	RedisPass   string
	Speed       int
}

//"engine_type = local, voice_name = xiaofeng, text_encoding = UTF8, tts_res_path = fo|res/tts/xiaofeng.jet;fo|res/tts/common.jet, sample_rate = 16000, speed = 50, volume = 50, pitch = 50, rdn = 2"
type TTSParams struct {
	EngineType    string   // 引擎类型。cloud：在线引擎，local：离线引擎
	VoiceName     string   // 使用在线引擎时指定的发音人，可在`控制台 -> 发音人授权管理`确认
	TTSResPath    []string // 合成资源所在路径。 fo|[file_info]|[offset]|[length]
	Speed         int      // 语速，取值范围：[0, 100]，默认 50
	Volume        int      // 音量，取值范围：[0, 100]，默认 50
	Pitch         int      // 音调，取值范围：[0, 100]，默认 50
	Rdn           int      // 合成音频数字发音
	Rcn           bool     // 中文发音
	TextEncoding  string   // 文本编码格式
	SampleRate    int      // 合成音频采样率
	BgSound       bool     // 背景音
	Aue           string   // 音频编码格式和压缩等级
	TTP           string   // 文本类型
	SpeedIncrease string   // 语速增强
	Effect        int      // 合成音效
}

type Speech struct {
	Id  string `json:"id"`
	Txt string `json:"txt"`
}

func New(opts *Options) *Server {
	// TODO set default
	return &Server{
		opts: opts,
	}
}

func (s *Server) Once(txt string, desPath string) error {
	log.Debug("tts:%s,login:%s", s.opts.TTSParams, s.opts.LoginParams)
	xf.SetTTSParams(s.opts.TTSParams)
	err := xf.Login(s.opts.LoginParams)

	if err != nil {
		return err
	}
	log.Debug("txt:%s,des_path:%s", txt, desPath)
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
