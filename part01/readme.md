# ctx
    Although context cancellation in Go is a versatile tool, there are
    a few things you need to keep in mind before proceeding. The most important of these is
    that the context can only be cancelled once. If you want to propagate multiple errors 
    in the same operation, then using contextual cancellation may not be the best option. 
    The use case for cancellation context is when you actually want to cancel an operation,
    not just notify downstream processes that an error occurred. Another thing to keep 
    in mind is that the same context instance should be passed to all functions and 
    goroutines you might want to cancel.

    Wrapping a context that already supports cancellation with WithTimeout or 
    WithCancel will create a variety of situations that may cause your context to 
    be canceled, and such secondary wrapping should be avoided.
    
    After the upper-level task is cancelled, all lower-level tasks will be cancelled;

    After the task of a certain layer in the middle is canceled, only the lower-level 
    task of the current task will be canceled, and it will not affect the upper-level 
    task and the task of the same level;

    Metadata can be shared between goroutines thread-safely; 
# points

    For the use of Context, Go officially mentioned the following points:

     Don't stuff Context into structs. The Context type is directly used as the first parameter of
     the function, and it is generally named ctx.

     Don't pass a nil context to the function, if you really don't know what to pass, the 
     standard library's TODO method prepares an emptyCtx for you.

     Don't stuff the types that should be used as function parameters into the context.
     The context should store some data shared by goroutines, such as server 
     information and so on.

    I think the first half of the first point is very far-fetched. 
    For example, in the official net/http package, the Context is placed in the 
    Request structure. The other points are indeed things to pay attention to.
