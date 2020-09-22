XFTTS := bin/xftts
DEF_SPEECH := "good morning"

build: clean $(XFTTS) clean-cache
$(XFTTS):
	mkdir -p bin out
	go build -work -ldflags "-s -w" -a -installsuffix cgo -o $(XFTTS) ./

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
	@find msc/ -type d -not \( -regex '.*/res.*' -o -regex '.*/$$' \) | xargs rm -rf

ifeq ($(OS),Windows_NT)
  # on Windows
else
  # on Unix/Linux
endif

.PHONY: clean clean-cache
