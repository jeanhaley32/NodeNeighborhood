import numpy as np
import pandas as pd

import datetime
import time

from PriorityList import *
from Crawler.Job import *

class Controller:
    
    def __init__(self):
        self.pl=PriorityList()
        return


    def initialize(self):
        bootstrap='0x00000000000000000000'
        self.pl.insert(Node(bootstrap))
        return
        
    def delegate(self):
        hid=self.pl.getHead()
        node=self.pl.getNode(hid)        
        self.pl.remove(hid) # possibly change this to a 'pull' or 'pop'
        jb=Job(node,'CNCT')
        return(jb)
        
    def ingest(self,jb):
        self.pl.insert(jb.node)
        for nid in jb.node.doc['neigh']:
            if self.pl.inList(nid):
                self.pl.updateNode(nid,'PONG')
            else:
                self.pl.insert(Node(nid))
        self.pl.updateNode(jb.node.id,'CNCT')
        return
