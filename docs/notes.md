always check if a goroutine has started or( if prev lines are blocking due to waiting for chan) to even get the chance to wait
also check if there is actual block to use goroutine,

jun 22: 2 pm - 5:40 pm:
- the job was not being finished at all, the program was killed before workers finished their job
- close the channel
- removed os.exit and now avoiding the infinite loop itself, which naturally accepts only one client

jun 24: 9:47:
- added redis queue using RPUSH and BLPOP