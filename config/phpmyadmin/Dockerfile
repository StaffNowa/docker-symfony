FROM phpmyadmin/phpmyadmin

MAINTAINER Vasilij Dusko <vasilij@prado.lt>

ARG PMA_MEMORY_LIMIT
ARG PMA_POST_MAX_SIZE
ARG PMA_UPLOAD_MAX_FILESIZE
ARG DEFAULT_TIMEZONE

RUN printf 'memory_limit = %s\n' ${PMA_MEMORY_LIMIT} >> $PHP_INI_DIR/conf.d/phpmyadmin-custom.ini && \
    printf 'post_max_size = %s\n' ${PMA_POST_MAX_SIZE} >> $PHP_INI_DIR/conf.d/phpmyadmin-custom.ini && \
    printf 'upload_max_filesize = %s\n' ${PMA_UPLOAD_MAX_FILESIZE} >> $PHP_INI_DIR/conf.d/phpmyadmin-custom.ini && \
    printf 'session.gc_maxlifetime = %s\n' 86400 >> $PHP_INI_DIR/conf.d/phpmyadmin-custom.ini && \
    printf '[PHP]\ndate.timezone = "%s"\n' ${DEFAULT_TIMEZONE} > $PHP_INI_DIR/conf.d/tzone.ini