# Ejemplo de uso del Data Source isard_users

Este ejemplo demuestra cómo usar el data source `isard_users` para buscar usuarios en Isard VDI.

## Configuración

Asegúrate de tener configurado el provider de Isard en tu archivo principal o en un archivo separado:

```terraform
provider "isard" {
  endpoint    = "isard.example.com"
  auth_method = "token"
  token       = "your-api-token"
}
```

## Uso

El archivo `main.tf` muestra varios ejemplos de cómo buscar y filtrar usuarios:

1. **Buscar por nombre**: Encuentra usuarios cuyo nombre contenga un texto específico
2. **Filtrar por categoría**: Obtiene todos los usuarios de una categoría
3. **Filtrar por rol**: Busca usuarios con un rol específico (admin, manager, user)
4. **Usuario específico**: Busca un usuario concreto y muestra su información
5. **Filtros combinados**: Aplica múltiples filtros simultáneamente

## Ejecutar el ejemplo

```bash
# Inicializar Terraform
terraform init

# Ver qué cambios se aplicarán (en este caso, solo lectura de datos)
terraform plan

# Aplicar la configuración
terraform apply
```

## Outputs

Los outputs mostrarán:
- IDs y nombres de usuarios encontrados
- Información detallada de usuarios específicos
- Listas de usuarios según diferentes criterios

## Casos de uso

### 1. Encontrar el ID de un usuario para asignarlo a un recurso

```terraform
data "isard_users" "john" {
  name_filter = "John Doe"
}

resource "isard_vm" "johns_desktop" {
  name        = "Desktop-${data.isard_users.john.users[0].username}"
  template_id = "template-123"
  # Otros atributos...
}
```

### 2. Listar todos los administradores

```terraform
data "isard_users" "admins" {
  role = "admin"
}

output "admin_list" {
  value = [for user in data.isard_users.admins.users : {
    name  = user.name
    email = user.email
  }]
}
```

### 3. Verificar usuarios activos en una categoría

```terraform
data "isard_users" "active_users" {
  category_id = "production"
  active      = true
}

output "active_user_count" {
  value = length(data.isard_users.active_users.users)
}
```
