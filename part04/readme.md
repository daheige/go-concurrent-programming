# Mutex vs Channel

In the above, we used two methods of mutex and channel to solve 
the data race problem of concurrent programs.
So when should we use a mutex, and when should we use a channel? 
The answer lies in the problem you are trying to solve.
If the problem you're trying to solve is better suited for a mutex, 
then go ahead and use a mutex. .
Use it if the question seems to fit better with the channel.

Most Go newbies try to use channels to solve all concurrency problems 
because it's a cool feature of the Go language.
this is not right. Language give us the option to use Mutex or Channel, 
there is nothing wrong with choosing either.

Typically, channels are used when goroutines need to communicate with each other,
and mutexes are used to ensure that only one goroutine can access critical parts 
of the code at a time.
Of the problem we solved above, I'd prefer to use a mutex because this problem doesn't 
require any communication between goroutines. It is only necessary to ensure that 
only one goroutine has
the right to use the shared variable at the same time.
The mutex was originally born to solve this problem, 
so using the mutex is a more natural choice.
