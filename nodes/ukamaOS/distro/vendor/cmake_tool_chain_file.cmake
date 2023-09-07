
#include_directories(${DESTDIR}/include)
include_directories(${CFLAGS})

set(CMAKE_C_FLAGS ${CFLAGS})
set(CMAKE_C_COMPILER ${XCC})
set(CMAKE_LD_FLAGS ${LDFLAGS})

set(CMAKE_TRY_COMPILE_TARGET_TYPE   STATIC_LIBRARY)
