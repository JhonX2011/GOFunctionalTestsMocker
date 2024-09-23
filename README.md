# Golang mock server

Simple library to building your rest mock APIs. 


### Description


With this library, you can start a mock server to respond to your application's requests.

#### Mock your rest API step by step


#### Step 1: Start Server

```go
    server, mocker := mock.New()
     //start http interface
     serverAddress := ":9999"
     go server.Run(serverAddress)
```

#### Step 2: Mock your http petitions


## Mock through code
```go
     mocker.When(
        mock.Request().
            URLPattern("^/users/.*").
            HeaderIsEqualTo("Content-Type", "application/json").
            ParamContains("api-version", "2").
            WithPriority(1).
            Build(),
    ).ThenReturn(
        mock.Response().
            WithStatus(200).
            WithBodyAsString(`{"id":12313, tags:["CBT"]}`).
            WithHeader("Content-Type", "application/json").
            Build(),
    )
```

## Mock through http

When the server mock is started, expose the following resource to add mock through http:

Endpoint: http://localhost:9999/mock/mapping

Method: **POST**

Header: **Content-Type: application/json**

Json Body:

```json

{
    "request": //request especification to match
    { 
        "url"://request url to match - optional
        { 
            "equal_to": "/inventories" //condition, posible values equal_to, pattern, contains 
        },
        "method": "GET", //request method to match - optional
        "headers":  //request header to match - optional
        {
            "Accept": // header name
            { 
                "contains": "xml" //condition, posible values equal_to, pattern, contains 
            }
        },
        "query_parameters": //request query parameters to match - optional 
        {
            "api-version": // the query param name 
            {
                "equal_to": "2" //condition, posible values equal_to, pattern, contains 
            }
        },
        "body": { //request body to match - optional
            "equal_to": "{ \"name\": \"any name\"}" //condition, posible values equal_to, pattern, contains   
            
        }
    },
    "response": {
        "status": 200, // response status to return
        "body": // response body to return
        {
            "key": "value"
        },
        "headers": // response headers to return
        {
            "Content-Type": "application/json"
        }
    }
}

```
