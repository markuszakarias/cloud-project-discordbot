# Discord self help bot, project in PROG2005 
For this project we have decided to create a self help discord bot! This bot acts as a smaller component in a series of moves discord has been making to make itself less of a gaming specific platform. Bot's such as this act as a smaller piece in a larger composition of non-gaming specific services being created for more casual users. There are more and more bots like this popping up every day on https://top.gg/ which is the prime website for all public discord bots to add to a server.

### Group information
**Group 23: Self-help**

**Members: Jørgen Eriksen, Elvis Arifagic, Markus Strømseth, Salvador Bascunan**  
  
# Installation and deployment guide
## Structure
- Installing golang
- Adding our deployed bot to your discord server 
- How to deploy and add the bot yourself to a discord server

### Installing golang
### Windows

Download the installer from:

- https://golang.org/doc/install

After downloading the installer, run it and install Golang.

Verify the installation by running the command, in the cmd `go version`

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

Last, verify your install by running the following command in the terminal: `go version`

### MacOS

Download the installer from:

- https://golang.org/doc/install

After downloading the installer, run it and install Golang.

Note!! The package installs the Go distribution to /usr/local/go. The package should put the /usr/local/go/bin directory in your PATH environment variable. You may need to restart any open Terminal sessions for the change to take effect.

Verify the installation by running the following command in the terminal: `go version`  

## Connect bot to your Discord Server

To use our Discord bot, open the following invite link and select the Discord server you want to add the bot to: https://discord.com/oauth2/authorize?client_id=836983652251336775&scope=bot

If you want create your own Bot with our code, then create a new application on https://discord.com/developers/, copy the token that is generated, and insert it in DC_TOKEN in the .env file. Then you can get take the invite link (https://discord.com/oauth2/authorize?client_id=YOUR_CLIENT_ID_HERE&scope=bot) from the discord website to add the bot to your discord server or to invite other people to use the bot.


