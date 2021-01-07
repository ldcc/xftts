package test

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	_ "xftts/routers"
	"xftts/xf"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	opts = &xf.Options{
		TTSParams: xf.TTSParams{
			EngineType:   "local",
			VoiceName:    "xiaoyan",
			TTSResPath:   "fo|res/tts/xiaoyan.jet;fo|res/tts/common.jet",
			Speed:        50,
			Volume:       50,
			Pitch:        50,
			Rdn:          2,
			SampleRate:   16000,
			TextEncoding: "UTF8",
		},
		LoginParams: xf.LoginParams{
			Appid: "5d57f7c2",
		},
	}
)

func init() {
	_ = xf.InitServer(opts)
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestGet is a sample to run an endpoint test
func TestSendGet(t *testing.T) {
	body := `{"txt": "请1号东风到内科门诊1号诊室就诊"}`
	r, _ := http.NewRequest("POST", "/xftts/Once", strings.NewReader(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	logs.Info("testing", "TestGet", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}
