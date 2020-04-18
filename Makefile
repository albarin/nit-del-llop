run-docker:
	docker build -t albarin/nit-del-llop .
	docker run albarin/nit-del-llop

build:
	go build -o bin/poster -v cmd/poster/*.go

run: build
	PORT=3000 \
	SECRET_TOKEN=ploc \
 	./bin/poster

deploy:
	git push heroku master -f