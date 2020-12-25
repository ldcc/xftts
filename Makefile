XFTTS := bin/xftts
BENCH_ONCE := bin/once

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

bench-once: clean-once
	mkdir -p bin out
	go test -work -ldflags "-s -w" -a -installsuffix cgo -c -o $(BENCH_ONCE) xftts/server
	#$(BENCH_ONCE) -test.v -test.bench BenchmarkOnce -test.run ^$$
	$(BENCH_ONCE) -test.v -test.run TestOnceN

clean: clean-cache
	@rm -f bin/*

clean-once:
	@rm -f $(BENCH_ONCE)

clean-cache:
	@rm -f msc/*.log msc/*.logcache
	@find msc/ -type d -not \( -regex '.*/res.*' -o -regex '.*/$$' \) | xargs rm -rf

ifeq ($(OS),Windows_NT)
  # on Windows
else
  # on Unix/Linux
endif

.PHONY: clean clean-cache
