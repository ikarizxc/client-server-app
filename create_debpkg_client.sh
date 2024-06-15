mkdir -p ./deb_pakage_client
mkdir -p ./deb_pakage_client/bin
mkdir -p ./deb_pakage_client/DEBIAN
touch  ./deb_pakage_client/DEBIAN/control

cat << EOF > ./deb_pakage_client/DEBIAN/control
Package: client
Version: 1.0
Architecture: amd64
Maintainer: ikarizxc
Description: client
EOF

cp ./client ./deb_pakage_client/bin/

dpkg-deb --build ./deb_pakage_client/ client.deb