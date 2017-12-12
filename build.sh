#!/bin/bash

SRCDIR=$(pwd)

tinydfs_setup_gopath() {
	cd "$SRCDIR"
	tinydfsGOPATH="${SRCDIR}/build"
	mkdir -p "${tinydfsGOPATH}/src"
	rsync -av --progress "${SRCDIR}/" "${tinydfsGOPATH}/src" --exclude=.git --exclude=build
	# preserve old gopath
	if [ -n "$GOPATH" ]; then
		GOPATH=":$GOPATH"
	fi
	export GOPATH=${tinydfsGOPATH}$GOPATH

	echo "TinyDFS: Inner GOPATH setup done."
}

tinydfs_install_dependencies() {
	# here goes the list of packages
	go get github.com/google/uuid

	echo "TinyDFS: Install package dependencies done."
}

tinydfs_build() {
	cd "${SRCDIR}/tinydfs"
	go build -o "${SRCDIR}/bin/tinydfs"

	echo "TinyDFS: Build done."
}

tinydfs_setup_gopath
tinydfs_install_dependencies
tinydfs_build


