# Terraform Provider para Isard VDI

Provider de Terraform para gestionar recursos en Isard VDI a través de su API v3.

## Características

- ✅ Creación, lectura y eliminación de desktops persistentes
- ✅ Listado de templates disponibles con filtrado por nombre
- ✅ Soporte para autenticación mediante token o formulario
- ✅ Configuración SSL flexible para desarrollo y producción

## Requisitos

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24 (para desarrollo)
- Instancia de Isard VDI con API v3 activa

## Instalación

### Desarrollo Local

1. Clonar el repositorio:
```bash
git clone https://github.com/prubioamzubiri/isard-terraform-provider.git
cd isard-terraform-provider
```

2. Compilar el provider:
```bash
go build -o terraform-provider-isard
```

3. Configurar el override local en `~/.terraformrc`:
```hcl
provider_installation {
  dev_overrides {
    "tknika/isard" = "/home/tu-usuario/source/isard-terraform-provider"
  }
  direct {}
}
```

### Instalación desde Registry (futuro)

```hcl
terraform {
  required_providers {
    isard = {
      source  = "registry.terraform.io/tknika/isard"
      version = "~> 1.0"
    }
  }
}
```

## Uso Básico

```hcl
# Configurar el provider
provider "isard" {
  endpoint     = "localhost"
  auth_method  = "form"
  cathegory_id = "default"
  username     = "admin"
  password     = "IsardVDI"
}

# Obtener templates disponibles
data "isard_templates" "ubuntu" {
  name_filter = "Ubuntu"
}

# Crear un desktop persistente
resource "isard_vm" "mi_desktop" {
  name        = "mi-desktop-terraform"
  description = "Desktop creado con Terraform"
  template_id = data.isard_templates.ubuntu.templates[0].id
}
```

## Documentación

- [Configuración del Provider](docs/index.md)
- [Resource: isard_vm](docs/resources/isard_vm.md)
- [Data Source: isard_templates](docs/data-sources/isard_templates.md)

## Ejemplos

Consulta el directorio [examples/](examples/) para ver ejemplos completos de uso.

## Desarrollo

### Compilar

```bash
go build -o terraform-provider-isard
```

### Ejecutar Tests

```bash
go test ./...
```

### Depuración

Para habilitar logs detallados:

```bash
export TF_LOG=DEBUG
terraform plan
```

## Contribuir

Las contribuciones son bienvenidas. Por favor:

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/mi-feature`)
3. Commit tus cambios (`git commit -am 'Agregar nueva característica'`)
4. Push a la rama (`git push origin feature/mi-feature`)
5. Abre un Pull Request

## Licencia

[Especificar licencia]

## Soporte

Para reportar bugs o solicitar features, abre un issue en GitHub.
