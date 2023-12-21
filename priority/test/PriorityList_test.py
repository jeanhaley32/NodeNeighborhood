import RandomNetwork
import PriorityList

nt=RandomNetwork.Network(100,0.001)
pl=PriorityList.PriorityList()
bootstrap='0x00000000000000000000'
pl.insert(PriorityList.Node(bootstrap))
pl.showPriority()