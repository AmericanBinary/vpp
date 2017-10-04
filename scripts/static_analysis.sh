#!/bin/bash

WHITELIST_CONTENT="^// DO NOT EDIT|^// File generated by|^// Code generated by|^// Automatically generated"
WHITELIST_ERRORS="should not use dot imports"

function static_analysis() {
  local TOOL="${@}"
  local PWD=$(pwd)

  local FILES=$(find "${PWD}" -mount -name "*.go" -type f -not -path "${PWD}/vendor/*" -exec grep -LE "${WHITELIST_CONTENT}"  {} +)


  local CMD=$(${TOOL} "${PWD}/cmd${SELECTOR}")
  local FLAVORS=$(${TOOL} "${PWD}/flavors${SELECTOR}")
  local PLUGINS=$(${TOOL} "${PWD}/plugins${SELECTOR}")

  local ALL="$CMD
$PLUGINS
$FLAVORS
"

  local OUT=$(echo "${ALL}" | grep -F "${FILES}" | grep -v "${WHITELIST_ERRORS}")
  if [[ ! -z $OUT ]] ; then
    echo "${OUT}" 1>&2
    exit 1
  fi
}
