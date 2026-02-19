---
page_title: "isardvdi_users Data Source - terraform-provider-isardvdi"
subcategory: ""
description: |-
  Fetches the list of users from Isard VDI with optional filtering capabilities.
---

# isardvdi_users (Data Source)

Este data source permite obtener información sobre usuarios en Isard VDI. Puedes filtrar usuarios por nombre, categoría, grupo, rol y estado activo.

## Ejemplo de Uso

```terraform
# Buscar usuarios por nombre
data "isardvdi_users" "admins" {
  name_filter = "admin"
}

# Obtener el ID del primer usuario encontrado
output "first_admin_id" {
  value = data.isardvdi_users.admins.users[0].id
}

# Filtrar usuarios por categoría
data "isardvdi_users" "category_users" {
  category_id = "default"
}

# Filtrar usuarios activos de un grupo específico
data "isardvdi_users" "active_group_users" {
  group_id = "default-default-default"
  active   = true
}

# Filtrar por rol
data "isardvdi_users" "managers" {
  role = "manager"
}

# Combinar múltiples filtros
data "isardvdi_users" "specific_users" {
  name_filter = "john"
  category_id = "default"
  role        = "user"
  active      = true
}
```

## Argumentos

* `name_filter` - (Opcional) Filtro para buscar usuarios por nombre. Realiza una búsqueda parcial sin distinción entre mayúsculas y minúsculas.
* `category_id` - (Opcional) ID de la categoría para filtrar usuarios.
* `group_id` - (Opcional) ID del grupo para filtrar usuarios.
* `role` - (Opcional) Rol del usuario para filtrar (admin, manager, advanced, user).
* `active` - (Opcional) Estado activo del usuario (true/false).

## Atributos

* `id` - Identificador del data source.
* `users` - Lista de usuarios que coinciden con los filtros aplicados. Cada usuario tiene los siguientes atributos:
  * `id` - ID del usuario.
  * `name` - Nombre completo del usuario.
  * `username` - Nombre de usuario.
  * `uid` - UID del usuario.
  * `email` - Email del usuario.
  * `active` - Indica si el usuario está activo.
  * `role` - Rol del usuario.
  * `category` - ID de la categoría del usuario.
  * `group` - ID del grupo del usuario.
  * `secondary_groups` - Lista de IDs de grupos secundarios.
  * `provider` - Proveedor de autenticación.
  * `email_verified` - Indica si el email está verificado.
  * `disclaimer_acknowledged` - Indica si se ha aceptado el disclaimer.
  * `role_name` - Nombre legible del rol.
  * `category_name` - Nombre legible de la categoría.
  * `group_name` - Nombre legible del grupo.

## Caso de Uso Común

Un caso común es buscar un usuario específico por nombre para obtener su ID y utilizarlo en otros recursos:

```terraform
# Buscar un usuario específico
data "isardvdi_users" "john_doe" {
  name_filter = "John Doe"
}

# Usar el ID del usuario en otro recurso
resource "isardvdi_vm" "user_desktop" {
  name        = "Desktop for ${data.isardvdi_users.john_doe.users[0].name}"
  template_id = "some-template-id"
  # Otros atributos...
}

# Mostrar información del usuario
output "user_info" {
  value = {
    id       = data.isardvdi_users.john_doe.users[0].id
    username = data.isardvdi_users.john_doe.users[0].username
    email    = data.isardvdi_users.john_doe.users[0].email
    role     = data.isardvdi_users.john_doe.users[0].role_name
  }
}
```

## Notas

* Si se proporciona `name_filter`, el data source utilizará el endpoint de búsqueda de la API, que es más eficiente para búsquedas específicas.
* Si no se proporciona `name_filter`, se obtendrán todos los usuarios y luego se aplicarán los filtros adicionales.
* Los filtros son acumulativos: si especificas múltiples filtros, solo se devolverán los usuarios que cumplan con todos los criterios.
* Si ningún usuario coincide con los filtros, la lista `users` estará vacía.
