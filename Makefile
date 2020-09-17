build: clean xftts clean-cache
	$(source ./xf/ldpath.sh)

XFTTS := bin/xftts
DEF_SPEECH := "你他妈的给我翻译翻译什么他妈的叫他妈的惊喜"

xftts:
	mkdir -p bin out
	go build -o $(XFTTS) ./
#-gcflags "-g -Wall -I./include" -ldflags "-L./libs/x64 -lmsc -lrt -ldl -lpthread"

# usage: make 你给我翻译翻译，什么他妈的叫他妈的惊喜/speech.wav
%.wav: $(XFTTS)
	@rm -f out/$(@F)
	@if [ $(*D) = "." ]; then\
		./bin/xftts -t $(DEF_SPEECH) -o out/$(@F);\
	else\
		./bin/xftts -t $(@D) -o out/$(@F);\
	fi

clean: clean-cache
	@rm -f bin/*

clean-cache:
	@rm -f msc/*.log msc/*.logcache
	find msc/ -type d -empty -delete

ifeq ($(OS),Windows_NT)
  # on Windows
else
  # on Unix/Linux
endif

.PHONY: clean clean-cache
