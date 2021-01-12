package test

import (
	"strconv"
	"testing"
	"time"
	"xftts/xf"
)

type fields struct {
	opts *xf.Options
}
type args struct {
	txt     string
	desPath string
}
type bench struct {
	name    string
	fields  fields
	args    args
	wantErr bool
}

var (
	defFields = fields{
		opts: &xf.Options{
			TTSParams: xf.TTSParams{
				EngineType:   "local",
				VoiceName:    "xiaomei",
				TTSResPath:   "fo|res/tts/xiaomei.jet;fo|res/tts/xiaoyan.jet;fo|res/tts/common.jet",
				Speed:        50,
				Volume:       50,
				Pitch:        50,
				Rdn:          2,
				SampleRate:   16000,
				TextEncoding: "UTF8",
			},
			LoginParams: xf.LoginParams{
				Appid:   "5ff5193f",
				WorkDir: "xf",
			},
			OutDir: "out/",
		}}
	defArgs = args{
		txt:     "请1号东风到内科门诊1号诊室就诊",
		desPath: "test",
	}
)

func BenchmarkOnce(b *testing.B) {
	bm := bench{
		name:   b.Name(),
		fields: defFields,
		args:   defArgs,
	}
	err := xf.InitServer(bm.fields.opts)
	if err != nil {
		b.Fatal(err)
	}

	i := 0
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i++
			desPath := bm.args.desPath + strconv.Itoa(i) + ".wav"
			if err := xf.TTSSrv.Once(bm.args.txt, desPath); (err != nil) != bm.wantErr {
				b.Errorf("Once() error = %v, wantErr %v", err, bm.wantErr)
			}
		}
	})

	err = xf.TTSSrv.Close()
	if err != nil {
		b.Fatal(err)
	}
}

func TestOnceN(t *testing.T) {
	type ret struct {
		int
		time.Duration
	}

	n := 1
	bm := bench{
		name:   t.Name(),
		fields: defFields,
		args:   defArgs,
	}
	err := xf.InitServer(bm.fields.opts)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	for i := 1; i <= n; i++ {
		start := time.Now()
		desPath := bm.args.desPath + strconv.Itoa(i) + ".wav"
		if err := xf.TTSSrv.Once(bm.args.txt, desPath); (err != nil) != bm.wantErr {
			t.Errorf("Once() error = %v, wantErr %v", err, bm.wantErr)
		}
		t.Log("\t执行次数：", i, "\t耗时：", time.Since(start))
	}
	elapsed := time.Since(start)
	t.Log("\t执行总次数：", n, "\t总耗时：", elapsed)

	err = xf.TTSSrv.Close()
	if err != nil {
		t.Fatal(err)
	}
}
