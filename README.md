# projectGroup23

Project in cloud course

### Clone repo and run locally

If you clone the repo you have to run the server locally, with these instructions on installation and running.

## Installation

### Windows

Download the installer from:

- https://golang.org/doc/install

After downloading the installer, run it and install Golang.

Verify the installation by running the command, in the cmd

```
go version
```

### Linux

If you have a previous version on Go installed, remove it before installing another.

Download the installer from:

- https://golang.org/doc/install

Extract the folder into `/usr/local`, which creates a Go tree in `/usr/local/go`. Do this by running the following command:

```
tar -C /usr/local -xzf go1.16.linux-amd64.tar.gz
```

Then add the `/usr/local/go/bin` to the PATH variable, this can be done by adding the following line to your `$HOME/.profile` or `/etc/profile`

```
export PATH=$PATH:/usr/local/go/bin
```

Note that you may need to restart your computer for it to take effect.

Last, verify your install by running the following command in the terminal:

```
go version
```

### MacOS

Download the installer from:

- https://golang.org/doc/install

After downloading the installer, run it and install Golang.

Note!! The package installs the Go distribution to /usr/local/go. The package should put the /usr/local/go/bin directory in your PATH environment variable. You may need to restart any open Terminal sessions for the change to take effect.

Verify the installation by running the following command in the terminal:

```
go version
```

## Run the program

After installing, clone the repo into a clean directory. Inside the repo directory, open a terminal and run:

```
go run .
```

Or to create an executable

```
go build .
```

To run the executable, run `./<executable>`

## Caching

We use the Read-Through strategy for caching and all writes to the database happens from API's. The application in case of a hit reads from cache, otherwise it reads from database and creates a cache on the read data.

![Read-Trough_Cache](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021-workspace/elvisa/projectgroup23/-/raw/master/assets/read_through_cache.PNG)

We decided to use BigCache as our cache provider, as our application only deals with read operations from the database and it fits best with regards to scaleability. The performance comparison shows this:

![Cache performance comparison](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021-workspace/elvisa/projectgroup23/-/raw/master/assets/cache_performance_comparison.PNG)

We added a wrapper module to easily abstract the call to the existing caching system which makes it easier to handle cache when calling functions.

https://github.com/josemiguelmelo/gocacheable
