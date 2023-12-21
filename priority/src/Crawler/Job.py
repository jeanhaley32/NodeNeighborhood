import datetime
import hashlib
import json


class Job:
    
    def __init__(self,node,task):
        self.created_time=datetime.datetime.now()
        self.start_time=None
        self.end_time=None
        self.node=node
        self.task=task
        self.jid=self.__getID__()
        self.status='created'
        
    def __getID__(self):
        string=self.created_time.__str__()+self.node.__str__()+self.task
        m = hashlib.sha256()
        m.update(string.encode('utf-8'))
        return(m.hexdigest())
    
    def __str__(self):
        jj={
            'job_id':self.jid,
            'created_time':self.created_time.__str__(),
            'start_time':self.start_time.__str__(),
            'end_time':self.end_time.__str__(),
            'task':self.task,
            'node':self.node.__json__(),
            'status':self.status
        }
        return(json.dumps(jj,indent=2))
