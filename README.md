# socks-http-proxy

## A http proxy over SOCKS5.

You can add domain names to `domain.file` file , this means that both HTTP and HTTPS requests for these domain will go through the SOCKS5.

```
http/https request
	 └── match domain name(domain.file)
			     ├── N
			     │   └── Local
			     └── Y
				 └── SOCKS5
```

## Installation

```shell
./build.sh
```

## Using it

```
./proxy -f ./cnf.json

export http_proxy=http://127.0.0.1:8888
export https_proxy=http://127.0.0.1:8888

curl https://en.wikipedia.org/wiki/Hello
```

