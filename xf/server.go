package xf

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

// 发音人资源
const (
	TTSResPath = "fo|res/tts/"
	ResSuffix  = ".jet"
	CommonRes  = ";" + TTSResPath + "common" + ResSuffix
)

var (
	TTSSrv *Server
)

type Server struct {
	sync.RWMutex
	opts *Options
}

func InitServer(opts *Options) error {
	if TTSSrv != nil {
		return nil
	}

	TTSSrv = &Server{opts: opts}
	return Login(opts.LoginParams.Format())
}

func (srv *Server) Close() error {
	return Logout()
}

func (srv *Server) Once(txt, desPath string, voiceName ...string) error {
	var paramCopy = srv.opts.TTSParams

	if len(voiceName) > 0 && len(voiceName[0]) > 0 {
		paramCopy.VoiceName = voiceName[0]
		paramCopy.TTSResPath = TTSResPath + voiceName[0] + ResSuffix + CommonRes
	}

	if paramCopy.TTSResPath == "" {
		paramCopy.TTSResPath = TTSResPath + srv.opts.VoiceName + ResSuffix + CommonRes
	}

	desPath = srv.opts.OutDir + desPath
	ttsParams := paramCopy.Format()

	srv.Lock()
	defer srv.Unlock()
	return TextToSpeech(txt, desPath, ttsParams)
}

func (srv *Server) GetOutPutDir() string {
	return srv.opts.OutDir
}

func (srv *Server) RemoveFile(desPath string) error {
	return os.Remove(srv.opts.OutDir + desPath)
}

type Options struct {
	TTSParams
	LoginParams

	OutDir    string // 音频输出目录
	BackupDir string // 备份目录，文件输出成功之后，再将文件复制到备份目录
	Level     int    // 音频生成速度级别，越快越耗CPU，级别1~10，数字越小速度越快
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
	WorkDir    string // msc 工作目录
	EngineMode string // 离线引擎启动模式
	XXXResPath string // fo|[path]|[offset]|[length]
}

func (p *LoginParams) Format() string {
	if p.Params != "" {
		return p.Params
	}

	var params = &p.Params
	appendParam("appid", p.Appid, params)
	appendParam("work_dir", p.WorkDir, params)
	appendParam("engine_start", p.EngineMode, params)
	appendParam(p.EngineMode+"_res_path", p.XXXResPath, params)
	return *params
}

func appendParam(field, value string, src *string) {
	if value == "" {
		return
	}

	if *src == "" {
		*src = fmt.Sprintf("%s=%s", field, value)
	} else {
		*src = fmt.Sprintf("%s=%s,%s", field, value, *src)
	}
}
