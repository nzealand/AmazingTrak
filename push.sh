#!/usr/bin/env bash
set -euo pipefail

# Make sure we are inside a git repo
git rev-parse --is-inside-work-tree >/dev/null 2>&1 || {
  echo "Not inside a git repository."
  exit 1
}

echo "Current changes:"
git status --short

if [ -z "$(git status --porcelain)" ]; then
  echo "No changes to commit."
  exit 0
fi

echo
read -r -p "Commit message: " msg

if [ -z "$msg" ]; then
  echo "Commit message cannot be empty."
  exit 1
fi

git add .
git commit -m "$msg"

# Push to existing upstream if configured, otherwise ask what to use
if git rev-parse --abbrev-ref --symbolic-full-name '@{u}' >/dev/null 2>&1; then
  git push
else
  branch="$(git branch --show-current)"
  read -r -p "No upstream set. Remote name [origin]: " remote
  remote="${remote:-origin}"
  git push -u "$remote" "$branch"
fi

echo "Done."
