# Adjust test task
The tool which makes http requests and prints the address of the request along with the MD5 hash of the response.

## Usage
1. Build app
```shell
    export GO111MODULE=on 
    go build -o myhttp cmd/walker
```

2. Run app
```shell
  ./myhttp google.com vk.com habr.com hh.ru facebook.com youtube.com 
  ./myhttp --parallel=1 google.com vk.com habr.com hh.ru facebook.com youtube.com
```

3. Test app
```shell
    export GO111MODULE=on 
    go test
    
```