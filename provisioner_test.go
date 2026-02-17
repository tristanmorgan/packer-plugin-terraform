package main

import (
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"testing"
)

func TestProvisioner_Impl(t *testing.T) {
	var raw any = &Provisioner{}
	if _, ok := raw.(packer.Provisioner); !ok {
		t.Fatalf("must be a Provisioner")
	}
}
