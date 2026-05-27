#!/usr/bin/env sh
# goccline installer — single-command setup.
#
#   curl -fsSL https://raw.githubusercontent.com/andhikapraa/goccline/main/install.sh | sh
#
# What this does:
#   1. Detects your platform (darwin/linux × arm64/amd64).
#   2. Downloads the matching prebuilt binary from the latest GitHub release.
#   3. Places it in /opt/homebrew/bin (macOS) or /usr/local/bin (Linux).
#   4. Writes a sensible default Config.toml to ~/.claude/goccline/.
#   5. Wires Claude Code's settings.json statusLine to the new binary
#      (with a backup of your existing config).
#
# Flags:
#   --no-wire     Don't touch ~/.claude/settings.json.
#   --prefix DIR  Install binary to DIR instead of the default.
#   --version V   Install a specific version (defaults to latest).
set -eu

REPO="andhikapraa/goccline"
VERSION="${GOCCLINE_VERSION:-latest}"
WIRE=1
PREFIX=""

while [ $# -gt 0 ]; do
  case "$1" in
    --no-wire) WIRE=0 ;;
    --prefix) PREFIX="$2"; shift ;;
    --prefix=*) PREFIX="${1#--prefix=}" ;;
    --version) VERSION="$2"; shift ;;
    --version=*) VERSION="${1#--version=}" ;;
    -h|--help)
      sed -n '2,18p' "$0" | sed 's/^# //;s/^#$//'
      exit 0 ;;
    *) echo "unknown flag: $1" >&2; exit 1 ;;
  esac
  shift
done

# ── Platform detection ───────────────────────────────────────────────────
os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"
case "$arch" in
  arm64|aarch64) arch="arm64" ;;
  x86_64|amd64)  arch="amd64" ;;
  *) echo "unsupported arch: $arch" >&2; exit 1 ;;
esac
case "$os" in
  darwin|linux) ;;
  *) echo "unsupported OS: $os (try go install or build from source)" >&2; exit 1 ;;
esac

target="goccline_${os}_${arch}"

# ── Default install prefix ───────────────────────────────────────────────
if [ -z "$PREFIX" ]; then
  if [ "$os" = "darwin" ] && [ -d /opt/homebrew/bin ] && [ -w /opt/homebrew/bin ]; then
    PREFIX="/opt/homebrew/bin"
  elif [ -w /usr/local/bin ]; then
    PREFIX="/usr/local/bin"
  else
    PREFIX="$HOME/.local/bin"
    mkdir -p "$PREFIX"
  fi
fi

# ── Resolve download URL ─────────────────────────────────────────────────
if [ "$VERSION" = "latest" ]; then
  url="https://github.com/${REPO}/releases/latest/download/${target}.tar.gz"
else
  url="https://github.com/${REPO}/releases/download/${VERSION}/${target}.tar.gz"
fi

# ── Download + extract ───────────────────────────────────────────────────
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

printf '\033[36m→\033[0m Downloading goccline (%s) from %s\n' "$target" "$url"
if ! curl -fsSL "$url" -o "$tmpdir/goccline.tar.gz"; then
  echo "" >&2
  echo "Couldn't fetch the prebuilt binary. Possible reasons:" >&2
  echo "  - No release yet for $VERSION × ${os}_${arch}" >&2
  echo "  - Network error" >&2
  echo "" >&2
  echo "Fallback: install with Go (requires Go 1.23+):" >&2
  echo "  go install github.com/${REPO}@${VERSION#v}" >&2
  exit 1
fi

tar -xzf "$tmpdir/goccline.tar.gz" -C "$tmpdir"

# Use install(1) where available (strips xattrs cleanly on macOS) — avoids
# the gotcha where ~/.local/bin downloads inherit com.apple.provenance and
# Gatekeeper SIGKILLs the binary inside Claude Code.
if command -v install >/dev/null 2>&1; then
  install -m 0755 "$tmpdir/goccline" "$PREFIX/goccline"
else
  cp "$tmpdir/goccline" "$PREFIX/goccline"
  chmod 0755 "$PREFIX/goccline"
fi

printf '\033[32m✓\033[0m Installed: %s/goccline\n' "$PREFIX"

# ── Default config ───────────────────────────────────────────────────────
cfg="$HOME/.claude/goccline/Config.toml"
mkdir -p "$(dirname "$cfg")"
if [ ! -f "$cfg" ]; then
  cat > "$cfg" <<'EOF'
# goccline — Claude Code statusline
# https://github.com/andhikapraa/goccline
# Edit freely; changes take effect on the next render.

[theme]
name = "catppuccin"  # catppuccin | garden | ocean | classic

[features]
show_commits = true
show_version = true
show_submodules = true

[prayer]
# location_mode = "gps"   # uncomment + `brew install corelocationcli` for
                           # auto-detect on macOS
latitude = 0.0
longitude = 0.0
method = 2                 # AlAdhan calc method (2 = ISNA)
school = 1                 # 1 = Hanafi, 0 = Shafi/Maliki/Hanbali
# location_label = ""      # auto-derived via reverse geocoding when empty
                           # AND coordinates are set

[display]
lines = 4

[display.line1]
components = ["repo_info", "session_info"]

[display.line2]
components = ["model_info", "session_mode", "version_info", "time_display"]

[display.line3]
components = ["cost_today", "cost_live", "burn_rate", "code_productivity"]

[display.line4]
components = ["context_window", "total_tokens", "usage_limits", "reset_timer"]

# To enable Islamic features, set [prayer] coordinates above and add:
# [display.line5]
# components = ["hijri_calendar", "wellness"]
# [display.line6]
# components = ["prayer_icon", "prayer_times_only"]
# separator = " "
EOF
  printf '\033[32m✓\033[0m Config:    %s\n' "$cfg"
else
  printf '\033[33m·\033[0m Config kept: %s (already exists)\n' "$cfg"
fi

# ── Wire into Claude Code ────────────────────────────────────────────────
if [ "$WIRE" -eq 1 ]; then
  settings="$HOME/.claude/settings.json"
  if [ -f "$settings" ]; then
    cp "$settings" "${settings}.bak.$(date +%s)"
    python3 - "$settings" "$PREFIX/goccline" <<'PY'
import json, sys
path, bin_path = sys.argv[1], sys.argv[2]
with open(path) as f:
    s = json.load(f)
s["statusLine"] = {
    "type": "command",
    "command": bin_path,
    "padding": 0,
    "refreshInterval": 10,
}
with open(path, "w") as f:
    json.dump(s, f, indent=2)
    f.write("\n")
PY
    printf '\033[32m✓\033[0m Wired:     %s (backup created)\n' "$settings"
  else
    mkdir -p "$(dirname "$settings")"
    cat > "$settings" <<EOF
{
  "statusLine": {
    "type": "command",
    "command": "$PREFIX/goccline",
    "padding": 0,
    "refreshInterval": 10
  }
}
EOF
    printf '\033[32m✓\033[0m Wired:     %s (created)\n' "$settings"
  fi
fi

# ── Done ────────────────────────────────────────────────────────────────
echo ""
echo "🎉 goccline is ready."
echo ""
echo "Next steps:"
echo "  1. Reload Claude Code (or just wait 10s for the next statusline refresh)"
echo "  2. Edit ${cfg} to taste"
if [ "$os" = "darwin" ]; then
  echo "  3. Optional — for fresh prayer coordinates while travelling:"
  echo "       brew install corelocationcli"
  echo "       # then set [prayer].location_mode = \"gps\" in Config.toml"
fi
echo ""
echo "Docs: https://github.com/${REPO}"
