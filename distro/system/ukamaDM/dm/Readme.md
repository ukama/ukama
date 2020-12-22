Use the following method to build client,server and bootstrap together.
### Server, Client, Bootstrap Build:
 * Create a build directory and change to that.
 * ``cmake [device-mangement git folder]``
 * ``make``
 * ``./server/server [Options]``
 * ``./client/client [Options]``
 * ``./bootstrap/bootstrap [Options]``

Alternatively you can also build client, server and bootstap individually using the following method.
### Server Build:
 * Create a build directory and change to that.
 * ``cmake [device-mangement git folder]/server``
 * ``make``
 * ``./server [Options]``

### Client Build:
 * Create a build directory and change to that.
 * ``cmake [device-mangement git folder]/client``
 * ``make``
 * ``./client [Options]``

 ### Bootstrap Build:
 * Create a build directory and change to that.
 * ``cmake [device-mangement git folder]/bootstrap``
 * ``make``
 * ``./bootstrap [Options]``