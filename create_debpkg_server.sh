mkdir -p ./deb_pakage_server
mkdir -p ./deb_pakage_server/bin
mkdir -p ./deb_pakage_server/DEBIAN
touch  ./deb_pakage_server/DEBIAN/control

cat << EOF > ./deb_pakage_server/DEBIAN/control
Package: server
Version: 1.0
Architecture: amd64
Maintainer: ikarizxc
Description: server
Depends: systemd
EOF

chmod 755 ./deb_pakage_server/DEBIAN/postinst
chmod 755 ./deb_pakage_server/DEBIAN/prerm

cp ./server ./deb_pakage_server/bin/

dpkg-deb --build ./deb_pakage_server/ server.deb