package commands

import (
	"docker-symfony/util"
	"fmt"
	"github.com/symfony-cli/console"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"
)

var envFile = util.GetCurrentDir() + "/.env"
var envDistFile = envFile + ".dist"

var envSecretFile = envFile + ".secret"

var startCmd = &console.Command{
	Category: "d4d",
	Name:     "start",
	Aliases:  []*console.Alias{{Name: "start"}},
	Usage:    "",
	Action: func(c *console.Context) error {
		doChecks()
		doBuildNginxConf()
		doBuildMySQLConf()
		doBuild()
		doBeforeStart()
		//start()

		return nil
	},
}

func doChecks() {
	if !util.IsCommandExist("docker") {
		fmt.Println("Docker not found. Please install it via https://docs.docker.com/install/")
		os.Exit(1)
	}

	if !util.IsCommandExist("docker-compose") {
		fmt.Println("Docker compose not found. Please install it via https://docs.docker.com/compose/install/")
		os.Exit(1)
	}

	if !util.FileExists(envFile) {
		if util.FileExists(envDistFile) {
			util.Copy(envDistFile, envFile)

			fileData, err := ioutil.ReadFile(envSecretFile)
			if err != nil {
				os.Exit(1)
			}

			fileString := string(fileData)
			util.AppendFile(envFile, fileString)
		} else {
			fmt.Sprintf("The %s file does not exist. Project setup will not work.", envDistFile)
			os.Exit(1)
		}
	}

	// Always validate user id and group id before start using .env file
	if runtime.GOOS != "darwin" {
		util.Sed("USER_ID=.*", fmt.Sprintf("USER_ID=%d", os.Getuid()), envFile)
		util.Sed("GROUP_ID=.*", fmt.Sprintf("GROUP_ID=%d", os.Getgid()), envFile)
	} else {
		util.Sed("USER_ID=.*", fmt.Sprintf("USER_ID=%d", os.Getuid()), envFile)
		util.Sed("GROUP_ID=.*", fmt.Sprintf("GROUP_ID=%d", os.Getuid()), envFile)
	}
	util.Sed("MYSQL_ROOT_PASSWORD=root", fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", util.GeneratePassword(20)), envSecretFile)
	util.Sed("MYSQL_PASSWORD=db_password", fmt.Sprintf("MYSQL_PASSWORD=%s", util.GeneratePassword(20)), envSecretFile)
	util.Sed("MONGODB_ROOT_PASSWORD=root", fmt.Sprintf("MONGODB_ROOT_PASSWORD=%s", util.GeneratePassword(20)), envSecretFile)
	util.Sed("MONGODB_PASSWORD=db_password", fmt.Sprintf("MONGODB_PASSWORD=%s", util.GeneratePassword(20)), envSecretFile)

	fileData, err := ioutil.ReadFile(envFile)
	if err != nil {
		os.Exit(1)
	}

	if !strings.Contains(string(fileData), "# .env.secret") {
		fileData, err := ioutil.ReadFile(envSecretFile)
		if err != nil {
			os.Exit(1)
		}

		fileString := string(fileData)
		util.AppendFile(envFile, fileString)
	}

	util.LoadEnvFile(envFile)

	sshKeyPath := os.Getenv("HOME") + "/.ssh"

	// Ensure all folders exists
	var dirs = []string{
		os.Getenv("PROJECT_PATH"),
		os.Getenv("SF_COMMUNITY_PATH"),
		os.Getenv("SF_COMMUNITY_PATH") + "/symfony",
		os.Getenv("SF_COMMUNITY_PATH") + "/recipes",
		os.Getenv("SF_COMMUNITY_PATH") + "/symfony-docs",
		"./.composer",
		"./.composer/cache",
		sshKeyPath,
		os.Getenv("MYSQL_DUMP_PATH"),
	}
	util.Mkdir(dirs, 0755)

	sshPrivateKey := sshKeyPath + "/id_rsa_test"
	sshPublicKey := sshPrivateKey + ".pub"
	if !util.FileExists(sshPrivateKey) {
		pubKey, privKey, _ := util.MakeSSHKeyPair()
		if pubKey != "" && privKey != "" {
			f, _ := os.Create(sshPublicKey)

			// get current user
			user, err := user.Current()
			if err != nil {
				log.Fatalf(err.Error())
			}

			f.WriteString(fmt.Sprintf("%s %s@d4d.lt\n", strings.Replace(pubKey, "\n", "", 1), user.Username))

			f2, _ := os.Create(sshPrivateKey)
			f2.WriteString(privKey)
		}
	}

	if !util.FileExists(sshKeyPath + "/known_hosts") {
		os.Create(sshKeyPath + "/known_hosts")
	}

	os.MkdirAll(os.Getenv("NGINX_SSL_PATH"), 0755)
	os.Mkdir(os.Getenv("NGINX_LOG_PATH"), 0755)
	os.Mkdir(os.Getenv("MYSQL_DATA_PATH"), 0755)
	os.Mkdir(os.Getenv("USER_CONFIG_PATH"), 0755)
	os.Mkdir(os.Getenv("MONGODB_LOG_PATH"), 0755)
	os.Mkdir(os.Getenv("MONGODB_DATA_PATH"), 0755)

	// Create a file if it does not exist
	util.CreateFileIfNotExists(os.Getenv("USER_CONFIG_PATH") + "/.bash_history")
	util.CreateFileIfNotExists(os.Getenv("USER_CONFIG_PATH") + "/.gitconfig")
	util.CreateFileIfNotExists(os.Getenv("USER_CONFIG_PATH") + "/.gitignore")

	if !util.FileExists(os.Getenv("USER_CONFIG_PATH") + "/.my.cnf") {
		data := fmt.Sprintf("[client]\nuser=%s\npassword=%s\n", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"))
		util.AppendFile(os.Getenv("USER_CONFIG_PATH")+"/.my.cnf", data)
	}

	if os.Getenv("NGINX_SSL") == "yes" {
		if !util.FileExists("user/nginx/ssl/d4d.pem") || !util.FileExists("user/nginx/ssl/d4d-key.pem") {
			// ./d4d mkcert ssl
		}
	}

	doNginxBuild()
	doPhpBuild()
	doPhpMyAdminBuild()

	// phpMyAdmin configuration
	pmaAuthType := "cookie"
	pmaMySQLUser := ""
	pmaMySQLPassword := ""

	if os.Getenv("PMA_AUTO_LOGIN") == "yes" {
		pmaAuthType = "config"
		pmaMySQLUser = "root"
		pmaMySQLPassword = os.Getenv("MYSQL_ROOT_PASSWORD")

		if os.Getenv("PMA_AUTO_LOGIN") == "no" {
			pmaMySQLUser = os.Getenv("MYSQL_USER")
			pmaMySQLPassword = os.Getenv("MYSQL_PASSWORD")
		}
	}

	util.Copy("config/phpmyadmin/config.user.inc.php.d4d", "config/phpmyadmin/config.user.inc.php")
	util.Sed("__AUTH_TYPE__", pmaAuthType, "config/phpmyadmin/config.user.inc.php")
	util.Sed("__MYSQL_USER__", pmaMySQLUser, "config/phpmyadmin/config.user.inc.php")
	util.Sed("__MYSQL_PASSWORD__", pmaMySQLPassword, "config/phpmyadmin/config.user.inc.php")
}

func doNginxBuild() {
	util.Copy("config/nginx/Dockerfile.build", "config/nginx/Dockerfile")

	util.Sed("__DEBIAN_VERSION__", os.Getenv("DEBIAN_VERSION"), "config/nginx/Dockerfile")

	if os.Getenv("NGINX_SSL") == "yes" {
		os.Mkdir("config/nginx/ssl", 0755)

		util.Copy("user/nginx/ssl/d4d.pem", "config/nginx/ssl/d4d.pem")
		util.Copy("user/nginx/ssl/d4d-key.pem", "config/nginx/ssl/d4d-key.pem")

		util.Sed("__D4D_SSL__", "COPY [\"ssl/d4d.pem\", \"ssl/d4d-key.pem\", \"/etc/nginx/ssl/\"]", "config/nginx/Dockerfile")

	} else {
		util.Sed("__D4D_SSL__", "", "config/nginx/Dockerfile")
	}
}

func doPhpBuild() {
	util.Copy("config/php/Dockerfile.build", "config/php/Dockerfile")

	packageList := []string{"gnupg2", "openssl", "git", "unzip", "libzip-dev", "nano", "libpng-dev", "libmagickwand-dev", "curl", "openssh-client", "less", "inkscape", "cron", "exiftool", "libicu-dev", "libmcrypt-dev", "libc-client-dev", "libkrb5-dev", "libssl-dev", "libxslt1-dev", "bash-completion"}
	peclInstall := []string{}
	phpExtInstall := []string{"pdo", "pdo_mysql", "opcache", "zip", "mysqli", "exif", "bcmath", "calendar", "intl", "soap", "sockets", "xsl"}
	phpExtEnable := []string{"mysqli", "calendar", "exif", "bcmath"}
	npmInstallGlobal := []string{}

	if os.Getenv("SF_CLI") == "yes" {
		packageList = append(packageList, "symfony-cli")
	}

	if os.Getenv("PHP_GD") == "yes" {
		phpExtInstall = append(phpExtInstall, "gd")

		util.Sed("__PHP_GD__", "&& docker-php-ext-configure gd --with-freetype-dir=/usr/include/ --with-jpeg-dir=/usr/include/", "config/php/Dockerfile")
	} else {
		util.Sed("__PHP_GD__", "", "config/php/Dockerfile")
	}

	if os.Getenv("PHP_IMAP") == "yes" {
		phpExtInstall = append(phpExtInstall, "imap")

		util.Sed("__PHP_IMAP__", "&& docker-php-ext-configure imap --with-kerberos --with-imap-ssl", "config/php/Dockerfile")
	} else {
		util.Sed("__PHP_IMAP__", "", "config/php/Dockerfile")
	}

	if os.Getenv("PHP_VERSION") != "8.2" {
		if os.Getenv("PHP_IMAGICK") == "yes" {
			peclInstall = append(peclInstall, "imagick")
			phpExtEnable = append(phpExtEnable, "imagick")
			util.Sed("__IMAGICK__", "", "config/php/Dockerfile")
		}
	}

	if os.Getenv("PHP_VERSION") == "5.6" || os.Getenv("PHP_VERSION") == "8.0" || os.Getenv("PHP_VERSION") == "8.1" || os.Getenv("PHP_VERSION") == "8.2" {
		util.Sed("&& echo 'extension=apc.so' >> \\$PHP_INI_DIR/conf.d/docker-php-ext-apcu.ini \\\\", "", "config/php/Dockerfile")
	}

	if os.Getenv("PHP_VERSION") != "5.6" && os.Getenv("PHP_VERSION") != "8.0" && os.Getenv("PHP_VERSION") != "8.1" && os.Getenv("PHP_VERSION") != "8.2" {
		peclInstall = append(peclInstall, "apcu", "apcu_bc")
		phpExtEnable = append(phpExtEnable, "apcu")
	}

	if os.Getenv("PHP_VERSION") == "5.6" {
		phpExtInstall = append(phpExtInstall, "mcrypt")
		phpExtEnable = append(phpExtEnable, "mcrypt")
	}

	if os.Getenv("PHP_VERSION") == "7.0" || os.Getenv("PHP_VERSION") == "7.1" {
		phpExtInstall = append(phpExtInstall, "mcrypt")
		phpExtEnable = append(phpExtEnable, "mcrypt")
	}

	if os.Getenv("PHP_VERSION") == "7.2" || os.Getenv("PHP_VERSION") == "7.3" || os.Getenv("PHP_VERSION") == "7.4" {
		peclInstall = append(peclInstall, "mcrypt")
		phpExtEnable = append(phpExtEnable, "mcrypt")
	}

	if os.Getenv("PHP_VERSION") == "8.0" || os.Getenv("PHP_VERSION") == "8.1" {
		peclInstall = append(peclInstall, "mcrypt")
		phpExtEnable = append(phpExtEnable, "mcrypt")
	}

	if os.Getenv("PHP_VERSION") == "7.4" {
		util.Sed("docker-php-ext-configure zip --with-libzip && ", "", "config/php/Dockerfile")
		util.Sed("docker-php-ext-configure gd --with-freetype-dir=/usr/include/ --with-jpeg-dir=/usr/include/ && ", "docker-php-ext-configure gd --with-freetype --with-jpeg && ", "config/php/Dockerfile")
	}

	if os.Getenv("PHP_VERSION") == "8.0" || os.Getenv("PHP_VERSION") == "8.1" || os.Getenv("PHP_VERSION") == "8.2" {
		util.Sed("docker-php-ext-configure zip --with-libzip", "docker-php-ext-configure zip", "config/php/Dockerfile")
		util.Sed("docker-php-ext-configure gd --with-freetype-dir=/usr/include/ --with-jpeg-dir=/usr/include/ &&", "docker-php-ext-configure gd --with-freetype --with-jpeg &&", "config/php/Dockerfile")

		if os.Getenv("PHP_IMAGICK") == "yes" {
			util.Sed("__IMAGICK__", "&& cd /tmp && git clone https://github.com/Imagick/imagick && cd imagick && phpize && ./configure && make && make install && echo extension=imagick.so > /usr/local/etc/php/conf.d/docker-php-ext-imagick.ini && rm -rf /tmp/imagick && cd /tmp", "config/php/Dockerfile")
		} else {
			util.Sed("__IMAGICK__", "", "config/php/Dockerfile")
		}
	}

	if os.Getenv("REDIS") == "yes" {
		peclInstall = append(peclInstall, "redis")
		phpExtEnable = append(phpExtEnable, "redis")
	}

	if os.Getenv("RABBITMQ") == "yes" {
		packageList = append(packageList, "librabbitmq-dev", "librabbitmq4")
		if os.Getenv("PHP_VERSION") != "8.0" && os.Getenv("PHP_VERSION") != "8.1" && os.Getenv("PHP_VERSION") != "8.2" {
			peclInstall = append(peclInstall, "amqp")
			util.Sed("__RABBIT_MQ__", "&& echo 'extension=amqp.so' >> $PHP_INI_DIR/conf.d/docker-php-ext-amqp.ini", "config/php/Dockerfile")
		} else {
			util.Sed("__RABBIT_MQ__", "&& cd /tmp && git clone https://github.com/php-amqp/php-amqp && cd php-amqp && phpize && ./configure && make && make install && echo 'extension=amqp.so' > /usr/local/etc/php/conf.d/docker-php-ext-amqp.ini && rm -rf /tmp/php-amqp && cd /tmp", "config/php/Dockerfile")
		}

	} else {
		util.Sed("__RABBIT_MQ__ \\\\", "", "config/php/Dockerfile")
	}

	if os.Getenv("MONGODB") == "yes" {
		peclInstall = append(peclInstall, "mongodb")
		util.Sed("__MONGODB__", "&& echo 'extension=mongodb.so' >> $$PHP_INI_DIR/conf.d/docker-php-ext-mongodb.ini", "config/php/Dockerfile")
	} else {
		util.Sed("__MONGODB__", "", "config/php/Dockerfile")
	}

	if os.Getenv("SUPERVISOR") == "yes" {
		packageList = append(packageList, "supervisor")
	}

	npmInstallGlobal = append(npmInstallGlobal, "npm", "grunt-cli", "yargs", "async", "sass", "gulp", "requirejs", "pm2", "uglify-js", "typescript", "eslint")

	if os.Getenv("XDEBUG") == "yes" {
		if os.Getenv("PHP_VERSION") == "5.6" {
			peclInstall = append(peclInstall, "xdebug-2.5.0")
		}

		if os.Getenv("PHP_VERSION") == "7.0" {
			peclInstall = append(peclInstall, "xdebug-2.7.0")
		}

		if os.Getenv("PHP_VERSION") == "7.1" {
			peclInstall = append(peclInstall, "xdebug-2.9.0")
		}

		if os.Getenv("PHP_VERSION") == "7.2" || os.Getenv("PHP_VERSION") == "7.3" || os.Getenv("PHP_VERSION") == "7.4" {
			peclInstall = append(peclInstall, "xdebug-3.1.0")
		}

		if os.Getenv("PHP_VERSION") == "8.0" || os.Getenv("PHP_VERSION") == "8.1" || os.Getenv("PHP_VERSION") == "8.2" {
			peclInstall = append(peclInstall, "xdebug-3.2.0")
		}

		phpExtEnable = append(phpExtEnable, "xdebug")

		util.Copy("config/php/conf.d/xdebug.d4d", "config/php/conf.d/xdebug.ini")

		util.Sed("__PHP_XDEBUG_CLIENT_PORT__", os.Getenv("XDEBUG_CLIENT_PORT"), "config/php/conf.d/xdebug.ini")
		util.Sed("__PHP_XDEBUG_START_WITH_REQUEST__", os.Getenv("XDEBUG_START_WITH_REQUEST"), "config/php/conf.d/xdebug.ini")

		if os.Getenv("XDEBUG_REMOTE_HOST") != "" {
			util.AppendFile("config/php/conf.d/xdebug.ini", fmt.Sprintf("\nxdebug.remote_host = %s", os.Getenv("XDEBUG_REMOTE_HOST")))
		}
		if os.Getenv("XDEBUG_REMOTE_CONNECT_BACK") != "" {
			util.AppendFile("config/php/conf.d/xdebug.ini", fmt.Sprintf("\nxdebug.remote_connect_back = %s", os.Getenv("XDEBUG_REMOTE_CONNECT_BACK")))
		}
	} else {
		util.Sed("__PHP_XDEBUG__", "", "config/php/Dockerfile")
	}

	if os.Getenv("IONCUBE") == "yes" {
		util.Sed("__PHP_IONCUBE__", "&& curl -fsSL 'https://downloads.ioncube.com/loader_downloads/ioncube_loaders_lin_x86-64.tar.gz' -o ioncube.tar.gz && mkdir -p /tmp/ioncube && tar -xvzf ioncube.tar.gz && mv ioncube/ioncube_loader_lin_${PHP_VERSION}.so `php-config --extension-dir` && rm -Rf ioncube.tar.gz ioncube && docker-php-ext-enable ioncube_loader_lin_${PHP_VERSION}", "config/php/Dockerfile")
	} else {
		util.Sed("__PHP_IONCUBE__ \\\\", "", "config/php/Dockerfile")
	}

	if os.Getenv("YARN") == "yes" {
		util.Sed("__YARN__", "&& apt-get remove -y cmdtest && curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - && echo \"deb https://dl.yarnpkg.com/debian/ stable main\" | tee /etc/apt/sources.list.d/yarn.list && apt-get update && apt-get install -y yarn", "config/php/Dockerfile")
	} else {
		util.Sed("__YARN__", "", "config/php/Dockerfile")
	}

	if os.Getenv("NODEJS") == "yes" {
		util.Sed("__NODEJS__", "&& mkdir -p /var/www/.npm && mkdir -p /var/www/html && printf '{\"name\": \"d4d\", \"version\": \"1.0.0\"}' > /var/www/html/package.json && chown -R $${USER_ID}:$${GROUP_ID} /var/www/.npm && chown -R $${USER_ID}:$${GROUP_ID} /var/www/html && mkdir -p /etc/apt/keyrings && curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg && NODE_MAJOR=$${NODE_JS_VERSION} && echo \"deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$${NODE_MAJOR}.x nodistro main\" | tee /etc/apt/sources.list.d/nodesource.list && apt-get update && apt-get install -y nodejs && npm install --location=global __NPM_INSTALL_GLOBAL__ \\\n    ", "config/php/Dockerfile")
	} else {
		util.Sed("__NODEJS__", "", "config/php/Dockerfile")
	}

	if os.Getenv("WKHTMLTOPDF") == "yes" {
		if os.Getenv("WKHTMLTOPDF_VERSION") == "0.12.3" {
			util.Sed("__WKHTMLTOPDF__", "&& curl -o wkhtmltox-${WKHTMLTOPDF_VERSION}_linux-generic-amd64.tar.xz -sL https://github.com/wkhtmltopdf/wkhtmltopdf/releases/download/${WKHTMLTOPDF_VERSION}/wkhtmltox-${WKHTMLTOPDF_VERSION}_linux-generic-amd64.tar.xz  && echo '9066ab2c7b2035c6eaa043d31aeb7260191e6c88 wkhtmltox-${WKHTMLTOPDF_VERSION}_linux-generic-amd64.tar.xz' | sha1sum -c - && tar -xvf wkhtmltox-${WKHTMLTOPDF_VERSION}_linux-generic-amd64.tar.xz && cp wkhtmltox/lib/* /usr/lib/ && cp wkhtmltox/bin/* /usr/bin/ && cp -r wkhtmltox/share/man/man1 /usr/share/man/ && chmod a+x /usr/bin/wkhtmltopdf && chmod a+x /usr/bin/wkhtmltoimage", "config/php/Dockerfile")
		}
		if os.Getenv("WKHTMLTOPDF_VERSION") == "0.12.4" {
			util.Sed("", "&& curl -o wkhtmltox-${WKHTMLTOPDF_VERSION}_linux-generic-amd64.tar.xz -sL https://github.com/wkhtmltopdf/wkhtmltopdf/releases/download/${WKHTMLTOPDF_VERSION}/wkhtmltox-${WKHTMLTOPDF_VERSION}_linux-generic-amd64.tar.xz  && echo '3f923f425d345940089e44c1466f6408b9619562 wkhtmltox-${WKHTMLTOPDF_VERSION}_linux-generic-amd64.tar.xz' | sha1sum -c - && tar -xvf wkhtmltox-${WKHTMLTOPDF_VERSION}_linux-generic-amd64.tar.xz && cp wkhtmltox/lib/* /usr/lib/ && cp wkhtmltox/bin/* /usr/bin/ && cp -r wkhtmltox/share/man/man1 /usr/share/man/ && chmod a+x /usr/bin/wkhtmltopdf && chmod a+x /usr/bin/wkhtmltoimage", "config/php/Dockerfile")
		}
		if os.Getenv("WKHTMLTOPDF_VERSION") == "0.12.5" {
			util.Sed("", "&& curl -o /tmp/wkhtmltox_${WKHTMLTOPDF_VERSION}.`echo $(lsb_release -cs)`_amd64.deb -sL https://github.com/wkhtmltopdf/wkhtmltopdf/releases/download/${WKHTMLTOPDF_VERSION}/wkhtmltox_${WKHTMLTOPDF_VERSION}-1.`echo $(lsb_release -cs)`_amd64.deb && apt-get --assume-yes install /tmp/wkhtmltox_${WKHTMLTOPDF_VERSION}.`echo $(lsb_release -cs)`_amd64.deb && rm /tmp/wkhtmltox_${WKHTMLTOPDF_VERSION}.`echo $(lsb_release -cs)`_amd64.deb && ln -s /usr/local/bin/wkhtmltopdf /usr/bin/wkhtmltopdf && ln -s /usr/local/bin/wkhtmltoimage /usr/bin/wkhtmltoimage", "config/php/Dockerfile")
		}
		if os.Getenv("WKHTMLTOPDF_VERSION") == "0.12.6" {
			util.Sed("", "&& curl -o /tmp/wkhtmltox_${WKHTMLTOPDF_VERSION}.`echo $(lsb_release -cs)`_amd64.deb -sL https://github.com/wkhtmltopdf/packaging/releases/download/${WKHTMLTOPDF_VERSION}-1/wkhtmltox_${WKHTMLTOPDF_VERSION}-1.`echo $(lsb_release -cs)`_amd64.deb && apt-get --assume-yes install /tmp/wkhtmltox_${WKHTMLTOPDF_VERSION}.`echo $(lsb_release -cs)`_amd64.deb && rm /tmp/wkhtmltox_${WKHTMLTOPDF_VERSION}.`echo $(lsb_release -cs)`_amd64.deb && ln -s /usr/local/bin/wkhtmltopdf /usr/bin/wkhtmltopdf && ln -s /usr/local/bin/wkhtmltoimage /usr/bin/wkhtmltoimage", "config/php/Dockerfile")
		}
	} else {
		util.Sed("__WKHTMLTOPDF__", "", "config/php/Dockerfile")
	}

	if os.Getenv("BLACKFIRE") == "yes" {
		util.Sed("__BLACKFIRE__", "&& curl -sS https://packages.blackfire.io/gpg.key | apt-key add - && echo \"deb http://packages.blackfire.io/debian any main\" | tee /etc/apt/sources.list.d/blackfire.list && apt-get update && apt-get install -y blackfire blackfire-php", "config/php/Dockerfile")
	} else {
		util.Sed("__BLACKFIRE__ \\\\", "", "config/php/Dockerfile")
	}

	util.Sed("__PACKAGE_LIST__", strings.Join(packageList, " "), "config/php/Dockerfile")
	util.Sed("__PHP_EXT_INSTALL__", strings.Join(phpExtInstall, " "), "config/php/Dockerfile")
	util.Sed("__PHP_EXT_ENABLE__", strings.Join(phpExtEnable, " "), "config/php/Dockerfile")
	util.Sed("__PECL_INSTALL__", strings.Join(peclInstall, " "), "config/php/Dockerfile")
	util.Sed("__NPM_INSTALL_GLOBAL__", strings.Join(npmInstallGlobal, " "), "config/php/Dockerfile")
	util.Sed("__CLEANUP__", "&& apt-get clean && rm -rf /var/lib/apt/lists/*", "config/php/Dockerfile")

	if os.Getenv("SF_CLI") == "yes" {
		util.Sed("__SYMFONY_CLI__", "echo \"deb [trusted=yes] https://repo.symfony.com/apt/ /\" | tee /etc/apt/sources.list.d/symfony-cli.list && \\", "config/php/Dockerfile")
	} else {
		util.Sed("__SYMFONY_CLI__", "\\", "config/php/Dockerfile")
	}

	if os.Getenv("PHP_VERSION") == "5.6" || os.Getenv("PHP_VERSION") == "7.0" {
		util.Sed("__CURL_INSECURE__", "-k", "config/php/Dockerfile")
	} else {
		util.Sed(" __CURL_INSECURE__", "", "config/php/Dockerfile")
	}

	util.Sed("    \n", "", "config/php/Dockerfile")
}

func doPhpMyAdminBuild() {
	util.Copy("config/phpmyadmin/Dockerfile.build", "config/phpmyadmin/Dockerfile")

	if runtime.GOOS != "darwin" {
		util.Sed("__PHP_MY_ADMIN__", "phpmyadmin/phpmyadmin", "config/phpmyadmin/Dockerfile")
	} else {
		util.Sed("__PHP_MY_ADMIN__", "arm64v8/phpmyadmin", "config/phpmyadmin/Dockerfile")
	}
}

func remove(haystack []string, needle string) []string {
	for index, value := range haystack {
		if value == needle {
			return append(haystack[:index], haystack[index+1:]...)
		}
	}
	return haystack
}

func doBuildNginxConf() {
	projectConfFile := "config/nginx/project.conf"

	if os.Getenv("NGINX_SSL") == "yes" {
		util.Copy("config/nginx/project-ssl.conf.default", projectConfFile)
	} else {
		util.Copy("config/nginx/project.conf.default", projectConfFile)
	}

	util.Sed("__INCLUDE__", "/etc/nginx/d4d/sf.conf", projectConfFile)
	util.Sed("__PHP_MAX_EXECUTION_TIME__", os.Getenv("PHP_MAX_EXECUTION_TIME"), projectConfFile)
	util.Sed("__NGINX_FASTCGI_BUFFERS__", os.Getenv("NGINX_FASTCGI_BUFFERS"), projectConfFile)
	util.Sed("__NGINX_FASTCGI_BUFFER_SIZE__", os.Getenv("NGINX_FASTCGI_BUFFER_SIZE"), projectConfFile)
	util.Sed("__PHP_UPLOAD_MAX_FILESIZE__", os.Getenv("PHP_UPLOAD_MAX_FILESIZE"), projectConfFile)

	util.Copy("config/nginx/d4d/pwa.conf.default", "config/nginx/d4d/pwa.conf")
	util.Sed("__SYMFONY_FRONT_CONTROLLER__", os.Getenv("SYMFONY_FRONT_CONTROLLER"), "config/nginx/d4d/pwa.conf")

	util.Copy("config/nginx/d4d/sf.conf.default", "config/nginx/d4d/sf.conf")
	util.Sed("__SYMFONY_FRONT_CONTROLLER__", os.Getenv("SYMFONY_FRONT_CONTROLLER"), "config/nginx/d4d/sf.conf")

	util.Copy("config/nginx/d4d/wp.conf.default", "config/nginx/d4d/wp.conf")
	util.Sed("__SYMFONY_FRONT_CONTROLLER__", os.Getenv("SYMFONY_FRONT_CONTROLLER"), "config/nginx/d4d/wp.conf")

	nginxIncludeCache := ""

	if os.Getenv("NGINX_CACHE") == "yes" {
		nginxIncludeCache = "include /etc/nginx/d4d/cache.conf;"
	}

	util.Sed("__INCLUDE_CACHE__", nginxIncludeCache, projectConfFile)
}

func doBuildMySQLConf() {
	util.Copy("config/mysql/d4d.cnf.d4d", "config/mysql/d4d.cnf")
	util.Sed("__MYSQL_MAX_ALLOWED_PACKET__", os.Getenv("MYSQL_MAX_ALLOWED_PACKET"), "config/mysql/d4d.cnf")
	util.Sed("__MYSQL_INNODB_LOG_FILE_SIZE__", os.Getenv("MYSQL_INNODB_LOG_FILE_SIZE"), "config/mysql/d4d.cnf")
	util.Sed("__MYSQL_WAIT_TIMEOUT__", os.Getenv("MYSQL_WAIT_TIMEOUT"), "config/mysql/d4d.cnf")
	util.Sed("__MYSQL_CHARACTER_SET_SERVER__", os.Getenv("MYSQL_CHARACTER_SET_SERVER"), "config/mysql/d4d.cnf")
	util.Sed("__MYSQL_COLLATION_SERVER__", os.Getenv("MYSQL_COLLATION_SERVER"), "config/mysql/d4d.cnf")
}

func doBuild() {
	util.Copy("docker/compose.yml", "docker-compose.yml")
	util.AppendFile("docker-compose.yml", util.FileGetContents("docker/php.yml"))

	if os.Getenv("DOCKER_ENV_PHP") == "yes" {
		util.AppendFile("docker-compose.yml", util.FileGetContents("docker/env/php.yml"))
	}

	if os.Getenv("MYSQL") == "yes" {
		if os.Getenv("MYSQL_INST") == "mysql" {
			util.AppendFile("docker-compose.yml", util.FileGetContents("docker/mysql.yml"))
		} else {
			if os.Getenv("MYSQL_INST") == "mariadb" {
				util.AppendFile("docker-compose.yml", util.FileGetContents("docker/mariadb.yml"))
			}
		}
		util.Sed("#php_depends_on", "depends_on:\r\n      - mysql", "docker-compose.yml")
	} else {
		util.Sed("#php_depends_on", "", "docker-compose.yml")
	}

	if os.Getenv("MAILHOG") == "yes" {
		util.AppendFile("docker-compose.yml", util.FileGetContents("docker/mailhog.yml"))
	}

	if os.Getenv("MAILPIT") == "yes" {
		util.AppendFile("docker-compose.yml", util.FileGetContents("docker/mailpit.yml"))
	}

	if os.Getenv("PMA") == "yes" && os.Getenv("MYSQL") == "yes" {
		util.AppendFile("docker-compose.yml", util.FileGetContents("docker/phpmyadmin.yml"))
	}

	if os.Getenv("REDIS") == "yes" {
		util.AppendFile("docker-compose.yml", util.FileGetContents("docker/redis.yml"))
	}

	if os.Getenv("RABBITMQ") == "yes" {
		util.AppendFile("docker-compose.yml", util.FileGetContents("docker/rabbitmq.yml"))
	}

	if os.Getenv("ELASTICSEARCH") == "yes" {
		util.AppendFile("docker-compose.yml", util.FileGetContents("docker/elasticsearch.yml"))
	}

	if os.Getenv("NGROK") == "yes" {
		util.AppendFile("docker-compose.yml", util.FileGetContents("docker/ngrok.yml"))
	}

	if os.Getenv("MONGODB") == "yes" {
		util.AppendFile("docker-compose.yml", util.FileGetContents("docker/mongodb.yml"))
	}

	if os.Getenv("ELK") == "yes" {
		util.AppendFile("docker-compose.yml", util.FileGetContents("docker/elk.yml"))
	}

	if os.Getenv("EXTERNAL_NETWORK") == "no" || os.Getenv("EXTERNAL_NETWORK") == "yes" {
		util.Sed("__NGINX_NETWORKS__", fmt.Sprintf("networks:\n      default:\n        aliases:\n          - %s", os.Getenv("PROJECT_DOMAIN_1")), "docker-compose.yml")
	}

	if os.Getenv("EXTERNAL_NETWORK") == "yes" {
		util.AppendFile("docker-compose.yml", util.FileGetContents("docker/network.yml"))
	}
}

func doBeforeStart() {
	envFile := util.GetCurrentDir() + "/.env"
	util.LoadEnvFile(envFile)

	if os.Getenv("CLEAN_NGINX_LOGS") == "yes" {
		if err := os.Truncate("var/logs/nginx/project_access.log", 0); err != nil {
			log.Printf("Failed to truncate: %v", err)
		}
	}

	if os.Getenv("CLEAN_SF_logs") == "yes" {
		os.RemoveAll("project/var/cache")
		os.RemoveAll("project/var/log")
	}
}
