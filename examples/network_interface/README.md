# Ejemplo: Network Interfaces con Control de Permisos

Este ejemplo muestra cómo crear y gestionar interfaces de red del sistema con diferentes configuraciones de permisos.

## Requisitos

- Privilegios de administrador en Isard VDI
- Terraform >= 1.0
- Provider de Isard configurado

## Ejemplos Incluidos

### 1. Interfaz Básica
Interfaz de red simple sin restricciones de permisos explícitas.

### 2. Interfaz Pública (Visible para Todos)
Interfaz configurada con permisos explícitos para ser visible a todos los usuarios del sistema usando listas vacías en el bloque `allowed`.

```hcl
allowed {
  roles      = []
  categories = []
  groups     = []
  users      = []
}
```

### 3. Interfaz Restringida por Rol
Interfaz que solo pueden usar usuarios con el rol de `admin`.

### 4. Interfaz Restringida por Categoría
Interfaz limitada a usuarios de una categoría específica (ej: marketing).

## Uso

1. Configurar el provider en `provider.tf`:

```hcl
provider "isardvdi" {
  endpoint     = "isard.example.com"
  auth_method  = "form"
  cathegory_id = "default"
  username     = "admin"
  password     = var.admin_password
}
```

2. Inicializar Terraform:

```bash
terraform init
```

3. Revisar el plan:

```bash
terraform plan
```

4. Aplicar la configuración:

```bash
terraform apply
```

## Control de Acceso

### Permitir Acceso a Todos

Para hacer una interfaz visible a **todos los usuarios**:

```hcl
allowed {
  roles      = []
  categories = []
  groups     = []
  users      = []
}
```

### Restringir por Rol

```hcl
allowed {
  roles      = ["admin", "manager"]
  categories = []
  groups     = []
  users      = []
}
```

### Restringir por Categoría

```hcl
allowed {
  roles      = []
  categories = ["categoria-id-1", "categoria-id-2"]
  groups     = []
  users      = []
}
```

### Restringir por Grupo

```hcl
allowed {
  roles      = []
  categories = []
  groups     = ["grupo-id-1"]
  users      = []
}
```

### Restringir por Usuario

```hcl
allowed {
  roles      = []
  categories = []
  groups     = []
  users      = ["user-id-1", "user-id-2"]
}
```

### Combinación de Restricciones

Las restricciones se evalúan con lógica OR. Un usuario tendrá acceso si cumple **cualquiera** de las condiciones:

```hcl
allowed {
  roles      = ["admin"]           # Admins SIEMPRE tienen acceso
  categories = ["especial"]        # O usuarios de categoría especial
  groups     = []
  users      = ["user-vip"]        # O este usuario específico
}
```

## Notas Importantes

1. **Solo administradores** pueden gestionar interfaces de red del sistema
2. El campo `net` debe corresponder con un bridge/red existente en el sistema
3. Para interfaces OVS con wireguard, usar `net = "4095"`
4. Al eliminar una interfaz, se desasigna automáticamente de todas las VMs
5. Las listas vacías `[]` significan "sin restricción en este nivel"
6. Si omites completamente el bloque `allowed`, la interfaz no tendrá restricciones explícitas

## Limpieza

Para eliminar todos los recursos creados:

```bash
terraform destroy
```

## Ver También

- [Documentación del recurso isardvdi_network_interface](../../docs/resources/isardvdi_network_interface.md)
- [Documentación del data source isardvdi_network_interfaces](../../docs/data-sources/isardvdi_network_interfaces.md)
