SUBDIRS := $(wildcard cloud/*/.)  $(wildcard bootstrap/*/.)  $(wildcard common/.) 

build: $(SUBDIRS)	
$(SUBDIRS):
	$(MAKE) -C $@

test: $(SUBDIRS)	
$(SUBDIRS):
	$(MAKE) -C $@

.PHONY: build test $(SUBDIRS)