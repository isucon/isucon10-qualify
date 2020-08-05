initssh:
	docker-compose -f docker-compose-ssh.yaml up --build

qemu/ssh:
	qemu-system-x86_64 \
    -m 2G -nographic \
	-net nic \
	-net user,hostfwd=tcp::10022-:22 \
    -drive file=./output-$(shell basename $@)/isucon10-qualify,format=qcow2 \
    -drive file=./register-ssh-config/user-data.img,format=raw

packer/%:
	rm -rf "./output-$(shell basename $@)"
	packer build -var 'env=$(shell basename $@)' packer.json

all/%: initssh/% packer/% qemu/ssh/%

