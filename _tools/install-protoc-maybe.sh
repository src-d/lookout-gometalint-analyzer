#!/bin/bash
#
# Checks Protobuff compiler version. 
# If absent, installs one in current dir from release.

PROTOC_VER="3.6.0"
OS="$(uname)"
if [[ "$OS" == "Darwin" ]]; then
	protoc_os="osx"
else
	protoc_os="linux"
fi

cur_ver="$(protoc --version | grep -o '[^ ]*$')"
if [[ "$cur_ver" == "$PROTOC_VER" ]]; then
	echo "Using protoc version $cur_ver"
else
	echo "Installing protoc version $PROTOC_VER"
	local protoc_zip = "protoc-$PROTOC_VER-$protoc_os-x86_64.zip"
	local url = "https://github.com/google/protobuf/releases/download/v$PROTOC_VER/$protoc_zip"

	wget "$url"
	if [[ "$?" -ne 0 ]]; then
		echo "Failed to download protoc release from $url"
		exit 2
	fi
	mkdir -p "./protoc"
	unzip -d "./protoc" "$protoc_zip"
	rm "$protoc_zip"
fi
