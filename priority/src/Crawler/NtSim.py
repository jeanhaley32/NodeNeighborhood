import datetime

import RandomNetwork
from PriorityList import *
from Crawler import *

class NtSim:
    def __init__(self,size):
        self.size=size
        self.nt=RandomNetwork.Network(self.size,query_time_scale=0.01)
        self.cntrl=Controller()
        self.wrkr=Worker(self.nt)
        return
    
    def Start(self):
        self.cntrl.initialize()
        return
    
    def Iterate(self):
        jb0=self.cntrl.delegate()
        #print(jb0)
        jb1=self.wrkr.executeJob(jb0)
        #print(jb1)
        self.cntrl.ingest(jb1)
        return
