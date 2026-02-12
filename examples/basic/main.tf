# Ejemplo básico de uso del provider Isard

terraform {
  required_providers {
    isardvdi = {
      source = "registry.terraform.io/tknika/isardvdi"
    }
  }
}

# Configuración del provider
provider "isardvdi" {
  endpoint     = "localhost"
  auth_method  = "form"
  cathegory_id = "default"
  username     = "admin"
  password     = "IsardVDI"
}

# Obtener templates disponibles
data "isardvdi_templates" "all" {}

# Crear un desktop simple
resource "isardvdi_vm" "ejemplo" {
  name        = "desktop-ejemplo"
  description = "Desktop de ejemplo creado con Terraform"
  template_id = data.isardvdi_templates.all.templates[0].id
}

# Outputs
output "desktop_id" {
  description = "ID del desktop creado"
  value       = isardvdi_vm.ejemplo.id
}

output "desktop_info" {
  description = "Información del desktop"
  value = {
    id          = isardvdi_vm.ejemplo.id
    name        = isardvdi_vm.ejemplo.name
    vcpus       = isardvdi_vm.ejemplo.vcpus
    memory      = isardvdi_vm.ejemplo.memory
    template_id = isardvdi_vm.ejemplo.template_id
  }
}

output "templates_disponibles" {
  description = "Lista de todos los templates disponibles"
  value       = data.isardvdi_templates.all.templates
}
