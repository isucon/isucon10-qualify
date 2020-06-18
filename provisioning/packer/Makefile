initimg:
	docker-compose -f docker-compose-img.yaml up --build
initssh:
	docker-compose -f docker-compose-ssh.yaml up --build

qemu/ssh:
	qemu-system-x86_64 \
    -m 2G -nographic \
	-net nic \
	-net user,hostfwd=tcp::10022-:22 \
    -drive file=./output/ubuntu,format=qcow2 \
    -drive file=./register-ssh-config/user-data.img,format=raw

packer:
	rm -rf "./output"
	packer build web.yaml

all: initssh packer qemu/ssh

qemu/tmp:
	rm ./original/tmp
	cp ./original/bionic-server-cloudimg-amd64.img ./original/tmp
	docker-compose -f docker-compose-ssh.yaml up --build
	qemu-system-x86_64 \
    -m 2G -nographic \
	-net nic \
	-net user,hostfwd=tcp::10022-:22 \
    -drive file=./original/tmp,format=qcow2 \
    -drive file=./register-ssh-config/user-data.img,format=raw
