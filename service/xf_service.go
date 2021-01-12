package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"xftts/models"
	"xftts/xf"
)

const (
	JYUT     = "jyut"
	MANDARIN = "mandarin"

	jyutprefix = "j_"
	mandprefix = "m_"
	mixprefix  = "mx_"
	wavsuffix  = ".wav"
)

type XfService struct{}

func NewXfService() *XfService {
	return new(XfService)
}

func (srv *XfService) MakeTTS(req *models.SpeechReq) (buf []byte, err error) {
	if len(req.Lang) == 0 {
		req.Lang = []string{JYUT}
	}

	var (
		desPath  string
		desPaths = make([]string, len(req.Lang))
	)

	for i, lang := range req.Lang {
		desPaths[i], err = srv.Once(req.Txt, lang)
		if err != nil {
			return
		}
		// TODO 缓存
	}

	// 合成语音
	desPath, err = srv.ConcatTTS(desPaths)
	if err != nil {
		err = fmt.Errorf("ffmpeg 合成语音失败，%v", err)
		return nil, err
	}

	buf, err = ioutil.ReadFile(desPath)
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

// 单次语音生成
func (srv *XfService) Once(txt, lang string) (desPath string, err error) {
	var (
		md5Sum    [16]byte
		hexSum    string
		prefix    string
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

	md5Sum = md5.Sum([]byte(txt))
	hexSum = hex.EncodeToString(md5Sum[:])
	desPath = prefix + hexSum + wavsuffix

	return desPath, xf.TTSSrv.Once(txt, desPath, voiveName)
}

// 多语种语音合成
func (srv *XfService) ConcatTTS(desPaths []string) (desPath string, err error) {
	// TODO
	return desPaths[0], err
}
