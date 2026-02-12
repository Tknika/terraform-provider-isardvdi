# Ejemplo usando autenticación por token

terraform {
  required_providers {
    isardvdi = {
      source = "registry.terraform.io/tknika/isardvdi"
    }
  }
}

# Configuración usando token JWT
provider "isardvdi" {
  endpoint     = "localhost"
  auth_method  = "token"
  cathegory_id = "default"
  token        = var.isard_token
}

# Obtener templates
data "isardvdi_templates" "all" {}

# Crear desktop
resource "isardvdi_vm" "token_auth_example" {
  name        = "desktop-token-auth"
  description = "Desktop creado usando autenticación por token"
  template_id = data.isardvdi_templates.all.templates[0].id
}

# Outputs
output "desktop_creado" {
  value = {
    id          = isardvdi_vm.token_auth_example.id
    name        = isardvdi_vm.token_auth_example.name
    description = isardvdi_vm.token_auth_example.description
  }
}
