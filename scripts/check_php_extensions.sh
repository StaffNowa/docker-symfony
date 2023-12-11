#!/bin/bash

extensions=("zip" "int" "opcache" "imagick" "gd" "amqp" "mongodb" "imap" "mcrypt" "redis" "xdebug")

is_extension_enabled() {
    php -m | grep -q "$1"
}

for extension in "${extensions[@]}"; do
    if is_extension_enabled "$extension"; then
        echo "$extension: ON"
    else
        echo "$extension: OFF"
    fi
done