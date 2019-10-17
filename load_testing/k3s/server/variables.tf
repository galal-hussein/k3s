variable "aws_access_key" {}

variable "aws_secret_key" {}

variable "db_instance_type" {
    default = ""
}
variable "k3s_ha" {
    default = 1
}
variable "server_instance_type" {
  # default = "c4.8xlarge"
}
variable "k3s_server_args" {
  default = ""
}
variable "prom_host" {
  default = ""
}
variable "graf_host" {
  default = ""
}
