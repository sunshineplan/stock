#! /bin/bash

installSoftware() {
    apt -qq -y install nginx
    apt -qq -y -t $(lsb_release -sc)-backports install golang-go
}

installMyIP() {
    curl -Lo- https://github.com/sunshineplan/mystocks-go/archive/v1.0.tar.gz | tar zxC /var/www
    mv /var/www/mystocks-go* /var/www/mystocks-go
    cd /var/www/mystocks-go
    go build
}

configMyIP() {
    read -p 'Please enter metadata server: ' server
    read -p 'Please enter VerifyHeader header: ' header
    read -p 'Please enter VerifyHeader value: ' value
    read -p 'Please enter unix socket(default: /run/mystocks-go.sock): ' unix
    [ -z $unix ] && unix=/var/www/mystocks-go/mystocks-go.sock
    read -p 'Please enter host(default: 127.0.0.1): ' host
    [ -z $host ] && host=127.0.0.1
    read -p 'Please enter port(default: 12345): ' port
    [ -z $port ] && port=12345
    read -p 'Please enter log path(default: /var/log/app/mystocks-go.log): ' log
    [ -z $log ] && log=/var/log/app/mystocks-go.log
    mkdir -p $(dirname $log)
    sed "s,\$server,$server," /var/www/mystocks-go/config.ini.default > /var/www/mystocks-go/config.ini
    sed -i "s/\$header/$header/" /var/www/mystocks-go/config.ini
    sed -i "s/\$value/$value/" /var/www/mystocks-go/config.ini
    sed -i "s,\$unix,$unix," /var/www/mystocks-go/config.ini
    sed -i "s,\$log,$log," /var/www/mystocks-go/config.ini
    sed -i "s/\$host/$host/" /var/www/mystocks-go/config.ini
    sed -i "s/\$port/$port/" /var/www/mystocks-go/config.ini
}

setupsystemd() {
    cp -s /var/www/mystocks-go/scripts/mystocks-go.service /etc/systemd/system
    systemctl enable mystocks-go
    service mystocks-go start
}

writeLogrotateScrip() {
    if [ ! -f '/etc/logrotate.d/app' ]; then
	cat >/etc/logrotate.d/app <<-EOF
		/var/log/app/*.log {
		    copytruncate
		    rotate 12
		    compress
		    delaycompress
		    missingok
		    notifempty
		}
		EOF
    fi
}

createCronTask() {
    cp -s /var/www/mystocks-go/scripts/mystocks-go.cron /etc/cron.monthly/mystocks-go
    chmod +x /var/www/mystocks-go/scripts/mystocks-go.cron
}

setupNGINX() {
    cp -s /var/www/mystocks-go/scripts/mystocks-go.conf /etc/nginx/conf.d
    sed -i "s/\$domain/$domain/" /var/www/mystocks-go/scripts/mystocks-go.conf
    sed -i "s,\$unix,$unix," /var/www/mystocks-go/scripts/mystocks-go.conf
    service nginx reload
}

main() {
    read -p 'Please enter domain:' domain
    installSoftware
    installMyIP
    configMyIP
    setupsystemd
    writeLogrotateScrip
    createCronTask
    setupNGINX
}

main
