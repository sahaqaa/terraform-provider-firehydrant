terraform {
  required_providers {
    firehydrant = {
      version = "0.1.0"
      source  = "firehydrant.com/twitch/firehydrant"
    }
  }
}

provider "firehydrant" {
}

data "firehydrant_service" "monolith" {
  id = "d975cf38-922e-4192-b366-dc632e73cd3d"
}

output "xxx" {
  value = data.firehydrant_service.monolith
}
