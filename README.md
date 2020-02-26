# echo-template
A RESTful backend template using [echo](https://github.com/labstack/echo)

## Usage
### docs
```shell
make doc
./build/doc
```
then visit [localhost:20002/index.html](http://localhost:20002/index.html) for docs.
### build
```shell
make
```

## Tools
### Generate rsa keypair for JWT
```
openssl genrsa > private.rsa
openssl rsa -in private.rsa -pubout > public.rsa
```