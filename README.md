# Observer
Implementation of the observer pattern for Golang based on different \"triggers\" which control when to notify the observers. 

```go get github.com/konimarti/observer```

## Notes on the implementations
Two type of observers are implemented suitable for different use cases:
* Channel-based observers accept new values through a ```chan interface{}``` channel, and
* Interval-based observers that collect new values in regular intervals from a ```func() interface{}``` function.

To create an new observer struct, a trigger has to be given as well. 
The trigger controls the behavior of the observer implementation and determines when to notify the observers. 
This allows for a greater flexibility with observer pattern so that it can be easily extended for different use cases and user-defined triggers. 

Calling the ```Subscription() chan interface{}``` method on the observers returns a channel that is used to notify the observers in case the triggers is activated.

## Example with channels

## Example with regular intervals

## Disclaimer

This package is still work-in-progress. Interfaces might still change substantially. It is also not recommend to use it in a production environment at this point.





