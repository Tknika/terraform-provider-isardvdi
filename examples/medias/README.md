# Ejemplo de Data Source isardvdi_medias

Este ejemplo demuestra cómo usar el data source `isardvdi_medias` para consultar y filtrar medios (ISOs, imágenes de disco) en Isard VDI.

## ¿Qué es un Medio en Isard VDI?

Los medios son archivos que se pueden adjuntar a máquinas virtuales:
- **ISOs**: Imágenes de instalación de sistemas operativos, drivers, herramientas
- **Disk Images**: Discos duros virtuales pre-configurados
- **Floppies**: Imágenes de disquete (raramente usado)

## Ejemplos de Filtrado

### 1. Todos los Medios
```hcl
data "isardvdi_medias" "all" {}
```
Obtiene todos los medios accesibles por el usuario.

### 2. Filtrar por Nombre
```hcl
data "isardvdi_medias" "ubuntu" {
  name_filter = "Ubuntu"  # Coincidencia parcial, insensible a mayúsculas
}
```
Encuentra todos los medios cuyo nombre contiene "Ubuntu".

### 3. Filtrar por Tipo
```hcl
data "isardvdi_medias" "isos" {
  kind = "iso"  # Valores: iso, disk, floppy
}
```
Obtiene solo medios de tipo ISO.

### 4. Filtrar por Estado
```hcl
data "isardvdi_medias" "ready" {
  status = "Downloaded"  # Solo medios completamente descargados
}
```
Estados comunes:
- `DownloadStarting` - Iniciando
- `Downloading` - En progreso
- `Downloaded` - Completado
- `Failed` - Error

### 5. Filtros Combinados
```hcl
data "isardvdi_medias" "ubuntu_ready" {
  name_filter = "ubuntu"
  kind        = "iso"
  status      = "Downloaded"
}
```
Solo ISOs de Ubuntu que estén completamente descargados.

### 6. Filtrar por Propietario
```hcl
data "isardvdi_medias" "my_medias" {
  user_id = "user-id-here"
}
```

### 7. Filtrar por Categoría/Grupo
```hcl
data "isardvdi_medias" "category_medias" {
  category_id = "default"
  group_id    = "group-id-here"
}
```

## Uso de los Resultados

### Obtener IDs
```hcl
output "media_ids" {
  value = data.isardvdi_medias.ubuntu.medias[*].id
}
```

### Crear un Mapa
```hcl
output "media_map" {
  value = {
    for media in data.isardvdi_medias.all.medias :
    media.name => media.id
  }
}
```

### Filtrado Post-Consulta
```hcl
output "large_isos" {
  value = [
    for media in data.isardvdi_medias.all.medias :
    media.name
    if media.kind == "iso" && media.status == "Downloaded"
  ]
}
```

### Acceso Condicional
```hcl
locals {
  ubuntu_iso = length(data.isardvdi_medias.ubuntu.medias) > 0 ? data.isardvdi_medias.ubuntu.medias[0] : null
}

output "ubuntu_iso_info" {
  value = local.ubuntu_iso != null ? {
    id   = local.ubuntu_iso.id
    name = local.ubuntu_iso.name
    url  = local.ubuntu_iso.url
  } : "No Ubuntu ISO found"
}
```

## Casos de Uso

### 1. Verificar Disponibilidad de Medios
Antes de crear VMs, verifica que los ISOs necesarios estén disponibles:

```hcl
data "isardvdi_medias" "required_iso" {
  name_filter = "Ubuntu 22.04"
  status      = "Downloaded"
}

# Usar en un check
check "iso_available" {
  assert {
    condition     = length(data.isardvdi_medias.required_iso.medias) > 0
    error_message = "Ubuntu 22.04 ISO no está disponible"
  }
}
```

### 2. Auditoría de Medios
Lista todos los medios y su estado:

```hcl
output "media_inventory" {
  value = {
    for media in data.isardvdi_medias.all.medias :
    media.name => {
      type     = media.kind
      status   = media.status
      category = media.category
      owner    = media.user
    }
  }
}
```

### 3. Identificar Descargas Pendientes
```hcl
output "downloads_in_progress" {
  value = [
    for media in data.isardvdi_medias.all.medias :
    {
      name   = media.name
      status = media.status
    }
    if contains(["DownloadStarting", "Downloading"], media.status)
  ]
}
```

### 4. Buscar por Múltiples Criterios
```hcl
# Combinando con data source de usuarios
data "isardvdi_users" "admin_user" {
  role = "admin"
}

data "isardvdi_medias" "admin_medias" {
  user_id = data.isardvdi_users.admin_user.users[0].id
  status  = "Downloaded"
}
```

## Ejecución

```bash
# Inicializar
terraform init

# Ver la consulta
terraform plan

# Aplicar (solo actualiza el state con los datos)
terraform apply

# Ver los outputs
terraform output

# Ver un output específico
terraform output ubuntu_isos
```

## Notas Importantes

1. **Permisos**: Los usuarios solo ven medios que les pertenecen o han sido compartidos con ellos. Los administradores ven todos.

2. **Filtro name_filter**: 
   - Es case-insensitive
   - Hace coincidencia parcial (substring)
   - "ubuntu" coincide con "Ubuntu 22.04 LTS"

3. **Lista Vacía ≠ Error**: Si no hay coincidencias, `medias` será una lista vacía, no un error.

4. **Campos Computados**: Todos los atributos en `medias` son de solo lectura (computed).

5. **Estados Transitorios**: Los estados como `Downloading` pueden cambiar entre ejecuciones de Terraform.

## Combinación con Recursos

Aunque actualmente no se puede adjuntar medios directamente a VMs en este provider, puedes usar el data source para:

- Documentar qué medios existen
- Verificar prerequisitos antes de crear recursos
- Generar inventarios y reportes
- Validar que los medios necesarios estén descargados

En futuras versiones, podrías usar los IDs así:

```hcl
data "isardvdi_medias" "drivers" {
  name_filter = "VirtIO Drivers"
}

resource "isardvdi_vm" "windows" {
  # ... other config ...
  # Hipotético en futuras versiones:
  # cdrom_media_id = data.isardvdi_medias.drivers.medias[0].id
}
```
