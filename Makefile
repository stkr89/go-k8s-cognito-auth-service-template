build:
	docker build -t stkr89/go-k8s-cognito-auth-microservice-template:latest .

build-push:
	docker build -t stkr89/go-k8s-cognito-auth-microservice-template:latest .
	docker push stkr89/go-k8s-cognito-auth-microservice-template:latest