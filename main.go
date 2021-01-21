package main

import (
	"flag"
	"fmt"
	"github.com/beego/beego/v2/server/web"
	"time"
	"xftts/cache"

	"github.com/beego/beego/v2/adapter/logs"
	_ "xftts/routers"
	"xftts/xf"
)

var usageStr = `
Usage: xftts [options]
语音合成参数选项:
    -tp <param>                 TTS合成参数[有默认值]
	-engine						引擎类型
	-voice						发音人
	-tts-res					离线资源所在路径
	-speed						语速
	-volume						音量
	-pitch						音调
	-rdn						合成音频数字发音
	-rate						合成音频采样率
	-encoding					文本编码格式
	-aue						音频编码格式和压缩等级
	-ttp						文本类型
	-inc						语速增强
讯飞 SDK 登录参数选项：
    -lp <param>                 登录参数[有默认值]
	-appid						登录参数
	-offmode					XF 提供的 SDK-Appid
	-xxx_res_path				离线引擎启动模式
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

var (
	opts     = &xf.Options{}
	txt      string
	out      string
	help     bool
	logFile  string
	logLevel string
)

func init() {
	// TTSParmas
	flag.StringVar(&opts.TTSParams.Params, "tp", "", "TTS合成参数")
	flag.StringVar(&opts.EngineType, "engine", "local", "引擎类型")
	flag.StringVar(&opts.VoiceName, "voice", "xiaomei", "发音人")
	flag.StringVar(&opts.TTSResPath, "tts-res", "", "离线资源所在路径")
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
	flag.StringVar(&opts.Appid, "appid", "5ff5193f", "XF 提供的 SDK-Appid")
	flag.StringVar(&opts.EngineMode, "offmode", "", "离线引擎启动模式")
	flag.StringVar(&opts.WorkDir, "work_dir", "xf", "msc 工作目录")
	flag.StringVar(&opts.XXXResPath, "xxx_res_path", "", "离线引擎所在路径。")

	// Options
	flag.StringVar(&opts.OutDir, "d", "out/", "音频输出目录")
	flag.StringVar(&opts.BackupDir, "b", "", "音频备份目录")
	flag.IntVar(&opts.Level, "level", 1, "音频生成速度级别，级别1~10，数字越小速度越快")

	//
	flag.StringVar(&txt, "t", "", "单次合成的文本")
	flag.StringVar(&out, "o", "speech.wav", "单次合成的输出路径")
	flag.StringVar(&logFile, "l", "logs/xftts.log", "日志输出路径")
	flag.StringVar(&logLevel, "ll", "debug", "日志输出级别")
	flag.BoolVar(&help, "h", false, "Help")

	flag.Parse()

	// 设置 logger
	err := logs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/logger.log", "separate":["error", "warning","info"],"level":3}`)
	if err != nil {
		logs.Error("日志配置失败", err)
	}
}

func main() {
	if help {
		fmt.Printf("%s\n", usageStr)
		return
	}

	err := xf.InitServer(opts)
	if err != nil {
		logs.Error(err)
		return
	}
	defer func() {
		err = xf.TTSSrv.Close()
		if err != nil {
			logs.Error(err)
		}
	}()

	if txt != "" && out != "" {
		logs.Info(fmt.Sprintf("合成文本:%s,输出:%s", txt, out))
		err = xf.TTSSrv.Once(txt, out)
		if err != nil {
			logs.Error(err)
		}
		return
	}

	timeout := time.Duration(web.AppConfig.DefaultInt("register_ip", 2)) * time.Hour
	cache.XfDump = cache.NewDump(timeout, func(key string) {
		err := xf.TTSSrv.RemoveFile(key)
		if err != nil {
			logs.Error(fmt.Errorf("删除缓存失败，%v", err))
		}
	})

	web.Run()
}
