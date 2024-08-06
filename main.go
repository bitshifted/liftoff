package main

import (
	"github.com/alecthomas/kong"
	"github.com/bitshifted/easycloud/log"
)

var CLI struct {
	TerraformPath   string `help:"Path to Terraform binary"`
	PlaybookBinPath string `help:"Path to ansible-playbook binary"`
}

func main() {
	log.Init()
	kong.Parse(&CLI)
}
