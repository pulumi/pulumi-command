resource "command_local_command" "random" {
  create = "openssl rand -hex 16"
}

output "output" {
  value = command_local_command.random.stdout
}
