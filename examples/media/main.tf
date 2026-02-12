terraform {
  required_providers {
    isardvdi = {
      source = "terraform.local/local/isard"
    }
  }
}

provider "isardvdi" {
  host     = "https://localhost"
  username = "admin"
  password = "IsardVDI"
  insecure = true
}

# Ejemplo 1: ISO público de Ubuntu
resource "isardvdi_media" "ubuntu_iso" {
  name        = "Ubuntu 22.04 LTS Desktop"
  url         = "https://releases.ubuntu.com/22.04/ubuntu-22.04.3-desktop-amd64.iso"
  kind        = "iso"
  description = "Ubuntu Desktop 22.04 LTS - Official ISO"
}

# Ejemplo 2: ISO restringido a roles específicos
resource "isardvdi_media" "windows_iso" {
  name        = "Windows Server 2022"
  url         = "https://example.com/path/to/windows-server-2022.iso"
  kind        = "iso"
  description = "Windows Server 2022 Evaluation - Solo para administradores"
  
  allowed {
    roles = ["admin", "advanced"]
  }
}

# Ejemplo 3: Drivers compartidos con grupos específicos
resource "isardvdi_media" "virtio_drivers" {
  name        = "VirtIO Drivers"
  url         = "https://fedorapeople.org/groups/virt/virtio-win/direct-downloads/stable-virtio/virtio-win.iso"
  kind        = "iso"
  description = "VirtIO drivers for Windows guests"
  
  allowed {
    roles  = ["admin", "advanced", "user"]
    groups = ["group-id-1", "group-id-2"]
  }
}

# Ejemplo 4: Imagen de disco para usuario específico
resource "isardvdi_media" "custom_disk" {
  name        = "Custom Application Disk"
  url         = "https://storage.example.com/custom-disk.qcow2"
  kind        = "disk"
  description = "Custom application pre-installed disk"
  
  allowed {
    users = ["user-id-123"]
  }
}

# Ejemplo 5: Disponible para toda una categoría
resource "isardvdi_media" "shared_iso" {
  name        = "Debian 12 NetInst"
  url         = "https://cdimage.debian.org/debian-cd/current/amd64/iso-cd/debian-12.5.0-amd64-netinst.iso"
  kind        = "iso"
  description = "Debian 12 Network Installer"
  
  allowed {
    categories = ["default"]
    roles      = ["user", "advanced"]
  }
}

# Outputs para ver los IDs creados
output "ubuntu_iso_id" {
  value       = isardvdi_media.ubuntu_iso.id
  description = "ID del medio Ubuntu ISO"
}

output "windows_iso_id" {
  value       = isardvdi_media.windows_iso.id
  description = "ID del medio Windows ISO"
}
