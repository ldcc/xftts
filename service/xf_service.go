package service

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/beego/beego/v2/adapter/logs"
	"io/ioutil"
	"os"
	"os/exec"
	"xftts/models"
	"xftts/xf"
)

const (
	JYUT     = "jyut"
	MANDARIN = "mandarin"

	jyutprefix = "j_"
	mandprefix = "m_"
	mixprefix  = "x_"
	wavsuffix  = ".wav"
	mp3suffix  = ".mp3"
)

type XfService struct{}

func NewXfService() *XfService {
	return new(XfService)
}

/**
 * 语音生成
 */
func (srv *XfService) MakeTTS(req *models.SpeechReq) (buf []byte, err error) {
	if len(req.Lang) == 0 {
		req.Lang = []string{JYUT}
	}

	var (
		mp3file string
		prefixs = make([]string, len(req.Lang))
		md5Sum  = md5.Sum([]byte(req.Txt))
		hexSum  = hex.EncodeToString(md5Sum[:])
	)

	for i, lang := range req.Lang {
		prefixs[i], err = srv.Once(req.Txt, lang, hexSum)
		if err != nil {
			return
		}
		// TODO 缓存
	}

	// 合成语音
	mp3file, err = srv.ConcatTTS(prefixs, hexSum)
	if err != nil {
		err = fmt.Errorf("ffmpeg 合成语音失败，%v", err)
		return nil, err
	}

	buf, err = ioutil.ReadFile(mp3file)
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
		return "", fmt.Errorf("不支持的语种")
	}

	desPath = prefix + hexSum + wavsuffix
	return prefix, xf.TTSSrv.Once(txt, desPath, voiveName)
}

/**
 * 多语种语音拼接
 * ffmpeg -i a.wav -i b.wav -filter_complex [0:0][1:0]concat=n=2:v=0:a=1[out] -map [out] c.mp3
 */
func (srv *XfService) ConcatTTS(prefixs []string, hexSum string) (fn string, err error) {
	var (
		fliter string
		done   = make(chan error)
		cmd    = exec.Command("ffmpeg")
		outDir = xf.TTSSrv.GetOutPutDir()
	)

	for i, px := range prefixs {
		fn = outDir + px + hexSum
		go func(fn string) {
			_err := srv.ConvertMp3(fn)
			if _err != nil {
				logs.Error(_err)
			}
			done <- _err
		}(fn)

		fliter += fmt.Sprintf("[%d:0]", i)
		cmd.Args = append(cmd.Args, "-i", fn+wavsuffix)
	}

	if len(prefixs) == 1 {
		select {
		case err = <-done:
			return
		}
	}

	fn = outDir + mixprefix + hexSum + mp3suffix
	fliter += fmt.Sprintf("concat=n=%d:v=0:a=1[out]", len(prefixs))
	cmd.Args = append(cmd.Args, "-filter_complex", fliter, "-map", "[out]", "-y", fn)
	cmd.Stderr = bytes.NewBuffer(nil)

	err = cmd.Run()
	if err != nil {
		return
	}

	return
}

func (srv *XfService) ConvertMp3(fn string) error {
	cmd := exec.Command("ffmpeg", "-i", fn+wavsuffix, fn+mp3suffix)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("转码 mp3 格式失败，%v", err)
	}
	err = os.Remove(fn + wavsuffix)
	if err != nil {
		return fmt.Errorf("清除原始 wav 文件失败，%v", err)
	}

	return nil
}
