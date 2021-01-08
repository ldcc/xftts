#-------打包环境--------
FROM golang:latest as builder

RUN mkdir -p /go/src/TextToSpeech
WORKDIR /go/src/TextToSpeech

COPY . /go/src/TextToSpeech

ENV LD_LIBRARY_PATH=/go/src/TextToSpeech/xf/libs/x64
ENV GOPROXY=https://goproxy.io,direct
ENV GOSUMDB=off
ENV GOOS=linux
ENV GOARCH=amd64

RUN make build

#-------运行环境--------
# 镜像文件
FROM golang:latest
# 维护者
MAINTAINER "412657308@qq.com"
LABEL "describe"="科大讯飞离线tts"

#创建工程文件夹
RUN mkdir -p /app && \
    mkdir -p /app/conf && \
    mkdir -p /app/logs && \
    mkdir -p /app/out && \
    mkdir -p /app/xf/libs/x64

# 拷贝当前目录代码到镜像
COPY ./msc/res /app/msc/res
COPY ./xf/include /app/xf/include
COPY ./xf/libs/x64 /app/xf/libs/x64
COPY ./conf/app.conf /app/conf/app.conf
COPY --from=builder /go/src/TextToSpeech/bin/xftts /app/xftts

#设置时区
COPY ./Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone
#环境变量
ENV LD_LIBRARY_PATH=/app/xf/libs/x64
ENV PATH /app/xftts:$PATH

WORKDIR /app

ENTRYPOINT ["./xftts"]