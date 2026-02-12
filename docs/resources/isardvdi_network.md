---
page_title: "isardvdi_network Resource - terraform-provider-isardvdi"
subcategory: ""
description: |-
  Manages a network in Isard VDI.
---

# Resource: isardvdi_network

Gestiona una red virtual de usuario en Isard VDI.

## Ejemplo de Uso

### Ejemplo Básico

```hcl
resource "isardvdi_network" "mi_red" {
  name        = "Red de Desarrollo"
  description = "Red virtual para equipo de desarrollo"
}
```

### Con Configuración Completa

```hcl
resource "isardvdi_network" "red_produccion" {
  name        = "Red Producción"
  description = "Red virtual para entorno de producción"
  model       = "virtio"
  qos_id      = "high-performance"
}
```

### Múltiples Redes

```hcl
resource "isardvdi_network" "redes_equipo" {
  count = 3
  
  name        = "Red Equipo ${count.index + 1}"
  description = "Red para equipo ${count.index + 1}"
  model       = "virtio"
  qos_id      = "unlimited"
}
```

## Argumentos

Los siguientes argumentos son soportados:

### Requeridos

- `name` - (Requerido) Nombre de la red. Debe ser único para el usuario.

### Opcionales

- `description` - (Opcional) Descripción de la red.
- `model` - (Opcional) Modelo de interfaz de red. Por defecto: `"virtio"`. Valores: `"virtio"`, `"e1000"`, `"rtl8139"`.
- `qos_id` - (Opcional) ID del perfil QoS de red a aplicar. Por defecto: `"unlimited"`.

## Atributos Exportados

Además de los argumentos anteriores, se exportan los siguientes atributos:

- `id` - ID único de la red en Isard VDI.
- `metadata_id` - ID de metadatos de la red (número grande, almacenado como string).

## Import

Las redes pueden ser importadas usando su ID:

```bash
terraform import isard_network.mi_red a1b2c3d4-e5f6-7890-abcd-ef1234567890
```

## Ciclo de Vida

### Create

Al crear una red:
1. Se crea la red usando `POST /api/v3/user/networks`
2. Se obtiene el ID de la red creada
3. Se lee la información completa incluyendo `metadata_id`

### Read

Al leer una red:
1. Se obtiene la información desde `GET /api/v3/user/networks/{id}`
2. Se actualizan todos los atributos

### Update

Al actualizar una red:
1. Se envían solo los campos modificados usando `PUT /api/v3/user/networks/{id}`
2. Se releen los valores actualizados

### Delete

Al eliminar una red:
1. Se elimina usando `DELETE /api/v3/user/networks/{id}`

## Notas Importantes

- Las redes virtuales de usuario son privadas y solo visibles para el usuario que las crea
- El `metadata_id` es un número grande (uint64) y se maneja como string para evitar overflow
- Asegúrate de que el perfil QoS especificado existe en el sistema
- El modelo `virtio` ofrece mejor rendimiento en la mayoría de los casos

## Uso con VMs

```hcl
# Crear una red
resource "isardvdi_network" "mi_red" {
  name        = "Red Desarrollo"
  description = "Red para VMs de desarrollo"
}

# Usar la red en una VM (a través de metadata_id)
# Nota: Las redes de usuario se asignan automáticamente,
# para control granular use isard_network_interface
```
