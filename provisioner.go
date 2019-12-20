package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/provisioner"
	"github.com/hashicorp/packer/template/interpolate"
)

type guestOSTypeConfig struct {
	runCommand     string
	installCommand string
	stagingDir     string
}

var guestOSTypeConfigs = map[string]guestOSTypeConfig{
	provisioner.UnixOSType: {
		runCommand: "cd {{.StagingDir}} \u0026\u0026 {{if .Sudo}}sudo {{end}}/usr/local/bin/terraform init \u0026\u0026 " +
			"{{if .Sudo}}sudo {{end}}/usr/local/bin/terraform apply -auto-approve",
		installCommand: "curl https://releases.hashicorp.com/terraform/{{.Version}}/terraform_{{.Version}}_linux_amd64.zip " +
			"-so /tmp/terraform.zip \u0026\u0026 " +
			"{{if .Sudo}}sudo {{end}}unzip -d /usr/local/bin/ /tmp/terraform.zip",
		stagingDir: "/tmp/packer-terraform",
	},
	provisioner.WindowsOSType: {
		runCommand: "cd {{.StagingDir}} \u0026\u0026 C:/Windows/Temp/terraform init \u0026\u0026 " +
			"C:/Windows/Temp/terraform apply -auto-approve",
		installCommand: "Invoke-WebRequest -Uri 'https://releases.hashicorp.com/terraform/{{.Version}}/terraform_{{.Version}}_windows_amd64.zip' " +
			"-OutFile 'C:/Windows/Temp/terraform.zip' \u0026\u0026 " +
			"Expand-Archive 'C:/Windows/Temp/terraform.zip' -DestinationPath 'C:/Windows/Temp/terraform'",
		stagingDir: "C:/Windows/Temp/packer-terraform",
	},
}

// Config struct containing variables
type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	Version        string `mapstructure:"version"`
	CodePath       string `mapstructure:"code_path"`
	RunCommand     string `mapstructure:"run_command"`
	InstallCommand string `mapstructure:"install_command"`
	StagingDir     string `mapstructure:"staging_dir"`
	PreventSudo    bool   `mapstructure:"prevent_sudo"`
	Variables      map[string]interface{}
	GuestOSType    string `mapstructure:"guest_os_type"`

	ctx interpolate.Context
}

// Provisioner is the interface to install and run Terraform
type Provisioner struct {
	config            Config
	guestOSTypeConfig guestOSTypeConfig
	guestCommands     *provisioner.GuestCommands
}

// RunTemplate for temp storage of interpolation vars
type RunTemplate struct {
	StagingDir string
	Sudo       bool
	Version    string
}

// ConfigSpec gets the FlatMapStructure for HCL Support.
func (p *Provisioner) ConfigSpec() hcldec.ObjectSpec { return p.config.FlatMapstructure().HCL2Spec() }

// Prepare parses the config and get everything ready
func (p *Provisioner) Prepare(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
	}, raws...)
	if err != nil {
		return err
	}

	if p.config.GuestOSType == "" {
		p.config.GuestOSType = provisioner.DefaultOSType
	}
	p.config.GuestOSType = strings.ToLower(p.config.GuestOSType)
	p.guestOSTypeConfig = guestOSTypeConfigs[p.config.GuestOSType]

	if p.config.StagingDir == "" {
		p.config.StagingDir = p.guestOSTypeConfig.stagingDir
	}

	_, err = os.Stat(p.config.CodePath)
	if err != nil {
		return fmt.Errorf("bad source '%s': %s", p.config.CodePath, err)
	}

	if p.config.Version == "" {
		p.config.Version = "0.12.16"
	}

	if p.config.InstallCommand == "" {
		p.config.InstallCommand = p.guestOSTypeConfig.installCommand
	}

	if p.config.RunCommand == "" {
		p.config.RunCommand = p.guestOSTypeConfig.runCommand
	}

	p.config.Variables, err = p.processVariables()
	if err != nil {
		return fmt.Errorf("Error processing Variables in JSON: %s", err)
	}

	return nil
}

// Provision does the work of installing Terraform and running it on the remote
func (p *Provisioner) Provision(_ context.Context, ui packer.Ui, comm packer.Communicator, _ map[string]interface{}) error {
	ui.Say("Provisioning with Terraform...")

	ui.Message("Uploading Code")
	if err := p.uploadDirectory(ui, comm, p.config.StagingDir, p.config.CodePath); err != nil {
		return fmt.Errorf("Error uploading code: %s", err)
	}

	if err := p.createTfvars(ui, comm); err != nil {
		return fmt.Errorf("Error generating tfvars: %s", err)
	}

	ui.Message("Installing Terraform")
	p.config.ctx.Data = &RunTemplate{
		StagingDir: p.config.StagingDir,
		Version:    p.config.Version,
		Sudo:       !p.config.PreventSudo,
	}
	command, err := interpolate.Render(p.config.InstallCommand, &p.config.ctx)
	if err != nil {
		return fmt.Errorf("Error rendering Template: %s", err)
	}
	if err := p.runCommand(ui, comm, command); err != nil {
		return fmt.Errorf("Error running Terraform: %s", err)
	}

	ui.Message("Running Terraform")
	command, err = interpolate.Render(p.config.RunCommand, &p.config.ctx)
	if err != nil {
		return fmt.Errorf("Error rendering Template: %s", err)
	}
	if err := p.runCommand(ui, comm, command); err != nil {
		return fmt.Errorf("Error installing Terraform: %s", err)
	}

	return nil
}

func (p *Provisioner) runCommand(ui packer.Ui, comm packer.Communicator, command string) error {
	var out, outErr bytes.Buffer
	cmd := &packer.RemoteCmd{
		Command: command,
		Stdin:   nil,
		Stdout:  &out,
		Stderr:  &outErr,
	}

	ctx := context.TODO()
	if err := cmd.RunWithUi(ctx, comm, ui); err != nil {
		return err
	}
	if cmd.ExitStatus() != 0 {
		return fmt.Errorf("non-zero exit status")
	}
	return nil
}

func (p *Provisioner) uploadDirectory(ui packer.Ui, comm packer.Communicator, dst string, src string) error {
	// Make sure there is a trailing "/" so that the directory isn't
	// created on the other side.
	if src[len(src)-1] != '/' {
		src = src + "/"
	}

	return comm.UploadDir(dst, src, nil)
}

func (p *Provisioner) processVariables() (map[string]interface{}, error) {
	jsonBytes, err := json.Marshal(p.config.Variables)
	if err != nil {
		panic(err)
	}

	// Process the bytes with the template processor
	p.config.ctx.Data = nil
	jsonBytesProcessed, err := interpolate.Render(string(jsonBytes), &p.config.ctx)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonBytesProcessed), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *Provisioner) createTfvars(ui packer.Ui, comm packer.Communicator) error {
	ui.Message("Creating tfvars file")

	template := "{{ range $key, $value := . }}" +
		"{{ $key }} = \"{{ $value }}\"\n" +
		"{{ end }}"

	p.config.ctx.Data = &p.config.Variables
	tfvarsData, err := interpolate.Render(template, &p.config.ctx)
	if err != nil {
		return err
	}

	// Upload the bytes
	remotePath := filepath.ToSlash(filepath.Join(p.config.StagingDir, "terraform.auto.tfvars"))
	if err := comm.Upload(remotePath, strings.NewReader(tfvarsData), nil); err != nil {
		return err
	}

	return nil
}
