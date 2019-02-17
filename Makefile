
#  __     __         _       _     _
#  \ \   / /_ _ _ __(_) __ _| |__ | | ___ ___
#   \ \ / / _` | '__| |/ _` | '_ \| |/ _ \ __|
#    \ V / (_| | |  | | (_| | |_) | |  __\__ \
#     \_/ \__,_|_|  |_|\__,_|_.__/|_|\___|___/
#

base_version = 0.1.0

ADDR := 192.168.0.20:3000
SECRET := 123456789
VERSION := $(base_version)-$(shell git rev-parse --short=7 HEAD)

#   _____                    _
#  |_   _|_ _ _ __ __ _  ___| |_ ___
#    | |/ _` | '__/ _` |/ _ \ __/ __|
#    | | (_| | | | (_| |  __/ |_\__ \
#    |_|\__,_|_|  \__, |\___|\__|___/
#                 |___/

b:
	go build -o ./build/bin/ethstats-server

start: b
	./build/bin/ethstats-server --secret ${SECRET} --addr ${ADDR}

docker-build:
	docker build -t eskoltech/ethstats-server:$(VERSION) .

docker-start: docker-build
	docker run -it -p 3000:3000 --name ethstats eskoltech/ethstats-server:$(VERSION) --secret $(SECRET) --addr $(ADDR)
