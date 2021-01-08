XFTTS := bin/xftts
BENCH_ONCE := bin/once
SEND_ONCE := bin/send

DEF_SPEECH := "good morning"

build: clean $(XFTTS) clean-cache
$(XFTTS):
	mkdir -p bin out
	go build -work -ldflags "-s -w" -a -installsuffix cgo -o $(XFTTS) ./

# usage: make say_your_words/speech.wav
%.wav: $(XFTTS)
	@rm -f out/$(@F)
	@if [ $(*D) = "." ]; then\
		$(XFTTS) -t $(DEF_SPEECH) -o out/$(@F);\
	else\
		$(XFTTS) -t $(@D) -o out/$(@F);\
	fi

serve: $(XFTTS)
	$(XFTTS)

$(SEND_ONCE):
	go test -work -ldflags "-s -w" -a -installsuffix cgo -c -o $(SEND_ONCE) xftts/tests

sent-once: $(SEND_ONCE)
	$(SEND_ONCE) -test.v -test.run TestSendGet

$(BENCH_ONCE):
	go test -work -ldflags "-s -w" -a -installsuffix cgo -c -o $(BENCH_ONCE) xftts/tests

bench-test: $(BENCH_ONCE)
	mkdir -p bin out
	$(BENCH_ONCE) -test.v -test.bench BenchmarkOnce -test.run TestOnceN #^$$

clean-bench: clean-cache
	@rm -f $(BENCH_ONCE)

clean-cache:
	@rm -f *.log
	@rm -f out/*.wav out/*.wav
	@rm -f msc/*.log msc/*.logcache
	@find msc/ -type d -not \( -regex '.*/res.*' -o -regex '.*/$$' \) | xargs rm -rf

clean: clean-cache
	@rm -f bin/*

.PHONY: clean clean-bench clean-cache
