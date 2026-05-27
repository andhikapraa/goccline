# goccline

Fast, native Claude Code statusline — built to replace the bash one that
took 4–5 seconds to render and timed out with cost components enabled.

```
~/Development/goccline (main) ✅ │ 🔗 a35a0566
🧠 Opus 4.7 (1M context) │ Style: default │ CC:2.1.152 │ GL:0.7.3 │ 🕐 19:47
📈 $6780 30d │ 📊 $1764 7d │ 📅 $174 today │ 💰 $83.39 │ 🔥 $55.82/hr │ Est: $279 (510M) │ +3436 / -427
🧠 61% (606K/1.0M) │ 1K in / 816K out │ ⏱ Used 5h:19% • 7d:15% │ ⏱ 5h:3h12m • 7d:Wed 3AM
🕌 10 Dhu al-Hijjah 1447 │ 🌓 │ May 27 2026 │ 📍 Tasikmalaya, Jawa Barat │ ☕ Coding 0m/45m
🕌 Fajr 04:51 ✓ │ Dhuhr 11:44 ✓ │ Asr 15:58 ✓ │ Maghrib 17:36 ✓ │ Isha 18:37 ✓ │ Fajr 04:50 (9h 50m)
```

**~400 ms** with a full 6-line render. Single Go binary, no runtime deps.

## Features

- 30+ components: repo info, model, cost (today/week/month/repo/live/burn-rate/projection), token usage, rate limits, Islamic prayer times, Hijri calendar, wellness timer, code productivity.
- Reads Claude Code's native `rate_limits` and `context_window` JSON fields (v2.1.80+); falls back to transcript parsing for older versions.
- Parallel JSONL parser dedup-by-request-id for cost calculations.
- Catppuccin / Garden / Ocean / Classic themes (truecolor).
- macOS GPS via CoreLocationCLI for fresh prayer locations when travelling. Falls back to manual lat/lon.
- Reverse-geocodes GPS coords via Nominatim (cached 24h) for the location label.
- Bash-statusline-compatible TOML config format — migrate by copying your existing `Config.toml`.

## Install

Requires Go 1.23+.

```bash
go install github.com/andhikapraa/goccline@latest
```

Or build from source:

```bash
git clone https://github.com/andhikapraa/goccline ~/src/goccline
cd ~/src/goccline
go build -o /opt/homebrew/bin/goccline .
```

Optional for GPS prayer times (macOS):

```bash
brew install corelocationcli
# Then enable Location Services for CoreLocationCLI in System Settings.
```

## Wire into Claude Code

Add to `~/.claude/settings.json`:

```json
{
  "statusLine": {
    "type": "command",
    "command": "/opt/homebrew/bin/goccline",
    "padding": 0,
    "refreshInterval": 10
  }
}
```

## Configure

`~/.claude/goccline/Config.toml`:

```toml
[theme]
name = "catppuccin"  # or "garden", "ocean", "classic"

[features]
show_commits = true

[prayer]
location_mode = "gps"   # or "manual"
latitude = -6.2088      # used when mode = "manual" or as GPS fallback
longitude = 106.8456
method = 2              # AlAdhan calc method (2 = ISNA)
school = 1              # 1 = Hanafi, 0 = Shafi/Maliki/Hanbali
# location_label = ""   # auto-derived via reverse geocoding when empty

[display]
lines = 6

[display.line1]
components = ["repo_info", "session_info"]

[display.line2]
components = ["model_info", "session_mode", "version_info", "time_display"]

[display.line3]
components = ["cost_repo", "cost_monthly", "cost_weekly", "cost_today",
              "cost_live", "burn_rate", "block_projection", "code_productivity"]

[display.line4]
components = ["context_window", "total_tokens", "usage_limits", "reset_timer"]

[display.line5]
components = ["hijri_calendar", "wellness"]

[display.line6]
components = ["prayer_icon", "prayer_times_only"]
separator = " "
```

## Available components

| Group | Components |
|---|---|
| Identity | `repo_info`, `session_info`, `model_info`, `session_mode`, `agent_display`, `version_info`, `version_display` |
| Git | `commits`, `submodules` |
| Time | `time_display` |
| Cost | `cost_live`, `cost_today`, `cost_daily`, `cost_weekly`, `cost_monthly`, `cost_repo`, `burn_rate`, `block_projection`, `code_productivity` |
| Tokens | `context_window`, `context_alert`, `total_tokens`, `token_usage` |
| Limits | `usage_limits`, `reset_timer` |
| Islamic | `prayer_icon`, `prayer_times`, `prayer_times_only`, `hijri_calendar`, `location_display` |
| Wellness | `wellness` |

## Acknowledgements

- Pricing table and JSONL parsing approach inspired by [backstabslash/goccc](https://github.com/backstabslash/goccc).
- Component layout and Islamic features ported from [rz1989s/claude-code-statusline](https://github.com/rz1989s/claude-code-statusline).
- Prayer times via [AlAdhan API](https://aladhan.com/prayer-times-api).
- Reverse geocoding via [Nominatim](https://nominatim.org/).
- Hijri calendar via [hablullah/go-hijri](https://github.com/hablullah/go-hijri).

## License

MIT
