terraform {
  required_providers {
    isard = {
      source = "terraform.local/local/isard"
    }
  }
}

provider "isard" {
  host     = "https://localhost"
  username = "admin"
  password = "IsardVDI"
  insecure = true
}

# Ejemplo 1: VM con un ISO específico
resource "isard_vm" "windows_with_drivers" {
  name        = "Windows VM con Drivers"
  description = "VM de Windows con drivers VirtIO"
  template_id = "windows-template-id"
  vcpus       = 4
  memory      = 8
  
  # Adjuntar ISO de drivers VirtIO
  isos = ["virtio-drivers-iso-id"]
}

# Ejemplo 2: VM con múltiples ISOs
resource "isard_vm" "lab_vm" {
  name        = "VM de Laboratorio"
  description = "VM con múltiples herramientas"
  template_id = "template-id"
  vcpus       = 2
  memory      = 4
  
  # Adjuntar múltiples ISOs
  isos = [
    "ubuntu-22.04-iso-id",
    "gparted-live-iso-id",
    "system-rescue-iso-id"
  ]
}

# Ejemplo 3: VM usando data source para buscar ISOs
data "isard_medias" "ubuntu_iso" {
  name_filter = "Ubuntu 22.04"
  kind        = "iso"
  status      = "Downloaded"
}

resource "isard_vm" "ubuntu_vm" {
  name        = "Ubuntu Desktop"
  description = "Ubuntu con ISO de instalación"
  template_id = "template-id"
  vcpus       = 2
  memory      = 4
  
  # Usar el primer ISO encontrado
  isos = length(data.isard_medias.ubuntu_iso.medias) > 0 ? [
    data.isard_medias.ubuntu_iso.medias[0].id
  ] : []
}

# Ejemplo 4: VM con floppy (raramente usado)
resource "isard_vm" "legacy_vm" {
  name        = "VM Legacy"
  description = "VM con floppy disk"
  template_id = "template-id"
  vcpus       = 1
  memory      = 2
  
  floppies = ["driver-floppy-id"]
}

# Ejemplo 5: VM con ISOs y floppies
resource "isard_vm" "complete_vm" {
  name        = "VM Completa"
  description = "VM con ISOs y floppies"
  template_id = "template-id"
  vcpus       = 2
  memory      = 4
  
  isos     = ["installation-iso-id", "tools-iso-id"]
  floppies = ["drivers-floppy-id"]
}

# Outputs
output "windows_vm_id" {
  value = isard_vm.windows_with_drivers.id
}

output "ubuntu_vm_id" {
  value = isard_vm.ubuntu_vm.id
}

output "ubuntu_iso_used" {
  value = length(data.isard_medias.ubuntu_iso.medias) > 0 ? data.isard_medias.ubuntu_iso.medias[0].name : "No ISO found"
}
