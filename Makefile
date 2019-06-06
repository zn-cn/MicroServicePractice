init:
    cd ..
	mv MicroServicePractice ${GOPATH}/src/Ethan/
	./pull.sh
	cd plugins
	docker-compose -f docker-compose.yml up -d
run:
	go run consignment/main.go &
	go run user/main.go &
	go run log/main.go &
	go run vessel/main.go &