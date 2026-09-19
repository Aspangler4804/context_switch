# Week of 09/09 - 09/16
This week we discussed alot, in particular we talked about virtualization and threads. We mentioned many ways to virtualize, and discussed different types of threads. 

One of the mentioned threads was OS threads, heavier threads that consume OS operating system attention, such that the OS schedules them and manages them. We also discussed Green threads, or user threads, that are managed by a runtime or userspace library, but essentially are invisible to the OS and can manage themselves. Context switching is highlighted as a benefit.

The green threads can operate at much more efficient rates, not worrying about blocking calls stopping the process, and can be efficiently created in Go with a goroutine. The class also mentioned a thread that lies between, such that a hybrid model that is managed by both. It turns out, the goroutine is more closesly aligned to this hybrid model, using an M:N scheduler to balance user space threads onto N OS threads.

I decided to implement a program that essentially throws heavy context switching operations at pools of goroutines and os bounded threads, to test and model the difference between such threads. 

Unfortunately, I decided to try and use reflection. For no reason. And I had to abandon it for the whole reason I wanted it, so the code now features reflection for no reason, when another structure could have done the job much more efficiently. Unfortunately most of my time was spent on reflection. Which is useless. Regardless, the program when ran outputs a slower time for OS threads than hybrid threads, modeling the discussion in class.

#Prerequisites
Go 1.22.2 or later
No external dependencies !
Uses only GO standard lib packages


Running
Clone the project
Run the program with go run .
Watch as your computer shows the time difference in threads as they complete a random order of tasks. With reflection. 
