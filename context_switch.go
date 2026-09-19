package main

import(
    "fmt"
    "runtime"
    "reflect"
    "sync"
    "time"
    "crypto/sha256"
)
//Go scheduler already does the M:N model, small pool of OS threads to handle 
//larger pool o fuser space threads

type tasks struct {
	//Arbitrary functions, string -> any
	functions map[string]any
	taken map[string]int //Map names to 
	mu sync.Mutex   //annoying
}

//Add something to registry
func(reg_ptr *tasks) Register(name string, fn any){
	//If empty,make
	if reg_ptr.functions == nil {
		reg_ptr.functions = make(map[string]any)
		reg_ptr.taken = make(map[string]int)
	}
	//Set the name and fn
	reg_ptr.functions[name] = fn
	reg_ptr.taken[name] = 0
}


//Spawn in user space threads- or goroutines- to handle tasks
func do_task(work *tasks){
	//Prevent races this is so beyond the scope fo this
	work.mu.Lock()
	
	var fn any
	valid_test := false
	//Find a new untaken tasks
	for name, took := range work.taken {
		if took != 50 {
			//Set the val up one
			work.taken[name] += 1
			//Grab it first
			fn = work.functions[name]
			valid_test = true
			break
		}
	}
	//Unlock before heavy duty
	work.mu.Unlock()

		//Never got any valid data
	if !valid_test {
		//All done
		return
	}
	
	//Execute the fn stuff through reflect magic 
	funct := reflect.ValueOf(fn)
	//Funct is actual concrete data of interface, 
	funct.Call(nil)
	
	//All are done
	do_task(work)	
} 



//AI helped me generate lists of random cpu intensive funtions
func primeSearch() {
	for n := 2; n < 200000; n++ {
		prime := true

		for d := 2; d*d <= n; d++ {
			if n%d == 0 {
				prime = false
				break
			}
		}

		_ = prime
	}
}
//And this one
func sha256Loop() {
	data := []byte("goroutine benchmark")

	for i := 0; i < 1_000_000; i++ {
		sum := sha256.Sum256(data)
		data = sum[:]
	}
}
//This one too
func fibonacci() {
	var fib func(int) int

	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		return fib(n-1) + fib(n-2)
	}

	_ = fib(38)
}
//Ok im making some functiosn that will context switch enough? to give 
//threads an issue. Or at least stimulate it 
func context_thread() {
	time.Sleep(50 * time.Millisecond)
	x:=0
	for i := 0; i < 1000000; i++{
		x+=1
	}
}
func sleep() {
	time.Sleep(500*time.Millisecond)
}

//Wrap over os ones so theyu can run like an os thread
func os_thread_work(work *tasks) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	do_task(work)
}

func main() {
	tasks := &tasks{}

	//Tried to take in arbitrary user input as a reflection function but
	//im dumb its not already compiled cant to that 
	//Now i just have a needlessly complicated reflection system for 
	//no reason

	//Take in a user defined function,
	//and also randomly apply a random number of functions
	//to the tasks table.
	// fmt.Print("Enter a function name: ")
	// //Go is pass by value, so give ref
	// _, err := fmt.Scan(&Fn_name)
	// fmt.Print("Enter a one lined function body: ")
	// _, err2 := fmt.Scan(&Fn_body)
	//
	// if err != nil & err2 != nil {
	// 	//Add the name and body to table
	// 	tasks.Register(Fn_name, Fn_body)
	// }
	
	//Add the ones in file
	tasks.Register("primeSearch", primeSearch)
	tasks.Register("sha256Loop", sha256Loop)
	tasks.Register("fibonacci", fibonacci)
	tasks.Register("context_thread", context_thread)

	//Now spin up 5 gothreads to handle ->
	var wg sync.WaitGroup
	
	start_go := time.Now()

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			do_task(tasks)
		}()
	}
	//Make sure all are done, get time
	wg.Wait()

	elapsed := time.Since(start_go)
	fmt.Println("The time for goroutines was : ", elapsed)
	for name := range tasks.taken {
		//Have to reset or the os will finish in nanoseconds funny
		tasks.taken[name] = 0
	}
	start_os := time.Now()

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			os_thread_work(tasks)
		}()
	}
		wg.Wait()
	fmt.Println("The time for os threads was: ", time.Since(start_os))


}
