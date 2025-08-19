# Liftoff

Liftoff is an open-source automation tool designed to streamline infrastructure provisioning and configuration management. It supports integration with Ansible, Terraform, and custom templates, making it suitable for DevOps workflows and cloud-native environments.

Liftoff supportd alternative cloud providers like Digital Ocean, Hetzner Cloud etc.

General idea is that user will provide a configuration file, and Liftoff will use provided data to provision cloud infrastructure using Terraform and configure applications using Ansible. An example configuration file might look like this:

```
template-repo: https://github.com/bitshifted/liftoff-templates.git
template-version: d5f3810dd52f046ad35c5565e4740285ea85a139
terraform:
  providers:
    - hcloud
ansible:
  inventory-file: inventory
  playbook-file: postgresql.yaml 
variables:
  default:
    hcloud_server_config:
      server_name: postgresql
      location: fsn1
      ssh_key_names: 
        - hcloud_playground
      enable_public_ipv6: false
      enable_public_ipv4: false
      labels:
        internal: true
        TTL: 3h
      keep_disk: true
      enable_backups: false
      private_networks:
        - infra-net
    ansible_user: ansible
    ansible_ssh_private_key: ~/.ssh/hcloud_playground
    ansible_host_key_checking: no
    ansible_ssh_public_keys:
      - "fromfile:~/.ssh/hcloud_playground.pub" 
    bastion_address: nat-gw.astrolabs.io
    gateway_ip: 10.1.0.1
    dns_servers: "8.8.8.8"
tags:
  stack: basic-infra
  application: nat-gateway
```

This will provision a server in Hetzner Cloud and install and configure Postgresql on it. Currently available templates can be found in [liftoff-template](https://github.com/bitshifted/liftoff-templates) directory.

## Features
- Infrastructure provisioning with Terraform
- Configuration management with Ansible
- Template processing for custom configuration files
- Extensible CLI for automation tasks
- Logging and error reporting

## Getting Started

Installation packages for Linux, Windows and Mac OS are avaialble in Releases page. Terraform and Ansible must be installed on the machine in order to use Liftoff.

Alternativaly, Docker image is available, which contains all necessary tools installed. You can use it if you don't want install addiitonal software.


## Usage
Run the CLI tool to provision infrastructure or configure environments:

```bash
./liftoff --config-file path/to/config.yaml setup
```

To cleanup provisioned infrastructure, run:

```bash
./liftoff --config-file path/to/config.yaml teardown
```

Additional options:

```
-terraform-path=STRING       Path to Terraform binary
--playbook-bin-path=STRING    Path to ansible-playbook binary
--config-file=STRING          Path to configuration file
--enable-debug                Enable debug logging

```

## Documentation
See the [docs](./docs) directory or the project wiki for detailed usage and configuration examples.

## License
This project is licensed under the MPL-2.0 License.
