# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2021-present, Ukama Inc.

include ../../../config.mk
PROJECT := $(notdir $(CURDIR))

# Directories
COMMONINCDIRS := ../common/inc
VENDORINCDIRS := ../../vendor/build/include/
SRCDIRS := src src/agents
INCDIRS := inc/ $(COMMONINCDIRS) $(VENDORINCDIRS)
BUILDDIR := build

#CAPP
CAPPEXE := $(BUILDDIR)/$(PROJECT)

# Version
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
SHASH=$(shell git describe --abbrev=4 --dirty --always --tags)

# Source extensions
SRCEXTS := c

# Header extensions
HDREXTS := h inc

#Cross compiler
ifdef XGCCPATH
XCC := $(XGCCPATH)$(XGCC)
XLD := $(XGCCPATH)$(XLD)
$(info XCC: $(XCC) XLD: $(XLD))
else
XCC := gcc
endif

# List of all recognized files found in the specified directories for test
SOURCES := $(foreach dir, $(SRCDIRS), $(foreach ext, $(SRCEXTS), $(wildcard $(dir)/*.$(ext))))
OBJECTS := $(foreach ext, $(SRCEXTS), $(patsubst %.$(ext), $(BUILDDIR)/%.o, $(filter %.$(ext), $(SOURCES))))
INCLUDES := $(foreach dir, $(INCDIRS), $(foreach ext, $(HDREXTS), $(wildcard $(dir)/*.$(ext))))

LDLIBS = -lprom -lpromhttp -lmicrohttpd -lpthread -ltoml -latomic
LDPATH = -L$(CURDIR)/../../vendor/build/lib -L$(CURDIR)/../../vendor/build/usr/local/lib 

LINKFLAG = --static
override CFLAGS += -g -Wall -fPIC -DHAVE_SYS_TIME_H -DGITHASH=\"$(SHASH)\"
override LDFLAGS =  $(LDLIBS) $(LDPATH)
INCFLAGS := $(INCDIRS:%=-I%)
DEPFLAGS := -MMD -MP
MEMCHECKREPORT := $(BUILDDIR)/memcheck.report
INCFLAGS += -I$(CURDIR)/../../vendor/build/include -I$(CURDIR)/../../vendor/build/usr/local/include

# Tools and flags
CPPLINT := cpplint
override CPPLINTFLAGS += --linelength=100 --filter=-build/header_guard,-runtime/references,-runtime/indentation_namespace,-build/namespaces --extensions=$(subst $( ),$(,),$(SRCEXTS)) --headers=$(subst $( ),$(,),$(HDREXTS))
CPPCHECK := cppcheck
override CPPCHECKFLAGS += --enable=style,warning,missingInclude
MEMCHECK := valgrind
override MEMCHECKFLAGS += --log-file=$(MEMCHECKREPORT) --track-origins=yes --leak-check=full --show-leak-kinds=all


# This makefile name
MAKEFILE := $(lastword $(MAKEFILE_LIST))

#$(info BUILDDIR $(BUILDDIR) CAPPEXE $(CAPPEXE) LDPATH $(LDPATH))

# Function to compile using $(CC) : (files: .c)
define compilecc
	@mkdir -p $(dir $1)
	@$(XCC) -c $2 -o $1 $(CFLAGS) $(INCFLAGS) -MT $1 -MF $(BUILDDIR)/$3.Td $(DEPFLAGS)
	@mv -f $(BUILDDIR)/$3.Td $(BUILDDIR)/$3.d && touch $1
	@echo CC: $1
endef

# Rules to build objects for each source file extension
$(BUILDDIR)/%.o: %.c $(BUILDDIR)/%.d $(MAKEFILE) $@
	$(call compilecc,$@,$<,$*)

$(BUILDDIR)/%.d: ;


.PHONY: all help run clean force cpplint cppcheck info list-headers list-sources list-objects debug

# Main target for building
all: $(CAPPEXE)
	@echo Done.

# Print commands
help:
	@echo "Some useful make targets:"
	@echo " make all          - Build entire project (modified sources only or dependents)"
	@echo " make run          - Build and launch excecutable immediately"
	@echo " make force        - Force rebuild of entire project (clean first)"
	@echo " make clean        - Remove all build output"
	@echo " make info         - Print out project configurations"
	@echo " make format-style  - Format your code to Ukama C Style"	
	@echo " make check-style   - Ukama C style checker tool"
	@echo " make memcheck     - Launch executable and does memory analysis."
	@echo " make cppcheck     - Static code analysis tool for the C and C++"
	@echo " make list-headers - Print out all recognized headers files"
	@echo " make list-sources - Print out all recognized sources files"
	@echo " make list-objects - Print out final objects"
	@echo ""

format-style:
	$(foreach src, $(SOURCES), $(call clang-format,$(src)))	
	@echo "Done"

check-style:
	$(eval CWD = $(shell pwd))
	echo $(CWD)
	@for src in $(SOURCES) ; do \
		echo "Checking format for $(CWD)/$$src"; \
		dif=`clang-format-9 "$(CWD)/$$src" | diff "$(CWD)/$$src" - | wc -l`; \
		if [ `echo $$dif` != `echo 0` ]; then \
			echo "clang-format: Fail";\
			echo "Err: $$src $$dif Lines to be modified."; \
			echo "Execute \" make format-style\".";\
			exit 1;\
		else\
			echo "clang-format: Pass.";\
		fi ;\
		echo "clang-tidy: ";\
		clang-tidy-10 --checks='-*,readability-identifier-naming' \
			-config="{CheckOptions: [ \
                    { key: readability-identifier-naming.NamespaceCase, value: lower_case },\
                    { key: readability-identifier-naming.ClassCase, value: lower_case  },\
                    { key: readability-identifier-naming.StructCase, value: lower_case  },\
                    { key: readability-identifier-naming.FunctionCase, value: lower_case },\
                    { key: readability-identifier-naming.VariableCase, value: lower_case },\
                    { key: readability-identifier-naming.GlobalConstantCase, value: lower_case }\
                    ]}" --quiet "$$src" -- $(INCFLAGS); \
	done
	@echo "Ukama Coding Style check pass..!!"

# Make sure Make do not delete included dependencies files
.PRECIOUS: $(BUILDDIR)/%.d

# Include the dependency files here (should not be before first target)
include $(wildcard $(foreach ext, $(SRCEXTS), $(patsubst %.$(ext), $(BUILDDIR)/%.d, $(filter %.$(ext), $(SOURCES)))))


# Compile binary if necessary, checks for modified files first
$(CAPPEXE): $(OBJECTS) $(MAKEFILE) $@
ifeq ($(filter-out %.c,$(SOURCES)),$(blank))
	@$(XCC) -o $@ $(OBJECTS) $(LINKFLAG) $(CFLAGS) $(LDFLAGS)
	@echo CC: $@ 
else
	$(CXX) -o $@ $(OBJECTS) $(CXXFLAGS) $(LDFLAGS)
	@echo CXX: $@
endif

# Launch excecutable, compile if necessary
run: $(CAPPEXE)
	@LD_LIBRARY_PATH=$(BUILDDIR) ./$(CAPPEXE)

# Memory check
memcheck: $(CAPPEXE)
	@LD_LIBRARY_PATH=$(BUILDDIR) $(MEMCHECK) $(MEMCHECK_REPORT) $(MEMCHECKFLAGS) ./$(EXCECUTABLE)

# Clean all build files
clean:
	@rm -rf $(CAPPEXE)
	@rm -rf $(BUILDDIR)
	@echo Cleaned.

# Force build of all files
force: clean all

# C++ style checker tool (following Google's C++ style guide)
cpplint:
	@$(CPPLINT) $(CPPLINTFLAGS) $(SOURCES) $(INCLUDES)

# Static code analysis tool for the C and C++
cppcheck:
	@$(CPPCHECK) $(CPPCHECKFLAGS) $(SOURCES) $(INCLUDES) $(INCFLAGS)

# Prints out project configurations
info:
	@echo Project: $(PROJECT)
	@echo Excecutable: $(EXCECUTABLE)
	@echo SourceDirs: $(SRCDIRS)
	@echo IncludeDirs: $(INCDIRS)
	@echo BuildDir: $(BUILDDIR)
	@echo CC: $(XCC)
	@echo CCFlags: $(CFLAGS)
	@echo LDFlags: $(LDFLAGS)
	@echo IncFlags: $(INCFLAGS)
	@echo DepFlags: $(DEPFLAGS)
	@echo CppLintFlags : $(CPPLINTFLAGS)
	@echo CppCheckFlags: $(CPPCHECKFLAGS)
	@echo MemCheckFlags: $(MEMCHECKFLAGS)


list-sources:
	$(foreach src, $(SOURCES), $(call print,$(src)))

list-headers:
	$(foreach hdr, $(INCLUDES), $(call print,$(hdr)))

list-objects:
	$(foreach obj, $(OBJECTS), $(call print,$(obj)))


# Debugging of this makefile, for development
debug:
	@echo SourceExts: $(SRCEXTS)
	@echo HeaderExts: $(HDREXTS)
#
# ### CLANG FORMAT ###
#
define clang-format
	@echo "Formatting $1"; 
	$(shell clang-format-9 -i $1;)
	$(shell clang-tidy-10 -checks='-*,readability-identifier-naming' \
		    -config="{CheckOptions: [ \
		    { key: readability-identifier-naming.NamespaceCase, value: lower_case },\
		    { key: readability-identifier-naming.ClassCase, value: lower_case  },\
		    { key: readability-identifier-naming.StructCase, value: lower_case  },\
		    { key: readability-identifier-naming.FunctionCase, value: lower_case },\
		    { key: readability-identifier-naming.VariableCase, value: lower_case },\
		    { key: readability-identifier-naming.GlobalConstantCase, value: lower_case }\
		    ]}" --fix "$1" -- $(INCFLAGS) )
endef

# 
# ### Utils ###
#
define print
	@echo $1

endef

# comma -> $(,)
, = ,
# blank -> $(blank)
blank =
# space -> $( )
space = $(blank) $(blank)
$(space) = $(space)
