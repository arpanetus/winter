package main

import (
	"github.com/steplems/winter/core"
	"os"
	"os/exec"
	"path"
	"plugin"
)

func createWinterDir(winterDir string) error {
	err := exec.Command("mkdir", winterDir).Run()
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
		err = os.RemoveAll(winterDir)
		if err != nil {
			log.Err("Couldn't remove temp .winter directory:", err)
			return err
		}
		return createWinterDir(winterDir)
	}
}

func winterBuild()  {
	pwd, err := os.Getwd()
	if err != nil {
		log.Err("Couldn't get execution path:", err)
		return
	}

	winterFile := path.Join(pwd, cli_config_file)
	winterDir := path.Join(pwd, cli_config_dir)
	winterFileSo := path.Join(winterDir, cli_config_file_so)

	err = updateWinterDir(winterDir)
	if err != nil {
		return
	}

	log.Info(winterFileSo)
	log.Info(winterFile)
	err = exec.Command("go", "build", "-buildmode=plugin", "-o", winterFileSo, winterFile).Run()
	if err != nil {
		log.Err("Couldn't build plugin from file winter.go:", err)
		return
	}

	winterPlugin, err := plugin.Open(winterFileSo)
	if err != nil {
		log.Err("Couldn't open winter.go file:", err)
		return
	}

	app, err := winterPlugin.Lookup(cli_app_config)
	if err != nil {
		log.Err("Couldn't get App config:", err)
		return
	}

	createMainFile(app.(*core.App), winterDir)
}

func createMainFile(config *core.App, buildPath string) {
	file, err := os.Create(path.Join(buildPath, "main.go"))
	if err != nil {
		log.Err("Couldn't create main.go file:", err)
		return
	}

	file.Write()
}