## Deploy and run the bot yourself (payment needed for full access)
1. Install docker for your system with the offical tutorial [docker installation](https://docs.docker.com/get-docker/)  

2. Create an account on all of the following api websites  
   a. https://spoonacular.com/food-api  
   b. https://newsapi.org/  
   c. https://openweathermap.org/api  

3. You will need service account credentials from firebase, first register as well a firebase account, follow these steps:  
   In the Firebase console, open Settings > Service Accounts.  
   Click Generate New Private Key, then confirm by clicking Generate Key.  
   Securely store the JSON file containing the key.  
   Rename the file to `firebasePrivateKey.json`  

4. Having all of the api-keys, replace the keys you have with all of the keys present in the `.env` file.  
   That would be the meal key, weather key and news key.  

5. Place your firebase key in database.go replacing the file which is already present  

6. Create an application on the developers page of discord [discord developers](https://discord.com/developers)  
   When the application is created, create a bot with the menu on the left and add that bot token to the `.env` file as well.  

7. If docker is not started now you need it running for the next part  

8. Now you can build the docker image, use the following command: `docker build -t discordbot -f Dockerfile`  

9. The image is now built and you can start the container using this command: `docker run --name your-container-name -d discordbot`  

Now the docker container is running and the bot is online ready for all the commands. From here there are multiple paths you can take for external deployment if you so wish. Our suggestion is to create a linux virtual machine on any platform (google cloud, azure, amazon) that provides those and run the docker container from there. The connection is based on the bot-token so as soon as that container is running the bot will run as well.

## Deployment

For this project we have decided to deploy the bot on Azure. We use the Azure Container Registry which store our docker image, and we use this together with Azure App Service which runs the application. The reason for doing is taking advantage of the freedom given in the task insofar as choosing our own technologies. Multiple people on the group felt like expolring the capabilities that Azure offers and we have found them to be easier to work with than others because of the container registry which makes the docker part very easy.


# Design choices and reflection 
## Architecture

For this we have created an [Architecture.md](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021-workspace/elvisa/projectgroup23/-/blob/master/ARCHITECTURE.md) which is a newer movement in software deveopment which says that anything over 2000 lines of code needs a document detailing the larger view of the project. Our intention is to be able to show that we understand what we have created and that less skilled developers are able to understand the way in which the system works.

## External dependencies

Going down the list of the functionality we wished to implement for this we needed some external dependencies.
Here is a list of all the direct external dependencies we have used. The team made sure that anything that we could do on our own
in a reasonable amount of lines of code we did write ourselves. Some things we had to use dependencies for (like mysql).

* https://github.com/josemiguelmelo/gocacheable

* https://github.com/tidwall/gjson

* https://github.com/denisenkom/go-mssqldb

* https://github.com/go-sql-driver/mysql

## Project technologies

- [x] -> Azure Deployment

- [x] -> Docker

- [x] -> Firebase storage

- [x] -> Azure SQL storage

- [x] -> Webhook functionality

- [x] -> Caching system

## Implemented Discord bot features

- [x] -> Command `!weather` with parameters

Be able to call the weather command with or without parameters. The default configuration will use the Get IP Location API to retrieve the location based on the ip address of the system that is running the application. Since this will in theory mean that the location would always be the same as the instance that is running the service. 

Because of this, we decided to add parameters to the command. Here you can pick a city location to retrieve the current weather for the specified location. Example:  

`!weather <city>` 


- [x] -> Command `!steamdeals` with parameters

Be able to call to !steamdeals with or without a parameter to change the amount of deals returned.


- [x] -> Command `!mealplan` with parameters

Call a mealplan to receive a breakfeast, dinner and snack meal plan.


- [x] -> Command `!newsletter` with parameters

Call !newsletter with or without parameter on location of the news. You can also give an optional number between 1-4 to change amount of headlines.


- [x] -> Command `!todo` with parameters

Be able to view, create, update and delete todotasks. These tasks are connected to your discord ID and will be stored in a azuresql database for persistence. This means that the tasks you created with your ID will still be there if the service reboots for whatever reason. There are different different commands for the different functions:

**View todo tasks**

`!todo mylist`

**Create todo task**

`!todo create This is my task`

**Update todo task description**

`!todo update <taskid> This is the update`

**Label task as finished**

`!todo finished <taskid>`

**Label task as inactive**

`!todo inactive <taskid>`

**Delete todo task**

`!todo delete <taskid>`


- [x] -> Command `!notifyweather` with parameters

**Register notification for a spesific city**

`!notifyweather <city>`
  
  
- [x] -> Command `!jokes` with parameters

**Get a random joke**

`!joke`

**Get all jokes that the user have made**

`!joke myjokes`

**Creates a joke**

`!joke create <joke text>`


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

## Original idea and what it turned into
This project for us started with wanting to create a travel data mash-up service but that proved to be very difficult. What made it difficult was not having access to the api's we needed. Real time flight data is hard to come by in a public api. We also could not find a simple way to find hiking trail information, some apis did have the information but in the wrong format (geo-data formats). The choice between this and self-help bot was always on the board for us, we went back and forth a bit on this but decided to give travl data the first try and then went for the bot instead. It also was very fun to work with a platform that you are pretty familiar with, in this case discord.

## What went well and what went bad
We utilized the branching system very well in this project such that we had development branches for each feature that were merged into master eventually. The balance between complexity, size and readability is something we feel that is done very well. With all the comments the project is very readable. Bad aspect of this project is the time that we spent on the caching system. It took a bit too much time to implement that functionality as there came up some more difficult bugs that took some time to solve.

## Hard aspects of the project.
As with anything we found it a bit hard to stop working and implementing new things. On several occasions thinking we are done we start talking about something we should add and losing a bit of time disscussing the idea. At some times during the course of the project we moved to quickly from feature to feature not focusing on finishing something right away.

## What have we learned
We as a group feel that we now are able to handle intermediate medium sized golang projects. The team feels very confindent in our ability to now program a larger multi-part application. Dockerizing an application with a multi-stage build is something we are now much more confident in our ability in having written a dockerfile and deployed it on Azure App Service.

## Total work cumulative by the group
Having not kept a strict tally on the work hours we would guess the total hours to be: `165`

