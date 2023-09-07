# Libidn2 README -- Introduction information

Libidn2 is a free software implementation of IDNA2008, Punycode and
Unicode TR46.  Its purpose is to encode and decode internationalized
domain names.

For technical reference, see:

 * [IDNA2008 Framework](https://tools.ietf.org/html/rfc5890)
 * [IDNA2008 Protocol](https://tools.ietf.org/html/rfc5891)
 * [IDNA2008 Unicode tables](https://tools.ietf.org/html/rfc5892)
 * [IDNA2008 Bidi rule](https://tools.ietf.org/html/rfc5893)
 * [Punycode](https://tools.ietf.org/html/rfc3492)
 * [Unicode IDNA Compatibility Processing](https://www.unicode.org/reports/tr46/)

The library contains functionality to convert internationalized domain
names to and from ASCII Compatible Encoding (ACE).

The API consists of two main functions, `idn2_to_ascii_8z` for
converting data from UTF-8 to ASCII Compatible Encoding (ACE), and
`idn2_to_unicode_8z8z` to convert ACE names into UTF-8 format. There
are several variations of these main functions, which accept UTF-32,
or input in the local system encoding. All functions assume
zero-terminated strings.

This library is backwards (API) compatible with the [libidn
library](https://www.gnu.org/software/libidn/).  Replacing the
`idna.h` header with `idn2.h` into a program is sufficient to switch
the application from IDNA2003 to IDNA2008 as supported by this
library.

Libidn2 is believed to be a complete IDNA2008 and TR46 implementation,
it contains an extensive test-suite, and is included in the continuous
fuzzing project
[OSS-Fuzz](https://bugs.chromium.org/p/oss-fuzz/issues/list?q=libidn2).

You can check the current test code coverage
[here](https://libidn.gitlab.io/libidn2/coverage/index.html) and the
current fuzzing code coverage
[here](https://libidn.gitlab.io/libidn2/fuzz-coverage/index.html).


# License

The installed C library libidn2 is dual-licensed under LGPLv3+|GPLv2+,
while the rest of the package is GPLv3+.  See the file
[COPYING](COPYING) for detailed information.


# Online docs

[API reference](https://libidn.gitlab.io/libidn2/reference/api-index-full.html)

[Manual](https://libidn.gitlab.io/libidn2/manual/libidn2.html)


# Obtaining the source

Software releases of libidn2 can be downloaded from
https://ftp.gnu.org/gnu/libidn/ and ftp://ftp.gnu.org/gnu/libidn/

Development of libidn2 is organized [through GitLab
website](https://gitlab.com/libidn/libidn2), and there is [an issue
tracker for reporting bugs](https://gitlab.com/libidn/libidn2/issues).


# Dependencies

To build Libidn2 you will need a POSIX shell to run ./configure and
the Unix make tool.

 * [Bash](https://www.gnu.org/software/bash/)
 * [Make](https://www.gnu.org/software/make/)

The shared libidn2 library uses GNU libunistring for Unicode
processing and GNU libiconv for character set conversion.  You should
install them before building and installing libidn2.  See the
following links for more information on these packages:

 * [Unistring](https://www.gnu.org/software/libunistring/)
 * [iconv](https://www.gnu.org/software/libiconv/)

Note that the iconv dependency is optional -- it is required for the
functions involving locale to UTF conversions -- but is recommended.

If you wish to build the project from version controlled sources,
rebuild all generated files (e.g., run autoreconf), or modify some
source code files, you will need to have additional tools installed.
None of the following tools are necessary if you build Libidn2 in the
usual way (i.e., ./configure && make).

 * [Automake](https://www.gnu.org/software/automake/)
 * [Autoconf](https://www.gnu.org/software/autoconf/)
 * [Libtool](https://www.gnu.org/software/libtool/)
 * [Gettext](https://www.gnu.org/software/gettext/)
 * [Texinfo](https://www.gnu.org/software/texinfo/)
 * [Gperf](https://www.gnu.org/software/gperf/)
 * [Gengetopt](https://www.gnu.org/software/gengetopt/)
 * [help2man](https://www.gnu.org/software/help2man/)
 * [Tar](https://www.gnu.org/software/tar/)
 * [Gzip](https://www.gnu.org/software/gzip/)
 * [Texlive & epsf](https://www.tug.org/texlive/) (for PDF manual)
 * [GTK-DOC](https://www.gtk.org/gtk-doc/) (for API manual)
 * [Git](https://git-scm.com/)
 * [Perl](https://www.cpan.org/)
 * [Valgrind](https://valgrind.org/) (optional)
 * [abi-compliance-checker](https://github.com/lvc/abi-compliance-checker)

The software is typically distributed with your operating system, and
the instructions for installing them differ.  Here are some hints:

Debian 10.x, Ubuntu 20.04:
```
apt-get install git autoconf automake libtool gettext autopoint gperf
apt-get install libunistring-dev valgrind gengetopt help2man
apt-get install texinfo git2cl gtk-doc-tools
apt-get install abi-compliance-checker abigail-tools
```


# Contributing

See [the contributing document](CONTRIBUTING.md).


# Estimating code coverage

Dependencies:
 * [lcov](https://github.com/linux-test-project/lcov) (for code coverage)

To test the code coverage of the test suite use the following:
```
$ ./configure --enable-code-coverage
$ make && make check && make code-coverage-capture
```

The current coverage report can be found [here](https://libidn.gitlab.io/libidn2/coverage/).


# Fuzzing

Libidn2 is being continuously fuzzed by [OSS-Fuzz](https://github.com/google/oss-fuzz).

Of course you can do local fuzzing on your own, see `fuzz/README.md` for instructions.

The code coverage of our fuzzers can be found [here](https://libidn.gitlab.io/libidn2/fuzz-coverage/).
