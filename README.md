# transaction-broadcasting

# Solution Design
- Simplicity with flexibility, that's the approach I am going to solve with the given requirement. As it is said in the requirements that this client module will be integrated into another program, the only thing that comes to my mind is that this must be a library which will be implemented by other services or programs. When it comes to a library, it must be generalized so that it will work anywhere like 'Write Once, Integrate Anywhere'. Also, it must allow flexibility for those who implement this library and let them manually configure and optimize with their business requirements and systems.

- The simplicity is achieved by providing out-of-the-box methods for users. At first, I degisned it to be just a single method handling broadcasting a transaction. However, that comes with less flexibility. Therefore, it divides the process of transaction broadcasting into smaller methods which can work independenty.

- For flexibility, separating the methods into smaller ones allows such ability. users who implement this module can do some busiess logic along the way or before calling the next method. This could be applying local logic or inserting some data that returns from each operation into a database. 

# Trade-Offs
- Even though this design allows flexibility, it still have a trade-off if a user would like to implement it as a middleware. This module is implemented and aimed for working in the business logic level or inside some modules. With such case, it still requires a user to extract the data of the transaction from the reqeust body.

- Allowing such flexibility means more human errors. For example, users must form a url of checking the status themselves. This could lead to an error suhc as no endpoint or not found. The reason of this design is for cases where the the endpoind is different from the broadcast and where the path is different.

# Asssumption
- Users may arrange the steps of the transaction broadcasting incorrectly. However, each method has a specific name of the parameters which the method requires. Even that, some of parameters shares the data types, but internal logic of each will handle that and returns the error. Therefore, that won't be a problem. 


# Scalability, Reliability, and Performace Aspects
## Scalability
- From my experience and the nature of a library, the library itself cannot achieve with only itself. Even if it can, it sill needs to be done with the code of a program which implements the library and also the infrastructure. Working together to achieve such scalability.

- To work that out, I integrate the module with a custom HTTP library (the third party module is used in this module). This will allow the users who implement this module can configure the HTTP request for their logic and optimize for the business and infrastructure capacity. These configurations include but not limited to timeout, connection pooling, and rate limits.

- I think achieving scalability needs four parts working together. The first part is the infrastructure, which monitors the traffic and decides when to scale up and scale down. With the action, the second part plays a role. The service which implenements the library must be designed with ability to handle increasing workloads. The third part which is embedded in the second part or in the service can be configured with the parameters for the HTTP requests so that it can help optimizing the sca;ability as well. Finally, the end service, the forth one. There could be a problem if other parts are scaled up, but left the end service handle all increasing traffic with the current resources.

- To sum up, the library cannnot ahcieve the scalability alone. That's why I add custom HTTP reqeust for a user to decide how to configure in this library.

## Reliability
- This could be done the HTTP library too. Inside HTTP module, there is a retry mechanism which can retrying the process of calling the endpoint until it succeeds or fails. This will prevent the prevent the program to crash and let it continue processing the transaction.

# Documentation
## How to Use it
### NewBroadcastService()
- the first method of this module
- instantiate the object of this module so that other methods can be accessed and implemented.

### WithRetryRequest()
- this method helps configure the retry machanism of the transaction broadcasting logic (not the HTTP library).
- if not provided, it will use the default values.

### WithCustomHTTPRequest()
- this method helps configure the HTTP request.
- if not provide, it will use the default value of the library.

### BroadcastTransaction()
- the method which will broadcast the transaction.
- it required the url of the endpoint along with the request of the transaction.
- the request has three fields. Symbol, Price, and Timestamp
- Inside this method, there is validation logic which checks the values of the reqeust.
- the return values are the txHash and the error

### MonitorTransaction()
- this method is for monitoring the status of the transaction.
- the user must form a url ans pass it to this method
- it returns txStatus and error.

### HandleStatus()
- this is for handling the action after receiving the status from MonitorTransaction()
- the logic inside this method is switch case whose cases are based on the provided status from the requirements
- Also, in this, there is a logic attached to case 'PENDING'. If the status is 'PENDING', the retry mechanism will be executed. (this will be explained in details later in Transaction Status Handling section below)
- it returns the final txStatus and error

## How to integrate it
### I make this a public repository (for my case). This could be done with other remote repositories which can store the module and allow the user can download and implement it
### Please refer to [GitHub Pages](https://github.com/paweenwatkwanja/raks-coin-exchange). In this example repository, how this module is integrated is shown.
### Briefly, if the project is implemented in Go, users can simle run 'go get https://github.com/paweenwatkwanja/transaction-broadcasting' and instantiate this object wherever they want to broadcast the transaction.

# Transaction Status Handling
## CONFIRMED
### For this status, the module exits the method right away and returns the 'CONFIRMED' string along with error with nil so that the user can handle what they would like to do next.
## FAILED
### For 'FAILED',  the module exits the method right away and returns the 'FAILED' string along with error so that the user can handle what they would like to do next.
## PENDING
### For this status, this means that the transaction is in process so the program must wait.
### To achieve this, there will be a method which will retry calling the service to check the status periodically. The retry counts and duration to wait can be configured by the users.
### Once the retry mechanism reaches the condition to stop, it will return the error saying that 'status is still pending' and exits.
### if the retry mechanism receives a new status that can be 'CONFIRMED', 'FAILED', or 'DNE', the method stops and return string of that status with its error, except for 'CONFIRMED' whose error is nil.
## DNE
### For 'DNE',  the module exits the method right away and returns the 'DNE' string along with error so that the user can handle what they would like to do next.