import RandomNetwork
from  PriorityList import *
from Crawler.Worker import *
from Crawler.Job import *

## Worker Test
nt=RandomNetwork.Network(100,1)
wrkr=Worker(nt)
nn=Node('0x00000000000000000000')
jb0=Job(nn,'CNCT')
print(jb0)
print('-----------------------')
jb1=wrkr.executeJob(jb0)
print(jb1)