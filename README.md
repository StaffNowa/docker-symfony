# Docker for Symfony (PHP-FPM - NGINX - MySQL)

Docker symfony gives you everything you need for developing Symfony application. This complete stack run with docker and docker-compose.

# Installation

## Install the latest Docker CE version
If you are running on Linux:
```
curl -fsSL https://get.docker.com -o get-docker.sh &&
sh get-docker.sh
```

If you would like to use Docker as a non-root user, you should now consider
adding your user to the "docker" group with something like:

```
sudo usermod -aG docker ${USER}
```

## Install the latest Docker Compose

If you are running on Linux:
```
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose && sudo chmod +x /usr/local/bin/docker-compose
```

# Configuration
1. Create a `.env` from the `.env.dist` file. Adapt it according to your symfony application

    ```bash
    cp .env.dist .env
    ```
   
2. Build / run containers
    ```
    ./d4d start
    ```  

3. Update your system host file (add symfony.local)
   ```bash
   # UNIX only: get containers IP address and update host (replace IP according to your configuration)
   $ sudo echo $(docker network inspect bridge | grep Gateway | grep -o -E '[0-9\.]+') "symfony.local" >> /etc/hosts
   ```
   
   **Note:** For **OS X**, please take a look [here](https://docs.docker.com/docker-for-mac/networking/).
   
4. Prepare Symfony app
    1. Get your logins to access MySQL server
    `./d4d passwd show`
    
    2. Update
  
        a) SF2, SF3: app/config/parameters.yml
          
        ```
        # ./project/app/config/parameters.yml
        parameters:
            database_host:     mysql
            database_port:     ~
            database_name:     db_name
            database_user:     db_user
            database_password: db_password (random password)
        ```
    
        b) SF4, SF5, SF6: .env
        ```
        DATABASE_URL=mysql://db_user:db_password@mysql:3306/db_name
        MAILER_URL=smtp://mailhog:1025
        ```
    3. Composer install & create database
        ```bash
        $ docker-compose exec php bash
        $ composer create-project symfony/website-skeleton my-project
            
        # Symfony 2
        $ sf doctrine:database:create
        $ sf doctrine:schema:update --force
        # Only if you have `doctrine/doctrine-fixtures-bundle` installed
        $ sf doctrine:fixtures:load --no-interaction
            
        # Symfony 3
        $ sf3 doctrine:database:create
        $ sf3 doctrine:schema:update --force
        # Only if you have `doctrine/doctrine-fixtures-bundle` installed
        $ sf3 doctrine:fixtures:load --no-interaction
    
        # Symfony 4
        $ sf4 doctrine:database:create
        $ sf4 doctrine:schema:update --force
        # Only if you have `doctrine/doctrine-fixtures-bundle` installed
        $ sf4 doctrine:fixtures:load --no-interaction
        
        # Symfony 5
        $ sf5 doctrine:database:create
        $ sf5 doctrine:schema:update --force
        # Only if you have `doctrine/doctrine-fixtures-bundle` installed
        $ sf5 doctrine:fixtures:load --no-interaction
       
        # Symfony 6
        $ sf6 doctrine:database:create
        $ sf6 doctrine:schema:update --force
        # Only if you have `doctrine/doctrine-fixtures-bundle` installed
        $ sf6 doctrine:fixtures:load --no-interaction
       ```
5. Enjoy :-)
    
    ## Usage
    
    Just run `./d4d start`, then:
    
    * Symfony app: visit [symfony.local](http://symfony.local)  
    * Symfony dev mode: visit [symfony.local/app_dev.php](http://symfony.local/app_dev.php)  
    * Logs (files location): logs/nginx and logs/symfony
    
    ## How it works?
    
    Have a look at the `docker-compose.yml` file, here are the `docker-compose` built images:
 
    * `nginx`: Nginx is one of the most popular web servers in the world and responsible for hosting some of the largest and highest-traffic sites on the internet. It is more resource-friendly than Apache in most cases and can be used as a web server or reverse proxy.
    * `php`: PHP is a popular general-purpose scripting language that is especially suited to web development.
    * `mysql:` MySQL is the most popular relational database management system.
    * `phpmyadmin`: phpMyAdmin was created so that users can interact with MySQL / MariaDB through a web interface.
    * `mailhog`: MailHog is an email testing tool for developers.
    * `redis`: Redis is an open source (BSD licensed), in-memory data structure store, used as a database, cache and message broker. It supports data structures such as strings, hashes, lists, sets, sorted sets with range queries, bitmaps, hyperloglogs, geospatial indexes with radius queries and streams. Redis has built-in replication, Lua scripting, LRU eviction, transactions and different levels of on-disk persistence, and provides high availability via Redis Sentinel and automatic partitioning with Redis Cluster. 
     
    This results in the following running containers:
    
    ```bash
    $ docker-compose ps
               Name                          Command               State                       Ports                     
    ---------------------------------------------------------------------------------------------------------------------
    docker-symfony_mailhog_1      MailHog                          Up      0.0.0.0:1025->1025/tcp, 0.0.0.0:8025->8025/tcp
    docker-symfony_mysql_1        docker-entrypoint.sh --cha ...   Up      0.0.0.0:3306->3306/tcp                        
    docker-symfony_nginx_1        nginx -g daemon off;             Up      0.0.0.0:80->80/tcp                            
    docker-symfony_php_1          docker-php-entrypoint php-fpm    Up      9000/tcp                                      
    docker-symfony_phpmyadmin_1   /docker-entrypoint.sh apac ...   Up      0.0.0.0:8080->80/tcp                          
    docker-symfony_redis_1        docker-entrypoint.sh redis ...   Up      6379/tcp          
    ```


# JetBrains support us! 

<p align="center">
  <a href="https://www.jetbrains.com" target="_blank">
      <img src="https://account.jetbrains.com/static/images/jetbrains-logo-inv.svg" width="120" title="hover text">
  </a>
</p>

# Digital Ocean support us!

<p align="center">
<a href="https://www.digitalocean.com/?refcode=c832ae92ce6c&utm_campaign=Referral_Invite&utm_medium=Referral_Program&utm_source=badge"><img src="https://web-platforms.sfo2.cdn.digitaloceanspaces.com/WWW/Badge%201.svg" alt="DigitalOcean Referral Badge" /></a>
</p>

# Blackfire.io support us!
<p align="center">
  <a href="https://www.blackfire.io" target="_blank">
      <img src="support/blackfire_io.png" width="120" title="hover text">
  </a>
</p>