#!/bin/bash

# shellcheck disable=SC2155
export LD_LIBRARY_PATH=$(pwd)/xf/libs/$(uname -m)

if [ -z "$1" ]; then
    exit 0
fi

case "$1" in
  help)
    echo "xftts 脚本
    xftts build - 编译可执行文件
    xftts clean - 清理缓存和日志文件
    xftts *.wav ['say_your_words'] - 合成语音（缺省语音）
    xftts 'say_your_words' [*.wav] - 合成语音（缺省文件名）"
    ;;
  build)
    make build;;
  clean)
    make clean;;
  *)
    make "$1"/"${2:-"speech.wav"}";;
esac
