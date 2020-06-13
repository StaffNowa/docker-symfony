ARG PHP_VERSION

FROM php:${PHP_VERSION}-fpm

ARG DEFAULT_TIMEZONE
ARG PHP_MEMORY_LIMIT
ARG PHP_MAX_EXECUTION_TIME
ARG PHP_UPLOAD_MAX_FILESIZE

ARG NODE_JS_VERSION

ARG USER_ID
ARG GROUP_ID

ARG XDEBUG_REMOTE_PORT
ARG XDEBUG_REMOTE_AUTOSTART
ARG XDEBUG_REMOTE_ENABLE
ARG XDEBUG_REMOTE_HOST
ARG XDEBUG_REMOTE_CONNECT_BACK

MAINTAINER Vasilij Dusko <vasilij@prado.lt>

RUN apt-get update && \
    apt-get install --no-install-recommends --no-install-suggests -y \
    gnupg1 openssl git unzip libzip-dev nano libpng-dev libmagickwand-dev curl xauth xvfb openssh-client less inkscape cron exiftool libicu-dev libmcrypt-dev librabbitmq-dev librabbitmq4 supervisor

# Install PHP extensions. Type docker-php-ext-install to see available extensions
RUN docker-php-ext-configure gd --with-freetype --with-jpeg && docker-php-ext-configure intl && docker-php-ext-install -j$(nproc) pdo pdo_mysql opcache zip gd mysqli exif bcmath calendar intl

RUN pecl install \
    imagick apcu apcu_bc mcrypt

RUN docker-php-ext-enable \
    imagick mysqli calendar exif bcmath apcu mcrypt

RUN ln -snf /usr/share/zoneinfo/${DEFAULT_TIMEZONE} /etc/localtime && echo ${DEFAULT_TIMEZONE} > /etc/timezone && \
    printf '[PHP]\ndate.timezone = "%s"\n' ${DEFAULT_TIMEZONE} > $PHP_INI_DIR/conf.d/tzone.ini && \
    printf '[CUSTOM]\nmemory_limit = "%s"\n' ${PHP_MEMORY_LIMIT} >> $PHP_INI_DIR/conf.d/custom.ini && \
    printf '\npost_max_size = %s\n' ${PHP_UPLOAD_MAX_FILESIZE} >> $PHP_INI_DIR/conf.d/custom.ini && \
    printf '\nupload_max_filesize = %s\n' ${PHP_UPLOAD_MAX_FILESIZE} >> $PHP_INI_DIR/conf.d/custom.ini && \
    printf '\nmax_execution_time = %s\n' ${PHP_MAX_EXECUTION_TIME} >> $PHP_INI_DIR/conf.d/custom.ini

RUN echo 'extension=apc.so' >> $PHP_INI_DIR/conf.d/docker-php-ext-apcu.ini

# Install composer
RUN php -r "copy('https://getcomposer.org/installer', 'composer-setup.php');"
RUN php composer-setup.php
RUN php -r "unlink('composer-setup.php');"
RUN mv composer.phar /usr/local/bin/composer
RUN mkdir -p /var/www/.composer
COPY composer/config.json /var/www/.composer/config.json
RUN chown -R ${USER_ID}:${GROUP_ID} /var/www/.composer && \
    chown ${USER_ID}:${GROUP_ID} /var/www

# npm & node
RUN mkdir -p /var/www/.npm && \
    chown -R ${USER_ID}:${GROUP_ID} /var/www/.npm && \
    printf 'Package: *\nPin: origin deb.nodesource.com\nPin-Priority: 600' > /etc/apt/preferences.d/nodejs && \
    curl -sL https://deb.nodesource.com/setup_${NODE_JS_VERSION} | bash && \
    apt-get install -y nodejs && \
    update-alternatives --install /usr/bin/node node /usr/bin/nodejs 10 && \
    npm install -g grunt-cli yargs async sass gulp requirejs pm2 uglify-js

# Yarn package manager
# bugfix: remove cmdtest to install yarn correctly.
RUN apt-get remove -y cmdtest && curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - && echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list && apt-get update && apt-get install -y yarn

# Install wkhtmltopdf
RUN curl -o wkhtmltox-0.12.4_linux-generic-amd64.tar.xz -sL https://github.com/wkhtmltopdf/wkhtmltopdf/releases/download/0.12.4/wkhtmltox-0.12.4_linux-generic-amd64.tar.xz  && echo '3f923f425d345940089e44c1466f6408b9619562 wkhtmltox-0.12.4_linux-generic-amd64.tar.xz' | sha1sum -c - && tar -xvf wkhtmltox-0.12.4_linux-generic-amd64.tar.xz && cp wkhtmltox/lib/* /usr/lib/ && cp wkhtmltox/bin/* /usr/bin/ && cp -r wkhtmltox/share/man/man1 /usr/share/man/ && chmod a+x /usr/bin/wkhtmltopdf && chmod a+x /usr/bin/wkhtmltoimage

RUN apt-get remove --purge -y gnupg1 && apt-get -y --purge autoremove && rm -rf /var/lib/apt/lists/*

# Copy xdebug and php config.
COPY conf.d/* /usr/local/etc/php/conf.d/

# Map user id from host user when it's provided
RUN if [ ! -z ${USER_ID} ] && [ ${USER_ID} -ne 0 ]; then usermod -u ${USER_ID} www-data; fi
RUN if [ ! -z ${GROUP_ID} ] && [ ${GROUP_ID} -ne 0 ]; then groupmod -g ${GROUP_ID} www-data; fi

RUN echo 'alias sf="php bin/console"' >> /var/www/.bashrc && \
    echo 'alias sf2="php app/console"' >> /var/www/.bashrc && \
    echo 'alias sf3="php bin/console"' >> /var/www/.bashrc && \
    echo 'alias sf4="php bin/console"' >> /var/www/.bashrc && \
    echo 'alias sf5="php bin/console"' >> /var/www.bashrc

# set default user and working directory
USER ${USER_ID}
WORKDIR /var/www/project