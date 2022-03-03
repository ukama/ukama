PROJECT := $(notdir $(CURDIR))
PLATFORMLIB := usys
BUILDSATICLIB = lib$(PLATFORMLIB).a
BUILDDYNAMICLIB = $(BUILDDIR)/lib$(PLATFORMLIB).so
PLATFORMLIBS = $(BUILDDYNAMICLIB)
TESTEXEC = $(BUILDDIR)/$(PROJECT)

VERSION = v0.0.1

UNITYROOT := ../tools/unity
SRCDIRS := sys/src log/src
INCDIRS := sys/inc log/inc $(UNITYROOT)/src
BUILDDIR := build
TESTDIRS := test $(UNITYROOT)/src
# Source extensions
SRCEXTS := c

# Header extensions
HDREXTS := h 

# List of all recognized files found in the specified directories for test
SOURCES := $(foreach dir, $(SRCDIRS), $(foreach ext, $(SRCEXTS), $(wildcard $(dir)/*.$(ext))))
OBJECTS := $(foreach ext, $(SRCEXTS), $(patsubst %.$(ext), $(BUILDDIR)/%.o, $(filter %.$(ext), $(SOURCES))))
INCLUDES := $(foreach dir, $(INCDIRS), $(foreach ext, $(HDREXTS), $(wildcard $(dir)/*.$(ext))))

TESTSRCS := $(foreach dir, $(TESTDIRS), $(foreach ext, $(SRCEXTS), $(wildcard $(dir)/*.$(ext))))
TESTOBJS := $(foreach ext, $(SRCEXTS), $(patsubst %.$(ext), $(BUILDDIR)/%.o, $(filter %.$(ext), $(TESTSRCS))))

$(info Sources :: $(SOURCES))
$(info Includes:: $(INCLUDES))
$(info Objects :: $(OBJECTS))

# Compilers and flags
ifdef XGCCPATH
CC = $(XGCCPATH)$(XGCC)
LD = $(XGCCPATH)$(XLD)
else
CC=gcc
endif

LDLIBS = -lpthread -lrt

override CFLAGS += -g -Wall -Wno-unused-variable -fPIC -DHAVE_SYS_TIME_H -DDMT_ABORT_NULL
override LDFLAGS +=  $(LDPATH) -L$(CURDIR)/$(BUILDDIR) $(LDLIBS) 
INCFLAGS := $(INCDIRS:%=-I%)
DEPFLAGS := -MMD -MP
MEMCHECKREPORT := $(BUILDDIR)/memcheck.report

$(info CC:: $(CC))
$(info LDFLAGS :: $(LDFLAGS))

# Tools and flags
CPPLINT := cpplint
override CPPLINTFLAGS += --linelength=100 --filter=-build/header_guard,-runtime/references,-runtime/indentation_namespace,-build/namespaces --extensions=$(subst $( ),$(,),$(SRCEXTS)) --headers=$(subst $( ),$(,),$(HDREXTS))
CPPCHECK := cppcheck
override CPPCHECKFLAGS += --enable=style,warning,missingInclude
MEMCHECK := valgrind
override MEMCHECKFLAGS += --log-file=$(MEMCHECKREPORT) --track-origins=yes --leak-check=full --show-leak-kinds=all

# This makefile name
MAKEFILE := $(lastword $(MAKEFILE_LIST))

# Function to compile using $(CC) : (files: .c)
define compilecc
	@mkdir -p $(dir $1)
	@$(CC) -c $2 -o $1 $(CFLAGS) $(INCFLAGS) -MT $1 -MF $(BUILDDIR)/$3.Td $(DEPFLAGS)
	@mv -f $(BUILDDIR)/$3.Td $(BUILDDIR)/$3.d && touch $1
	@echo CC: $1
endef

# Rules to build objects for each source file extension
$(BUILDDIR)/%.o: %.c $(BUILDDIR)/%.d $(MAKEFILE) $@
	$(call compilecc,$@,$<,$*)

$(BUILDDIR)/%.d: ;


.PHONY: all help run clean force cpplint cppcheck info list-headers list-sources list-objects debug container $(PLATFORMLIBS) 

# Main target for building
all: $(PLATFORMLIBS) $(TESTEXEC)
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
	@echo "Source file Done."
	#$(foreach inc, $(INCLUDES), $(call clang-format,$(inc)))
	#@echo "Include file Done."


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
                    { key: readability-identifier-naming.NamespaceCase, value: camelBack },\
                    { key: readability-identifier-naming.ClassCase, value: CamelCase  },\
                    { key: readability-identifier-naming.StructCase, value: CamelCase  },\
                    { key: readability-identifier-naming.FunctionCase, value: lower_case },\
                    { key: readability-identifier-naming.VariableCase, value: camelBack },\
                    { key: readability-identifier-naming.TypedefCase, value: CamelCase },\
                    { key: readability-identifier-naming.GlobalConstantCase, value: camelBack },\
                    { key: readability-braces-around-statements.ShortStatementLines, value: 0}\
                    ]}" --quiet "$$src" -- $(INCFLAGS); \
	done
	@echo "Ukama Coding Style check pass..!!"

# Make sure Make do not delete included dependencies files
.PRECIOUS: $(BUILDDIR)/%.d

# Include the dependency files here (should not be before first target)
include $(wildcard $(foreach ext, $(SRCEXTS), $(patsubst %.$(ext), $(BUILDDIR)/%.d, $(filter %.$(ext), $(SOURCES)))))

$(PLATFORMLIBS): $(OBJECTS)
	@echo "Building dynamic lib" 
	$(CC) -o $@ $^ $(LDPATH) $(LDLIBS) -shared


# Compile binary if necessary, checks for modified files first
$(TESTEXEC):$(TESTOBJS) $(PLATFORMLIBS) $(MAKEFILE) $@
ifeq ($(filter-out %.c,$(TESTSRCS)),$(blank))
	@$(CC) -o $@ $(TESTOBJS) $(LINKFLAG) $(CFLAGS) $(LDFLAGS) -l$(PLATFORMLIB)
	@echo CC: $@ 
else
	$(CXX) -o $@ $(TESTOBJS) $(CXXFLAGS) $(LDFLAGS) -l$(PLATFORMLIBS)
	@echo CXX: $@
endif

run: $(TESTEXEC)
	@LD_LIBRARY_PATH=$(BUILDDIR) ./$(TESTEXEC)
	
# Clean all build files
clean:
	rm -rf $(PLATFORMLIBS)
	rm -rf $(BUILDDIR)
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
	@echo CC: $(CC)
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
		    { key: readability-identifier-naming.NamespaceCase, value: camelBack },\
		    { key: readability-identifier-naming.ClassCase, value: CamelCase  },\
		    { key: readability-identifier-naming.StructCase, value: CamelCase  },\
		    { key: readability-identifier-naming.FunctionCase, value: lower_case },\
		    { key: readability-identifier-naming.VariableCase, value: camelBack },\
		    { key: readability-identifier-naming.TypedefCase, value: CamelCase },\
		    { key: readability-identifier-naming.GlobalConstantCase, value: camelBack },\
		    { key: readability-braces-around-statements.ShortStatementLines, value: 0}\
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
