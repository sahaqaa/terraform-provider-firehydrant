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

data "firehydrant_service" "qqq" {
  id = "qqq"
}
output "qqq" {
  value = data.firehydrant_service.qqq.name
}

# data "firehydrant_service" "lake-service" {
#   id = "lake-service"
# }

# output "lake-service" {
#   value = data.firehydrant_service.lake-service.name
# }

resource "firehydrant_service" "puddles-service" {
  name        = "Puddles!!!"
  description = "new puddles"
}
