
.PHONY: all

PLUGS := $(wildcard *.go)
OBJS := $(PLUGS:go=so)

all: $(OBJS)

$(OBJS): %.so: %.go
	go build -buildmode=c-shared -a -o ../build/plugins/$@ $<
