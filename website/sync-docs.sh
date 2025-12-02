#!/bin/sh
# Sync docs to Docusaurus: add frontmatter + fix links
# Files in /docs: 01-quickstart.md, 02-api.md, etc.
set -e

# Detect script location and make paths absolute
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
ROOT_DIR=$(cd "$SCRIPT_DIR/.." && pwd)

src="$ROOT_DIR/docs"
dst="$SCRIPT_DIR/docs"
root="$ROOT_DIR"

echo "🔄 Syncing docs..."
echo ""

# Create destination directory if it doesn't exist
mkdir -p "$dst"

# Process all numbered docs
for f in "$src"/[0-9][0-9]-*.md; do
  [ -f "$f" ] || continue
  base=$(basename "$f")
  # Extract position number (strip leading zeros to avoid octal interpretation)
  pos_raw="${base%%-*}"
  # Remove leading zeros: 01 -> 1, 08 -> 8, 10 -> 10
  pos=$(echo "$pos_raw" | sed 's/^0*//')
  [ -z "$pos" ] && pos=0  # Handle case of "00"
  out="${base#*-}"                            # Output filename without NN- prefix
  title=$(sed -n 's/^# //p;q' "$f")          # Extract first # title

  # Generate file: frontmatter + content (skip line 1) + fix links
  # Strip NN- from (NN-file.md) → (file.md)
  { printf -- '---\nsidebar_position: %s\ntitle: %s\n---\n\n' "$pos" "$title"
    sed '1d; s|(\([0-9][0-9]-\)\([^)]*\.md\))|(\2)|g' "$f"
  } > "$dst/$out"

  echo "✅ $out (pos: $pos)"
done

# CONTRIBUTING.md from root (position 12)
if [ -f "$root/CONTRIBUTING.md" ]; then
  title=$(sed -n 's/^# //p;q' "$root/CONTRIBUTING.md")
  { printf -- '---\nsidebar_position: 12\ntitle: %s\n---\n\n' "$title"
    sed '1d' "$root/CONTRIBUTING.md"
  } > "$dst/contributing.md"
  echo "✅ contributing.md (pos: 12)"
fi

# TODO.md from root (position 13)
if [ -f "$root/TODO.md" ]; then
  title=$(sed -n 's/^# //p;q' "$root/TODO.md")
  { printf -- '---\nsidebar_position: 13\ntitle: %s\n---\n\n' "$title"
    sed '1d; s|CONTRIBUTING\.md|contributing.md|g' "$root/TODO.md"
  } > "$dst/todo.md"
  echo "✅ todo.md (pos: 13)"
fi

# Copy OpenAPI specs to static/ (symlinks don't work in Docker)
if [ -f "$ROOT_DIR/internal/api/docs/swagger.yaml" ]; then
  cp "$ROOT_DIR/internal/api/docs/swagger.yaml" "$SCRIPT_DIR/static/openapi.yaml"
  echo "✅ openapi.yaml → static/"
fi

if [ -f "$ROOT_DIR/internal/api/docs/swagger.json" ]; then
  cp "$ROOT_DIR/internal/api/docs/swagger.json" "$SCRIPT_DIR/static/openapi.json"
  echo "✅ openapi.json → static/"
fi

echo ""
echo "✨ Done! /docs/*.md → /website/docs/*.md"
