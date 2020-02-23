FROM phpmyadmin/phpmyadmin

MAINTAINER Vasilij Dusko <vasilij@prado.lt>

ARG PHP_MEMORY_LIMIT
ARG PHP_POST_MAX_SIZE
ARG PHP_UPLOAD_MAX_FILESIZE
ARG DEFAULT_TIMEZONE

RUN printf 'memory_limit = %s\n' ${PHP_MEMORY_LIMIT} >> $PHP_INI_DIR/conf.d/phpmyadmin-misc.ini && \
    printf 'post_max_size = %s\n' ${PHP_POST_MAX_SIZE} >> $PHP_INI_DIR/conf.d/phpmyadmin-custom.ini && \
    printf 'upload_max_filesize = %s\n' ${PHP_UPLOAD_MAX_FILESIZE} >> $PHP_INI_DIR/conf.d/phpmyadmin-custom.ini && \
    printf 'session.gc_maxlifetime = %s\n' 86400 >> $PHP_INI_DIR/conf.d/phpmyadmin-custom.ini && \
    printf '[PHP]\ndate.timezone = "%s"\n' ${DEFAULT_TIMEZONE} > $PHP_INI_DIR/conf.d/tzone.ini