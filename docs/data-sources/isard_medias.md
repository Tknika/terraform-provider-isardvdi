# isard_medias Data Source

Obtiene información sobre los medios (ISOs, imágenes de disco) disponibles en Isard VDI.

## Ejemplo de Uso

```hcl
# Obtener todos los medios
data "isard_medias" "all" {}

# Buscar medios por nombre
data "isard_medias" "ubuntu" {
  name_filter = "Ubuntu"
}

# Filtrar por tipo de medio
data "isard_medias" "isos" {
  kind = "iso"
}

# Filtrar por estado
data "isard_medias" "downloaded" {
  status = "Downloaded"
}

# Filtrar por categoría y grupo
data "isard_medias" "category_medias" {
  category_id = "default"
  group_id    = "584a8dcb-abb0-446d-ba7c-d05fdb57f051"
}

# Filtrar por tipo y nombre
data "isard_medias" "ubuntu_isos" {
  name_filter = "ubuntu"
  kind        = "iso"
  status      = "Downloaded"
}

# Obtener medios de un usuario específico
data "isard_medias" "user_medias" {
  user_id = "local-default-admin-admin"
}
```

## Esquema de Argumentos

Todos los argumentos son opcionales y actúan como filtros:

- `name_filter` (String) - Filtra medios por nombre (búsqueda insensible a mayúsculas, coincidencia parcial).
- `kind` (String) - Filtra por tipo de medio. Valores: `iso`, `disk`, `floppy`.
- `status` (String) - Filtra por estado del medio. Valores comunes:
  - `DownloadStarting` - Iniciando descarga
  - `Downloading` - Descargando
  - `Downloaded` - Descarga completa
  - `Failed` - Descarga fallida
- `category_id` (String) - Filtra por ID de categoría.
- `group_id` (String) - Filtra por ID de grupo.
- `user_id` (String) - Filtra por ID de usuario propietario.

## Atributos de Referencia

- `id` (String) - Identificador del data source.
- `medias` (List of Object) - Lista de medios que coinciden con los filtros. Cada medio tiene:
  - `id` (String) - ID único del medio.
  - `name` (String) - Nombre del medio.
  - `description` (String) - Descripción del medio.
  - `url` (String) - URL de descarga del medio.
  - `kind` (String) - Tipo de medio (iso, disk, floppy).
  - `status` (String) - Estado del medio.
  - `user` (String) - ID del usuario propietario.
  - `category` (String) - ID de la categoría.
  - `group` (String) - ID del grupo.
  - `icon` (String) - Icono del medio.
  - `path` (String) - Ruta del archivo en el servidor.

## Uso con Outputs

```hcl
data "isard_medias" "ubuntu" {
  name_filter = "Ubuntu"
  kind        = "iso"
}

# Mostrar todos los IDs
output "ubuntu_iso_ids" {
  value = data.isard_medias.ubuntu.medias[*].id
}

# Mostrar nombres y estados
output "ubuntu_isos_status" {
  value = [
    for media in data.isard_medias.ubuntu.medias : {
      name   = media.name
      status = media.status
      id     = media.id
    }
  ]
}
```

## Combinación con Recursos

Puedes usar el data source para referenciar medios existentes:

```hcl
# Buscar un ISO específico
data "isard_medias" "virtio_drivers" {
  name_filter = "VirtIO"
  kind        = "iso"
}

# Usar el ID en otro recurso (cuando se soporte adjuntar medios a VMs)
resource "isard_vm" "windows_vm" {
  name        = "Windows VM"
  template_id = "template-id"
  # En futuras versiones: media_ids = [data.isard_medias.virtio_drivers.medias[0].id]
}
```

## Filtros Múltiples

Los filtros son acumulativos (AND lógico):

```hcl
# Solo medios que cumplan TODAS las condiciones
data "isard_medias" "filtered" {
  name_filter = "Ubuntu"      # Y nombre contiene "Ubuntu"
  kind        = "iso"         # Y tipo es ISO
  status      = "Downloaded"  # Y está completamente descargado
  category_id = "default"     # Y pertenece a categoría default
}
```

## Tipos de Medio (kind)

- **iso** - Imágenes ISO (sistemas operativos, drivers, herramientas)
- **disk** - Imágenes de disco (qcow2, raw, vmdk)
- **floppy** - Imágenes de disquete (raramente usado)

## Estados Comunes (status)

- **DownloadStarting** - Proceso de descarga iniciando
- **Downloading** - Descarga en progreso
- **Downloaded** - Descarga completada exitosamente
- **Failed** - Descarga falló
- **Deleting** - Medio siendo eliminado

## Notas

- Sin filtros, devuelve todos los medios accesibles por el usuario autenticado
- Los administradores ven todos los medios del sistema
- Los usuarios solo ven medios propios o compartidos con ellos
- El filtro `name_filter` realiza búsqueda parcial insensible a mayúsculas
- Si no hay coincidencias, `medias` será una lista vacía (no genera error)
