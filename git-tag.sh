#!/bin/bash
set -euxo pipefail
# git-tag.sh
GIT_VER=$(./hugo-docs-i18n.exe version -g)
# GIT_VER=`./hugo-docs-i18n.exe version -g`
VER_MSG=$(./hugo-docs-i18n.exe version -m)
# VER_MSG=`./hugo-docs-i18n.exe version -m`

echo "${GIT_VER}" "${VER_MSG}"

