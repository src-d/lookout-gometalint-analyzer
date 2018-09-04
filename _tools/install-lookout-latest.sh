#!/bin/bash
#
# Installs latest lookout SDK binary
# Depends on GNU grep
set -x

#TODO(bzz): check local cached version

curl -s --connect-timeout 5 \
    --max-time 10 \
    --retry 5 \
    --retry-delay 0 \
    --retry-max-time 40\
    "https://api.github.com/repos/src-d/lookout/releases/latest" \
  |& tee -a ../lookout-install.log \
  | grep -oP '"browser_download_url": "\K(.*)(?=")' \
  | grep linux \
  | wget -qi -

if [[ "${PIPESTATUS[0]}" -ne 0 || "${PIPESTATUS[1]}" -ne 0 || "${PIPESTATUS[2]}" -ne 0 || "${PIPESTATUS[4]}" -ne 0 ]]; then
  echo "Unable download latest lookout SDK release" >&2
  exit 2
fi

if ! tar -xvzf lookout-sdk_*.tar.gz ; then
  echo "Unable to extract lookout release archive" >&2
  exit 2
fi

if ! mv lookout-sdk_*/lookout* .. ; then
  echo "Unable to move lookout-sdk binary file" >&2
  exit 2
fi

rm -rf lookout-sdk_*
