import datetime
from PriorityList import *
from Crawler import *

class Worker:

    def __init__(self,nt):
        self.nt=nt
        return

    def executeJob(self,jb):
        jb.start_time=datetime.datetime.now()
        jb.status='started'
        nid=jb.node.id
        ndoc=self.nt.query(nid)
        jb.node.doc=ndoc
        jb.end_time=datetime.datetime.now()
        jb.status='completed'
        return(jb)
