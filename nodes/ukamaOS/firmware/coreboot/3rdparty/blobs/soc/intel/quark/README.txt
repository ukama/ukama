These binaries are the result of compiling the QuarkFsp sources and
provided as a convenience since not everybody wants to setup a working
edk2 tree.
Since the sources, as well as the edk2 sources, are BSD-licensed,
redistribution is not an issue.

These binaries are untested and come with no warranty!

Instructions to build your own binaries, using the coreboot toolchain:

$ git clone http://github.com/tianocore/edk2
$ cd edk2
$ git checkout a5cd3bb037cf87ecda0a5c8cd8a3eda722591b70
$ git clone https://review.gerrithub.io/LeeLeahy/quarkfsp QuarkFspPkg
$ (cd QuarkFspPkg; patch -p1 -i $path/to/this/directory/QuarkFsp.patch)
$ . edksetup.sh
$ cat $path/to/your/coreboot/toolchain/share/edk2config/tools_def.txt >> Conf/tools_def.txt

$ # builds the debug images
$ QuarkFspPkg/BuildFsp2_0.sh -d32
$ QuarkFspPkg/BuildFsp2_0Pei.sh -d32

$ # builds the release images
$ QuarkFspPkg/BuildFsp2_0.sh -r32
$ QuarkFspPkg/BuildFsp2_0Pei.sh -r32
