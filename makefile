DockerComposeStart: ImageBuild
	docker compose up

ImageBuild: GitCloneRep
	cd ./PostCommentAuthorization && docker build -t migrationcommentpost:1 -f dockerfilemigration . && docker build -t backendpostcomment:1 .

GitCloneRep: MakeSSO
	git clone https://github.com/Senkoker/PostCommentAuthorization

MakeSSO: GitCloneRepSSO
	cd ./SSO_service && make StartDockerCompose && cd ..

GitCloneRepSSO:
	git clone https://github.com/Senkoker/SSO_service

Down:
	docker compose down