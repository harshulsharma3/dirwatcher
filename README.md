# Directory Watcher


 **Install**

`git clone https://github.com/harshulsharma3/dirwatcher`


**Setup dependencies**

`go mod tidy`


**Run the app from source location**

`go run main.go`

**APIs**

*POST API To start task with below Json as payload*

`localhost:8080/config`

payload ->

{

    "directory": "C:/Users/Lenovo/Desktop/New folder/testing",
    
    "interval": 5,
    
    "magic_string":"hello"  
}


*GET API to get information for current running directory*

`localhost:8080/results`



*Post API to stop the process*

`localhost:8080/stop`
