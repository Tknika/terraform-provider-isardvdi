variable "isard_endpoint" {
  description = "Endpoint del servidor Isard VDI"
  type        = string
  default     = "isard.example.com"
}

variable "auth_method" {
  description = "Método de autenticación (token o form)"
  type        = string
  default     = "token"
}

variable "isard_token" {
  description = "Token de API para autenticación"
  type        = string
  sensitive   = true
}

variable "search_username" {
  description = "Nombre del usuario a buscar"
  type        = string
  default     = "admin"
}

variable "group_id" {
  description = "ID del grupo para filtrar usuarios"
  type        = string
  default     = "default-default-default"
}
