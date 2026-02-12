terraform {
  required_providers {
    isardvdi = {
      source = "tknika/isardvdi"
    }
  }
}

provider "isardvdi" {
  endpoint    = var.isard_endpoint
  auth_method = var.auth_method
  token       = var.isard_token
}

# Ejemplo 1: Buscar usuarios por nombre
data "isardvdi_users" "search_by_name" {
  name_filter = "admin"
}

output "users_by_name" {
  description = "Usuarios que contienen 'admin' en su nombre"
  value = [for user in data.isardvdi_users.search_by_name.users : {
    id       = user.id
    name     = user.name
    username = user.username
    email    = user.email
  }]
}

# Ejemplo 2: Obtener todos los usuarios de una categoría
data "isardvdi_users" "category_users" {
  category_id = "default"
}

output "category_users_count" {
  description = "Número de usuarios en la categoría default"
  value       = length(data.isardvdi_users.category_users.users)
}

# Ejemplo 3: Buscar usuarios por rol
data "isardvdi_users" "managers" {
  role = "manager"
}

output "manager_list" {
  description = "Lista de managers"
  value = [for user in data.isardvdi_users.managers.users : {
    name         = user.name
    email        = user.email
    category     = user.category_name
    group        = user.group_name
  }]
}

# Ejemplo 4: Buscar un usuario específico
data "isardvdi_users" "specific_user" {
  name_filter = var.search_username
}

output "specific_user_details" {
  description = "Detalles del usuario buscado"
  value = length(data.isardvdi_users.specific_user.users) > 0 ? {
    id                     = data.isardvdi_users.specific_user.users[0].id
    name                   = data.isardvdi_users.specific_user.users[0].name
    username               = data.isardvdi_users.specific_user.users[0].username
    email                  = data.isardvdi_users.specific_user.users[0].email
    role                   = data.isardvdi_users.specific_user.users[0].role_name
    category               = data.isardvdi_users.specific_user.users[0].category_name
    group                  = data.isardvdi_users.specific_user.users[0].group_name
    active                 = data.isardvdi_users.specific_user.users[0].active
    email_verified         = data.isardvdi_users.specific_user.users[0].email_verified
  } : null
}

# Ejemplo 5: Filtrar usuarios activos de un grupo específico
data "isardvdi_users" "active_group_users" {
  group_id = var.group_id
  active   = true
}

output "active_users_in_group" {
  description = "Usuarios activos en el grupo especificado"
  value = [for user in data.isardvdi_users.active_group_users.users : {
    name     = user.name
    username = user.username
    role     = user.role_name
  }]
}

# Ejemplo 6: Combinar múltiples filtros
data "isardvdi_users" "filtered_users" {
  category_id = "default"
  role        = "user"
  active      = true
}

output "filtered_users_summary" {
  description = "Resumen de usuarios filtrados"
  value = {
    total_count = length(data.isardvdi_users.filtered_users.users)
    users       = [for user in data.isardvdi_users.filtered_users.users : user.username]
  }
}

# Ejemplo 7: Obtener IDs de todos los usuarios para otros usos
output "all_user_ids" {
  description = "Lista de todos los IDs de usuarios encontrados"
  value = [for user in data.isardvdi_users.search_by_name.users : user.id]
}

# Ejemplo 8: Verificar si un usuario existe
output "user_exists" {
  description = "Verifica si el usuario buscado existe"
  value = length(data.isardvdi_users.specific_user.users) > 0
}
