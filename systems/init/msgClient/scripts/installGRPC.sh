#!/bin/sh -x

# Build and install GRPC with sharedlibs enable to .ukama/local

export UKAMA_DEP_ROOT=$HOME/.ukama/
export UKAMA_SRC_DIR=$UKAMA_DEP_ROOT/src/
export UKAMA_INSTALL_DIR=$UKAMA_DEP_ROOT/local/

mkdir -p $UKAMA_SRC_DIR
mkdir -p $UKAMA_INSTALL_DIR

export PATH="$UKAMA_INSTALL_DIR/bin:$PATH"
CWD=`pwd`

# remove existing grpc src
rm -rf $UKAMA_SRC_DIR/grpc

# Clone grpc ./
git clone --recurse-submodules -b v1.46.3 --depth 1 --shallow-submodules \
	https://github.com/grpc/grpc $UKAMA_SRC_DIR/grpc

cd $UKAMA_SRC_DIR/grpc
mkdir -p cmake/build
cd cmake/build
cmake -DBUILD_SHARED_LIBS=ON \
	  -DgRPC_INSTALL=ON \
      -DgRPC_BUILD_TESTS=OFF \
      -DCMAKE_INSTALL_PREFIX=$UKAMA_INSTALL_DIR ../../
make
make install
cd $CWD

#rm -rf $UKAMA_SRC_DIR/grpc
