# Ejemplo de Deployments en Isard VDI

Este ejemplo muestra cómo crear y gestionar deployments en Isard VDI usando Terraform.

## ¿Qué es un Deployment?

Un deployment en Isard VDI permite crear múltiples desktops a partir de una plantilla para diferentes usuarios, grupos o categorías. Es ideal para:

- Proporcionar desktops a un equipo o departamento completo
- Crear entornos de desarrollo/prueba para grupos específicos
- Distribuir escritorios virtuales a estudiantes o participantes de un curso
- Gestionar desktops con configuraciones específicas para diferentes roles

## Características del Ejemplo

Este ejemplo incluye varios tipos de deployments:

1. **Deployment Básico para Grupo**: Crea desktops para todos los miembros de un grupo específico
2. **Deployment con Hardware Personalizado**: Configura recursos específicos (CPU, RAM) para desktops de alto rendimiento
3. **Deployment para Usuarios Específicos**: Asigna desktops solo a usuarios concretos con permisos personalizados
4. **Deployment por Categoría**: Crea desktops para todos los usuarios de una categoría
5. **Deployment con Red Personalizada**: Configura interfaces de red específicas
6. **Deployment Multi-Grupo**: Comparte desktops entre múltiples grupos

## Uso

1. Copia el archivo `terraform.tfvars.example` a `terraform.tfvars`:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

2. Edita `terraform.tfvars` y configura tus valores:
   ```hcl
   isard_token = "tu-token-aqui"
   template_id = "tu-template-uuid-aqui"
   ```

3. Inicializa Terraform:
   ```bash
   terraform init
   ```

4. Revisa el plan de ejecución:
   ```bash
   terraform plan
   ```

5. Aplica la configuración:
   ```bash
   terraform apply
   ```

## Variables Requeridas

- `isard_token`: Token de autenticación para Isard VDI
- `template_id`: ID de la plantilla que se usará para crear los desktops

## Personalización

Puedes personalizar los deployments modificando:

- **Recursos de Hardware**: Ajusta `vcpus` y `memory` según tus necesidades
- **Visibilidad**: Cambia `visible` para controlar si los desktops son visibles inmediatamente
- **Permisos**: Modifica `user_permissions` para definir qué pueden hacer los usuarios
- **Criterios de Asignación**: Usa `users`, `groups`, `categories` o `roles` en el bloque `allowed`

## Outputs

El ejemplo proporciona outputs con los IDs de los deployments creados, útiles para:
- Referencias en otros módulos de Terraform
- Integración con scripts externos
- Documentación de la infraestructura

## Limpieza

Para eliminar todos los recursos creados:

```bash
terraform destroy
```

**Advertencia**: Esto eliminará permanentemente todos los deployments y sus desktops asociados.

## Notas Importantes

- Al crear un deployment, Isard VDI creará automáticamente un desktop para cada usuario que coincida con los criterios especificados
- Cambiar el `template_id` forzará la recreación del deployment
- Algunos cambios pueden requerir que los desktops estén detenidos
- La eliminación de un deployment elimina también todos sus desktops asociados
