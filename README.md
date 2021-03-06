# Xf-TextToSpeech

## TODO-list

- [x] 粤语语音生成
- [x] 单个语音生成
- [x] 多语语音合成
- [x] 添加语音缓存功能
- [x] 使用 beego 部署 http 服务
- [x] 提供单个语音生成的 http 接口
- [x] 提供多语语音合成的 http 接口
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

## wave 文件编码格式

该规范用于对接讯飞接口输出的 pcm 数据流，可用于编写音频文件的生成：[xf/tts.h][tts.h]

![wave-header][wave-header]

### RIFF Chunk

|标签名|大小|取值|描述|
|:---|:---:|:---|---|
|ChunkID|4|"RIFF"|用作 RIFF 资源交换文件标识|
|ChunkSize|4|0~2^32|描述了除 ChunkID、ChunkSize 外的 Chunk 总字节数，为文件总大小 `-8` 字节|
|Format|4|"WAVE"|wave 文件标识，当为 `"WAVE"` 时还至少需要两个 Sub Chunks：__Format Chunk__、__Data Chunk__|

### Format Chunk

|标签名|大小|值|描述|
|:---|:---:|:---|---|
|SubChunk1ID|4|"fmt "|最后一位为空白字符，音频编码格式 Chunk 的标识符|
|SubChunk1Size|4|0 ~ 2^32|Format Chunk 的大小，一般为 16|
|AudioFormat|2|0 or 1|表示音频数据的格式，值 1 表示数据为线性 PCM 编码|
|NumChannels|2|1 or 2|表示音频声道数，值 1 为单声道， 2 是双声道|
|SampleRate|4|-|采样率，每秒从连续信号中提取并组成离散信号的采样个数，用赫兹（Hz）表示|
|ByteRate|4|-|比特率，波形文件每秒的字节数，`ByteRate = SampleRate * NumChannels * BitsPerSample / 8`|
|BlockAlign|2|-| 数据块对齐单位，单次采样的字节大小，`BlockAlign = NumChannels * BitsPerSample / 8`|
|BitsPerSample|2|8*2^n|采样位数或采样深度，一般使用 16|

### Data Chunk

|标签名|大小|值|描述|
|:---|:---:|:---|---|
|SubChunk2ID|4|"data"|音频数据 Chunk 的标识符|
|SubChunk2Size|4|0~2^32|Data Chunk 的大小，PCM 音频数据域的长度|
|data|-|-|PCM 音频数据流|

生成 wave 头文件后拼接 xf 的 pcm 数据流即可生成 .wav 文件了。

## mepg layer3 文件编码格式

MP3 文件分为 3 部分：TAG_V2(ID3V2)，Frames，TAG_V1(ID3V1)：

|Layer|描述|
|:---|:---|
|ID3V2|包含作曲，专辑等信息，不定长，ID3V1的扩展|
|Frames|音频的序列帧，不定长|
|ID3V1|包含作者，作曲，专辑等信息，固定 128 字节|

### ID3V2

ID3V2 由 10 字节的标签头、若干标签帧以及一些其它可选项（以下按拼接顺序排列）组成：

|标签|描述|
|:---|:---|
|Header|标签头，固定 10 字节|
|Extended Header|拓展头，不定长，可选项，默认不添加|
|Frames|标签帧，不定长|
|Padding|空白填充，不定长，可选项，默认不添加|
|Footer|脚注，固定 10 字节，可选项，默认不添加|

标签头由以下部分组成：

|Tag|大小|值|描述|
|:---|:---|---|---|
|ID|3|"ID3"|ID3v2 的标识符|
|Version|2|0x04 00|前一个字节为 ID3v2的主版本，后一个字节为副版本|
|Flags|1|0b11110000|目前只有 4 个标识位，从左到右分别是：帧不同步、扩展头、实验版标识、脚注|
|Size|4|-|标签大小，包括标签头的 10 个字节和所有的标签帧的大小|

标签帧的结构：

![chapter-frame][chapter-frame]


## 资料来源

- [1] <http://www.topherlee.com/software/pcm-tut-wavformat.html>
- [2] <https://mutagen-specs.readthedocs.io/en/latest/id3/id3v2.4.0-structure.html>
- [3] <https://id3.org/id3v2-chapters-1.0>

[voice-man]: [https://console.xfyun.cn/services/tts]
[qtts-api]: xf/libs/qtts.png
[wave-header]: xf/libs/wave-header.png
[chapter-frame]: xf/libs/chapter-frame.png
[tts.h]: xf/tts.h
