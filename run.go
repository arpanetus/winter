package main

import (
	"github.com/steplems/winter/core"
	"os"
	"os/exec"
	"path"
	"plugin"
)

func createWinterDir(winterDir string) error {
	err := os.Mkdir(winterDir, os.ModeDir)
	if err != nil {
		log.Err("Couldn't create temp .winter directory:", err)
		return err
	}
	return nil
}

func updateWinterDir(winterDir string) error {
	if _, err := os.Stat(winterDir); os.IsNotExist(err) {
		return createWinterDir(winterDir)
	} else {
		err = os.Remove(winterDir)
		if err != nil {
			log.Err("Couldn't remove temp .winter directory:", err)
			return err
		}
		return createWinterDir(winterDir)
	}
}

func winterRun()  {
	pwd, err := os.Getwd()
	if err != nil {
		log.Err("Can not get execution path:", err)
		return
	}

	winterFile := path.Join(pwd, cli_config_file)
	winterDir := path.Join(pwd, cli_config_dir)
	winterFileSo := path.Join(winterDir, cli_config_file_so)

	err = updateWinterDir(winterDir)
	if err != nil {
		return
	}

	err = exec.Command("go", "build", "-buildmode=plugin", "-o", winterFileSo, winterFile).Run()
	if err != nil {
		log.Err("Couldn't build plugin from file winter.go:", err)
		return
	}

	winterPlugin, err := plugin.Open(winterFileSo)
	if err != nil {
		log.Err("Can not open winter.go file:", err)
		return
	}

	app, err := winterPlugin.Lookup(cli_app_config)
	if err != nil {
		log.Err("Can not get App config:", err)
		return
	}

	config := app.(core.App)
	log.Info(config)
}
