XFTTS := bin/xftts
DEF_SPEECH := "你他妈的给我翻译翻译什么他妈的叫他妈的惊喜"

build: clean xftts clean-cache
xftts:
	mkdir -p bin out
	go build -a -installsuffix cgo -o $(XFTTS) ./

# usage: make say_your_words/speech.wav
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
