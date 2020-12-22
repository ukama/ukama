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
	
#Clean
clean:
	@echo Starting cleaning process.
	for dir in $(BUILDSUBDIRS); do \
		$(MAKE) -C $$dir -f Makefile $@; \
  	done
	rm -rf $(ROOTFSPATH)
	rm -rf *.img

#Distclean

distclean:
	@echo Starting distclean process.
	for dir in $(BUILDSUBDIRS); do \
                $(MAKE) -C $$dir -f Makefile $@; \
	done
	rm -rf $(ROOTFSPATH)
	rm -rf *.img

# Print commands
help:
	@echo "Some useful make targets:"
	@echo " make TARGETBOARD=< board >                      - Build entire project with required toochain for board (anode|cnode|hnode)"
	@echo " make initramfs TARGETBOARD=< board >            - Build initramfs image for required board (anode|cnode|hnode)"
	@echo " make distclean                                  - Remove all the build files, configs, patches etc."
	@echo " make clean                                      - Remove all build output"
	@echo " ** All build artifacts are stored under _ukamafs directory. ** "
