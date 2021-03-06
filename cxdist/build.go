package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func build() {
	// checkout code
	if err := cloneRepo(sourceGit, flagBranch, buildDir); err != nil {
		log.Fatal(err)
	}

	// build it
	if err := buildApp(buildDir, flagVersion); err != nil {
		log.Fatal(err)
	}
}

func cloneRepo(repo, branch, dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	if _, err := cmd("git", "clone", "-b", branch, repo, dir); err != nil {
		return err
	}
	return nil
}

func buildApp(dir, ver string) error {
	if _, err := cmd("goxc", "-include=''", "-bc=linux,windows,darwin,freebsd", "-pv="+ver, "-d="+publishDir, "-main-dirs-exclude=Godeps,testdata,_project,vendor,cxdist", "-n=cx", "-wd="+dir, "-tasks-=go-test,go-vet"); err != nil {
		//if _, err := cmd("gox", "-ldflags=\"-X main.VERSION="+ver+"\"", "-osarch=\"darwin/386 darwin/amd64 linux/386 linux/amd64 windows/386 windows/amd64\"", "-output=\""+publishDir+"/{{.Dir}}_{{.OS}}_{{.Arch}}\""); err != nil {
		return err
	}
	return nil
}

func cmd(arg ...string) ([]byte, error) {
	log.Println(strings.Join(arg, " "))
	cmd := exec.Command(arg[0], arg[1:]...)
	cmd.Stderr = os.Stderr
	return cmd.Output()
}
