#!/bin/bash

# shellcheck disable=SC2155
export LD_LIBRARY_PATH=$(pwd)/xf/libs/$(uname -m)

if [ -z "$1" ]; then
    exit 0
fi

case "$1" in
  help)
    echo "xftts 脚本
    - build 编译可执行文件
    - clean 清理缓存和日志文件
    - 'words' [*.wav] 合成语音"
    ;;
  build)
    make build;;
  clean)
    make clean;;
  *)
    make "$1"/"${2:-"speech.wav"}";;
esac
