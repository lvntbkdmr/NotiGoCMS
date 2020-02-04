#!/bin/bash
set -u -e -o pipefail

# setup git
git config --global user.email "lvntbkdmr@gmail.com"
git config --global user.name "Levent Bekdemir"
git config --global github.user "lvntbkdmr"
git config --global github.token "${GITHUB_TOKEN}"

git rm -r --cached blog

git submodule add --force https://github.com/lvntbkdmr/blog blog

./NotiGoCMS

cd blog
git status
git add *
git status
now=`date "+%Y-%m-%d %a"`

# "git commit" returns 1 if there's nothing to commit, so don't report this as failed build
set +e
git commit -am "ci: update from notion on ${now}"
if [ "$?" -ne "0" ]; then
    echo "nothing to commit"
    exit 0
fi
set -e
git push || true
cd ../