package main

import (
	"io"
	"log"
	"net/http"
	"os/exec"

	"os"
	"runtime"
)

var ext string

type Config struct {
	Input        string `json:"Input"`
	Flag         string `json:"Flag"`
	DefaultValue string `json:"DefaultValue"`
}

func main() {
	if err := GetAliyunCliPkg(); err != nil {
		log.Fatal(err)
	}
	if err := ConfigAliyunCli(); err != nil {
		log.Fatal(err)
	}
}

func ConfigAliyunCli() error {
	var configs = []Config{
		{Input: os.Getenv("MOD"), Flag: "--mode", DefaultValue: "AK"},
		{Input: os.Getenv("PROFILE"), Flag: "--profile", DefaultValue: ""},
		{Input: os.Getenv("LANGUAGE"), Flag: "--language", DefaultValue: "zh"},
		{Input: os.Getenv("REGION"), Flag: "--region", DefaultValue: ""},
		{Input: os.Getenv("CONFIG-PATH"), Flag: "--config-path", DefaultValue: ""},
		{Input: os.Getenv("ACCESS-KEY-ID"), Flag: "--access-key-id", DefaultValue: ""},
		{Input: os.Getenv("ACCESS-KEY-SECRET"), Flag: "--access-key-secret", DefaultValue: ""},
		{Input: os.Getenv("STS-TOKEN"), Flag: "--sts-token", DefaultValue: ""},
		{Input: os.Getenv("RAM-ROLE-NAME"), Flag: "--ram-role-name", DefaultValue: ""},
		{Input: os.Getenv("RAM-ROLE-ARN"), Flag: "--ram-role-arn", DefaultValue: ""},
		{Input: os.Getenv("ROLE-SESSION-NAME"), Flag: "--role-session-name", DefaultValue: ""},
		{Input: os.Getenv("PRIVATE-KEY"), Flag: "--private-key", DefaultValue: ""},
		{Input: os.Getenv("KEY-PAIR-NAME"), Flag: "--key-pair-name", DefaultValue: ""},
		{Input: os.Getenv("READ-TIMEOUT"), Flag: "--read-timeout", DefaultValue: ""},
		{Input: os.Getenv("CONNECT-TIMEOUT"), Flag: "--connect-timeout", DefaultValue: ""},
		{Input: os.Getenv("RETRY-COUNT"), Flag: "--retry-count", DefaultValue: ""},
		{Input: os.Getenv("SKIP-SECURE-VERIFY"), Flag: "--skip-secure-verify", DefaultValue: ""},
		{Input: os.Getenv("EXPIRED-SECONDS"), Flag: "--expired-seconds", DefaultValue: ""},
		{Input: os.Getenv("SECURE"), Flag: "--secure", DefaultValue: ""},
		{Input: os.Getenv("FORCE"), Flag: "--force", DefaultValue: ""},
		{Input: os.Getenv("ENDPOINT"), Flag: "--endpoint", DefaultValue: ""},
		{Input: os.Getenv("VERSION"), Flag: "--version", DefaultValue: ""},
		{Input: os.Getenv("HEADER"), Flag: "--header", DefaultValue: ""},
		{Input: os.Getenv("BODY"), Flag: "--body", DefaultValue: ""},
		{Input: os.Getenv("PAGER"), Flag: "--pager", DefaultValue: ""},
		{Input: os.Getenv("OUTPUT"), Flag: "--output", DefaultValue: ""},
		{Input: os.Getenv("WAITER"), Flag: "--waiter", DefaultValue: ""},
		{Input: os.Getenv("DRYRUN"), Flag: "--dryrun", DefaultValue: ""},
		{Input: os.Getenv("QUIET"), Flag: "--quiet", DefaultValue: ""},
	}
	var config string
	for _, v := range configs {
		if v.Input != "" {
			config = config + " " + v.Flag + " " + v.Input
			continue
		}
		config = config + " " + v.Flag + " " + v.DefaultValue
	}
	cmd := exec.Command("./aliyun", "configure", "set", config)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func GetAliyunCliPkg() error {
	version := os.Getenv("VERSION")
	if version == "" {
		version = "3.0.55"
	}
	var system string
	switch runtime.GOOS {
	case "win32":
		system = "windows"
		ext = "zip"
	case "darwin":
		system = "macosx"
		ext = "tgz"
	case "linux":
		system = "linux"
		ext = "tgz"
	}
	var url = "https://github.com/aliyun/aliyun-cli/releases/download/v" + version + "/aliyun-cli-" + system + "-" + version + "-amd64." + ext
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if ext == "tgz" {
		file, err := os.Create("./aliyun.tgz")
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return err
		}
		cmd := exec.Command("tar", "-zxvf", "./aliyun.tgz", "-C", ".")
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	}
	file, err := os.Create("./aliyun.zip")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	cmd := exec.Command("unzip", "aliyun.zip", "-d", "cli")
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
