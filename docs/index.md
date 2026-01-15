# Isard Provider

El provider de Isard permite gestionar recursos en Isard VDI mediante Terraform.

## Configuración del Provider

### Ejemplo de Uso

```hcl
provider "isard" {
  endpoint     = "localhost"
  auth_method  = "form"
  cathegory_id = "default"
  username     = "admin"
  password     = "IsardVDI"
}
```

### Autenticación con Token

```hcl
provider "isard" {
  endpoint     = "mi-servidor.isard.com"
  auth_method  = "token"
  token        = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  cathegory_id = "default"
}
```

## Argumentos

Los siguientes argumentos son soportados:

### Requeridos

- `endpoint` - (Requerido) El hostname o IP del servidor Isard VDI (sin protocolo, se usa HTTPS automáticamente)
- `auth_method` - (Requerido) Método de autenticación. Valores aceptados: `"form"` o `"token"`
- `cathegory_id` - (Requerido) ID de la categoría en Isard VDI

### Opcionales según método de autenticación

#### Para `auth_method = "form"`

- `username` - (Requerido) Nombre de usuario para autenticación
- `password` - (Requerido) Contraseña para autenticación

#### Para `auth_method = "token"`

- `token` - (Requerido) Token JWT de Isard VDI

Nota: Opcionalmente se puede especificar `token` junto con `auth_method = "form"` para usar el token directamente en las llamadas API después de la autenticación inicial.

## Configuración SSL

El provider está configurado para omitir la verificación de certificados SSL (desarrollo). Para entornos de producción, se recomienda modificar el código para validar certificados.

## Variables de Entorno

Puedes usar variables de entorno en lugar de especificar credenciales directamente:

```hcl
provider "isard" {
  endpoint     = var.isard_endpoint
  auth_method  = var.isard_auth_method
  username     = var.isard_username
  password     = var.isard_password
  cathegory_id = var.isard_category
}
```

```bash
export TF_VAR_isard_endpoint="localhost"
export TF_VAR_isard_auth_method="form"
export TF_VAR_isard_username="admin"
export TF_VAR_isard_password="IsardVDI"
export TF_VAR_isard_category="default"
```

## Recursos y Data Sources

- [Resource: isard_vm](resources/isard_vm.md) - Gestión de desktops persistentes
- [Data Source: isard_templates](data-sources/isard_templates.md) - Consulta de templates disponibles
