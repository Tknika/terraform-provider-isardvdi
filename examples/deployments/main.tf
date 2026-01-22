terraform {
  required_providers {
    isard = {
      source = "tknika/isard"
    }
  }
}

provider "isard" {
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
resource "isard_deployment" "group_deployment" {
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
resource "isard_deployment" "performance_deployment" {
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
resource "isard_deployment" "user_deployment" {
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
resource "isard_deployment" "category_deployment" {
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
resource "isard_deployment" "network_deployment" {
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
resource "isard_deployment" "multi_group_deployment" {
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

# Outputs para mostrar información de los deployments creados
output "group_deployment_id" {
  description = "ID del deployment del grupo DevOps"
  value       = isard_deployment.group_deployment.id
}

output "performance_deployment_id" {
  description = "ID del deployment de alto rendimiento"
  value       = isard_deployment.performance_deployment.id
}

output "user_deployment_id" {
  description = "ID del deployment de usuarios VIP"
  value       = isard_deployment.user_deployment.id
}

output "category_deployment_id" {
  description = "ID del deployment de estudiantes"
  value       = isard_deployment.category_deployment.id
}
