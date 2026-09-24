terraform {
  required_providers {
    isardvdi = {
      source = "tknika/isardvdi"
    }
  }
}

provider "isardvdi" {
  endpoint    = "isard.example.com"
  auth_method = "token"
  token       = var.isard_token
}

variable "isard_token" {
  description = "Token de autenticación para Isard VDI"
  type        = string
  sensitive   = true
}

variable "template_id" {
  description = "ID de la plantilla a utilizar"
  type        = string
  default     = "your-template-uuid"
}

# Deployment que arranca automáticamente sus desktops y espera a que reporten
# la IP de wireguard antes de terminar el apply. Las VMs quedan encendidas al
# finalizar; el proveedor no las apaga ni orquesta ningún provisionamiento
# posterior (por ejemplo, Ansible).
resource "isardvdi_deployment" "with_ip" {
  name         = "Deployment con IP wireguard"
  description  = "Desktops accesibles por IP wireguard nada más aplicar"
  template_id  = var.template_id
  desktop_name = "Desktop automatizable por SSH gracias a su IP"
  visible      = true

  wait_for_ip             = true
  ip_wait_timeout_seconds = 300

  allowed = {
    roles      = []
    categories = []
    groups     = ["students-group-id"]
    users      = []
  }
}

# Este output se puede consumir, por ejemplo, para generar un inventory de
# Ansible fuera de Terraform (plantilla externa, script, etc.). El proveedor
# no genera el inventory ni ejecuta Ansible por sí mismo.
output "desktops" {
  description = "Desktops del deployment con su IP (wireguard), user y user_name"
  value       = isardvdi_deployment.with_ip.desktops
}
