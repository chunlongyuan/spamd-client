docker_compose_files=-f docker-compose.yml

up:
	docker-compose $(docker_compose_files) up -d

down:
	docker-compose $(docker_compose_files) down --remove-orphans


ps:
	docker-compose $(docker_compose_files) ps 

restart:
	$(MAKE) down up



clean:
	docker system prune -f


