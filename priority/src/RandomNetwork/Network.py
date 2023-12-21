import hashlib
import random
import json
import time

from RandomNetwork import *

class Network:
    
    def __init__(self,size,query_time_scale=0.001):
        self.size=size
        self.query_time_scale=query_time_scale
        self.nt=generateNetwork(self.size)
        return
        
    def query(self,nid):
        tt=self.query_time_scale*random.randint(20,100)/100
        time.sleep(tt)
        res=None
        if nid in self.nt:
            res=self.nt[nid]
        return(res)
    
    def printNode(self,nid):
        json_formatted_str=json.dumps(self.nt[nid], indent=2)
        print(json_formatted_str)
        return
        
    def printNodeIds(self,s,e):
        ll=list(self.nt.keys())[s:e]
        for nid in ll:
            print(nid)
        return

    def printNodes(self,s,e):
        ll=list(self.nt.keys())[s:e]
        for nid in ll:
            self.printNode(nid)
        return
