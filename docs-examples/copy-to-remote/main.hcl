variable "server_public_ip" {
  type = string
}

variable "user_name" {
  type = string
}

variable "private_key" {
  type      = string
  sensitive = true
}

variable "payload" {
  type        = string
  description = "Local source directory or archive to copy."
}

variable "dest_dir" {
  type        = string
  description = "Destination directory on the remote host."
}

locals {
  archive = fileArchive(var.payload)
  conn = {
    host        = var.server_public_ip
    user        = var.user_name
    private_key = var.private_key
  }
}

resource "command_remote_copytoremote" "copy" {
  connection  = local.conn
  source      = local.archive
  remote_path = var.dest_dir
}

resource "command_remote_command" "find" {
  connection = local.conn
  create     = "find ${var.dest_dir}/${var.payload} | sort"
  triggers   = [local.archive]

  depends_on = [command_remote_copytoremote.copy]
}

output "remoteContents" {
  value = command_remote_command.find.stdout
}
