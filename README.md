# Observer
Implementation of the observer pattern for Golang based on different \"triggers\" which control when to notify the observers. 

```go get github.com/konimarti/observer```

## Notes on the implementations
Two type of observers are implemented suitable for different use cases:
* Channel-based observers accept new values through a ```chan interface{}``` channel, and
* Function-based observers that collect new values in regular intervals from a ```func() interface{}``` function.

To create an new observer struct, a trigger has to be given as well. 
The trigger controls the behavior of the observer implementation and determines when to notify the observers. 
This allows for a greater flexibility with observer pattern so that it can be easily extended for different use cases and user-defined triggers. 

Calling the ```Subscribe() Subscriber``` method on the observer returns a subscriber responds to events and provides the current values.

## Example with channels

## Example with anonymous functions

## Credits

The implementation of the observer pattern was inspired by ```github.com/imkira/go-observer```.

## Disclaimer

This package is still work-in-progress. Interfaces might still change substantially. It is also not recommend to use it in a production environment at this point.





