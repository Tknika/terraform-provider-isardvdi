# Ejemplo de Media

Este ejemplo demuestra cómo gestionar archivos multimedia (ISOs, imágenes de disco) en Isard VDI usando el resource `isard_media`.

## Configuración

Los ejemplos incluyen:

1. **ISO Público**: Ubuntu 22.04 LTS disponible para todos
2. **ISO Restringido por Roles**: Windows Server solo para administradores y usuarios avanzados
3. **Drivers Compartidos por Grupos**: VirtIO drivers disponibles para grupos específicos
4. **Disco Personalizado por Usuario**: Imagen de disco para un usuario específico
5. **ISO por Categoría**: Debian disponible para toda una categoría

## Requisitos Importantes

### URLs HTTPS Obligatorias

Todas las URLs **deben usar HTTPS**. El servidor Isard VDI rechazará URLs HTTP:

```hcl
# ✅ CORRECTO
url = "https://releases.ubuntu.com/22.04/ubuntu-22.04-desktop-amd64.iso"

# ❌ INCORRECTO - Será rechazado
url = "http://example.com/image.iso"
```

### Accesibilidad de URLs

- Las URLs deben ser accesibles desde el servidor Isard VDI (no desde tu máquina local)
- Isard verificará que puede acceder a la URL antes de crear el medio
- Si la URL no es accesible, la creación fallará

## Tipos de Medio (kind)

- `iso`: Imágenes ISO (instalación de SO, drivers)
- `disk`: Imágenes de disco (qcow2, raw, vmdk)
- `floppy`: Imágenes de disquete (raramente usado)

## Permisos (allowed)

El bloque `allowed` es opcional. Si no se especifica:
- Solo el creador puede ver el medio
- Los administradores siempre ven todos los medios

Opciones de permisos:

```hcl
allowed {
  # Por roles
  roles = ["admin", "advanced", "user"]
  
  # Por categorías (IDs)
  categories = ["category-id-1", "category-id-2"]
  
  # Por grupos (IDs)
  groups = ["group-id-1", "group-id-2"]
  
  # Por usuarios específicos (IDs)
  users = ["user-id-1", "user-id-2"]
}
```

Los permisos son **aditivos**: un usuario solo necesita coincidir con uno de los criterios.

## Uso

```bash
# Inicializar Terraform
terraform init

# Ver el plan de ejecución
terraform plan

# Aplicar la configuración
terraform apply

# Destruir los recursos
terraform destroy
```

## Notas

- La creación de medios es **asíncrona**. El proveedor espera un momento después de la solicitud inicial.
- Los medios comienzan en estado `DownloadStarting` y progresan a través de varios estados.
- Cambios en `name`, `url`, o `kind` requieren **recrear** el recurso.
- Si un medio está siendo usado por VMs activas, la eliminación puede fallar.

## Ejemplo con Data Source de Usuarios

Puedes combinar con el data source `isard_users` para asignar permisos:

```hcl
data "isard_users" "admins" {
  role = "admin"
}

resource "isard_media" "restricted_iso" {
  name = "Restricted Software"
  url  = "https://example.com/software.iso"
  kind = "iso"
  
  allowed {
    users = data.isard_users.admins.users[*].id
  }
}
```
