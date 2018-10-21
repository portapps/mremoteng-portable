//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"os"
	"path/filepath"

	. "github.com/portapps/portapps"
)

func init() {
	Papp.ID = "mremoteng-portable"
	Papp.Name = "mRemoteNG"
	Init()
}

func main() {
	Papp.AppPath = AppPathJoin("app")
	Papp.DataPath = CreateFolder(AppPathJoin("data"))
	Papp.Process = PathJoin(Papp.AppPath, "mRemoteNG.exe")
	Papp.Args = []string{"/cons:" + FormatWindowsPath(PathJoin(Papp.DataPath, "confCons.xml"))}
	Papp.WorkingDir = Papp.AppPath

	// mRemoteNG.settings check
	dSettings := PathJoin(Papp.DataPath, "mRemoteNG.settings")
	aSettings := PathJoin(Papp.AppPath, "mRemoteNG.settings")
	if !Exists(dSettings) {
		if err := WriteToFile(dSettings, `<?xml version="1.0" encoding="utf-8"?><settings/>`); err != nil {
			Log.Error("Cannot write to mRemoteNG.settings: ", err)
		}
	}
	os.Remove(aSettings)
	if err := os.Symlink(dSettings, aSettings); err != nil {
		Log.Error("Cannot create symlink to mRemoteNG.settings: ", err)
	}

	// extApps.xml exists in data ? Create symlink and remove old one
	dExtApps := PathJoin(Papp.DataPath, "extApps.xml")
	aExtApps := PathJoin(Papp.AppPath, "extApps.xml")
	if Exists(dExtApps) {
		os.Remove(aExtApps)
		if err := os.Symlink(dExtApps, aExtApps); err != nil {
			Log.Error("Cannot create symlink to extApps.xml: ", err)
		}
	}

	// pnlLayout.xml exists in data ? Create symlink and remove old one
	dPnlLayout := PathJoin(Papp.DataPath, "pnlLayout.xml")
	aPnlLayout := PathJoin(Papp.AppPath, "pnlLayout.xml")
	if Exists(dPnlLayout) {
		os.Remove(aPnlLayout)
		if err := os.Symlink(dPnlLayout, aPnlLayout); err != nil {
			Log.Error("Cannot create symlink to pnlLayout.xml: ", err)
		}
	}

	Launch(os.Args[1:])

	// confCons.xml copy back on close if not exists in data
	dConfCons := PathJoin(Papp.DataPath, "confCons.xml")
	aConfCons := PathJoin(Papp.AppPath, "confCons.xml")
	if !Exists(dConfCons) && Exists(aConfCons) {
		if err := CopyFile(aConfCons, dConfCons); err != nil {
			Log.Error("Cannot copy confCons.xml: ", err)
		}
		os.Remove(aConfCons)
		if err := os.Symlink(dConfCons, aConfCons); err != nil {
			Log.Error("Cannot create symlink to confCons.xml: ", err)
		}
	}
	oldConfConsFiles, _ := filepath.Glob(PathJoin(Papp.AppPath, "confCons*"))
	for _, oldConfConsFile := range oldConfConsFiles {
		if err := os.Remove(oldConfConsFile); err != nil {
			Log.Error("Cannot remove old confCons file: ", err)
		}
	}

	// extApps.xml handling on close
	if !Exists(dExtApps) && Exists(aExtApps) {
		if err := CopyFile(aExtApps, dExtApps); err != nil {
			Log.Error("Cannot copy extApps.xml: ", err)
		}
		os.Remove(aExtApps)
		if err := os.Symlink(dExtApps, aExtApps); err != nil {
			Log.Error("Cannot create symlink to extApps.xml: ", err)
		}
	}

	// pnlLayout.xml handling on close
	if !Exists(dPnlLayout) && Exists(aPnlLayout) {
		if err := CopyFile(aPnlLayout, dPnlLayout); err != nil {
			Log.Error("Cannot copy pnlLayout.xml: ", err)
		}
		os.Remove(aPnlLayout)
		if err := os.Symlink(dPnlLayout, aPnlLayout); err != nil {
			Log.Error("Cannot create symlink to pnlLayout.xml: ", err)
		}
	}
}
