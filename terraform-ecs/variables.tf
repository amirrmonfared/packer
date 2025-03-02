variable "region" {
  type        = string
  default     = "eu-west-3"
  description = "AWS region for all resources"
}

variable "vpc_cidr" {
  type        = string
  default     = "10.0.0.0/16"
  description = "CIDR block for the VPC"
}

variable "public_subnet_cidr" {
  type        = string
  default     = "10.0.1.0/24"
  description = "CIDR block for the single public subnet"
}

variable "docker_image" {
  type        = string
  default     = "ghcr.io/amirrmonfared/packer/packer:master"
  description = "Docker image to deploy"
}

variable "service_name" {
  type        = string
  default     = "packer-service"
  description = "Name of the ECS service"
}
