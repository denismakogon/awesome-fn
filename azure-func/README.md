Azure function with an HTTP trigger on Fn
=========================================

Purpose
-------
This is an experimental work being done to show how to integrate Azure Function written in Node.JS in Fn Project.

Assumption
----------
For simpler and faster prototyping Azure introduced what they called "local development mode", where a developer gets 
a Docker container with the web server pre-installed for the particular runtime (two available: Node.JS, .Net)

My personal assumption that is derived from my own experience is that Azure actually runs a web server for your functions in their production public cloud
and here's why:

 - with Azure Functions V2, Azure restricts a number of programming languages within a single serverless application
 - in Azure Functions, you pay for an idle, you have to turn off your serverless application in order to keep yourself away from charges

Okay, so, my assumption is that in Azure Function, a serverless application is nothing but a **container** or a **virtual machine**.

How it works?
-------------
I did the following:

 - used an instructions how to create a Node.JS function with an HTTP Trigger in Azure and run it locally
 - wrote the Golang function that acts as a proxy between Fn and tiny .Net web server inside of a function

So, when a user submits the request to Fn it will start a Go function that will start a tiny .Net web server and act as a proxy.

Known issues
------------
Cold start isn't going anywhere. .Net web server needs 1 second to start.

Decisions made
--------------
You may ask why am i not including .Net web server start into a supervisor.d or similar things.
Well, supervisor.d requires Python to be installed, but in order to keep container small, it's simpler to start .Net server in a background before starting Go function dispatching.

How to use?
-----------

Do the following:
```bash
cd init-images/
docker build -t azure-func:http-trigger-node -f node.dockerfile .
fn init --init-image=reverser-init http-trigger-nodejs
cd azure-func/
fn --verbose deploy --app azure --local
fn create trigger azure http-trigger-nodejs http-trigger-nodejs
```

List the triggers:
```bash
fn list t azure
FUNCTION        NAME                 ID                              TYPE    SOURCE                ENDPOINT
azure-node      http-trigger-nodejs  01CTGB3C2SNG8G00RZJ000000W      http    /http-trigger-nodejs  http://localhost:8080/t/azure/http-trigger-nodejs
```

Call a function:
```bash
echo -e '{"name":"Denis"} | fn invoke azure http-trigger-nodejs'
```

The result suppose to be:
```bash
Hello Denis
```

TODO
----
Add .Net runtime
