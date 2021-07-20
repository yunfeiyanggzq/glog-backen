#!/bin/bash

echo "
╔═╗┌┬┐┌─┐┬─┐┌┬┐  ┌─┐┬ ┬┌─┐┌─┐┬┌─
╚═╗ │ ├─┤├┬┘ │   │  ├─┤├┤ │  ├┴┐
╚═╝ ┴ ┴ ┴┴└─ ┴   └─┘┴ ┴└─┘└─┘┴ ┴
"
# 解压位置   程序名字  程序md5
# $1        $2         $3

# shellcheck disable=SC2164
cd "$1"

# shellcheck disable=SC2164
checkFUnc() {
  code=0
  go build
   # shellcheck disable=SC2181
   if [ $? -ne 0 ]; then
    echo "============================================================================================"
    echo "go build failed"
    echo "============================================================================================"
    code=1
  else
    echo "============================================================================================"
    echo "go build success"
    echo "============================================================================================"
  fi

  # shellcheck disable=SC2034
  build_hash_code=$(md5 "$1")
  # shellcheck disable=SC1009
  if [ "$build_hash_code" = "$2" ]; then
    echo "============================================================================================"
    echo "md5 equal"
    echo "============================================================================================"
  else
    echo "============================================================================================"
    echo "md5 not equal"
    echo "============================================================================================"
    code=1
  fi

  #检查英语语法
  # shellcheck disable=SC2038
  find . -name "*" | xargs misspell -error
  # shellcheck disable=SC2181
  if [ $? -ne 0 ]; then
    echo "============================================================================================"
    echo "use opensource tool client9/misspell to correct commonly misspelled English words failed"
    echo "============================================================================================"
    code=1
  else
    echo "============================================================================================"
    echo "use opensource tool client9/misspell to correct commonly misspelled English words success"
    echo "============================================================================================"
  fi

  #检查md格式
  find ./ -name "*.md" | grep -v commandline | grep -v .github | grep -v swagger | grep -v api | xargs mdl -r ~MD009,~MD002,~MD007,~MD010,~MD013,~MD024,~MD026,~MD029,~MD033,~MD034,~MD036,~MD04
  # shellcheck disable=SC2181
  if [ $? -ne 0 ]; then
    echo "============================================================================================"
    echo "use markdownLint v0.5.0 to lint markdown file failed"
    echo "============================================================================================"
  else
    echo "============================================================================================"
    echo "use markdownLint v0.5.0 to lint markdown file success"
    echo "============================================================================================"
  fi

  #检查链接
  set +e
  for name in $(find . -name \*.md | grep -v CHANGELOG); do
    # shellcheck disable=SC2086
    if [ -f $name ]; then
      markdown-link-check -q -v $name
      # shellcheck disable=SC2181
      if [ $? -ne 0 ]; then
        echo "============================================================================================"
        echo "use markdown-link-check to check links in markdown files failed"
        echo "============================================================================================"
      else
        echo "============================================================================================"
        echo "use markdown-link-check to check links in markdown files success"
        echo "============================================================================================"
      fi
    fi
  done

  #检查yaml格式
  # shellcheck disable=SC2038
  find . -name "*.yml" | xargs yamllint -d "{extends: default, rules: {line-length: disable}}"
  # shellcheck disable=SC2181
  if [ $? -ne 0 ]; then
    echo "============================================================================================"
    echo "use yaml lint tool to check the format of all yaml files failed"
    echo "============================================================================================"
  else
    echo "============================================================================================"
    echo "use yaml lint tool to check the format of all yaml files success"
    echo "============================================================================================"
  fi

  #检查shell
  # shellcheck disable=SC2038
  find . -name "*.sh" | xargs shellcheck
  # shellcheck disable=SC2181
  if [ $? -ne 0 ]; then
    echo "============================================================================================"
    echo "use ShellCheck to check the validates of shell scripts in pouch repo failed"
    echo "============================================================================================"
  else
    echo "============================================================================================"
    echo "use ShellCheck to check the validates of shell scripts in pouch repo success"
    echo "============================================================================================"
  fi

  #检查go格式
  golangci-lint run
  # shellcheck disable=SC2181
  if [ $? -ne 0 ]; then
    echo "============================================================================================"
    echo "use golangci-lint to check go format failed"
    echo "============================================================================================"
  else
    echo "============================================================================================"
    echo "use golangci-lint to check go format success"
    echo "============================================================================================"
  fi

  return $code
}

checkFUnc "$2" "$3"
# shellcheck disable=SC2181
if [ $? -ne 0 ]; then
  echo "00失败00"
else
  echo "11成功11"
fi
