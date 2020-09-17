package main

import (
	"flag"
	"fmt"
	"github.com/imroc/log"
	"os"
	"strings"
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
    -s <digit>                  合成速度级别(1-10),数值越小速度越快，越耗CPU[默认为1]
    -r <addr>                   redis连接地址
    -rp <pass>                  redis密码
合成选项:
    -t <text>                	待合成的文本
    -o <file>               	音频输出路径
日志选项:
    -l <file>                   日志输出路径[默认./xftts.log]
    -ll <level>                 日志输出级别(debug,info,warn,error)
其他:
    -h                          查看帮助 
`

func main() {
	var (
		opts     = &server.Options{}
		txt      string
		out      string
		help     bool
		logFile  string
		logLevel string
	)

	// TODO slim TTSParmas
	flag.StringVar(&opts.TTSParams, "tp", "engine_type = local, voice_name = xiaofeng, text_encoding = UTF8, tts_res_path = fo|res/tts/xiaofeng.jet;fo|res/tts/common.jet, sample_rate = 16000, speed = 50, volume = 50, pitch = 50, rdn = 2", "TTS合成参数")
	flag.StringVar(&opts.LoginParams, "lp", "appid = 5d57f7c2, work_dir = .", "登录参数")
	flag.StringVar(&opts.RedisAddr, "r", ":6379", "redis连接地址")
	flag.StringVar(&opts.RedisPass, "rp", "", "redis连接密码")
	flag.StringVar(&opts.OutDir, "d", "", "音频输出目录")
	flag.StringVar(&opts.BackupDir, "b", "", "音频保存目录")
	flag.IntVar(&opts.Speed, "s", 1, "合成速度")

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

	err := configureLog(logFile, logLevel)
	if err != nil {
		log.Debug("日志配置失败:%v", err)
		return
	}

	var srv = server.New(opts)

	log.Debug("合成文本:%q,输出:%s", txt, out)
	err = srv.Once(txt, out)
	if err != nil {
		log.Debug("%v", err)
	}
}

func configureLog(logFile, logLevel string) error {
	var debug bool

	switch strings.ToLower(logLevel) {
	case "debug":
		debug = true
	case "info":
		debug = false
	}

	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	log.SetDebug(debug)
	log.SetOutput(file)
	return nil
}
