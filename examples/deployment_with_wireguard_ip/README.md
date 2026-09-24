# Deployment con IP wireguard (para Ansible u otros consumidores externos)

Este ejemplo muestra cómo hacer que Terraform arranque automáticamente los
desktops de un `isardvdi_deployment` y espere a que reporten su IP de
wireguard antes de terminar el `apply`, exponiéndola como salida.

## Qué hace

- Crea un `isardvdi_deployment` con `wait_for_ip = true`.
- Terraform arranca los desktops del deployment y espera hasta
  `ip_wait_timeout_seconds` (por defecto 300 s) a que cada uno reporte IP.
- El atributo `desktops` queda relleno con `id`, `name`, `ip`, `user` y
  `user_name` de cada desktop.
- El output `desktops` expone esa lista para que pueda ser consumida por
  fuera de Terraform (por ejemplo, para generar un inventory de Ansible).

## Fuera de alcance

El proveedor **no** ejecuta ninguna herramienta de provisionamiento. Tampoco
apaga las VMs automáticamente: quedan encendidas tras el apply para permitir
el acceso (SSH, etc.).

## Uso

```bash
terraform init
terraform apply \
  -var="isard_token=tu-token-aqui" \
  -var="template_id=tu-template-uuid-aqui"

terraform output -json desktops
```
