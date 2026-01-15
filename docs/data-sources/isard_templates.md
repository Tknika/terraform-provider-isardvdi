# Data Source: isard_templates

Obtiene la lista de templates disponibles para el usuario autenticado en Isard VDI.

## Ejemplo de Uso

### Obtener Todos los Templates

```hcl
data "isard_templates" "todos" {}

output "lista_completa" {
  value = data.isard_templates.todos.templates
}
```

### Filtrar Templates por Nombre

```hcl
data "isard_templates" "ubuntu" {
  name_filter = "Ubuntu"
}

output "templates_ubuntu" {
  value = data.isard_templates.ubuntu.templates
}
```

### Usar con Resource

```hcl
data "isard_templates" "ubuntu" {
  name_filter = "Ubuntu Desktop"
}

resource "isard_vm" "mi_desktop" {
  name        = "desktop-desde-datasource"
  description = "Usa el primer template Ubuntu encontrado"
  template_id = data.isard_templates.ubuntu.templates[0].id
}
```

### Filtrado Case-Insensitive

```hcl
# Estos tres ejemplos devolverán los mismos resultados
data "isard_templates" "ubuntu1" {
  name_filter = "ubuntu"
}

data "isard_templates" "ubuntu2" {
  name_filter = "Ubuntu"
}

data "isard_templates" "ubuntu3" {
  name_filter = "UBUNTU"
}
```

### Verificar que Exista al Menos un Template

```hcl
data "isard_templates" "windows" {
  name_filter = "Windows"
}

# Fallar si no hay templates
locals {
  template_count = length(data.isard_templates.windows.templates)
  template_id    = local.template_count > 0 ? data.isard_templates.windows.templates[0].id : null
}

resource "isard_vm" "windows_desktop" {
  count = local.template_count > 0 ? 1 : 0
  
  name        = "desktop-windows"
  description = "Desktop Windows"
  template_id = local.template_id
}
```

## Argumentos

Los siguientes argumentos son soportados:

### Opcionales

- `name_filter` - (Opcional) Filtro para buscar templates por nombre. La búsqueda es case-insensitive y busca coincidencias parciales (substring). Si no se especifica, devuelve todos los templates disponibles.

## Atributos Exportados

- `id` - ID del data source (siempre es `"templates"`).
- `templates` - Lista de objetos template. Cada objeto template contiene:
  - `id` - ID único del template.
  - `name` - Nombre del template.
  - `category` - ID de la categoría a la que pertenece el template.
  - `group` - ID del grupo al que pertenece el template.
  - `user_id` - ID del usuario propietario del template.
  - `icon` - Nombre del icono asociado al template.
  - `description` - Descripción del template.
  - `enabled` - Boolean indicando si el template está habilitado.
  - `status` - Estado actual del template (ej: `"Stopped"`).
  - `desktop_size` - Tamaño del disco del desktop en bytes.

## Comportamiento del Filtrado

El filtrado por nombre funciona de la siguiente manera:

1. **Sin filtro:** Devuelve todos los templates disponibles para el usuario
2. **Con filtro:** Devuelve solo los templates cuyo nombre contenga el string especificado
3. **Case-insensitive:** No distingue entre mayúsculas y minúsculas
4. **Substring match:** Busca el filtro como parte del nombre, no requiere coincidencia exacta

### Ejemplos de Filtrado

Si tienes estos templates:
- "Ubuntu Desktop 22.04"
- "Windows 10 Pro"
- "Debian Server"
- "Ubuntu Server 20.04"

```hcl
# Devuelve: Ubuntu Desktop 22.04, Ubuntu Server 20.04
data "isard_templates" "ubuntu" {
  name_filter = "Ubuntu"
}

# Devuelve: Ubuntu Desktop 22.04
data "isard_templates" "desktop" {
  name_filter = "Desktop"
}

# Devuelve: Ubuntu Desktop 22.04, Ubuntu Server 20.04, Debian Server
data "isard_templates" "servers" {
  name_filter = "server"
}

# Devuelve: todos los templates
data "isard_templates" "todos" {}
```

## Ejemplos Adicionales

### Seleccionar Template Específico

```hcl
data "isard_templates" "ubuntu" {
  name_filter = "Ubuntu"
}

locals {
  # Buscar un template específico por nombre exacto
  ubuntu_22_template = [
    for t in data.isard_templates.ubuntu.templates : t
    if t.name == "Ubuntu Desktop 22.04"
  ][0]
}

resource "isard_vm" "mi_vm" {
  name        = "mi-desktop-ubuntu-22"
  template_id = local.ubuntu_22_template.id
}
```

### Listar Información de Templates

```hcl
data "isard_templates" "todos" {}

output "templates_info" {
  value = {
    for template in data.isard_templates.todos.templates :
    template.name => {
      id          = template.id
      description = template.description
      size_gb     = template.desktop_size / 1024 / 1024 / 1024
      enabled     = template.enabled
      status      = template.status
    }
  }
}
```

### Crear Desktops para Cada Template

```hcl
data "isard_templates" "desarrollo" {
  name_filter = "Dev"
}

resource "isard_vm" "dev_desktops" {
  for_each = { 
    for idx, tmpl in data.isard_templates.desarrollo.templates : 
    idx => tmpl 
  }
  
  name        = "desktop-${each.value.name}-${each.key}"
  description = "Desktop basado en ${each.value.name}"
  template_id = each.value.id
}
```

### Validación de Templates

```hcl
data "isard_templates" "ubuntu" {
  name_filter = "Ubuntu"
}

# Validar que al menos un template fue encontrado
resource "null_resource" "validacion" {
  count = length(data.isard_templates.ubuntu.templates) == 0 ? 1 : 0
  
  provisioner "local-exec" {
    command = "echo 'ERROR: No se encontraron templates Ubuntu' && exit 1"
  }
}
```

## Notas Importantes

1. **Permisos:** Solo devuelve templates a los que el usuario autenticado tiene acceso.

2. **Performance:** La API devuelve todos los templates y el filtrado se hace localmente. Para grandes cantidades de templates, considera usar filtros específicos.

3. **Templates Deshabilitados:** El data source devuelve todos los templates, incluyendo los deshabilitados. Verifica el campo `enabled` si necesitas filtrar solo habilitados:

```hcl
locals {
  templates_habilitados = [
    for t in data.isard_templates.todos.templates : t
    if t.enabled
  ]
}
```

4. **Actualizaciones:** El data source se ejecuta en cada `terraform plan` o `terraform apply`, por lo que siempre obtendrás la lista actual de templates.

## Limitaciones Conocidas

1. No se puede filtrar por otros campos (categoría, grupo, estado) - solo por nombre
2. No se puede ordenar los resultados
3. No hay paginación para grandes cantidades de templates
