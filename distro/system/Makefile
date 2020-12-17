# Build System apps [Will be replaced by containers]
#ALLDIR := $(sort $(dir $(wildcard */.)))
#FDIR = $(wildcard $(ALLDIR)/Make*)
#FLTRDIR := $(shell find ./*/ -maxdepth 1 -type f -name 'Makefile' -printf '%h\0' | sort -zu | awk -F'./' '{$$1=$$1; print $$4}') 
CURMAKE := $(abspath $(lastword $(MAKEFILE_LIST)))
CURPATH := $(dir $(CURMAKE))

SYSSUBDIRS = $(shell find ./*/ -maxdepth 1 -type f -name 'Makefile' -printf '%h'| awk -F './' '{ for (i=1; i<=NF; i++) printf $$i" " }')

#Add a flag with supported architecures so no package which is specfic for X86 get build for arm.

.PHONY: subdirs $(SYSSUBDIRS) clean info

ifndef ROOTFSPATH
override ROOTFSPATH = $(CURPATH)_ukamafs
endif

#Export variables
export

#Build sub-directories
subdirs: $(SYSSUBDIRS)
$(SYSSUBDIRS):
	$(MAKE) -C $@

#Clean sub-directories
clean:
	@echo "Cleaning System build."
	rm -rf $(CURPATH)_ukamafs
	for dir in $(SYSSUBDIRS); do \
                $(MAKE) -C $$dir -f Makefile $@; \
        done
info: 
	@echo $(SYSSUBDIRS)
	#@echo $(ALLDIR)
	#@echo $(FDIR)
	#@echo $(shell find ./*/ -maxdepth 1 -type f -name 'Makefile' -printf '%h\0' | sort -zu | awk -F './' '{ for (i=1; i<=NF; i++) printf $$i" " }')
	#@echo $(shell find . -type f -name 'Makefile' -printf '%h' | sort -zu | sed -z 's/$$/ /')
	#@echo $(shell find ./*/ -type f -name 'Makefile' -printf '%h'| awk -F './' '{ for (i=1; i<=NF; i++) printf $$i" " }')


