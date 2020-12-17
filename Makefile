#Makefile UkamaDistro 

UKAMAOSMAKE := $(abspath $(lastword $(MAKEFILE_LIST)))
UKAMAOSDIR := $(dir $(UKAMAOSMAKE))

ANODEBOARD = anode
CNODEBOARD = cnode
ENODEBOARD = enode
BUILDSUBDIRS = firmware os distro 

#Final UKAMAOS path and name
ROOTFSDIR := _ukamafs
ROOTFSPATH := $(UKAMAOSDIR)$(ROOTFSDIR)

#Script for initramfs
RAMFSSCRIPT := $(UKAMAOSDIR)distro/scripts/mkrootfs.sh

#BUILDSUBDIRS = firmware/build os/build distro/build 
.PHONY: subdirs $(BUILDSUBDIRS) initramfs clean

#Export variables
export 

#Build sub-directories
subdirs: $(BUILDSUBDIRS)
$(BUILDSUBDIRS):
	$(MAKE) -C $@

#Initramfs
initramfs:
	@echo Creating initramfs for Ukama Distro.
	(chmod +x $(RAMFSSCRIPT))
	($(RAMFSSCRIPT) -p $(ROOTFSPATH) -u $(TARGETBOARD))
	
#Clean sub-directories
clean:
	@echo Starting cleaning process.
	for dir in $(BUILDSUBDIRS); do \
		$(MAKE) -C $$dir -f Makefile $@; \
  	done
	rm -rf $(ROOTFSPATH)
	rm -rf *.img
