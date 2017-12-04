#!/bin/bash
set -e
if [[ "$1" == "" ]]
then
  echo -e "\033[31m 请输入注释！ \033[0m"
  exit 1
fi
branch_name=`git branch | sed -n '/*/s///p'`
git add -A
git commit -m $1
git push origin $branch_name

