#!/bin/bash
#set -u -e -o pipefail

GITHUB_ACTOR="lvntbkdmr"

# setup git
git config --global user.email "lvntbkdmr@gmail.com"
git config --global user.name "Levent Bekdemir"
git config --global github.user "${GITHUB_ACTOR}"
git config --global github.token "${GH_TOKEN_OVERRIDE}"

git add -u
git reset update_notigocms.sh
now=`date "+%Y-%m-%d %a"`

# "git commit" returns 1 if there's nothing to commit, so don't report this as failed build
set +e
git commit -am "ci: update from notion on ${now}"
if [ "$?" -ne "0" ]; then
    echo "nothing to commit"
    exit 0
fi
set -e
git push "https://${GITHUB_ACTOR}:${GH_TOKEN_OVERRIDE}@github.com/lvntbkdmr/NotiGoCMS.git" master || true
