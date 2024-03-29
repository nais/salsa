package php

import (
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/build/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComposerLockJsonParsing(t *testing.T) {
	got, err := ComposerDeps(composerLockContents)
	assert.NoError(t, err)
	want := map[string]build.Dependency{}
	want["guzzlehttp/guzzle"] = test.Dependency("guzzlehttp/guzzle", "7.4.1", "sha1", "")
	want["guzzlehttp/promises"] = test.Dependency("guzzlehttp/promises", "1.5.1", "sha1", "")
	want["guzzlehttp/psr7"] = test.Dependency("guzzlehttp/psr7", "2.1.0", "sha1", "")
	want["nikic/fast-route"] = test.Dependency("nikic/fast-route", "v1.3.0", "sha1", "")
	want["psr/container"] = test.Dependency("psr/container", "2.0.2", "sha1", "")
	want["psr/http-client"] = test.Dependency("psr/http-client", "1.0.1", "sha1", "")
	want["psr/http-factory"] = test.Dependency("psr/http-factory", "1.0.1", "sha1", "")
	want["psr/http-message"] = test.Dependency("psr/http-message", "1.0.1", "sha1", "")
	want["psr/http-server-handler"] = test.Dependency("psr/http-server-handler", "1.0.1", "sha1", "")
	want["psr/http-server-middleware"] = test.Dependency("psr/http-server-middleware", "1.0.1", "sha1", "")
	want["psr/log"] = test.Dependency("psr/log", "1.1.4", "sha1", "")
	want["ralouphie/getallheaders"] = test.Dependency("ralouphie/getallheaders", "3.0.3", "sha1", "")
	want["slim/slim"] = test.Dependency("slim/slim", "4.9.0", "sha1", "")
	want["symfony/deprecation-contracts"] = test.Dependency("symfony/deprecation-contracts", "v2.5.0", "sha1", "")

	test.AssertEqual(t, got, want)
}

const composerLockContents = `{
    "_readme": [
        "This file locks the dependencies of your project to a known state",
        "Read more about it at https://getcomposer.org/doc/01-basic-usage.md#installing-dependencies",
        "This file is @generated automatically"
    ],
    "content-hash": "f74ba492654fa5b3c63c100054f2e588",
    "packages": [
        {
            "name": "guzzlehttp/guzzle",
            "version": "7.4.1",
            "source": {
                "type": "git",
                "url": "https://github.com/guzzle/guzzle.git",
                "reference": "ee0a041b1760e6a53d2a39c8c34115adc2af2c79"
            },
            "dist": {
                "type": "zip",
                "url": "https://api.github.com/repos/guzzle/guzzle/zipball/ee0a041b1760e6a53d2a39c8c34115adc2af2c79",
                "reference": "ee0a041b1760e6a53d2a39c8c34115adc2af2c79",
                "shasum": ""
            },
            "require": {
                "ext-json": "*",
                "guzzlehttp/promises": "^1.5",
                "guzzlehttp/psr7": "^1.8.3 || ^2.1",
                "php": "^7.2.5 || ^8.0",
                "psr/http-client": "^1.0",
                "symfony/deprecation-contracts": "^2.2 || ^3.0"
            },
            "provide": {
                "psr/http-client-implementation": "1.0"
            },
            "require-dev": {
                "bamarni/composer-bin-plugin": "^1.4.1",
                "ext-curl": "*",
                "php-http/client-integration-tests": "^3.0",
                "phpunit/phpunit": "^8.5.5 || ^9.3.5",
                "psr/log": "^1.1 || ^2.0 || ^3.0"
            },
            "suggest": {
                "ext-curl": "Required for CURL handler support",
                "ext-intl": "Required for Internationalized Domain Name (IDN) support",
                "psr/log": "Required for using the Log middleware"
            },
            "type": "library",
            "extra": {
                "branch-alias": {
                    "dev-master": "7.4-dev"
                }
            },
            "autoload": {
                "psr-4": {
                    "GuzzleHttp\\": "src/"
                },
                "files": [
                    "src/functions_include.php"
                ]
            },
            "notification-url": "https://packagist.org/downloads/",
            "license": [
                "MIT"
            ],
            "authors": [
                {
                    "name": "Graham Campbell",
                    "email": "hello@gjcampbell.co.uk",
                    "homepage": "https://github.com/GrahamCampbell"
                },
                {
                    "name": "Michael Dowling",
                    "email": "mtdowling@gmail.com",
                    "homepage": "https://github.com/mtdowling"
                },
                {
                    "name": "Jeremy Lindblom",
                    "email": "jeremeamia@gmail.com",
                    "homepage": "https://github.com/jeremeamia"
                },
                {
                    "name": "George Mponos",
                    "email": "gmponos@gmail.com",
                    "homepage": "https://github.com/gmponos"
                },
                {
                    "name": "Tobias Nyholm",
                    "email": "tobias.nyholm@gmail.com",
                    "homepage": "https://github.com/Nyholm"
                },
                {
                    "name": "Márk Sági-Kazár",
                    "email": "mark.sagikazar@gmail.com",
                    "homepage": "https://github.com/sagikazarmark"
                },
                {
                    "name": "Tobias Schultze",
                    "email": "webmaster@tubo-world.de",
                    "homepage": "https://github.com/Tobion"
                }
            ],
            "description": "Guzzle is a PHP HTTP client library",
            "keywords": [
                "client",
                "curl",
                "framework",
                "http",
                "http client",
                "psr-18",
                "psr-7",
                "rest",
                "web service"
            ],
            "support": {
                "issues": "https://github.com/guzzle/guzzle/issues",
                "source": "https://github.com/guzzle/guzzle/tree/7.4.1"
            },
            "funding": [
                {
                    "url": "https://github.com/GrahamCampbell",
                    "type": "github"
                },
                {
                    "url": "https://github.com/Nyholm",
                    "type": "github"
                },
                {
                    "url": "https://tidelift.com/funding/github/packagist/guzzlehttp/guzzle",
                    "type": "tidelift"
                }
            ],
            "time": "2021-12-06T18:43:05+00:00"
        },
        {
            "name": "guzzlehttp/promises",
            "version": "1.5.1",
            "source": {
                "type": "git",
                "url": "https://github.com/guzzle/promises.git",
                "reference": "fe752aedc9fd8fcca3fe7ad05d419d32998a06da"
            },
            "dist": {
                "type": "zip",
                "url": "https://api.github.com/repos/guzzle/promises/zipball/fe752aedc9fd8fcca3fe7ad05d419d32998a06da",
                "reference": "fe752aedc9fd8fcca3fe7ad05d419d32998a06da",
                "shasum": ""
            },
            "require": {
                "php": ">=5.5"
            },
            "require-dev": {
                "symfony/phpunit-bridge": "^4.4 || ^5.1"
            },
            "type": "library",
            "extra": {
                "branch-alias": {
                    "dev-master": "1.5-dev"
                }
            },
            "autoload": {
                "psr-4": {
                    "GuzzleHttp\\Promise\\": "src/"
                },
                "files": [
                    "src/functions_include.php"
                ]
            },
            "notification-url": "https://packagist.org/downloads/",
            "license": [
                "MIT"
            ],
            "authors": [
                {
                    "name": "Graham Campbell",
                    "email": "hello@gjcampbell.co.uk",
                    "homepage": "https://github.com/GrahamCampbell"
                },
                {
                    "name": "Michael Dowling",
                    "email": "mtdowling@gmail.com",
                    "homepage": "https://github.com/mtdowling"
                },
                {
                    "name": "Tobias Nyholm",
                    "email": "tobias.nyholm@gmail.com",
                    "homepage": "https://github.com/Nyholm"
                },
                {
                    "name": "Tobias Schultze",
                    "email": "webmaster@tubo-world.de",
                    "homepage": "https://github.com/Tobion"
                }
            ],
            "description": "Guzzle promises library",
            "keywords": [
                "promise"
            ],
            "support": {
                "issues": "https://github.com/guzzle/promises/issues",
                "source": "https://github.com/guzzle/promises/tree/1.5.1"
            },
            "funding": [
                {
                    "url": "https://github.com/GrahamCampbell",
                    "type": "github"
                },
                {
                    "url": "https://github.com/Nyholm",
                    "type": "github"
                },
                {
                    "url": "https://tidelift.com/funding/github/packagist/guzzlehttp/promises",
                    "type": "tidelift"
                }
            ],
            "time": "2021-10-22T20:56:57+00:00"
        },
        {
            "name": "guzzlehttp/psr7",
            "version": "2.1.0",
            "source": {
                "type": "git",
                "url": "https://github.com/guzzle/psr7.git",
                "reference": "089edd38f5b8abba6cb01567c2a8aaa47cec4c72"
            },
            "dist": {
                "type": "zip",
                "url": "https://api.github.com/repos/guzzle/psr7/zipball/089edd38f5b8abba6cb01567c2a8aaa47cec4c72",
                "reference": "089edd38f5b8abba6cb01567c2a8aaa47cec4c72",
                "shasum": ""
            },
            "require": {
                "php": "^7.2.5 || ^8.0",
                "psr/http-factory": "^1.0",
                "psr/http-message": "^1.0",
                "ralouphie/getallheaders": "^3.0"
            },
            "provide": {
                "psr/http-factory-implementation": "1.0",
                "psr/http-message-implementation": "1.0"
            },
            "require-dev": {
                "bamarni/composer-bin-plugin": "^1.4.1",
                "http-interop/http-factory-tests": "^0.9",
                "phpunit/phpunit": "^8.5.8 || ^9.3.10"
            },
            "suggest": {
                "laminas/laminas-httphandlerrunner": "Emit PSR-7 responses"
            },
            "type": "library",
            "extra": {
                "branch-alias": {
                    "dev-master": "2.1-dev"
                }
            },
            "autoload": {
                "psr-4": {
                    "GuzzleHttp\\Psr7\\": "src/"
                }
            },
            "notification-url": "https://packagist.org/downloads/",
            "license": [
                "MIT"
            ],
            "authors": [
                {
                    "name": "Graham Campbell",
                    "email": "hello@gjcampbell.co.uk",
                    "homepage": "https://github.com/GrahamCampbell"
                },
                {
                    "name": "Michael Dowling",
                    "email": "mtdowling@gmail.com",
                    "homepage": "https://github.com/mtdowling"
                },
                {
                    "name": "George Mponos",
                    "email": "gmponos@gmail.com",
                    "homepage": "https://github.com/gmponos"
                },
                {
                    "name": "Tobias Nyholm",
                    "email": "tobias.nyholm@gmail.com",
                    "homepage": "https://github.com/Nyholm"
                },
                {
                    "name": "Márk Sági-Kazár",
                    "email": "mark.sagikazar@gmail.com",
                    "homepage": "https://github.com/sagikazarmark"
                },
                {
                    "name": "Tobias Schultze",
                    "email": "webmaster@tubo-world.de",
                    "homepage": "https://github.com/Tobion"
                },
                {
                    "name": "Márk Sági-Kazár",
                    "email": "mark.sagikazar@gmail.com",
                    "homepage": "https://sagikazarmark.hu"
                }
            ],
            "description": "PSR-7 message implementation that also provides common utility methods",
            "keywords": [
                "http",
                "message",
                "psr-7",
                "request",
                "response",
                "stream",
                "uri",
                "url"
            ],
            "support": {
                "issues": "https://github.com/guzzle/psr7/issues",
                "source": "https://github.com/guzzle/psr7/tree/2.1.0"
            },
            "funding": [
                {
                    "url": "https://github.com/GrahamCampbell",
                    "type": "github"
                },
                {
                    "url": "https://github.com/Nyholm",
                    "type": "github"
                },
                {
                    "url": "https://tidelift.com/funding/github/packagist/guzzlehttp/psr7",
                    "type": "tidelift"
                }
            ],
            "time": "2021-10-06T17:43:30+00:00"
        },
        {
            "name": "nikic/fast-route",
            "version": "v1.3.0",
            "source": {
                "type": "git",
                "url": "https://github.com/nikic/FastRoute.git",
                "reference": "181d480e08d9476e61381e04a71b34dc0432e812"
            },
            "dist": {
                "type": "zip",
                "url": "https://api.github.com/repos/nikic/FastRoute/zipball/181d480e08d9476e61381e04a71b34dc0432e812",
                "reference": "181d480e08d9476e61381e04a71b34dc0432e812",
                "shasum": ""
            },
            "require": {
                "php": ">=5.4.0"
            },
            "require-dev": {
                "phpunit/phpunit": "^4.8.35|~5.7"
            },
            "type": "library",
            "autoload": {
                "psr-4": {
                    "FastRoute\\": "src/"
                },
                "files": [
                    "src/functions.php"
                ]
            },
            "notification-url": "https://packagist.org/downloads/",
            "license": [
                "BSD-3-Clause"
            ],
            "authors": [
                {
                    "name": "Nikita Popov",
                    "email": "nikic@php.net"
                }
            ],
            "description": "Fast request router for PHP",
            "keywords": [
                "router",
                "routing"
            ],
            "support": {
                "issues": "https://github.com/nikic/FastRoute/issues",
                "source": "https://github.com/nikic/FastRoute/tree/master"
            },
            "time": "2018-02-13T20:26:39+00:00"
        },
        {
            "name": "psr/container",
            "version": "2.0.2",
            "source": {
                "type": "git",
                "url": "https://github.com/php-fig/container.git",
                "reference": "c71ecc56dfe541dbd90c5360474fbc405f8d5963"
            },
            "dist": {
                "type": "zip",
                "url": "https://api.github.com/repos/php-fig/container/zipball/c71ecc56dfe541dbd90c5360474fbc405f8d5963",
                "reference": "c71ecc56dfe541dbd90c5360474fbc405f8d5963",
                "shasum": ""
            },
            "require": {
                "php": ">=7.4.0"
            },
            "type": "library",
            "extra": {
                "branch-alias": {
                    "dev-master": "2.0.x-dev"
                }
            },
            "autoload": {
                "psr-4": {
                    "Psr\\Container\\": "src/"
                }
            },
            "notification-url": "https://packagist.org/downloads/",
            "license": [
                "MIT"
            ],
            "authors": [
                {
                    "name": "PHP-FIG",
                    "homepage": "https://www.php-fig.org/"
                }
            ],
            "description": "Common Container Interface (PHP FIG PSR-11)",
            "homepage": "https://github.com/php-fig/container",
            "keywords": [
                "PSR-11",
                "container",
                "container-interface",
                "container-interop",
                "psr"
            ],
            "support": {
                "issues": "https://github.com/php-fig/container/issues",
                "source": "https://github.com/php-fig/container/tree/2.0.2"
            },
            "time": "2021-11-05T16:47:00+00:00"
        },
        {
            "name": "psr/http-client",
            "version": "1.0.1",
            "source": {
                "type": "git",
                "url": "https://github.com/php-fig/http-client.git",
                "reference": "2dfb5f6c5eff0e91e20e913f8c5452ed95b86621"
            },
            "dist": {
                "type": "zip",
                "url": "https://api.github.com/repos/php-fig/http-client/zipball/2dfb5f6c5eff0e91e20e913f8c5452ed95b86621",
                "reference": "2dfb5f6c5eff0e91e20e913f8c5452ed95b86621",
                "shasum": ""
            },
            "require": {
                "php": "^7.0 || ^8.0",
                "psr/http-message": "^1.0"
            },
            "type": "library",
            "extra": {
                "branch-alias": {
                    "dev-master": "1.0.x-dev"
                }
            },
            "autoload": {
                "psr-4": {
                    "Psr\\Http\\Client\\": "src/"
                }
            },
            "notification-url": "https://packagist.org/downloads/",
            "license": [
                "MIT"
            ],
            "authors": [
                {
                    "name": "PHP-FIG",
                    "homepage": "http://www.php-fig.org/"
                }
            ],
            "description": "Common interface for HTTP clients",
            "homepage": "https://github.com/php-fig/http-client",
            "keywords": [
                "http",
                "http-client",
                "psr",
                "psr-18"
            ],
            "support": {
                "source": "https://github.com/php-fig/http-client/tree/master"
            },
            "time": "2020-06-29T06:28:15+00:00"
        },
        {
            "name": "psr/http-factory",
            "version": "1.0.1",
            "source": {
                "type": "git",
                "url": "https://github.com/php-fig/http-factory.git",
                "reference": "12ac7fcd07e5b077433f5f2bee95b3a771bf61be"
            },
            "dist": {
                "type": "zip",
                "url": "https://api.github.com/repos/php-fig/http-factory/zipball/12ac7fcd07e5b077433f5f2bee95b3a771bf61be",
                "reference": "12ac7fcd07e5b077433f5f2bee95b3a771bf61be",
                "shasum": ""
            },
            "require": {
                "php": ">=7.0.0",
                "psr/http-message": "^1.0"
            },
            "type": "library",
            "extra": {
                "branch-alias": {
                    "dev-master": "1.0.x-dev"
                }
            },
            "autoload": {
                "psr-4": {
                    "Psr\\Http\\Message\\": "src/"
                }
            },
            "notification-url": "https://packagist.org/downloads/",
            "license": [
                "MIT"
            ],
            "authors": [
                {
                    "name": "PHP-FIG",
                    "homepage": "http://www.php-fig.org/"
                }
            ],
            "description": "Common interfaces for PSR-7 HTTP message factories",
            "keywords": [
                "factory",
                "http",
                "message",
                "psr",
                "psr-17",
                "psr-7",
                "request",
                "response"
            ],
            "support": {
                "source": "https://github.com/php-fig/http-factory/tree/master"
            },
            "time": "2019-04-30T12:38:16+00:00"
        },
        {
            "name": "psr/http-message",
            "version": "1.0.1",
            "source": {
                "type": "git",
                "url": "https://github.com/php-fig/http-message.git",
                "reference": "f6561bf28d520154e4b0ec72be95418abe6d9363"
            },
            "dist": {
                "type": "zip",
                "url": "https://api.github.com/repos/php-fig/http-message/zipball/f6561bf28d520154e4b0ec72be95418abe6d9363",
                "reference": "f6561bf28d520154e4b0ec72be95418abe6d9363",
                "shasum": ""
            },
            "require": {
                "php": ">=5.3.0"
            },
            "type": "library",
            "extra": {
                "branch-alias": {
                    "dev-master": "1.0.x-dev"
                }
            },
            "autoload": {
                "psr-4": {
                    "Psr\\Http\\Message\\": "src/"
                }
            },
            "notification-url": "https://packagist.org/downloads/",
            "license": [
                "MIT"
            ],
            "authors": [
                {
                    "name": "PHP-FIG",
                    "homepage": "http://www.php-fig.org/"
                }
            ],
            "description": "Common interface for HTTP messages",
            "homepage": "https://github.com/php-fig/http-message",
            "keywords": [
                "http",
                "http-message",
                "psr",
                "psr-7",
                "request",
                "response"
            ],
            "support": {
                "source": "https://github.com/php-fig/http-message/tree/master"
            },
            "time": "2016-08-06T14:39:51+00:00"
        },
        {
            "name": "psr/http-server-handler",
            "version": "1.0.1",
            "source": {
                "type": "git",
                "url": "https://github.com/php-fig/http-server-handler.git",
                "reference": "aff2f80e33b7f026ec96bb42f63242dc50ffcae7"
            },
            "dist": {
                "type": "zip",
                "url": "https://api.github.com/repos/php-fig/http-server-handler/zipball/aff2f80e33b7f026ec96bb42f63242dc50ffcae7",
                "reference": "aff2f80e33b7f026ec96bb42f63242dc50ffcae7",
                "shasum": ""
            },
            "require": {
                "php": ">=7.0",
                "psr/http-message": "^1.0"
            },
            "type": "library",
            "extra": {
                "branch-alias": {
                    "dev-master": "1.0.x-dev"
                }
            },
            "autoload": {
                "psr-4": {
                    "Psr\\Http\\Server\\": "src/"
                }
            },
            "notification-url": "https://packagist.org/downloads/",
            "license": [
                "MIT"
            ],
            "authors": [
                {
                    "name": "PHP-FIG",
                    "homepage": "http://www.php-fig.org/"
                }
            ],
            "description": "Common interface for HTTP server-side request handler",
            "keywords": [
                "handler",
                "http",
                "http-interop",
                "psr",
                "psr-15",
                "psr-7",
                "request",
                "response",
                "server"
            ],
            "support": {
                "issues": "https://github.com/php-fig/http-server-handler/issues",
                "source": "https://github.com/php-fig/http-server-handler/tree/master"
            },
            "time": "2018-10-30T16:46:14+00:00"
        },
        {
            "name": "psr/http-server-middleware",
            "version": "1.0.1",
            "source": {
                "type": "git",
                "url": "https://github.com/php-fig/http-server-middleware.git",
                "reference": "2296f45510945530b9dceb8bcedb5cb84d40c5f5"
            },
            "dist": {
                "type": "zip",
                "url": "https://api.github.com/repos/php-fig/http-server-middleware/zipball/2296f45510945530b9dceb8bcedb5cb84d40c5f5",
                "reference": "2296f45510945530b9dceb8bcedb5cb84d40c5f5",
                "shasum": ""
            },
            "require": {
                "php": ">=7.0",
                "psr/http-message": "^1.0",
                "psr/http-server-handler": "^1.0"
            },
            "type": "library",
            "extra": {
                "branch-alias": {
                    "dev-master": "1.0.x-dev"
                }
            },
            "autoload": {
                "psr-4": {
                    "Psr\\Http\\Server\\": "src/"
                }
            },
            "notification-url": "https://packagist.org/downloads/",
            "license": [
                "MIT"
            ],
            "authors": [
                {
                    "name": "PHP-FIG",
                    "homepage": "http://www.php-fig.org/"
                }
            ],
            "description": "Common interface for HTTP server-side middleware",
            "keywords": [
                "http",
                "http-interop",
                "middleware",
                "psr",
                "psr-15",
                "psr-7",
                "request",
                "response"
            ],
            "support": {
                "issues": "https://github.com/php-fig/http-server-middleware/issues",
                "source": "https://github.com/php-fig/http-server-middleware/tree/master"
            },
            "time": "2018-10-30T17:12:04+00:00"
        },
        {
            "name": "psr/log",
            "version": "1.1.4",
            "source": {
                "type": "git",
                "url": "https://github.com/php-fig/log.git",
                "reference": "d49695b909c3b7628b6289db5479a1c204601f11"
            },
            "dist": {
                "type": "zip",
                "url": "https://api.github.com/repos/php-fig/log/zipball/d49695b909c3b7628b6289db5479a1c204601f11",
                "reference": "d49695b909c3b7628b6289db5479a1c204601f11",
                "shasum": ""
            },
            "require": {
                "php": ">=5.3.0"
            },
            "type": "library",
            "extra": {
                "branch-alias": {
                    "dev-master": "1.1.x-dev"
                }
            },
            "autoload": {
                "psr-4": {
                    "Psr\\Log\\": "Psr/Log/"
                }
            },
            "notification-url": "https://packagist.org/downloads/",
            "license": [
                "MIT"
            ],
            "authors": [
                {
                    "name": "PHP-FIG",
                    "homepage": "https://www.php-fig.org/"
                }
            ],
            "description": "Common interface for logging libraries",
            "homepage": "https://github.com/php-fig/log",
            "keywords": [
                "log",
                "psr",
                "psr-3"
            ],
            "support": {
                "source": "https://github.com/php-fig/log/tree/1.1.4"
            },
            "time": "2021-05-03T11:20:27+00:00"
        },
        {
            "name": "ralouphie/getallheaders",
            "version": "3.0.3",
            "source": {
                "type": "git",
                "url": "https://github.com/ralouphie/getallheaders.git",
                "reference": "120b605dfeb996808c31b6477290a714d356e822"
            },
            "dist": {
                "type": "zip",
                "url": "https://api.github.com/repos/ralouphie/getallheaders/zipball/120b605dfeb996808c31b6477290a714d356e822",
                "reference": "120b605dfeb996808c31b6477290a714d356e822",
                "shasum": ""
            },
            "require": {
                "php": ">=5.6"
            },
            "require-dev": {
                "php-coveralls/php-coveralls": "^2.1",
                "phpunit/phpunit": "^5 || ^6.5"
            },
            "type": "library",
            "autoload": {
                "files": [
                    "src/getallheaders.php"
                ]
            },
            "notification-url": "https://packagist.org/downloads/",
            "license": [
                "MIT"
            ],
            "authors": [
                {
                    "name": "Ralph Khattar",
                    "email": "ralph.khattar@gmail.com"
                }
            ],
            "description": "A polyfill for getallheaders.",
            "support": {
                "issues": "https://github.com/ralouphie/getallheaders/issues",
                "source": "https://github.com/ralouphie/getallheaders/tree/develop"
            },
            "time": "2019-03-08T08:55:37+00:00"
        },
        {
            "name": "slim/slim",
            "version": "4.9.0",
            "source": {
                "type": "git",
                "url": "https://github.com/slimphp/Slim.git",
                "reference": "44d3c9c0bfcc47e52e42b097b6062689d21b904b"
            },
            "dist": {
                "type": "zip",
                "url": "https://api.github.com/repos/slimphp/Slim/zipball/44d3c9c0bfcc47e52e42b097b6062689d21b904b",
                "reference": "44d3c9c0bfcc47e52e42b097b6062689d21b904b",
                "shasum": ""
            },
            "require": {
                "ext-json": "*",
                "nikic/fast-route": "^1.3",
                "php": "^7.3 || ^8.0",
                "psr/container": "^1.0 || ^2.0",
                "psr/http-factory": "^1.0",
                "psr/http-message": "^1.0",
                "psr/http-server-handler": "^1.0",
                "psr/http-server-middleware": "^1.0",
                "psr/log": "^1.1 || ^2.0 || ^3.0"
            },
            "require-dev": {
                "adriansuter/php-autoload-override": "^1.2",
                "ext-simplexml": "*",
                "guzzlehttp/psr7": "^2.0",
                "laminas/laminas-diactoros": "^2.8",
                "nyholm/psr7": "^1.4",
                "nyholm/psr7-server": "^1.0",
                "phpspec/prophecy": "^1.14",
                "phpspec/prophecy-phpunit": "^2.0",
                "phpstan/phpstan": "^0.12.99",
                "phpunit/phpunit": "^9.5",
                "slim/http": "^1.2",
                "slim/psr7": "^1.5",
                "squizlabs/php_codesniffer": "^3.6"
            },
            "suggest": {
                "ext-simplexml": "Needed to support XML format in BodyParsingMiddleware",
                "ext-xml": "Needed to support XML format in BodyParsingMiddleware",
                "php-di/php-di": "PHP-DI is the recommended container library to be used with Slim",
                "slim/psr7": "Slim PSR-7 implementation. See https://www.slimframework.com/docs/v4/start/installation.html for more information."
            },
            "type": "library",
            "autoload": {
                "psr-4": {
                    "Slim\\": "Slim"
                }
            },
            "notification-url": "https://packagist.org/downloads/",
            "license": [
                "MIT"
            ],
            "authors": [
                {
                    "name": "Josh Lockhart",
                    "email": "hello@joshlockhart.com",
                    "homepage": "https://joshlockhart.com"
                },
                {
                    "name": "Andrew Smith",
                    "email": "a.smith@silentworks.co.uk",
                    "homepage": "http://silentworks.co.uk"
                },
                {
                    "name": "Rob Allen",
                    "email": "rob@akrabat.com",
                    "homepage": "http://akrabat.com"
                },
                {
                    "name": "Pierre Berube",
                    "email": "pierre@lgse.com",
                    "homepage": "http://www.lgse.com"
                },
                {
                    "name": "Gabriel Manricks",
                    "email": "gmanricks@me.com",
                    "homepage": "http://gabrielmanricks.com"
                }
            ],
            "description": "Slim is a PHP micro framework that helps you quickly write simple yet powerful web applications and APIs",
            "homepage": "https://www.slimframework.com",
            "keywords": [
                "api",
                "framework",
                "micro",
                "router"
            ],
            "support": {
                "docs": "https://www.slimframework.com/docs/v4/",
                "forum": "https://discourse.slimframework.com/",
                "irc": "irc://irc.freenode.net:6667/slimphp",
                "issues": "https://github.com/slimphp/Slim/issues",
                "rss": "https://www.slimframework.com/blog/feed.rss",
                "slack": "https://slimphp.slack.com/",
                "source": "https://github.com/slimphp/Slim",
                "wiki": "https://github.com/slimphp/Slim/wiki"
            },
            "funding": [
                {
                    "url": "https://opencollective.com/slimphp",
                    "type": "open_collective"
                },
                {
                    "url": "https://tidelift.com/funding/github/packagist/slim/slim",
                    "type": "tidelift"
                }
            ],
            "time": "2021-10-05T03:00:00+00:00"
        },
        {
            "name": "symfony/deprecation-contracts",
            "version": "v2.5.0",
            "source": {
                "type": "git",
                "url": "https://github.com/symfony/deprecation-contracts.git",
                "reference": "6f981ee24cf69ee7ce9736146d1c57c2780598a8"
            },
            "dist": {
                "type": "zip",
                "url": "https://api.github.com/repos/symfony/deprecation-contracts/zipball/6f981ee24cf69ee7ce9736146d1c57c2780598a8",
                "reference": "6f981ee24cf69ee7ce9736146d1c57c2780598a8",
                "shasum": ""
            },
            "require": {
                "php": ">=7.1"
            },
            "type": "library",
            "extra": {
                "branch-alias": {
                    "dev-main": "2.5-dev"
                },
                "thanks": {
                    "name": "symfony/contracts",
                    "url": "https://github.com/symfony/contracts"
                }
            },
            "autoload": {
                "files": [
                    "function.php"
                ]
            },
            "notification-url": "https://packagist.org/downloads/",
            "license": [
                "MIT"
            ],
            "authors": [
                {
                    "name": "Nicolas Grekas",
                    "email": "p@tchwork.com"
                },
                {
                    "name": "Symfony Community",
                    "homepage": "https://symfony.com/contributors"
                }
            ],
            "description": "A generic function and convention to trigger deprecation notices",
            "homepage": "https://symfony.com",
            "support": {
                "source": "https://github.com/symfony/deprecation-contracts/tree/v2.5.0"
            },
            "funding": [
                {
                    "url": "https://symfony.com/sponsor",
                    "type": "custom"
                },
                {
                    "url": "https://github.com/fabpot",
                    "type": "github"
                },
                {
                    "url": "https://tidelift.com/funding/github/packagist/symfony/symfony",
                    "type": "tidelift"
                }
            ],
            "time": "2021-07-12T14:48:14+00:00"
        }
    ],
    "packages-dev": [],
    "aliases": [],
    "minimum-stability": "stable",
    "stability-flags": [],
    "prefer-stable": false,
    "prefer-lowest": false,
    "platform": [],
    "platform-dev": [],
    "plugin-api-version": "2.2.0"
}`

func TestBuildPhp(t *testing.T) {
	tests := []test.IntegrationTest{
		{
			Name:      "find build file and parse output",
			BuildType: BuildComposer(),
			WorkDir:   "testdata",
			BuildPath: "testdata/composer.lock",
			Cmd:       "composer.lock",
			Want: test.Want{
				Key:     "guzzlehttp/guzzle",
				Version: "7.4.1",
				Algo:    "sha1",
				Digest:  "",
			},
		},
		{
			Name:         "cant find build file",
			BuildType:    BuildComposer(),
			WorkDir:      "testdata/whatever",
			Error:        true,
			ErrorMessage: "could not find match, reading dir open testdata/whatever: no such file or directory",
		},
	}

	test.Run(t, tests)
}
