# Ejemplo de uso del recurso isard_qos_net

# Crear un perfil de QoS de red
resource "isard_qos_net" "standard" {
  name             = "Standard Network QoS"
  description      = "Perfil estándar de QoS para red"
  average_download = 10240  # 10 MB/s en KB/s
  average_upload   = 5120   # 5 MB/s en KB/s
  peak_download    = 20480  # 20 MB/s en KB/s
  peak_upload      = 10240  # 10 MB/s en KB/s
  burst_download   = 102400 # 100 MB en KB
  burst_upload     = 51200  # 50 MB en KB
}

# Output para ver la información del QoS
output "qos_info" {
  value = {
    id               = isard_qos_net.standard.id
    name             = isard_qos_net.standard.name
    average_download = isard_qos_net.standard.average_download
    average_upload   = isard_qos_net.standard.average_upload
  }
}
