// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package config

type AnsibleConfig struct {
	InventoryFile string `yaml:"inventory-file"`
	PlaybookFile  string `yaml:"playbook-file"`
}
