# Docker Symfony (PHP-FPM - NGINX - MySQL)

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

If you are running on Windows:
```
https://hub.docker.com/editions/community/docker-ce-desktop-windows
```

## Install the latest Docker Machine

If you are running on Linux:
```
base=https://github.com/docker/machine/releases/download/v0.16.0 &&
curl -L $base/docker-machine-$(uname -s)-$(uname -m) >/tmp/docker-machine &&
install /tmp/docker-machine /usr/local/bin/docker-machine
```

If you are running on Windows:
```
$ base=https://github.com/docker/machine/releases/download/v0.16.0 &&
  mkdir -p "$HOME/bin" &&
  curl -L $base/docker-machine-Windows-x86_64.exe > "$HOME/bin/docker-machine.exe" &&
  chmod +x "$HOME/bin/docker-machine.exe"
```

## Install the latest Docker Compose

If you are running on Linux:
```
curl -L "https://github.com/docker/compose/releases/download/1.23.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose && chmod +x /usr/local/bin/docker-compose
```

If you are running on Windows:
Docker for Windows and Docker Toolbox already include Compose with other Docker apps, so most Windows users do not need to install Compose separately.

# Configuration
1. Create a `.env` from the `.env.dist` file. Adapt it according to your symfony application

    ```bash
    cp .env.dist .env
    ```

2. Build / run containers with (with and without detached mode)
    ```
    docker-compose build
    docker-compose up -d
    ```

3. Update your system host file (add symfony.local)
   ```bash
   # UNIX only: get containers IP address and update host (replace IP according to your configuration) (on Windows, edit C:\Windows\System32\drivers\etc\hosts)
   $ sudo echo $(docker network inspect bridge | grep Gateway | grep -o -E '[0-9\.]+') "symfony.local" >> /etc/hosts
   ```
   
   **Note:** For **OS X**, please take a look [here](https://docs.docker.com/docker-for-mac/networking/) and for **Windows** read [this](https://docs.docker.com/docker-for-windows/#/step-4-explore-the-application-and-run-examples) (4th step).
   
4. Prepare Symfony app
    
    1. Create a new user
        ```
        docker-compose exec db bash
        mysql -u root -p
        
        mysql> use mysql;
        mysql> CREATE USER 'db_user'@'%' IDENTIFIED BY 'db_password';
        mysql> GRANT ALL PRIVILEGES ON db_name.* TO 'db_user'@'%';
        mysql> FLUSH PRIVILEGES;
        ```
    2. Update
  
        a) SF2, SF3: app/config/parameters.yml
          
        ```
        # ./project/app/config/parameters.yml
        parameters:
            database_host:     db
            database_port:     ~
            database_name:     db_name
            database_user:     db_user
            database_password: db_password
        ```
    
        b) SF4: .env
        ```
        DATABASE_URL=mysql://db_user:db_password@172.18.0.4:3306/db_name
        ```
    3. Composer install & create database
        ```bash
        $ docker-compose exec php bash
        $ composer create-project symfony/website-skeleton my-project
            
        # Symfony2
        $ sf doctrine:database:create
        $ sf doctrine:schema:update --force
        # Only if you have `doctrine/doctrine-fixtures-bundle` installed
        $ sf doctrine:fixtures:load --no-interaction
            
        # Symfony3
        $ sf3 doctrine:database:create
        $ sf3 doctrine:schema:update --force
        # Only if you have `doctrine/doctrine-fixtures-bundle` installed
        $ sf3 doctrine:fixtures:load --no-interaction
    
        # Symfony4
        $ sf4 doctrine:database:create
        $ sf4 doctrine:schema:update --force
        # Only if you have `doctrine/doctrine-fixtures-bundle` installed
        $ sf4 doctrine:fixtures:load --no-interaction
        ```
5. Enjoy :-)
    
    ## Usage
    
    Just run `docker-compose up -d`, then:
    
    * Symfony app: visit [symfony.local](http://symfony.local)  
    * Symfony dev mode: visit [symfony.local/app_dev.php](http://symfony.local/app_dev.php)  
    * Logs (files location): logs/nginx and logs/symfony
    
    ## How it works?
    
    Have a look at the `docker-compose.yml` file, here are the `docker-compose` built images:
    
    * `db`: This is the MySQL database container,
    * `php`: This is the PHP-FPM container in which the application volume is mounted,
    * `nginx`: This is the Nginx webserver container in which application volume is mounted too,
    
    This results in the following running containers:
    
    ```bash
    $ docker-compose ps
                  Name                             Command               State                    Ports                  
    ---------------------------------------------------------------------------------------------------------------------
    symfony_db_1_f0586075033b           docker-entrypoint.sh mysqld      Up      0.0.0.0:3306->3306/tcp, 33060/tcp       
    symfony_nginx_1_3c244ea6ff7b        nginx -g daemon off;             Up      0.0.0.0:443->443/tcp, 0.0.0.0:80->80/tcp
    symfony_php_1_916d0314f3e0          docker-php-entrypoint php-fpm    Up      9000/tcp                                
    symfony_phpmyadmin_1_a5ce79ef63bd   /run.sh supervisord -n -j  ...   Up      0.0.0.0:8080->80/tcp, 9000/tcp 
    ```
