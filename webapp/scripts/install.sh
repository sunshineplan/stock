#! /bin/bash

installSoftware() {
    apt -qq -y install nginx
    apt -qq -y -t $(lsb_release -sc)-backports install golang-go
}

installStock() {
    curl -Lo- https://github.com/sunshineplan/stock/archive/v1.0.tar.gz | tar zxC /var/www
    mv /var/www/stock* /var/www/stock
    cd /var/www/stock/webapp
    go build -ldflags "-s -w" -o stock
}

configStock() {
    read -p 'Please enter metadata server: ' server
    read -p 'Please enter VerifyHeader header: ' header
    read -p 'Please enter VerifyHeader value: ' value
    read -p 'Please enter unix socket(default: /run/stock.sock): ' unix
    [ -z $unix ] && unix=/run/stock.sock
    read -p 'Please enter host(default: 127.0.0.1): ' host
    [ -z $host ] && host=127.0.0.1
    read -p 'Please enter port(default: 12345): ' port
    [ -z $port ] && port=12345
    read -p 'Please enter log path(default: /var/log/app/stock.log): ' log
    [ -z $log ] && log=/var/log/app/stock.log
    mkdir -p $(dirname $log)
    sed "s,\$server,$server," /var/www/stock/webapp/config.ini.default > /var/www/stock/webapp/config.ini
    sed -i "s/\$header/$header/" /var/www/stock/webapp/config.ini
    sed -i "s/\$value/$value/" /var/www/stock/webapp/config.ini
    sed -i "s,\$unix,$unix," /var/www/stock/webapp/config.ini
    sed -i "s,\$log,$log," /var/www/stock/webapp/config.ini
    sed -i "s/\$host/$host/" /var/www/stock/webapp/config.ini
    sed -i "s/\$port/$port/" /var/www/stock/webapp/config.ini
}

setupsystemd() {
    cp -s /var/www/stock/webapp/scripts/stock.service /etc/systemd/system
    systemctl enable stock
    service stock start
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
    cp -s /var/www/stock/webapp/scripts/stock.cron /etc/cron.monthly/stock
    chmod +x /var/www/stock/webapp/scripts/stock.cron
}

setupNGINX() {
    cp -s /var/www/stock/webapp/scripts/stock.conf /etc/nginx/conf.d
    sed -i "s/\$domain/$domain/" /var/www/stock/webapp/scripts/stock.conf
    sed -i "s,\$unix,$unix," /var/www/stock/webapp/scripts/stock.conf
    service nginx reload
}

main() {
    read -p 'Please enter domain:' domain
    installSoftware
    installStock
    configStock
    setupsystemd
    writeLogrotateScrip
    createCronTask
    setupNGINX
}

main
