---
page_title: "isard_qos_net Resource - terraform-provider-isard"
subcategory: ""
description: |-
  Manages QoS network configuration in Isard VDI.
---

# Resource: isard_qos_net

Gestiona un perfil de Quality of Service (QoS) para redes en Isard VDI. **Requiere privilegios de administrador.**

## Ejemplo de Uso

### Ejemplo Básico

```hcl
resource "isard_qos_net" "standard" {
  id          = "standard-qos"
  name        = "QoS Estándar"
  description = "Perfil QoS estándar para redes"
}
```

### Con Límites de Ancho de Banda

```hcl
resource "isard_qos_net" "limited" {
  id               = "limited-bandwidth"
  name             = "Ancho de Banda Limitado"
  description      = "Limita el tráfico de red a 100Mbps"
  inbound_average  = 102400  # 100 Mbps en KB/s
  inbound_peak     = 153600  # 150 Mbps peak
  inbound_burst    = 204800  # 200 Mbps burst
  outbound_average = 102400
  outbound_peak    = 153600
  outbound_burst   = 204800
}
```

### Perfil Alta Performance

```hcl
resource "isard_qos_net" "high_performance" {
  id               = "high-perf"
  name             = "Alta Performance"
  description      = "Sin límites para aplicaciones críticas"
  inbound_average  = 1048576  # 1 Gbps
  outbound_average = 1048576  # 1 Gbps
}
```

### Uso con Red de Usuario

```hcl
# Crear perfil QoS
resource "isard_qos_net" "dev_team" {
  id          = "dev-team-qos"
  name        = "QoS Equipo Dev"
  description = "200 Mbps para equipo de desarrollo"
  inbound_average  = 204800  # 200 Mbps
  outbound_average = 204800
}

# Usar en una red
resource "isard_network" "dev_network" {
  name        = "Red Desarrollo"
  description = "Red con QoS controlado"
  qos_id      = isard_qos_net.dev_team.id
}
```

## Argumentos

Los siguientes argumentos son soportados:

### Requeridos

- `id` - (Requerido) ID único del perfil QoS. No se puede modificar después de creado.
- `name` - (Requerido) Nombre descriptivo del perfil.

### Opcionales

- `description` - (Opcional) Descripción del perfil.
- `inbound_average` - (Opcional) Ancho de banda promedio de entrada en KB/s.
- `inbound_peak` - (Opcional) Ancho de banda pico de entrada en KB/s.
- `inbound_burst` - (Opcional) Ancho de banda burst de entrada en KB/s.
- `outbound_average` - (Opcional) Ancho de banda promedio de salida en KB/s.
- `outbound_peak` - (Opcional) Ancho de banda pico de salida en KB/s.
- `outbound_burst` - (Opcional) Ancho de banda burst de salida en KB/s.

## Atributos Exportados

Los mismos que los argumentos de entrada.

## Import

Los perfiles QoS pueden ser importados usando su ID:

```bash
terraform import isard_qos_net.standard standard-qos
```

## Ciclo de Vida

### Create

Al crear un perfil QoS:
1. Se valida que el usuario tenga privilegios de administrador
2. Se crea usando `POST /api/v3/admin/table/add/qos_net`
3. Se verifica la creación exitosa

### Read

Al leer un perfil QoS:
1. Se obtiene la lista completa desde `GET /api/v3/admin/table/qos_net`
2. Se busca el perfil específico por ID

### Update

Al actualizar un perfil QoS:
1. Se envían solo los campos modificados usando `PUT /api/v3/admin/table/update/qos_net`
2. Se releen los valores actualizados

### Delete

Al eliminar un perfil QoS:
1. Se elimina usando `DELETE /api/v3/admin/table/qos_net/{id}`

## Notas Importantes

- **Solo administradores** pueden gestionar perfiles QoS
- Los valores de ancho de banda están en **KB/s** (kilobytes por segundo)
- Para convertir de Mbps a KB/s: `Mbps * 1024 / 8`
  - Ejemplo: 100 Mbps = 100 * 1024 / 8 = 12800 KB/s
- Si no se especifican límites, el tráfico no está restringido
- El perfil `unlimited` es el valor por defecto del sistema

## Conversión de Unidades

### Mbps a KB/s

```
KB/s = (Mbps × 1024) / 8
```

Ejemplos:
- 10 Mbps = 1,280 KB/s
- 100 Mbps = 12,800 KB/s
- 1000 Mbps (1 Gbps) = 128,000 KB/s

### Valores Comunes

```hcl
# 10 Mbps
inbound_average = 1280

# 50 Mbps
inbound_average = 6400

# 100 Mbps
inbound_average = 12800

# 500 Mbps
inbound_average = 64000

# 1 Gbps
inbound_average = 128000
```

## Tipos de Límites

### Average (Promedio)
El ancho de banda promedio sostenido. Este es el límite principal.

### Peak (Pico)
El ancho de banda máximo que se puede alcanzar temporalmente.

### Burst (Ráfaga)
Cuántos bytes pueden enviarse en una ráfaga antes de aplicar el límite promedio.

## Ejemplo Completo

```hcl
# Perfil para desarrollo (100 Mbps)
resource "isard_qos_net" "dev" {
  id               = "qos-dev"
  name             = "Desarrollo"
  description      = "100 Mbps para desarrollo"
  inbound_average  = 12800
  inbound_peak     = 25600
  inbound_burst    = 51200
  outbound_average = 12800
  outbound_peak    = 25600
  outbound_burst   = 51200
}

# Perfil para producción (1 Gbps)
resource "isard_qos_net" "prod" {
  id               = "qos-prod"
  name             = "Producción"
  description      = "1 Gbps para producción"
  inbound_average  = 128000
  outbound_average = 128000
}

# Red con QoS de desarrollo
resource "isard_network" "dev_net" {
  name   = "Red Desarrollo"
  qos_id = isard_qos_net.dev.id
}

# Interfaz con QoS de producción
resource "isard_network_interface" "prod_bridge" {
  id     = "prod-bridge"
  name   = "Bridge Producción"
  net    = "br-prod"
  kind   = "bridge"
  qos_id = isard_qos_net.prod.id
}
```
