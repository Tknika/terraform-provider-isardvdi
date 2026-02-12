---
page_title: "Provider: Isard"
description: |-
  The Isard provider is used to interact with Isard VDI resources.
---

# Isard Provider

El provider de Isard permite gestionar recursos en Isard VDI mediante Terraform.

## Configuración del Provider

### Ejemplo de Uso

```hcl
provider "isard" {
  endpoint         = "localhost"
  auth_method      = "form"
  cathegory_id     = "default"
  username         = "admin"
  password         = "IsardVDI"
  ssl_verification = false  # Para desarrollo local con certificados autofirmados
}
```

### Autenticación con Token

```hcl
provider "isard" {
  endpoint     = "mi-servidor.isard.com"
  auth_method  = "token"
  token        = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  cathegory_id = "default"
  # ssl_verification por defecto es true (recomendado para producción)
}
```

### Producción con SSL Validado

```hcl
provider "isard" {
  endpoint         = "isard.empresa.com"
  auth_method      = "form"
  cathegory_id     = "default"
  username         = var.isard_username
  password         = var.isard_password
  ssl_verification = true  # Valida certificados SSL (default)
}
```

## Argumentos

Los siguientes argumentos son soportados:

### Requeridos

- `endpoint` - (Requerido) El hostname o IP del servidor Isard VDI (sin protocolo, se usa HTTPS automáticamente)
- `auth_method` - (Requerido) Método de autenticación. Valores aceptados: `"form"` o `"token"`
- `cathegory_id` - (Requerido) ID de la categoría en Isard VDI

### Opcionales

- `ssl_verification` - (Opcional) Habilita la verificación de certificados SSL. Establece a `false` para deshabilitar la verificación SSL (útil para desarrollo con certificados autofirmados). Por defecto: `true`. **Recomendación:** Mantener en `true` para entornos de producción.

### Opcionales según método de autenticación

#### Para `auth_method = "form"`

- `username` - (Requerido) Nombre de usuario para autenticación
- `password` - (Requerido) Contraseña para autenticación

#### Para `auth_method = "token"`

- `token` - (Requerido) Token JWT de Isard VDI

Nota: Opcionalmente se puede especificar `token` junto con `auth_method = "form"` para usar el token directamente en las llamadas API después de la autenticación inicial.

## Configuración SSL

Por defecto, el provider valida los certificados SSL del servidor Isard VDI (`ssl_verification = true`), lo cual es la configuración recomendada para entornos de producción.

Para entornos de desarrollo con certificados autofirmados, puedes deshabilitar la verificación SSL:

```hcl
provider "isard" {
  endpoint         = "localhost"
  auth_method      = "form"
  cathegory_id     = "default"
  username         = "admin"
  password         = "IsardVDI"
  ssl_verification = false  # Solo para desarrollo
}
```

**Advertencia de Seguridad:** Deshabilitar la verificación SSL (`ssl_verification = false`) hace que las conexiones sean vulnerables a ataques man-in-the-middle. Solo debe usarse en entornos de desarrollo controlados.

## Variables de Entorno

Puedes usar variables de entorno en lugar de especificar credenciales directamente:

```hcl
provider "isard" {
  endpoint         = var.isard_endpoint
  auth_method      = var.isard_auth_method
  username         = var.isard_username
  password         = var.isard_password
  cathegory_id     = var.isard_category
  ssl_verification = var.isard_ssl_verification
}
```

```bash
export TF_VAR_isard_endpoint="localhost"
export TF_VAR_isard_auth_method="form"
export TF_VAR_isard_username="admin"
export TF_VAR_isard_password="IsardVDI"
export TF_VAR_isard_category="default"
export TF_VAR_isard_ssl_verification="false"  # Solo para desarrollo
```

## Recursos y Data Sources

### Recursos

- [Resource: isard_vm](resources/isard_vm.md) - Gestión de desktops persistentes
- [Resource: isard_deployment](resources/deployment.md) - Gestión de deployments
- [Resource: isard_media](resources/isard_media.md) - Gestión de medios (ISOs y floppies)
- [Resource: isard_network](resources/isard_network.md) - Gestión de redes virtuales de usuario
- [Resource: isard_network_interface](resources/isard_network_interface.md) - Gestión de interfaces de red del sistema
- [Resource: isard_qos_net](resources/isard_qos_net.md) - Gestión de perfiles QoS de red

### Data Sources

- [Data Source: isard_templates](data-sources/isard_templates.md) - Consulta de templates disponibles
- [Data Source: isard_users](data-sources/isard_users.md) - Consulta de usuarios del sistema
- [Data Source: isard_medias](data-sources/isard_medias.md) - Consulta de medios (ISOs y floppies)
- [Data Source: isard_network_interfaces](data-sources/isard_network_interfaces.md) - Consulta de interfaces de red del sistema
- [Data Source: isard_groups](data-sources/isard_groups.md) - Consulta de grupos del sistema
