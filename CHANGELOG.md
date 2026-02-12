# Changelog

Todos los cambios notables de este proyecto serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/es/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/lang/es/).

## [0.2.1] - 2026-02-12

### Cambiado
- **BREAKING CHANGE**: Provider renombrado de `isard` a `isardvdi`
  - Todos los recursos cambian de `isard_*` a `isardvdi_*`
  - Todos los data sources cambian de `isard_*` a `isardvdi_*`
  - El provider ahora se declara como `provider "isardvdi"` en lugar de `provider "isard"`
  - El source del provider es ahora `registry.terraform.io/tknika/isardvdi`
- Actualizado módulo Go de `terraform-provider-isard` a `terraform-provider-isardvdi`
- Actualizada toda la documentación con el nuevo nombre

### Migración desde v0.2.0

Para migrar desde la versión 0.2.0, necesitas actualizar tu código Terraform:

1. Actualiza el bloque `required_providers`:
```hcl
terraform {
  required_providers {
    isardvdi = {  # antes: isard
      source = "registry.terraform.io/tknika/isardvdi"  # antes: tknika/isard
      version = "~> 0.2.1"
    }
  }
}
```

2. Actualiza el bloque del provider:
```hcl
provider "isardvdi" {  # antes: provider "isard"
  endpoint     = "localhost"
  auth_method  = "form"
  cathegory_id = "default"
  username     = "admin"
  password     = "IsardVDI"
}
```

3. Actualiza todos los recursos y data sources:
```hcl
# Antes:
resource "isard_vm" "example" { ... }
data "isard_templates" "all" { ... }

# Ahora:
resource "isardvdi_vm" "example" { ... }
data "isardvdi_templates" "all" { ... }
```

4. Ejecuta:
```bash
terraform state replace-provider registry.terraform.io/tknika/isard registry.terraform.io/tknika/isardvdi
terraform init -upgrade
```

## [0.2.0] - 2026-01-XX

### Agregado
- Soporte inicial para gestión de VMs persistentes
- Soporte para deployments
- Soporte para medios (ISOs y floppies)
- Soporte para redes virtuales
- Soporte para interfaces de red del sistema
- Soporte para perfiles QoS de red
- Data sources para templates, usuarios, grupos, medias e interfaces de red
- Autenticación mediante token JWT y formulario
- Configuración SSL flexible

### Seguridad
- Soporte para verificación SSL configurable

