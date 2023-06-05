terraform {
	required_providers {
		docker = {
			source = "kreuzwerker/docker"
			version = "3.0.2"
		}
	}
}

provider "docker" {
	 host = "unix:///var/run/docker.sock"
}

provider "kubernetes" {
	config_path = "~/.kube/config"
	config_context = "docker-desktop"
}

resource "kubernetes_namespace" "passman" {
	metadata {
		name = "passman"
	}
}

resource "kubernetes_deployment" "passman-server" {
	metadata {
		name = "passman-server"
		labels = {
			app = "passman-server"
		}
	}

	spec {
		replicas = 1
		selector {
			match_labels = {
				app = "passman-server"
			}
		}
		template {
			metadata {
				name = "passman-server"
				labels = {
					app = "passman-server"
				}
			}
			spec {
				container {
					name = "passman-server"
					image = "passman:latest"
					image_pull_policy = "Never"
					volume_mount {
						mount_path = "/var/passman"
						name = kubernetes_persistent_volume.passman-data.metadata.0.name
					}
				}

				volume {
					name = kubernetes_persistent_volume.passman-data.metadata.0.name
					persistent_volume_claim {
						claim_name = kubernetes_persistent_volume_claim.passman-data.metadata.0.name
					}
				}
			}
		}
	}
}

resource "kubernetes_service" "passman-server" {
	metadata {
		name = "passman-server"
	}

	spec {
		selector = {
			app = kubernetes_deployment.passman-server.metadata.0.labels.app
		}

		port {
			port = 8080
			target_port = 8080
		}

		type = "LoadBalancer"
	}
}

resource "kubernetes_persistent_volume" "passman-data" {
	metadata {
		name = "passman-data"
	}

	spec {
		capacity = {
			storage = "2Gi"
		}
		access_modes = ["ReadWriteMany"]
		persistent_volume_source {
			host_path {
				path = "/Users/divij/workspace/local/go/passwords/data"
			}
		}
		storage_class_name = "manual"
	}
}

resource "kubernetes_persistent_volume_claim" "passman-data" {
	metadata {
		name = "passman-data"
	}

	spec {
		access_modes = ["ReadWriteMany"]
		resources {
			requests = {
				storage = "2Gi"
			}
		}
		volume_name = kubernetes_persistent_volume.passman-data.metadata.0.name
		storage_class_name = "manual"
	}
}
