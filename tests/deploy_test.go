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

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestGet is a sample to run an endpoint test
func TestSendGet(t *testing.T) {
	err := xf.InitServer(defFields.opts)
	if err != nil {
		t.Fatal(err)
	}

	body := `{"txt": "请1号东风到内科门诊1号诊室就诊", "lang":["jyut"]}`
	r, _ := http.NewRequest("POST", "/xftts/make-tts", strings.NewReader(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	logs.Info("testing", "TestGet", "Code", w.Code)
	if w.Code != 200 {
		logs.Info(w.Body)
	}

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})

	err = xf.TTSSrv.Close()
	if err != nil {
		t.Fatal(err)
	}
}
