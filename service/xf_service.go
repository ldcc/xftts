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
		desPaths = make([]string, len(req.Lang))
		//desPath  string
	)

	for i, lang := range req.Lang {
		desPaths[i], err = srv.Once(req.Txt, lang)
		if err != nil {
			return nil, err
		}
		// TODO 缓存
	}

	// TODO 合并

	buf, err = ioutil.ReadFile(desPaths[0])
	if err != nil {
		err = fmt.Errorf("获取文件失败，%v", err)
		return
	}

	if len(buf) == 0 {
		err = fmt.Errorf("资源长度为空，%v", err)
		return
	}

	return nil, err
}

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
