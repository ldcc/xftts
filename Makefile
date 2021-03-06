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
		$(XFTTS) -t $(DEF_SPEECH) -o $(@F);\
	else\
		$(XFTTS) -t $(@D) -o $(@F);\
	fi

serve: $(XFTTS)
	./$(XFTTS)

sent-once:
	go test -work -ldflags "-s -w" -a -installsuffix cgo -c -o $(SEND_ONCE) xftts/tests
	$(SEND_ONCE) -test.v -test.run TestSendGet

$(BENCH_ONCE):
	go test -work -ldflags "-s -w" -a -installsuffix cgo -c -o $(BENCH_ONCE) xftts/tests

bench: $(BENCH_ONCE)
	mkdir -p bin out
	$(BENCH_ONCE) -test.v -test.bench BenchmarkOnce -test.run ^$$
test: $(BENCH_ONCE)
	mkdir -p bin out
	$(BENCH_ONCE) -test.v -test.run TestOnceN

MSC := xf/msc
clean-cache:
	@rm -f out/*
	@rm -f $(MSC)/*.log $(MSC)/*.logcache
	@find $(MSC)/ -type d -not \( -regex '.*/res.*' -o -regex '.*/$$' \) | xargs rm -rf

clean: clean-cache
	@rm -f bin/*

.PHONY: clean clean-cache
