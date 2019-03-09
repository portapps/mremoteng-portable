//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/portapps/portapps"
	"github.com/portapps/portapps/pkg/utl"
)

var (
	app *App
)

func init() {
	var err error

	// Init app
	if app, err = New("mremoteng-portable", "mRemoteNG"); err != nil {
		Log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	utl.CreateFolder(app.DataPath)
	app.Process = utl.PathJoin(app.AppPath, "mRemoteNG.exe")
	app.Args = []string{
		fmt.Sprintf("/cons:%s", utl.FormatWindowsPath(utl.PathJoin(app.DataPath, "confCons.xml"))),
	}

	// mRemoteNG.settings check
	dSettings := utl.PathJoin(app.DataPath, "mRemoteNG.settings")
	aSettings := utl.PathJoin(app.AppPath, "mRemoteNG.settings")
	if !utl.Exists(dSettings) {
		if err := utl.WriteToFile(dSettings, `<?xml version="1.0" encoding="utf-8"?><settings/>`); err != nil {
			Log.Fatal().Err(err).Msg("Cannot write to mRemoteNG.settings")
		}
	}
	_ = os.Remove(aSettings)
	if err := os.Symlink(dSettings, aSettings); err != nil {
		Log.Fatal().Err(err).Msg("Cannot create symlink to mRemoteNG.settings")
	}

	// extApps.xml exists in data ? Create symlink and remove old one
	dExtApps := utl.PathJoin(app.DataPath, "extApps.xml")
	aExtApps := utl.PathJoin(app.AppPath, "extApps.xml")
	if utl.Exists(dExtApps) {
		_ = os.Remove(aExtApps)
		if err := os.Symlink(dExtApps, aExtApps); err != nil {
			Log.Fatal().Err(err).Msg("Cannot create symlink to extApps.xml")
		}
	}

	// pnlLayout.xml exists in data ? Create symlink and remove old one
	dPnlLayout := utl.PathJoin(app.DataPath, "pnlLayout.xml")
	aPnlLayout := utl.PathJoin(app.AppPath, "pnlLayout.xml")
	if utl.Exists(dPnlLayout) {
		_ = os.Remove(aPnlLayout)
		if err := os.Symlink(dPnlLayout, aPnlLayout); err != nil {
			Log.Fatal().Err(err).Msg("Cannot create symlink to pnlLayout.xml")
		}
	}

	// On exit
	defer func() {
		// confCons.xml copy back on close if not exists in data
		dConfCons := utl.PathJoin(app.DataPath, "confCons.xml")
		aConfCons := utl.PathJoin(app.AppPath, "confCons.xml")
		if !utl.Exists(dConfCons) && utl.Exists(aConfCons) {
			if err := utl.CopyFile(aConfCons, dConfCons); err != nil {
				Log.Error().Err(err).Msg("Cannot copy confCons.xml")
			}
			_ = os.Remove(aConfCons)
			if err := os.Symlink(dConfCons, aConfCons); err != nil {
				Log.Error().Err(err).Msg("Cannot create symlink to confCons.xml")
			}
		}
		oldConfConsFiles, _ := filepath.Glob(utl.PathJoin(app.AppPath, "confCons*"))
		for _, oldConfConsFile := range oldConfConsFiles {
			if err := os.Remove(oldConfConsFile); err != nil {
				Log.Error().Err(err).Msg("Cannot remove old confCons file")
			}
		}

		// extApps.xml handling on close
		if !utl.Exists(dExtApps) && utl.Exists(aExtApps) {
			if err := utl.CopyFile(aExtApps, dExtApps); err != nil {
				Log.Error().Err(err).Msg("Cannot copy extApps.xml")
			}
			_ = os.Remove(aExtApps)
			if err := os.Symlink(dExtApps, aExtApps); err != nil {
				Log.Error().Err(err).Msg("Cannot create symlink to extApps.xml")
			}
		}

		// pnlLayout.xml handling on close
		if !utl.Exists(dPnlLayout) && utl.Exists(aPnlLayout) {
			if err := utl.CopyFile(aPnlLayout, dPnlLayout); err != nil {
				Log.Error().Err(err).Msg("Cannot copy pnlLayout.xml")
			}
			_ = os.Remove(aPnlLayout)
			if err := os.Symlink(dPnlLayout, aPnlLayout); err != nil {
				Log.Error().Err(err).Msg("Cannot create symlink to pnlLayout.xml")
			}
		}
	}()

	app.Launch(os.Args[1:])
}
