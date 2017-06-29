pkg_origin=opsline
pkg_name=swizzle
pkg_version=1.0
pkg_maintainer="Chetan Sarva <csarva@opsline.com>"
pkg_license=()
pkg_source=https://foo.com/dummy/v${pkg_version}/${pkg_name}-${pkg_version}
pkg_shasum=
pkg_deps=()
pkg_build_deps=(core/go core/git core/musl core/gcc)
pkg_bin_dirs=(bin)

pkg_go_path="github.com/opsline/swizzle"

if [[ -f .BUILD/env ]]; then
  source .BUILD/env
fi

if [[ -n "$SNAPSHOT_VERSION" ]]; then
  pkg_version="${pkg_version}-${SNAPSHOT_VERSION}"
fi

do_build_go_static() {
  go get -a -installsuffix cgo -ldflags "-linkmode external -extldflags \"-static\" -s -w" $pkg_go_path/...
}

do_build_go() {
  go get $pkg_go_path/...
}

do_install_go() {
  mkdir -p $pkg_prefix/bin/
  cp -v $GOPATH/bin/* $pkg_prefix/bin/
}

do_begin() {
  export GOPATH="${HAB_CACHE_SRC_PATH}/go"
}

do_clean() {
  do_default_clean
  rm -rf $GOPATH $HAB_CACHE_SRC_PATH/package
}

do_prepare() {
  export CC=$(hab pkg path core/musl)/bin/musl-gcc
  dynamic_linker="$(pkg_path_for musl)/lib/ld-musl-x86_64.so.1"
  export LDFLAGS="$LDFLAGS -Wl,--dynamic-linker=$dynamic_linker"
}

do_download() {
  if [[ -n "$SNAPSHOT" ]]; then
    local ACTUAL_PKG_FILENAME
    ACTUAL_PKG_FILENAME=$(ls /src/.BUILD/ | grep $pkg_name | grep '.tgz$')
    pkg_shasum=$(cat /src/.BUILD/$ACTUAL_PKG_FILENAME.sha256)
    cp -v "/src/.BUILD/$ACTUAL_PKG_FILENAME" "$HAB_CACHE_SRC_PATH/$pkg_filename"
  fi
}

do_verify() {
  if [[ -n "$SNAPSHOT" ]]; then
    do_default_verify
    return $?
  fi
  return 0 # skip for non-snapshot builds (building from source tree)
}

# Whether snapshot or source, create `package` dir
do_unpack() {
  if [[ -n "$SNAPSHOT" ]]; then
    cd $HAB_CACHE_SRC_PATH
    tar xzf $pkg_filename
  else
    cp -a /src $HAB_CACHE_SRC_PATH/package
  fi

  mkdir -p $GOPATH/src/$(dirname $pkg_go_path)
  mv $HAB_CACHE_SRC_PATH/package $GOPATH/src/$pkg_go_path
}

do_build() {
  do_build_go_static
}

do_install() {
  do_install_go
}
