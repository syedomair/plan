# [Plan Application](https://plans-admin.herokuapp.com)
It is an example application, which has two parts 
 * [Plan API](https://plans-api.herokuapp.com): A backend API application developed with Golang. 
 * [Plan Admin Dashboard](https://plans-admin.herokuapp.com): A Frontend application developed with ReactJs, Redux, and Material-UI.

Test user email/password  admin@gmail.com/123456

## Plan API Table of Contents
* [Concurrent Code](#concurrent-code)
* [Webhook](#webhook)
* [CURL Demo](#curl-demo)
* [Batch Processing](#batch-processing)
* [JWT Feature](#jwt-feature)
* [Common Lib](#common-lib)
* [Test](#test)
* [File Structure](#file-structure)




## Concurrent Code
When the new plan is added to the system, the following things happens asynchronously in the seperate goroutines, which resulted in a non blocking response.   
 * Increment the stat table plan counter.
 * Create a notification record for the plan in the notification table. 


https://github.com/syedomair/plan-api/blob/master/services/plans-svc/plan/repository.go#L72

## Webhook 
When the Plan is created or updated the notification message is send out asynchronously in the seperate goroutines. Here are the steps: 
 * Create a notification message queue
 * Send out a POST notification message to the httpbin site (for testing), which consume queuing element with the span of 2 seconds.  
 * Create a record in the notification_log, which store the success/failure of the POST notification message. 
 
https://github.com/syedomair/plan-api/blob/master/services/plans-svc/plan/repository.go#L150
https://github.com/syedomair/plan-api/blob/master/lib/functions.go#L18

## Curl Demo
curl -X GET https://plans-api.herokuapp.com/ping

{"data":{"response":"pong"},"result":"success"}



curl -X POST https://plans-api.herokuapp.com/login -d '{"email":"admin@gmail.com", "password":"123456"}'

{"data":{"token":"eyJhbGciOiJI............."},"result":"success"}



curl -X POST https://plans-api.herokuapp.com/plans -H 'Token:eyJhbGciOiJ.............' -d '{"title":"Test Plan", "status":"1", "cost":"1111", "validity":"30"}'

{"data":{"plan_id":"2eef6db1-4c6b-452a-aee5-ec702ad9341f"},"result":"success"}



curl -X GET https://plans-api.herokuapp.com/plans -H 'Token:eyJhbGciOiJIUzI1NiIsI.....................' 

{"data":{"count":"4","limit":"100","list":[{"id":"2eef6db1-4c6b-452a-aee5-ec702ad9341f","title":"Test Plan","status":1,"validity":30,"cost":1111,"created_at":"2019-04-02T16:09:19Z","updated_at":"2019-04-02T16:09:19Z"},{"id":"d107aa5c-9995-47b2-b34a-203ad655b621","title":"Monthly Plan","status":1,"validity":30,"cost":9999,"created_at":"","updated_at":""},{"id":"c9de5200-dbad-44b8-b5fc-ab1381730de7","title":"Weekly Plan","status":1,"validity":7,"cost":4999,"created_at":"","updated_at":""},{"id":"7be49965-fc69-4d15-ae53-aebcd7367402","title":"Daily Plan","status":1,"validity":1,"cost":1999,"created_at":"","updated_at":""}],"offset":"0"},"result":"success"}



## Batch Processing
Batch API is for processing complicated time consuming task.
 
When GET /users/batch is called, it simply trigger an asynchronous goroutine function and returns a batch_task_id instantaneously.
 
Then user call GET batch-tasks/{batch_task_id} if the batch task is not completed, it returns the the status only. If the task is completed, it returns the data. 


https://github.com/syedomair/plan-api/blob/master/services/users-svc/user/repository.go#L56	
https://github.com/syedomair/plan-api/blob/master/services/users-svc/user/repository.go#L25

## JWT Feature
On successful login, the system generate the JWT token with the following payload. 
```
	type Claims struct {
		CurrentUserId string `json:"current_user_id"`
		FirstName     string `json:"first_name"`
		LastName      string `json:"last_name"`
		Email         string `json:"email"`
		CreatedAt     int64  `json:"created_at"`
		jwt.StandardClaims
	}
```

## Common Lib
It contains code for common payload validation, common query string validation and  common response 


https://github.com/syedomair/plan-api/blob/master/lib/common_service.go


## Test 
There are multiple testing options available. 
 * Check Makefile to see different options
 * There is a standalone client application, which tests every API endpoint in the application.   

https://github.com/syedomair/plan-api/blob/master/test/testURL/urlTest.go
 * Test with Postman environment and collection files are available for local as well as Heroku server   

https://github.com/syedomair/plan-api/tree/master/test/postman 




## File Structure
```
├── database
│   ├── create_db.sql
│   ├── README
│   └── table_def.sql
├── lib
│   ├── auth.go
│   ├── common_service.go
│   ├── connection.go
│   ├── enum.go
│   ├── functions.go
│   └── logger.go
├── LICENSE
├── main.go
├── Makefile
├── models
│   ├── BatchTask.go
│   ├── Notification.go
│   ├── NotificationLog.go
│   ├── ParamConf.go
│   ├── Plan.go
│   ├── PlanMessage.go
│   ├── StatUserRegPerMonth.go
│   ├── TestResponse.go
│   ├── User.go
│   └── UserLogin.go
├── README.md
├── router.go
├── services
│   ├── plans-svc
│   │   └── plan
│   │       ├── handler.go
│   │       ├── handler_test.go
│   │       ├── repository.go
│   │       └── repository_test.go
│   ├── plnmsgs-svc
│   │   └── plnmsg
│   │       ├── handler.go
│   │       ├── handler_test.go
│   │       ├── repository.go
│   │       └── repository_test.go
│   ├── public-svc
│   │   └── public
│   │       ├── handler.go
│   │       ├── handler_test.go
│   │       ├── repository.go
│   │       └── repository_test.go
│   ├── stats-svc
│   │   └── stat
│   │       ├── handler.go
│   │       ├── handler_test.go
│   │       └── repository.go
│   └── users-svc
│       └── user
│           ├── handler.go
│           ├── handler_test.go
│           ├── repository.go
│           └── repository_test.go
├── test
│   ├── postman
│   │   ├── Heroku.postman_environment.json
│   │   ├── Local.postman_environment.json
│   │   └── PLAN_API_Test.postman_collection.json
│   ├── testdata
│   │   └── test_data.go
│   └── testURL
│       └── urlTest.go

```
