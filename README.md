# Xf-TextToSpeech


## TODO-list

- [x] 粤语语音生成
- [x] 单个语音生成
- [ ] 多语语音合成
- [ ] 添加语音缓存功能
- [x] 使用 beego 部署 http 服务
- [x] 提供单个语音生成的 http 接口
- [ ] 提供多语语音合成的 http 接口
- [x] docker 快速部署
- [x] 提供性能基准测试

## Usage

### 编译

使用 Make Tool

`make build`

或者

`source xftts.sh build`

### 语音合成

使用 Make Tool

`make say_your_words/xxx.wav`

或者

`./xftts.sh 'say_your_words' [*.wav]`

## API 调用流程

![qtts-api][qtts-api]

## 参数说明

### qtts.h::MSPLogin

```c
int MSPAPI MSPLogin (const char * 	usr,
                     const char * 	pwd,
                     const char * 	params 
                    )
```

#### params

|业务|参数|名称|说明|
|:---:|---|---|---|
|通用|appid         |应用ID    |SDK申请成功后获取到的appid。申请SDK请前往讯飞开放平台，此参数必须传入|
|离线|engine_start  |离线引擎启动  |启动离线引擎，支持参数，ivw：唤醒，asr：识别|
|离线|[xxx]_res_path|离线引擎资源路径|设置ivw、asr引擎离线资源路径，详细格式如下：fo&#124;[path]&#124;[offset]&#124;[length]&#124<br>示例如下,单个资源路径：<br>ivw_res_path=fo&#124;res/ivw/wakeupresource.jet，<br>多个资源路径：asr_res_path=fo&#124;res/asr/common.jet;fo&#124;res/asr/sms.jet|

### msp_cmn.h::QTTSSessionBegin

```c
const char* MSPAPI QTTSSessionBegin(const char * params,
                                    int *        errorCode 
                                    )
```

#### params

|业务|参数|名称|说明|
|:---:|---|---|---|
|通用|engine_type     |引擎类型       |可取值：<br>cloud：在线引擎<br>local：离线引擎，默认为cloud|
|通用|voice_name      |发音人        |不同的发音人代表了不同的音色，<br>如男声、女声、童声等，具体参数值请到[发音人授权管理][voice-man]确认|
|通用|speed           |语速         |合成音频对应的语速，<br>取值范围：[0,100]，数值越大语速越快。<br>默认值：50|
|通用|volume          |音量         |合成音频的音量，<br>取值范围：[0,100]，数值越大音量越大。<br>默认值：50|
|通用|pitch           |语调         |合成音频的音调，<br>取值范围：[0,100]，数值越大音调越高。<br>默认值：50|
|离线|tts_res_path    |合成资源路径     |合成资源所在路径，支持fo 方式参数设置，对应格式如下：<br> fo&#124;[path]&#124;[offset]&#124;[length] <br>（1）若是合并资源，则只需传入一个资源路径，如：fo&#124;combined.jet&#124;0&#124;1024 <br>（2）若是分离资源，则需传两个资源路径，如：fo&#124;common.jet&#124;0&#124;1024;fo&#124; xiaoyan.jet&#124;0&#124;1024|
|通用|rdn             |数字发音       |合成音频数字发音，支持参数，<br>0 数值优先,<br>1 完全数值,<br>2 完全字符串，<br>3 字符串优先，<br>默认值：0|
|离线|rcn             |中文发音       |支持参数：<br>0：表示发音为yao<br>1：表示发音为yi<br>默认值：0|
|通用|text_encoding   |文本编码格式     |合成文本编码格式，支持参数，GB2312，GBK，BIG5，UNICODE，GB18030，UTF8|
|通用|sample_rate     |合成音频采样率    |合成音频采样率，支持参数，16000，8000，默认为16000|
|在线|background_sound|背景音        |合成音频中的背景音，支持参数，<br>0：无背景音乐，<br>1：有背景音乐|
|在线|aue             |音频编码格式和压缩等级|码算法：raw；speex；speex-wb；ico<br>编码等级：raw：不进行解压缩<br>speex系列：0-10；<br>默认为speex-wb;7<br>speex对应sample_rate=8000 <br>speex-wb对应sample_rate=16000<br>ico对应sample_rate=16000|
|在线|ttp             |文本类型       |合成文本类型，支持参数，<br>text: 普通格式文本<br>cssml：cssml 格式文本<br>默认值：text|
|离线|speed_increase  |语速增强       |通过设置此参数控制合成音频语速基数，取值范围，<br>1：正常 2：2 倍语速 4：4 倍语速|
|离线|effect          |合成音效       |合成音频的音效，取值范围，<br>0 无音效，1 忽远忽近，2 回声，3 机器人，4 合唱，5 水下，6 混响，7 阴阳怪气|

[voice-man]: [https://console.xfyun.cn/services/tts]
[qtts-api]: xf/libs/qtts.png