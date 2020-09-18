package main

import (
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	"xftts/server"
)

var usageStr = `
Usage: xftts [options]
讯飞语音参数选项:
    -tp <param>                 TTS合成参数[有默认值]
    -lp <param>                 登录参数[有默认值]
合成服务模式选项:
    -d <dir>                    音频保存的目录
    -b <dir>                    音频备份的目录
    -t <text>                	待合成的文本
    -o <file>               	音频输出路径
日志选项:
    -l <file>                   日志输出路径[默认./xftts.log]
    -ll <level>                 日志输出级别(debug,info,warn,error)
其他:
    -h                          查看帮助
`

const (
	DefTTSResPath = "fo|res/tts/xiaofeng.jet;fo|res/tts/common.jet"
)

func main() {
	var (
		opts     = &server.Options{}
		txt      string
		out      string
		help     bool
		logFile  string
		logLevel string
	)

	// TTSParmas
	flag.StringVar(&opts.TTSParams.Params, "tp", "", "TTS合成参数")
	flag.StringVar(&opts.EngineType, "engine", "local", "引擎类型")
	flag.StringVar(&opts.VoiceName, "voice", "xiaofeng", "在线引擎发音人")
	flag.StringVar(&opts.TTSResPath, "tts-res", DefTTSResPath, "离线资源所在路径")
	flag.IntVar(&opts.Speed, "speed", 50, "语速")
	flag.IntVar(&opts.Volume, "volume", 50, "音量")
	flag.IntVar(&opts.Pitch, "pitch", 50, "音调")
	flag.IntVar(&opts.Rdn, "rdn", 2, "合成音频数字发音")
	flag.IntVar(&opts.SampleRate, "rate", 16000, "合成音频采样率")
	flag.StringVar(&opts.TextEncoding, "encoding", "UTF8", "文本编码格式")
	flag.StringVar(&opts.Aue, "aue", "", "音频编码格式和压缩等级")
	flag.StringVar(&opts.TTP, "ttp", "", "文本类型")
	flag.StringVar(&opts.SpeedIncrease, "inc", "", "语速增强")

	// LoginParams
	flag.StringVar(&opts.LoginParams.Params, "lp", "", "登录参数")
	flag.StringVar(&opts.Appid, "appid", "5d57f7c2", "XF 提供的 SDK-Appid")
	flag.StringVar(&opts.EngineMode, "offmode", "", "离线引擎启动模式")
	flag.StringVar(&opts.XXXResPath, "xxx_res_path", "", "离线引擎所在路径。")

	// Options
	flag.StringVar(&opts.OutDir, "d", "", "音频输出目录")
	flag.StringVar(&opts.BackupDir, "b", "", "音频保存目录")
	flag.StringVar(&txt, "t", "", "单次合成的文本")
	flag.StringVar(&out, "o", "out/speech.wav", "单次合成的输出路径")
	flag.StringVar(&logFile, "l", "logs/xftts.log", "日志输出路径")
	flag.StringVar(&logLevel, "ll", "debug", "日志输出级别")
	flag.BoolVar(&help, "h", false, "Help")

	flag.Parse()

	if help {
		fmt.Printf("%s\n", usageStr)
		return
	}

	err := configureLog()
	if err != nil {
		logs.Error(fmt.Errorf("日志配置失败:%v", err))
		return
	}

	var srv = server.New(opts)

	logs.Info(fmt.Sprintf("合成文本:%s,输出:%s", txt, out))
	err = srv.Once(txt, out)
	if err != nil {
		logs.Info(fmt.Errorf("%v", err))
	}
}

func configureLog() (err error) {
	//设置logger
	err = logs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/logger.log", "separate":["error", "warning","info"],"level":3}`)
	return
}
