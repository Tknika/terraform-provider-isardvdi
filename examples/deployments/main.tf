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

# Ejemplo 1: Deployment básico para un grupo
resource "isardvdi_deployment" "group_deployment" {
  name         = "Deployment Equipo DevOps"
  description  = "Desktops para el equipo de DevOps"
  template_id  = var.template_id
  desktop_name = "Desktop DevOps"
  visible      = false

  allowed {
    groups = ["devops-group-uuid"]
  }
}

# Ejemplo 2: Deployment con hardware personalizado
resource "isardvdi_deployment" "performance_deployment" {
  name         = "Deployment Alto Rendimiento"
  description  = "Deployment con recursos mejorados para tareas intensivas"
  template_id  = var.template_id
  desktop_name = "Desktop Performance"
  visible      = true
  
  vcpus      = 4
  memory     = 16.0

  allowed {
    groups = ["performance-group-uuid"]
  }
}

# Ejemplo 3: Deployment para usuarios específicos
resource "isardvdi_deployment" "user_deployment" {
  name         = "Deployment Usuarios VIP"
  description  = "Desktops para usuarios específicos"
  template_id  = var.template_id
  desktop_name = "Desktop VIP"
  visible      = true

  vcpus  = 2
  memory = 8.0

  allowed {
    users = [
      "user-uuid-1",
      "user-uuid-2",
      "user-uuid-3"
    ]
  }
  
  user_permissions = ["viewer", "desktop_start", "desktop_stop"]
}

# Ejemplo 4: Deployment para una categoría completa
resource "isardvdi_deployment" "category_deployment" {
  name         = "Deployment Estudiantes"
  description  = "Desktops para todos los estudiantes de la institución"
  template_id  = var.template_id
  desktop_name = "Desktop Estudiante"
  visible      = true

  allowed {
    categories = ["students-category-uuid"]
  }
}

# Ejemplo 5: Deployment con configuración de red personalizada
resource "isardvdi_deployment" "network_deployment" {
  name         = "Deployment Red Avanzada"
  description  = "Deployment con interfaces de red personalizadas"
  template_id  = var.template_id
  desktop_name = "Desktop Network"
  
  vcpus  = 2
  memory = 4.0
  
  interfaces = [
    "interface-uuid-1",
    "interface-uuid-2"
  ]

  allowed {
    groups = ["network-admin-group-uuid"]
  }
}

# Ejemplo 6: Deployment para múltiples grupos
resource "isardvdi_deployment" "multi_group_deployment" {
  name         = "Deployment Multi-Grupo"
  description  = "Deployment compartido entre varios grupos"
  template_id  = var.template_id
  desktop_name = "Desktop Compartido"
  visible      = true

  allowed {
    groups = [
      "group-uuid-1",
      "group-uuid-2",
      "group-uuid-3"
    ]
  }
}

# Ejemplo 7: Deployment con viewers personalizados (solo web)
resource "isardvdi_deployment" "web_viewers_deployment" {
  name         = "Deployment Acceso Web"
  description  = "Deployment con solo viewers basados en navegador"
  template_id  = var.template_id
  desktop_name = "Desktop Web"
  visible      = true
  
  # Solo viewers web, sin necesidad de instalar clientes
  viewers = ["browser_vnc", "browser_rdp"]

  allowed {
    groups = ["remote-workers-uuid"]
  }
}

# Ejemplo 8: Deployment con SPICE para máximo rendimiento
resource "isardvdi_deployment" "spice_deployment" {
  name         = "Deployment SPICE"
  description  = "Deployment optimizado para alto rendimiento con SPICE"
  template_id  = var.template_id
  desktop_name = "Desktop SPICE"
  visible      = true
  
  vcpus  = 4
  memory = 8.0
  
  # Solo SPICE para máximo rendimiento
  viewers = ["file_spice"]

  allowed {
    groups = ["design-team-uuid"]
  }
}

# Ejemplo 9: Deployment con todos los viewers disponibles
resource "isardvdi_deployment" "full_access_deployment" {
  name         = "Deployment Acceso Completo"
  description  = "Deployment con todos los métodos de acceso disponibles"
  template_id  = var.template_id
  desktop_name = "Desktop Full"
  visible      = true
  
  # Todos los viewers disponibles
  viewers = [
    "browser_rdp",
    "browser_vnc",
    "file_rdpgw",
    "file_rdpvpn",
    "file_spice"
  ]

  allowed {
    users = ["admin-user-uuid"]
  }
}

# Outputs para mostrar información de los deployments creados
output "group_deployment_id" {
  description = "ID del deployment del grupo DevOps"
  value       = isardvdi_deployment.group_deployment.id
}

output "performance_deployment_id" {
  description = "ID del deployment de alto rendimiento"
  value       = isardvdi_deployment.performance_deployment.id
}

output "user_deployment_id" {
  description = "ID del deployment de usuarios VIP"
  value       = isardvdi_deployment.user_deployment.id
}

output "category_deployment_id" {
  description = "ID del deployment de estudiantes"
  value       = isardvdi_deployment.category_deployment.id
}

output "web_viewers_deployment_id" {
  description = "ID del deployment con viewers web"
  value       = isardvdi_deployment.web_viewers_deployment.id
}

output "spice_deployment_id" {
  description = "ID del deployment con SPICE"
  value       = isardvdi_deployment.spice_deployment.id
}
