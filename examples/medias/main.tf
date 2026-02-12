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

# Ejemplo 1: Obtener todos los medios
data "isardvdi_medias" "all" {}

# Ejemplo 2: Buscar ISOs de Ubuntu
data "isardvdi_medias" "ubuntu_isos" {
  name_filter = "Ubuntu"
  kind        = "iso"
}

# Ejemplo 3: Obtener solo medios descargados
data "isardvdi_medias" "downloaded" {
  status = "Downloaded"
}

# Ejemplo 4: Medios de una categoría específica
data "isardvdi_medias" "category_medias" {
  category_id = "default"
}

# Ejemplo 5: Filtro combinado - ISOs de Windows descargados
data "isardvdi_medias" "windows_isos_ready" {
  name_filter = "Windows"
  kind        = "iso"
  status      = "Downloaded"
}

# Ejemplo 6: Medios de disk (imágenes de disco)
data "isardvdi_medias" "disk_images" {
  kind = "disk"
}

# Ejemplo 7: Medios de un usuario específico
data "isardvdi_medias" "my_medias" {
  user_id = "local-default-admin-admin"
}

# Outputs para visualizar los resultados

output "total_medias" {
  value       = length(data.isardvdi_medias.all.medias)
  description = "Número total de medios disponibles"
}

output "all_media_names" {
  value       = data.isardvdi_medias.all.medias[*].name
  description = "Nombres de todos los medios"
}

output "ubuntu_isos" {
  value = [
    for media in data.isardvdi_medias.ubuntu_isos.medias : {
      id     = media.id
      name   = media.name
      status = media.status
      url    = media.url
    }
  ]
  description = "Información de ISOs de Ubuntu"
}

output "downloaded_medias_by_kind" {
  value = {
    for media in data.isardvdi_medias.downloaded.medias :
    media.kind => media.name...
  }
  description = "Medios descargados agrupados por tipo"
}

output "windows_isos_ids" {
  value       = data.isardvdi_medias.windows_isos_ready.medias[*].id
  description = "IDs de ISOs de Windows listos para usar"
}

# Ejemplo de uso condicional
output "first_ubuntu_iso_id" {
  value       = length(data.isardvdi_medias.ubuntu_isos.medias) > 0 ? data.isardvdi_medias.ubuntu_isos.medias[0].id : "No Ubuntu ISOs found"
  description = "ID del primer ISO de Ubuntu encontrado"
}

# Ejemplo: Listar medios en descarga
output "downloading_medias" {
  value = [
    for media in data.isardvdi_medias.all.medias :
    {
      name   = media.name
      status = media.status
    }
    if media.status == "Downloading" || media.status == "DownloadStarting"
  ]
  description = "Medios que están siendo descargados actualmente"
}

# Ejemplo: Medios por categoría y grupo
output "category_group_summary" {
  value = {
    for media in data.isardvdi_medias.category_medias.medias :
    media.name => {
      id       = media.id
      kind     = media.kind
      category = media.category
      group    = media.group
    }
  }
  description = "Resumen de medios de la categoría"
}
