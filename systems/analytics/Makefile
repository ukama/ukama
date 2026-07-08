SUBDIRS := ingest analysis aggregator api-gateway

.PHONY: $(SUBDIRS) build test gen clean

build:
	for dir in $(SUBDIRS); do \
		$(MAKE) -C $$dir build || exit 1; \
	done

test:
	for dir in $(SUBDIRS); do \
		$(MAKE) -C $$dir test || exit 1; \
	done

gen:
	for dir in $(SUBDIRS); do \
		$(MAKE) -C $$dir gen || exit 1; \
	done

clean:
	for dir in $(SUBDIRS); do \
		$(MAKE) -C $$dir clean || exit 1; \
	done
