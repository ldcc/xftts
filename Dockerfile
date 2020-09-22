# 镜像文件
FROM golang:latest
# 维护者
MAINTAINER "412657308@qq.com"
LABEL "describe"="科大讯飞离线tts"

#创建工程文件夹
RUN mkdir -p /app && \
    mkdir -p /app/logs && \
    mkdir -p /app/out && \
    mkdir -p /app/xf/libs/x64

# 拷贝当前目录代码到镜像
COPY ./msc/res /app/msc/res
COPY ./xftts /app/xftts
COPY ./xf/include /app/xf/include
COPY ./xf/libs/x86_64 /app/xf/libs/x86_64

#设置时区
COPY ./Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone
#环境变量
ENV LD_LIBRARY_PATH=/app/xf/libs/x86_64
ENV PATH /app/xftts:$PATH

WORKDIR /app
