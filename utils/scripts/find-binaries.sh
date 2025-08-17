#!/usr/bin/env bash
# detect binary files in your repo using file --mime-encoding

echo "Scanning for binary files in the repo…"

binary_files=()

# Loop over all tracked files
while IFS= read -r -d '' file; do
  # Ask file(1) for the encoding; e.g. "utf-8" vs. "binary"
  encoding=$(file --mime-encoding -b "$file")
  if [[ $encoding == binary ]]; then
    binary_files+=("$file")
  fi
done < <(git ls-files -z)

if [ ${#binary_files[@]} -eq 0 ]; then
  echo "✅ No binary files detected."
else
  echo "⚠️  Binary files detected:"
  printf '  %s\n' "${binary_files[@]}"
  exit 1
fi

