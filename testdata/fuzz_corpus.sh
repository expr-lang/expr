#!/bin/bash

inputfile=./fuzz_corpus.txt
zipname=fuzz_expr_seed_corpus.zip

if [ ! -f "$inputfile" ]; then
  echo "Error: File $inputfile not found!"
  exit 1
fi

lineno=1
while IFS= read -r line; do
  echo "$line" >"file_${lineno}.txt"
  ((lineno++))
done <"$inputfile"

zip "$zipname" file_*.txt

rm file_*.txt

echo "Done!"
