// Code generated by "mapstructure-to-hcl2 -type Config"; DO NOT EDIT.

package main

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

// FlatConfig is an auto-generated flat version of Config.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatConfig struct {
	PackerBuildName     *string           `mapstructure:"packer_build_name" cty:"packer_build_name" hcl:"packer_build_name"`
	PackerBuilderType   *string           `mapstructure:"packer_builder_type" cty:"packer_builder_type" hcl:"packer_builder_type"`
	PackerCoreVersion   *string           `mapstructure:"packer_core_version" cty:"packer_core_version" hcl:"packer_core_version"`
	PackerDebug         *bool             `mapstructure:"packer_debug" cty:"packer_debug" hcl:"packer_debug"`
	PackerForce         *bool             `mapstructure:"packer_force" cty:"packer_force" hcl:"packer_force"`
	PackerOnError       *string           `mapstructure:"packer_on_error" cty:"packer_on_error" hcl:"packer_on_error"`
	PackerUserVars      map[string]string `mapstructure:"packer_user_variables" cty:"packer_user_variables" hcl:"packer_user_variables"`
	PackerSensitiveVars []string          `mapstructure:"packer_sensitive_variables" cty:"packer_sensitive_variables" hcl:"packer_sensitive_variables"`
	Version             *string           `mapstructure:"version" cty:"version" hcl:"version"`
	CodePath            *string           `mapstructure:"code_path" cty:"code_path" hcl:"code_path"`
	RunCommand          *string           `mapstructure:"run_command" cty:"run_command" hcl:"run_command"`
	InstallCommand      *string           `mapstructure:"install_command" cty:"install_command" hcl:"install_command"`
	StagingDir          *string           `mapstructure:"staging_dir" cty:"staging_dir" hcl:"staging_dir"`
	PreventSudo         *bool             `mapstructure:"prevent_sudo" cty:"prevent_sudo" hcl:"prevent_sudo"`
	VariableString      *string           `mapstructure:"variable_string" cty:"variable_string" hcl:"variable_string"`
	GuestOSType         *string           `mapstructure:"guest_os_type" cty:"guest_os_type" hcl:"guest_os_type"`
}

// FlatMapstructure returns a new FlatConfig.
// FlatConfig is an auto-generated flat version of Config.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*Config) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatConfig)
}

// HCL2Spec returns the hcl spec of a Config.
// This spec is used by HCL to read the fields of Config.
// The decoded values from this spec will then be applied to a FlatConfig.
func (*FlatConfig) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"packer_build_name":          &hcldec.AttrSpec{Name: "packer_build_name", Type: cty.String, Required: false},
		"packer_builder_type":        &hcldec.AttrSpec{Name: "packer_builder_type", Type: cty.String, Required: false},
		"packer_core_version":        &hcldec.AttrSpec{Name: "packer_core_version", Type: cty.String, Required: false},
		"packer_debug":               &hcldec.AttrSpec{Name: "packer_debug", Type: cty.Bool, Required: false},
		"packer_force":               &hcldec.AttrSpec{Name: "packer_force", Type: cty.Bool, Required: false},
		"packer_on_error":            &hcldec.AttrSpec{Name: "packer_on_error", Type: cty.String, Required: false},
		"packer_user_variables":      &hcldec.AttrSpec{Name: "packer_user_variables", Type: cty.Map(cty.String), Required: false},
		"packer_sensitive_variables": &hcldec.AttrSpec{Name: "packer_sensitive_variables", Type: cty.List(cty.String), Required: false},
		"version":                    &hcldec.AttrSpec{Name: "version", Type: cty.String, Required: false},
		"code_path":                  &hcldec.AttrSpec{Name: "code_path", Type: cty.String, Required: false},
		"run_command":                &hcldec.AttrSpec{Name: "run_command", Type: cty.String, Required: false},
		"install_command":            &hcldec.AttrSpec{Name: "install_command", Type: cty.String, Required: false},
		"staging_dir":                &hcldec.AttrSpec{Name: "staging_dir", Type: cty.String, Required: false},
		"prevent_sudo":               &hcldec.AttrSpec{Name: "prevent_sudo", Type: cty.Bool, Required: false},
		"variable_string":            &hcldec.AttrSpec{Name: "variable_string", Type: cty.String, Required: false},
		"guest_os_type":              &hcldec.AttrSpec{Name: "guest_os_type", Type: cty.String, Required: false},
	}
	return s
}
