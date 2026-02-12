# isardvdi_media Resource

Gestiona archivos multimedia (ISOs, imágenes de disco) en Isard VDI.

## Ejemplo de Uso

```hcl
# ISO básico
resource "isardvdi_media" "ubuntu_iso" {
  name        = "Ubuntu 22.04 LTS"
  url         = "https://releases.ubuntu.com/22.04/ubuntu-22.04-desktop-amd64.iso"
  kind        = "iso"
  description = "Ubuntu Desktop 22.04 LTS ISO"
}

# Media con permisos específicos
resource "isardvdi_media" "restricted_iso" {
  name        = "Windows Server ISO"
  url         = "https://example.com/windows-server.iso"
  kind        = "iso"
  description = "Windows Server 2022 ISO"
  
  allowed {
    roles      = ["admin", "advanced"]
    categories = ["default"]
    groups     = ["group-id-1", "group-id-2"]
    users      = ["user-id-1"]
  }
}
```

## Esquema de Argumentos

### Atributos Requeridos

- `name` (String) - Nombre del medio. **Requiere reemplazo** si se cambia. Máximo 50 caracteres.
- `url` (String) - URL del archivo multimedia. **Debe ser HTTPS**. **Requiere reemplazo** si se cambia.
- `kind` (String) - Tipo de medio. Valores permitidos: `iso`, `disk`, `floppy`. **Requiere reemplazo** si se cambia.

### Atributos Opcionales

- `description` (String) - Descripción del medio. Máximo 255 caracteres.
- `allowed` (Block) - Define quién puede usar este medio. Si no se especifica, solo el creador tiene acceso.
  - `roles` (List of String) - Lista de roles permitidos (ej: "admin", "advanced", "user").
  - `categories` (List of String) - Lista de IDs de categorías permitidas.
  - `groups` (List of String) - Lista de IDs de grupos permitidos.
  - `users` (List of String) - Lista de IDs de usuarios permitidos.

### Atributos de Solo Lectura

- `id` (String) - ID único del medio en Isard VDI.

## Notas Importantes

### Requisitos de URL

- **Solo se permiten URLs HTTPS**. Las URLs HTTP serán rechazadas por la API de Isard.
- La URL debe ser accesible desde el servidor Isard VDI.
- El servidor Isard verificará la accesibilidad de la URL antes de aceptar la creación.

### Comportamiento de Creación

- La creación es **asíncrona**. El proveedor espera 2 segundos después de la solicitud inicial y luego busca el medio creado por nombre.
- El medio comienza en estado `DownloadStarting` y progresa a través de varios estados de descarga.

### Actualizaciones

- Los campos `name`, `url`, y `kind` **requieren reemplazo** del recurso (destruir y recrear).
- Actualmente, **no se soportan actualizaciones in-place** del medio. Cualquier cambio en `description` o `allowed` también requerirá reemplazo del recurso.

### Eliminación

- La eliminación es inmediata y permanente.
- Si el medio está siendo usado por VMs activas, la eliminación puede fallar.

## Tipos de Medio

- `iso` - Imagen ISO (típicamente para instalación de sistemas operativos)
- `disk` - Imagen de disco (discos duros virtuales, drivers)
- `floppy` - Imagen de disquete (raramente usado)

## Permisos (allowed)

El bloque `allowed` controla qué usuarios pueden ver y usar este medio cuando creen o modifiquen máquinas virtuales. Si no se especifica:

- El medio solo es visible para el usuario que lo creó
- Los administradores siempre pueden ver todos los medios

Al especificar `allowed`, puedes compartir el medio con:
- Roles específicos (todos los usuarios con ese rol)
- Categorías completas (todos los usuarios en esa categoría)
- Grupos específicos (todos los usuarios en esos grupos)
- Usuarios individuales (solo esos usuarios específicos)

Los permisos son **aditivos**: un usuario solo necesita coincidir con uno de los criterios para tener acceso.
