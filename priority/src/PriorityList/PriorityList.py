import datetime
import json

from PriorityList import *

class Node:

    def __init__(self,nid):
        tt=datetime.datetime.now()
        ut=datetime.datetime.timestamp(tt)
        self.id=nid
        self.doc={'id':nid,'loc':None,'neigh':[]}
        self.creation=ut
        self.last_connect=None
        self.last_pong=ut
        self.priority=ut
        self.log={'CNCT':0,'PONG':0}
        return
    
    def loadDoc(self,nodeDoc):
        assert self.id==nodeDoc['id'],"Id mismatch"
        self.doc=nodeDoc
        return
    
    def connect(self):
        tt=datetime.datetime.now()
        ut=datetime.datetime.timestamp(tt)
        self.last_connect=ut
        self.priority=ut
        self.log['CNCT']+=1
        return
        
    def pong(self):
        tt=datetime.datetime.now()
        ut=datetime.datetime.timestamp(tt)
        self.last_pong=ut
        self.priority=0.7*ut+0.3*self.priority  # (ut-self.priority)*0.7+self.priority
        self.log['PONG']+=1
        return
    
    def __str__(self):
        jj={
            'node_id':self.id,
            'priority':self.priority.__str__(),
            'last_pong':self.last_pong.__str__(),
            'last_connect':self.last_connect.__str__(),
            'doc':self.doc
        }
        return(json.dumps(jj,indent=2))
 
    def __json__(self):
        return(self.doc)
    

class PriorityList:
    
    def __init__(self):
        self.plist={}
        return

    def inList(self,nid):
        return(nid in self.plist)

    def insert(self,node):
        self.plist[node.id]=node
        return
    
    def remove(self,nid):
        del(self.plist[nid])
        
    def updateNode(self,nid,update):
        if update=='CNCT':
            self.plist[nid].connect()
        elif update=='PONG':
            self.plist[nid].pong()
        else:
            raise Exception("Unknown update")
        return
        
    def getHead(self):
        tmp={self.plist[nid].priority:self.plist[nid].id for nid in self.plist}
        return(tmp[min(tmp)])

    def getTail(self):
        tmp={self.plist[nid].priority:self.plist[nid].id for nid in self.plist}
        return(tmp[max(tmp)])

    def getNode(self,nid):
        return(self.plist[nid])

    def showPriority(self):
        print('length=',len(self.plist))
        tmp={self.plist[nid].priority:nid for nid in self.plist}
        tmp=dict(sorted(tmp.items()))
        for pp in tmp:
            nid=tmp[pp]
            print(pp,nid,self.plist[nid].log)
        return
