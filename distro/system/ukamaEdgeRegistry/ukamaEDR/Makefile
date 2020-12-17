PROJECT := $(notdir $(CURDIR))
EXCECUTABLE = $(BUILDDIR)/$(PROJECT)
BUILDSATICLIB = lib/ubsp/build/libubsp.a
BUILDDYNAMICLIB = lib/ubsp/build/libubsp.so
BUILDLIBS = $(BUILDSATICLIB)
CJSONLIB = lib/ubsp/third-party/libcjson.a
UBSPLIB = $(BUILDSATICLIB)
# Directories specification for LIBS
SRCDIRS := src utils registry
INCDIRS := $(CURDIR) 
BUILDDIR := build

# Directories specification for Executable
TESTSRCDIRS := $(SRCDIRS) test

# Source extensions
SRCEXTS := c

# Header extensions
HDREXTS := h inc

# List of all recognized files found in the specified directories for test
SOURCES := $(foreach dir, $(SRCDIRS), $(foreach ext, $(SRCEXTS), $(wildcard $(dir)/*.$(ext))))
OBJECTS := $(foreach ext, $(SRCEXTS), $(patsubst %.$(ext), $(BUILDDIR)/%.o, $(filter %.$(ext), $(SOURCES))))

# List of all recognized files found in the specified directories for libs
TESTSOURCES := $(foreach dir, $(TESTSRCDIRS), $(foreach ext, $(SRCEXTS), $(wildcard $(dir)/*.$(ext))))
TESTOBJECTS := $(foreach ext, $(SRCEXTS), $(patsubst %.$(ext), $(BUILDDIR)/%.o, $(filter %.$(ext), $(TESTSOURCES))))

INCLUDES := $(foreach dir, $(INCDIRS), $(foreach ext, $(HDREXTS), $(wildcard $(dir)/*.$(ext))))

# Compilers and flags
CC := gcc
LDLIBS = -lubsp -lcjson -lpthread -lm
LDPATH = -L$(CURDIR)/lib/ubsp/build -L$(CURDIR)/lib/ubsp/third-party
LINKFLAG = --static
override CFLAGS += -g -Wall -Wno-unused-variable -fPIC -DHAVE_SYS_TIME_H
override LDFLAGS +=  $(LDPATH) $(LDLIBS) 
INCFLAGS := $(INCDIRS:%=-I%)
DEPFLAGS := -MMD -MP
MEMCHECKREPORT := $(BUILDDIR)/memcheck.report

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


.PHONY: all help run clean force cpplint cppcheck info list-headers list-sources list-objects debug

# Main target for building
all: $(UBSPLIB) $(EXCECUTABLE)
	@echo Done.

#UBSP lib
$(UBSPLIB):
	@echo "Building UBSP libarary."
	cd lib/ubsp && make CC=$(CC) LD=$(LD)
	
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
$(EXCECUTABLE): $(UBSPLIB) $(OBJECTS) $(MAKEFILE) $@
ifeq ($(filter-out %.c,$(SOURCES)),$(blank))
	@$(CC) -o $@ $(OBJECTS) $(LINKFLAG) $(CFLAGS) $(LDFLAGS)
	@echo CC: $@ 
else
	$(CXX) -o $@ $(OBJECTS) $(CXXFLAGS) $(LDFLAGS)
	@echo CXX: $@
endif

ifdef LIBS
$(BUILDSATICLIB): $(LIBOBJECTS)
	ar rcs $@ $^
        
$(BUILDDYNLIB): $(LIBOBJECTS)
	$(CC) -o $@ $^ $(LDPATH) $(LDLIBS) -shared
endif

# Launch excecutable, compile if necessary
run: $(EXCECUTABLE)
	@LD_LIBRARY_PATH=$(BUILDDIR) ./$(EXCECUTABLE)

# Memory check
memcheck: $(EXCECUTABLE)
	@LD_LIBRARY_PATH=$(BUILDDIR) $(MEMCHECK) $(MEMCHECK_REPORT) $(MEMCHECKFLAGS) ./$(EXCECUTABLE)

# Clean all build files
clean:
	@rm -rf $(EXCECUTABLE)
	@rm -rf $(BUILDDIR)
	@cd lib/ubsp && make $@
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
