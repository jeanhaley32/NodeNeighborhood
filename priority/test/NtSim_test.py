import RandomNetwork
from PriorityList import *
from Crawler import *

ntsim=NtSim(10000)
ntsim.Start()
for ii in range(500):
    ntsim.Iterate()
    #print(len(ntsim.cntrl.pl.pList))
ntsim.cntrl.pl.showPriority()
