#!/bin/bash

# shellcheck disable=SC2155
export LD_LIBRARY_PATH=$(pwd)/xf/libs/$(uname -m)

XFTTS=bin/xftts
case "$1" in
  help)
    echo "xftts 脚本
    xftts build - 编译可执行文件
    xftts clean - 清理缓存和日志文件
    xftts *.wav ['say_your_words'] - 合成语音（缺省语音）
    xftts 'say_your_words' [*.wav] - 合成语音（缺省文件名）"
    ;;
  build)
    rm -f $XFTTS
    go build -ldflags "-s -w" -a -installsuffix cgo -o $XFTTS ./
    ;;
  clean)
    rm -f logs/*.log msc/*.log msc/*.logcache
    find msc/ -type d -not \( -regex '.*/res.*' -o -regex '.*/$$' \) -print0 | xargs -0 rm -rf
    ;;
  *.wav)
    DEF_SPEECH="happy 30"
    rm -f out/"$1"
    if [ -n "$2" ]; then
      ./$XFTTS -t "$2" -o out/"$1";\
      exit 0
    fi
    ./$XFTTS -t "$DEF_SPEECH" -o out/"$1";\
    ;;
  *)
    if [[ "$2" = *'.wav' ]]; then
      rm -f out/"$2"
      ./$XFTTS -t "$1" -o out/"$2"
      exit 0
    fi
    ./$XFTTS -t "$1" -o out/speech.wav;\
    ;;
esac