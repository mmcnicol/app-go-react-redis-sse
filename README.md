# app-go-react-redis-sse

an "eager loading" demo using "server sent events"


## ui
* a connection to the backend is created, invoking a handler to eager load 3 types of data.
* event listeners are created for each type of data, to receive messages from the backend.
* for each type of data, a message is received to indicate whether or not is was fetched successfully from mockapis.
* if a type of data was fetched successfully, the user can click a corresponding button to call a backend handler to get the data from the redis cache and display it in the ui.
* the backend informs the ui to close the connection, once the outcome of the 3 calls is known.


## backend 
* a handler
  * uses goroutines to fetch 3 types of data by making calls to mockapis, concurrently.
  * as soon as the outcome of a call is known, a message is sent to the ui so that it can indicate success or failure.
  * a WaitGroup is used so that the backend can inform the ui to close the connection, once the outcome of the 3 calls is known.
  * when a call is successful, the data is cached in redis.
* a handler
  * returns requested type of data from the redis cache

