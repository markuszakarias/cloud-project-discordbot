# Project in PROG 2005-2021, NTNU

Group 23: Self-help

Members: Jørgen Eriksen, Elvis Arifagic, Markus Strømseth, Salvador Bascunan

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

## Connect bot to Discord Server

To use our Discord bot, open the following invite link and select the Discord server you want to add the bot to: https://discord.com/oauth2/authorize?client_id=836983652251336775&scope=bot

If you want create your own Bot with our code, then create a new application on https://discord.com/developers/, copy the token that is generated, and insert it in DC_TOKEN in the .env file. Then you can get take the invite link (https://discord.com/oauth2/authorize?client_id=YOUR_CLIENT_ID_HERE&scope=bot) from the discord website to add the bot to your discord server or to invite other people to use the bot.

## Deployment

We have decided to deploy the bot on Azure. We use the Azure Container Registry which store our docker image, and we use this together with Azure App Service which runs the application.

## Architecture

image maybe?

## Introduction

Something about the thought process in the project. 

## Storage

Firestore and azuresql..

## External dependencies

* https://github.com/josemiguelmelo/gocacheable

* https://github.com/tidwall/gjson

## Project technologies (unsure about title)

- [x] -> Azure Deployment

- [x] -> Docker

- [x] -> Firebase storage

- [x] -> Azure SQL storage

- [x] -> Webhook functionality

- [x] -> Caching system

## Discord bot features

- [x] -> Command `!weather` with parameters

Be able to call the weather command with or without parameters. The default configuration will use the Get IP Location API to retrieve the location based on the ip address of the system that is running the application. Since this will in theory mean that the location would always be the same as the instance that is running the service. 

Because of this, we decided to add parameters to the command. Here you can pick a city location to retrieve the current weather for the specified location. Example:

```discord
!weather Gjøvik
```

- [x] -> Command `!steamdeals` with parameters

Be able to call to steamdeals...

- [x] -> Command `!mealplan` with parameters

Be able to call to mealplan...

- [x] -> Command `!newsletter` with parameters

Be able to call to newsletter...

- [x] -> Command `!todo` with parameters

Be able to view, create, update and delete todotasks. These tasks are connected to your discord ID and will be stored in a azuresql database for persistence. This means that the tasks you created with your ID will still be there if the service reboots for whatever reason. There are different different commands for the different functions:

**View todo tasks**

```
!todo mylist
```

**Create todo task**

```
!todo create This is my task
```


**Update todo task description**

```
!todo update <taskid> This is the update
```

**Label task as finished**

```
!todo finished <taskid>
```

**Label task as inactive**

```
!todo inactive <taskid>
```

**Delete todo task**

```
!todo delete <taskid>
```

- [_] -> Command `!notifyweather` with parameters

**Register notification for a spesific city**

```
!notifyweather <city>
```

- [x] -> Command `!jokes` with parameters

**Get a random joke**

```
!joke
```

**Get all jokes that the user have made**

```
!joke myjokes
```

**Creates a joke**

```
!joke create <joke text>
```

- [x] -> Command `!help` without parameters

The self-help bot has a help command so the user can view and get instruction on how the different commands work. The `!help` command, without any parameters, will print a helper message giving a small introduction to the discord bot and give further instructions on how to get information about the other commands.

- [x] -> Command `!help todo`

- [x] -> Command `!help newsletter`

- [x] -> Command `!help mealplan`

- [x] -> Command `!help weather`

- [x] -> Command `!help notifyweather`


## Caching

We use the Read-Through strategy for caching and all writes to the database happens from API's. The application in case of a hit reads from cache, otherwise it reads from database and creates a cache on the read data.

![Read-Trough_Cache](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021-workspace/elvisa/projectgroup23/-/raw/master/assets/read_through_cache.PNG)

We decided to use BigCache as our cache provider, as our application only deals with read operations from the database and it fits best with regards to scaleability. The performance comparison shows this:

![Cache performance comparison](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021-workspace/elvisa/projectgroup23/-/raw/master/assets/cache_performance_comparison.PNG)

We added a wrapper module to easily abstract the call to the existing caching system which makes it easier to handle cache when calling functions.

https://github.com/josemiguelmelo/gocacheable


## How Webhook (notification) works
The webhook is used as a functionality to notify selected user with the weather details of the day, every day at 8am. When the server starts it sends a notification at once, and then calculates the time uintil 8am. So after the first notification gets send at 8am, then the webhook wil run every 24 hour after that. The webhook information is stored in Cloud Firestore, and has a collection where each document represent one webhook/user. The document has only two fiels, City (city to get weather information from) and UserId (uniqe user id in Discord). We know we could have UserId as document id, but if the project where to expand in a way that users could have multiple webhooks, then that solution would not have worked.

