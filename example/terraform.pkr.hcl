packer {
  required_plugins {
    terraform = {
      version = ">= 0.0.2"
      source  = "github.com/tristanmorgan/terraform"
    }
    docker = {
      version = "~> 1"
      source  = "github.com/hashicorp/docker"
    }
  }
}

source "docker" "amazon" {
  commit = true
  image  = "amazonlinux:2023"
}

build {
  sources = ["source.docker.amazon"]

  provisioner "shell" {
    inline = [
      "yum install -y unzip"
    ]
  }

  provisioner "terraform" {
    code_path    = "./tfcode"
    prevent_sudo = "true"
    variable_string = jsonencode({
      consul_server_node = false
      nomad_alt_url      = "https://example.com"
    })
  }

  post-processor "docker-tag" {
    repository = "tristanmorgan/packer-tf-test"
    tags       = ["latest"]
  }
}
