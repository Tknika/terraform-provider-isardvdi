terraform {
  required_providers {
    isard = {
      source = "terraform.local/local/isard"
    }
  }
}

provider "isard" {
  host     = "https://localhost"
  username = "admin"
  password = "IsardVDI"
  insecure = true
}

# Ejemplo 1: Deployment con ISO específico
resource "isard_deployment" "training_with_tools" {
  name         = "Training Environment"
  description  = "Entorno de formación con herramientas"
  template_id  = "template-id"
  desktop_name = "Training Desktop"
  visible      = true
  
  vcpus  = 2
  memory = 4
  
  # Adjuntar ISO de herramientas
  isos = ["tools-iso-id"]
  
  allowed {
    roles = ["user"]
  }
}

# Ejemplo 2: Deployment con múltiples ISOs
resource "isard_deployment" "lab_deployment" {
  name         = "Lab Deployment"
  description  = "Deployment de laboratorio con múltiples ISOs"
  template_id  = "template-id"
  desktop_name = "Lab Desktop"
  visible      = true
  
  vcpus  = 4
  memory = 8
  
  # Múltiples ISOs para diferentes propósitos
  isos = [
    "ubuntu-iso-id",
    "drivers-iso-id",
    "tools-iso-id"
  ]
  
  allowed {
    roles  = ["advanced", "user"]
    groups = ["students-group-id"]
  }
}

# Ejemplo 3: Deployment usando data source
data "isard_medias" "installation_media" {
  name_filter = "Windows Server"
  kind        = "iso"
  status      = "Downloaded"
}

data "isard_medias" "virtio_drivers" {
  name_filter = "VirtIO"
  kind        = "iso"
  status      = "Downloaded"
}

resource "isard_deployment" "windows_deployment" {
  name         = "Windows Server Deployment"
  description  = "Deployment de Windows Server con drivers"
  template_id  = "windows-template-id"
  desktop_name = "Windows Server"
  visible      = true
  
  vcpus  = 4
  memory = 16
  
  # Combinar múltiples data sources
  isos = concat(
    length(data.isard_medias.installation_media.medias) > 0 ? [data.isard_medias.installation_media.medias[0].id] : [],
    length(data.isard_medias.virtio_drivers.medias) > 0 ? [data.isard_medias.virtio_drivers.medias[0].id] : []
  )
  
  allowed {
    roles = ["admin", "advanced"]
  }
}

# Ejemplo 4: Deployment condicional basado en disponibilidad de media
data "isard_medias" "required_iso" {
  name_filter = "Required Software"
  kind        = "iso"
  status      = "Downloaded"
}

resource "isard_deployment" "conditional_deployment" {
  count = length(data.isard_medias.required_iso.medias) > 0 ? 1 : 0
  
  name         = "Conditional Deployment"
  description  = "Solo se crea si el ISO requerido está disponible"
  template_id  = "template-id"
  desktop_name = "Desktop"
  visible      = true
  
  vcpus  = 2
  memory = 4
  
  isos = [data.isard_medias.required_iso.medias[0].id]
  
  allowed {
    roles = ["user"]
  }
}

# Outputs
output "training_deployment_id" {
  value = isard_deployment.training_with_tools.id
}

output "windows_deployment_info" {
  value = {
    id              = isard_deployment.windows_deployment.id
    name            = isard_deployment.windows_deployment.name
    isos_attached   = length(isard_deployment.windows_deployment.isos)
  }
}

output "conditional_created" {
  value = length(isard_deployment.conditional_deployment) > 0 ? "Yes" : "No - Required ISO not available"
}
