FROM phpmyadmin/phpmyadmin

MAINTAINER Vasilij Dusko <vasilij@prado.lt>

ARG PMA_CLIENT_MAX_BODY_SIZE
ARG PMA_MEMORY_LIMIT
ARG PMA_POST_MAX_SIZE
ARG PMA_UPLOAD_MAX_FILESIZE
ARG DEFAULT_TIMEZONE

RUN apk add tzdata

RUN sed -i 's#client_max_body_size .*M#client_max_body_size '"${PMA_CLIENT_MAX_BODY_SIZE}"'#g' /etc/nginx.conf && \
    sed -i 's#memory_limit = .*M#memory_limit = '"${PMA_MEMORY_LIMIT}"'#g' ${PHP_INI_DIR}/conf.d/php-phpmyadmin.ini && \
    sed -i 's#post_max_size = .*M#post_max_size = '"${PMA_POST_MAX_SIZE}"'#g' ${PHP_INI_DIR}/conf.d/php-phpmyadmin.ini && \
    sed -i 's#upload_max_filesize = .*M#upload_max_filesize = '"${PMA_UPLOAD_MAX_FILESIZE}"'#g' ${PHP_INI_DIR}/conf.d/php-phpmyadmin.ini && \
    printf '[PHP]\ndate.timezone = "%s"\n' ${DEFAULT_TIMEZONE} > $PHP_INI_DIR/conf.d/tzone.ini