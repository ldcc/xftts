package server

import (
	"strconv"
	"testing"
)

func BenchmarkOnce(b0 *testing.B) {
	type fields struct {
		opts *Options
	}
	type args struct {
		txt     string
		desPath string
	}
	benchmarks := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Performance",
			fields: struct{ opts *Options }{
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
				}},
			args: args{
				txt:     "请1号东风到内科门诊1号诊室就诊",
				desPath: "out/test",
			},
		},
	}

	//ct := make(chan time.Time)
	for _, bm := range benchmarks {
		b0.Run(bm.name, func(b1 *testing.B) {
			//start := time.Now()
			for i := 0; i < b0.N; i++ {
				go func(i int) {
					s := NewServer(bm.fields.opts)
					desPath := bm.args.desPath + strconv.Itoa(i) + ".mp3"
					if err := s.Once(bm.args.txt, desPath); (err != nil) != bm.wantErr {
						b1.Errorf("Once() error = %v, wantErr %v", err, bm.wantErr)
					}
					//ct
					//elapsed := time.Since(start)
				}(i)
			}

		})
	}

	//b0.Log("执行次数：", b0.N, "总耗时：", elapsed)
}
