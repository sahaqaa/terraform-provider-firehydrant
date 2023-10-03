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

# data "firehydrant_service" "monolith" {
#   id = "d975cf38-922e-4192-b366-dc632e73cd3d"
# }

# data "firehydrant_service" "lake-service" {
#   id = "lake-service"
# }

# output "monolith" {
#   value = data.firehydrant_service.monolith.name
# }

# output "lake-service" {
#   value = data.firehydrant_service.lake-service.name
# }

# resource "firehydrant_service" "puddles-service" {
#   name        = "Puddles!!!"
#   description = "new puddles"
# }
