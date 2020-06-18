# setup get cloud-image-ubuntu and its sha256sum
rm -rf original
mkdir original
curl -O "https://cloud-images.ubuntu.com/bionic/current/bionic-server-cloudimg-amd64.img"
mv bionic-server-cloudimg-amd64.img ./original
curl -O "https://cloud-images.ubuntu.com/bionic/current/SHA256SUMS"
mv SHA256SUMS ./original
