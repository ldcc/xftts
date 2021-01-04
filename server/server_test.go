package server

import (
	"strconv"
	"testing"
	"time"
)

type fields struct {
	opts *Options
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
		opts: &Options{
			TTSParams: TTSParams{
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
			LoginParams: LoginParams{
				Appid: "5d57f7c2",
			},
		}}
	defArgs = args{
		txt:     "请1号东风到内科门诊1号诊室就诊",
		desPath: "out/test",
	}
)

func BenchmarkOnce(b *testing.B) {
	bm := bench{
		name:   b.Name(),
		fields: defFields,
		args:   defArgs,
	}
	srv, err := NewServer(bm.fields.opts)
	if err != nil {
		b.Fatal(err)
	}

	i := 0
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i++
			desPath := bm.args.desPath + strconv.Itoa(i) + ".mp3"
			if err := srv.Once(bm.args.txt, desPath); (err != nil) != bm.wantErr {
				b.Errorf("Once() error = %v, wantErr %v", err, bm.wantErr)
			}
		}
	})

	err = srv.Close()
	if err != nil {
		b.Fatal(err)
	}
}

func TestOnceN(t *testing.T) {
	type ret struct {
		int
		time.Duration
	}

	n := 100
	bm := bench{
		name:   t.Name(),
		fields: defFields,
		args:   defArgs,
	}
	srv, err := NewServer(bm.fields.opts)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	for i := 1; i <= n; i++ {
		start := time.Now()
		desPath := bm.args.desPath + strconv.Itoa(i) + ".mp3"
		if err := srv.Once(bm.args.txt, desPath); (err != nil) != bm.wantErr {
			t.Errorf("Once() error = %v, wantErr %v", err, bm.wantErr)
		}
		t.Log("执行次数：", i, "\t耗时：", time.Since(start))
	}
	elapsed := time.Since(start)
	t.Log("执行总次数：", n, "总耗时：", elapsed)

	err = srv.Close()
	if err != nil {
		t.Fatal(err)
	}
}
