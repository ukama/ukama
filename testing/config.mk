# Binary version
ifdef RELEASE
BIN_VER:=$$(git describe)
else
BIN_VER:=0.0.$(BUILD_NUMBER)
endif
