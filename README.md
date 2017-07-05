#Metar

##Steps to Run

1. Install golang
###Mac
```bash
    $ brew install go
```
###Linux

```bash
    $ sudo apt-get install golang
```

2. Setup go paths

```bash
    $ mkdir $HOME/go
    $ export GOPATH=$HOME/go
    $ export PATH=$PATH:$GOPATH/bin
```

Add it in your .bashrc if you want these configs for ever


3. Get this package

```bash
    $ go get github.com/gokultp/metar
```

4. Traverse to the package

```bash
    $ cd $GOPATH/src/github.com/gokultp/metar
```

5. Set ENV (redis url)

```bash
    $ export REDIS_URL=localhost:6379
```

6. Compile and Run
 
 either 

```bash
    $ go run main.go
```

or

```bash
    $ go build main.go
    $ ./main
```



