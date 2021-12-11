CHANGELOG for "Docker for Symfony"
======================

* 1.0.6 (2020-01-26)
    * FEATURE   [PHP] Added ability to use imagecreatefromjpeg
    * FEATURE   [PHP] Added ability to change PHP version using command ./docker set php_version 5.6 / 7.0 / 7.2 / 7.3 / 7.4
    * FEATURE   [PHP] Added ability enable / disable xdebug using command ./docker set xdebug yes / no
    * FEATURE   [PHP] Added ability enable / disable ionCube using command ./docker set ioncube yes / no
    * BUGFIX    [PHP] Returned mcrypt extension into d4d.lt installation script

* 1.0.7 (2020-02-16)
    * BUGFIX    [PHP]   Remove xdebug from pecl installation if disabled
    * FEATURE   [MYSQL] Import MySQL database copy without using a password form the console 
    * FEATURE   [PHP]   Added ability to change max_execution_time in file .env
    * FEATURE   [NGINX] Added ability to put rewrite rules into file config/nginx/rewrite/project.conf
    * FEATURE   [NGINX] Added ability select SF vs SF + PWA configuration
    * FEATURE   [PMA]   Added phpMyAdmin auto login functionality
    
* 1.0.8 (2020-04-19)
    * FEATURE   [D4D]   Added ability enable / disable redis in docker-compose.yml file
    * FEATURE   [D4D]   Added ability enable / disable rabbitmq in docker-compose.yml file
    
* 1.0.9 (2020-06-12)
    * FEATURE   [D4D]   Added option to select project type Symfony (SF), Symfony + PWA (SF_PWA), WordPress (WP)

* 1.0.10 (2020-07-22)
    * BUGFIX    [PHP]           Fix  missing redis module
    * FEATURE   [SUPERVISOR]    Added ability to switch on / off supervisor in the file .env
    * FEATURE   [ES]            Added ability to use elastic search.
    * FEATURE   [NGROK]         Added ability to use remote access.

* 1.0.11 (2020-08-26)
    * FEATURE   [D4D]           Layer optimization.
    * FEATURE   [D4D]           Added ability select MySQL vs MariaDB.
    * FEATURE   [MONGODB]       Added ability to use MongoDB.
    
* 1.0.12 (2020-11-23)
    * FEATURE   [MKCERT]        Added ability to have local developer SSL certificate.
    * FEATURE   [NGROK]         Added ability to have ngrok auth token in the file .env
    * FEATURE   [NGINX]         Updated Debian OS + nginx version.
    * FEATURE   [PHP]           The PHP 7.4 is set as default version.
    * FEATURE   [COMPOSER]      Added ability to select composer version 1.x vs 2.x
    * FEATURE   [WKHTMLTOPDF]   Added ability to select version 0.12.3, 0.12.4, 0.12.5 or 0.12.6

* 1.0.13 (2020-12-25)
  * FEATURE     [PHP]           PHP 8 is available now!
  * FEATURE     [PHP]           Added imagick changes required to support PHP 8.0
  * FEATURE     [D4D]           Added phpMyAdmin auto login functionality as root user.
  * FEATURE     [D4D]           Added ability to store .env.git without credentials (.env.secret.dist)
  * FEATURE     [D4D]           Added ability to run multiple dockers on the same time (network name external-d4d and needed to have nginx reverse proxy docker container)
  * FEATURE     [D4D]           Added ability to see running dockers IP address, version
  * FEATURE     [D4D]           Added ability to update main docker images from repository
  
* 1.0.14 (2021-02-19)
  * FEATURE     [D4D]           Added ability to see the main docker information (docker, docker-machine, docker-compose versions), main software versions, ip addresses every container
  * FEATURE     [PHP]           PHP 8 - RabbitMQ amqp extension support
  * FEATURE     [NGINX]         nginx version update
  * FEATURE     [COMPOSER]      Default composer version 2.x
  * FEATURE     [PHP]           Changed default setting for "post_max_size" and "upload_max_filesize"

* 1.0.15 (2021-04-21)
  * FEATURE     [PHP]           PHP - enabled XSL extension.
  * FEATURE     [PHP]           Added ability contributing to Symfony.
  * FEATURE     [D4D]           Added ability run default container on ./d4d start.
  * FEATURE     [NGINX]         NGINX version update.
  * FEATURE     [PHP]           Blackfire Agent v1 to v2 upgrade.

* 1.0.16 (2021-07-18)
  * FEATURE     [PHP]           Added PHP-CS-Fixer implementation.
  * FEATURE     [D4D]           Updated docker-compose version.
  * BUGFIX      [D4D]           NGINX + PWA network issue bug fix.
  * BUGFIX      [PHP]           PHP configuration bugfix.
  * FEATURE     [PHP]           Added Local PHP Security Checker implementation.
  * FEATURE     [NGINX]         NGINX version update.
  * FEATURE     [COMPOSER]      Composer version update.
  * FEATURE     [MySQL]         Default MySQL version as MariaDB.
  * FEATURE     [D4D]           Added webpack-analyzer.
  * FEATURE     [NGINX]         NGINX version update.
  * FEATURE     [D4D]           Added Symfony CLI implementation.

* 1.0.17 (2021-10-05)
  * FEATURE     [RabbitMQ]      Updated RabbitMQ version.
  * FEATURE     [D4D]           Updated Symfony CLI version.
  * FEATURE     [NGINX]         Updated Nginx version.
  * BUGFIX      [D4D]           MySQL 8.0 version not shown issue fix.
  * FEATURE     [D4D]           Docker code optimization.
  * FEATURE     [MariaDB]       Updated MariaDB version.
  * FEATURE     [D4D]           Updated composer version.
  * FEATURE     [PHP]           Added eslint implementation.
  * FEATURE     [PHP]           Added PHPStan integration.
  * FEATURE     [NGINX]         Nginx layer optimization.
  * FEATURE     [D4D]           mysql_dump.sh allow non-root user.

* 1.0.18 ()
  * FEATURE     [COMPOSER]      Updated composer version.
  * FEATURE     [D4D]           Added Deployer.org implementation.
  * FEATURE     [D4D]           Removed docker-machine from installation script.
  * BUGFIX      [D4D]           Fix container`s IP addresses.
  * FEATURE     [D4D]           Added docker container backup implementation.
  * FEATURE     [RabbitMQ]      Updated RabbitMQ version.
  * FEATURE     [PHP]           Updated Symfony CLI version.
  * FEATURE     [D4D]           Updated D4D installation script.
  * FEATURE     [NGINX]         Updated Nginx version.
  * FEATURE     [PHP]           Updated PHPStan version.
  * FEATURE     [PHP]           Updated default PHP version.
  * FEATURE     [PHP]           Updated Local PHP Security Checker version.