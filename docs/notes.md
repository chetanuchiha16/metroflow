always check if a goroutine has started or( if prev lines are blocking due to waiting for chan) to even get the chance to wait
also check if there is actual block to use goroutine,

jun 22: 2 pm:
the job was not being finished at all, the program was killed before workers finished their job