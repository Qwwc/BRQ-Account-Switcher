<p align="center">
  <img src="build/appicon.png" width="128" alt="BRQ">
</p>

<h1 align="center">BRQ Account Switcher</h1>

<p align="center">A fast account switcher for Steam, Epic Games, Riot Games, Ubisoft, Discord, Battle.net and more, with BRQ's violet look and animations.</p>

## Based on TcNo Account Switcher

BRQ Account Switcher is a modified version of
[TcNo Account Switcher](https://github.com/TCNOco/TcNo-Acc-Switcher) by TroubleChute (Wesley Pyburn),
licensed under the **GNU General Public License v3.0**. BRQ is independent and is not affiliated with
or endorsed by TcNo or TroubleChute. The original README is kept as [README.upstream.md](README.upstream.md).

Under the GPL-3.0, this program and its source come with no warranty, and anyone who receives a build
is entitled to the complete corresponding source under the same licence (see [LICENSE](LICENSE)).

### What BRQ changes

- Name, logo, icons and default "BRQ Violet" theme taken from the BRQ logo.
- Chamfered, glowing platform tiles with an animated energy frame, staggered entrance, hover lift and light sweep.
- Living violet background, animated title-bar light and metallic logo shine.
- Own data folder (`%AppData%\BRQ Account Switcher`) and own `brq:` link scheme, so it never clashes with an installed TcNo.
- Telemetry, crash uploads, feedback, Discord status and the self-updater are switched off, because they point at upstream servers.
  Set `UpdateRepo` / `TelemetryEnabled` in `internal/api/api.go` and the Discord `clientID` in `internal/discordrpc/manager.go` to turn them on with your own services.

All animations follow the in-app "animations" setting and Windows' reduced-motion preference.

## Building from source (Windows)

Requirements: Go 1.27.3+, Node.js 24, pnpm 11.0.8 and Wails CLI 3.1.0
(`go install github.com/wailsapp/wails/v3/cmd/wails3@v3.1.0`).

```
wails3 task windows:build
```

The app is written to `bin/BRQ-Acc-Switcher.exe`.

BRQ helper scripts:

- `node tools/brq/gen_logo.mjs` regenerates the vector logo (`frontend/public/img/BRQ_Logo_Flat.svg`, `frontend/src/lib/brqLogoPaths.ts`).
- `go run ./tools/brq/icons` rebuilds the app/tray icons from `build/brq_logo_source.webp`.

## Licence

GPL-3.0. See [LICENSE](LICENSE), [OPEN_SOURCE_LICENSES.txt](OPEN_SOURCE_LICENSES.txt) and [DISCLAIMER.md](DISCLAIMER.md).
