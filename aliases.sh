#!/usr/bin/env bash

function copy_wasm_script() {
    cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" "public/assets"
}

$@