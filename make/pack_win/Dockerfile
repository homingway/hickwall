FROM monokrome/wine
MAINTAINER oliveagle@gmail.com

MKDIR /tmp/innosetup
COPY isetup-5.5.5.exe /tmp/innosetup/isetup-5.5.5.exe
RUN xvfb-run -a wine /tmp/innosetup/isetup-5.5.5.exe /silent && rm /tmp/innosetup/isetup-5.5.5.exe
COPY ./iscc_wrapper.sh /root/bin/iscc
