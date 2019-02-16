ARG PHP_VERSION

FROM php:${PHP_VERSION}-fpm

ARG PHP_TIMEZONE
ARG PHP_MEMORY_LIMIT
ARG USER_ID
ARG GROUP_ID

MAINTAINER Vasilij Dusko <vasilij@prado.lt>

RUN apt-get update && \
    apt-get install --no-install-recommends --no-install-suggests -y gnupg1 \
    openssl \
    git \
    unzip \
    libzip-dev \
    nano \
    libpng-dev \
    libmagickwand-dev \
    curl

# Install PHP extensions. Type docker-php-ext-install to see available extensions
RUN docker-php-ext-configure zip --with-libzip && \
    docker-php-ext-install -j$(nproc) \
    pdo \
    pdo_mysql \
    opcache \
    zip \
    gd

RUN apt-get install -y libicu-dev \
    && docker-php-ext-configure intl \
    && docker-php-ext-install intl

RUN pecl install \
    xdebug \
    imagick

RUN docker-php-ext-enable \
    xdebug \
    imagick

# Set timezone
RUN ln -snf /usr/share/zoneinfo/${PHP_TIMEZONE} /etc/localtime && echo ${PHP_TIMEZONE} > /etc/timezone \
&& printf '[PHP]\ndate.timezone = "%s"\n' ${PHP_TIMEZONE} > $PHP_INI_DIR/conf.d/tzone.ini \
&& "date"

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
RUN curl -sL https://deb.nodesource.com/setup_9.x | bash && \
    apt-get install -y nodejs npm && \
    update-alternatives --install /usr/bin/node node /usr/bin/nodejs 10

# bugfix: remove cmdtest to install yarn correctly.
RUN apt-get remove -y cmdtest

# yarn package manager
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - && \
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list && \
    apt-get update && \
    apt-get install -y yarn

RUN apt-get remove --purge -y gnupg1 && apt-get -y --purge autoremove && rm -rf /var/lib/apt/lists/*

# Copy xdebug and php config.
#COPY conf.d/* /usr/local/etc/php/conf.d/

# Map user id from host user when it's provided
RUN if [ ! -z ${USER_ID} ]; then usermod -u ${USER_ID} www-data; fi
RUN if [ ! -z ${GROUP_ID} ]; then groupmod -g ${GROUP_ID} www-data; fi

RUN echo 'alias sf="php bin/console"' >> /var/www/.bashrc && \
    echo 'alias sf2="php app/console"' >> /var/www/.bashrc && \
    echo 'alias sf3="php bin/console"' >> /var/www/.bashrc && \
    echo 'alias sf4="php bin/console"' >> /var/www.bashrc

# set default user and working directory
USER ${USER_ID}
WORKDIR /var/www/project