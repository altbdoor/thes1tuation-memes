#!/bin/bash

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd $script_dir

find_cmd="find"
sed_cmd="sed"

if [[ $(uname) == "Darwin" ]]; then
    find_cmd="gfind"
    sed_cmd="gsed"
fi

original_ifs=$IFS
IFS=$'\n'

html_files=$("$find_cmd" ../_site/ -type f -name '*.html' -print)
for file in $html_files; do
    $sed_cmd -i -E 's/(<a\s+|<a$)/& rel="noopener noreferrer" /g' $file
done

IFS=$original_ifs
