package service

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os/exec"

	"xftts/cache"
	"xftts/models"
	"xftts/xf"
)

const (
	JYUT     = "jyut"
	MANDARIN = "mandarin"

	jyutprefix = "j_"
	mandprefix = "m_"

	wavsuffix = ".wav"
	mp3suffix = ".mp3"
)

type XfService struct {
	dump cache.DumpFile
}

func NewXfService() *XfService {
	srv := new(XfService)
	srv.dump = cache.XfDump
	return srv
}

/**
 * 语音生成
 */
func (srv *XfService) MakeTTS(req *models.SpeechReq) (buf []byte, err error) {
	if len(req.Lang) == 0 {
		req.Lang = []string{JYUT}
	}

	var (
		mixfile string
		prefixs = make([]string, len(req.Lang))
		md5Sum  = md5.Sum([]byte(req.Txt))
		hexSum  = hex.EncodeToString(md5Sum[:])
	)

	// 语音生成
	for i, lang := range req.Lang {
		prefixs[i], err = srv.Once(req.Txt, lang, hexSum)
		if err != nil {
			return
		}
	}

	// 语音合成
	mixfile, err = srv.ConcatTTS(prefixs, hexSum)
	if err != nil {
		err = fmt.Errorf("ffmpeg 合成语音失败，%v", err)
		return
	}

	buf, err = ioutil.ReadFile(mixfile)
	if err != nil {
		err = fmt.Errorf("获取文件失败，%v", err)
		return
	}

	if len(buf) == 0 {
		err = fmt.Errorf("资源长度为空，%v", err)
		return
	}

	return buf, err
}

/**
 * 单次语音生成
 * 返回当此生成的语音类型
 * 语音类型用文件名前缀 `prefix` 作为标识
 */
func (srv *XfService) Once(txt, lang, hexSum string) (prefix string, err error) {
	var (
		desPath   string
		voiveName string
	)

	switch lang {
	case MANDARIN:
		prefix = mandprefix
		voiveName = "xiaoyan"
	case JYUT:
		prefix = jyutprefix
		voiveName = "xiaomei"
	default:
		err = fmt.Errorf("不支持的语种")
		return
	}

	// 检查 mp3 缓存
	desPath = prefix + hexSum
	if srv.dump.Lookup(desPath) == nil {
		// 生成 tts
		err = xf.TTSSrv.Once(txt, desPath, voiveName)
		if err != nil {
			return
		}

		// 转码 mp3
		err = srv.ConvertMp3(desPath)
		if err != nil {
			return
		}

		// 缓存 mp3
		srv.dump.Extend(desPath)
	}

	return
}

/**
 * 多语种语音拼接
 * ffmpeg -i a -i b -filter_complex [0:0][1:0]concat=n=2:v=0:a=1[out] -map [out] c
 */
func (srv *XfService) ConcatTTS(prefixs []string, hexSum string) (mixfile string, err error) {
	var (
		fliter    string
		desPath   string
		mixprefix string

		cmd    = exec.Command("ffmpeg")
		outDir = xf.TTSSrv.GetOutPutDir()
	)

	for i, px := range prefixs {
		mixprefix += px
		fliter += fmt.Sprintf("[%d:0]", i)
		cmd.Args = append(cmd.Args, "-i", outDir+px+hexSum)
	}

	// 检查 mixfile 缓存
	desPath = mixprefix + hexSum
	mixfile = outDir + desPath
	if srv.dump.Lookup(desPath) == nil {
		fliter += fmt.Sprintf("concat=n=%d:v=0:a=1[out]", len(prefixs))
		cmd.Args = append(cmd.Args, "-filter_complex", fliter, "-map", "[out]", "-y", "-f", "mp3", mixfile)
		cmd.Stderr = bytes.NewBuffer(nil)

		err = cmd.Run()
		if err != nil {
			err = fmt.Errorf("拼接语音失败，%v", err)
			return
		}
		// 添加 mix 缓存
		srv.dump.Extend(desPath)
	}

	return
}

/**
 * wav 转码 mp3 格式
 * ffmpeg -i `hex` -y -f mp3 `hex`
 */
func (srv *XfService) ConvertMp3(desPath string) error {
	fn := xf.TTSSrv.GetOutPutDir() + desPath
	cmd := exec.Command("ffmpeg", "-i", fn, "-y", "-f", "mp3", fn)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("转码 mp3 格式失败，%v", err)
	}
	//go func() {
	//	err = os.Remove(fn + wavsuffix)
	//	if err != nil {
	//		logs.Warning(fmt.Errorf("清除原始 wav 文件失败，%v", err))
	//	}
	//}()

	return err
}
