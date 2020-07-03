FROM \
    karalabe/xgo-latest

MAINTAINER \
    Bill Zissimopoulos <billziss at navimatics.com>

# add 32-bit and 64-bit architectures and install 7zip
RUN \
    dpkg --add-architecture i386 && \
    dpkg --add-architecture amd64 && \
    apt-get update && \
    apt-get install -y --no-install-recommends p7zip-full

# install OSXFUSE
RUN \
    wget -q -O osxfuse.dmg \
        https://github.com/osxfuse/osxfuse/releases/download/osxfuse-3.8.3/osxfuse-3.8.3.dmg && \
    7z e osxfuse.dmg 0.hfs && \
    7z e 0.hfs "FUSE for macOS/Extras/FUSE for macOS 3.8.3.pkg" && \
    7z e "FUSE for macOS 3.8.3.pkg" Core.pkg/Payload && \
    7z e Payload && \
    7z x Payload~ -o/tmp && \
    cp -R /tmp/usr/local/include/osxfuse /usr/local/include && \
    cp /tmp/usr/local/lib/libosxfuse_i64.2.dylib /usr/local/lib/libosxfuse.dylib

# install LIBFUSE
RUN \
    apt-get update -y && \
    apt-get install -y --no-install-recommends libfuse-dev:i386 && \
    apt-get install -y --no-install-recommends libfuse-dev:amd64 && \
    apt-get download libfuse-dev:i386 && \
    dpkg -x libfuse-dev*i386*.deb /

# install WinFsp-FUSE
RUN \
    wget -q -O winfsp.zip \
        https://github.com/billziss-gh/winfsp/archive/release/1.2.zip && \
    7z e winfsp.zip 'winfsp-release-1.2/inc/fuse/*' -o/usr/local/include/winfsp

ENV \
    OSXCROSS_NO_INCLUDE_PATH_WARNINGS 1
