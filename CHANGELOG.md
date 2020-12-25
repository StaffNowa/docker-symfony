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