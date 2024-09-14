FROM php:8.3 AS builder

# Install Composer
RUN apt update && apt install git subversion zip unzip -y && \
    curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer --version=2.1.12

# Install third party plugins with Composer
WORKDIR /app
COPY composer.json composer.json
RUN composer install

FROM wordpress:6.6.1-php8.3 as runner

# Install locales
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y locales && \
    sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && \
    dpkg-reconfigure --frontend=noninteractive locales && \
    update-locale LANG=en_US.UTF-8
ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US:en \
    LC_ALL=en_US.UTF-8

# WordPress config
COPY --chown=www-data:www-data wp-config.php .htaccess /var/www/html/

# Add third party plugins to the image
COPY --from=builder --chown=www-data:www-data /app/plugins /var/www/html/wp-content/plugins

# Add custom themes and plugins
COPY --chown=www-data:www-data themes/ /var/www/html/wp-content/themes/
COPY --chown=www-data:www-data plugins/ /var/www/html/wp-content/plugins/

# Custom init script
COPY entrypoint.sh /usr/local/bin/apache2-custom.sh
COPY init.sh /usr/local/bin/init.sh
RUN chmod +x /usr/local/bin/apache2-custom.sh && chmod +x /usr/local/bin/init.sh

CMD ["apache2-custom.sh"]
