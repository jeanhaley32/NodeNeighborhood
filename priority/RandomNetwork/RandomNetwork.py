#################################################################################################################
##   Random Network Generator
#################################################################################################################

import hashlib
import random

import numpy as np
from scipy.stats import poisson

def randomId():
    tt=random.randint(0,10000)
    string=tt.__str__()
    m = hashlib.sha1()
    m.update(string.encode('utf-8'))
    return('0x'+m.hexdigest()[:20])

def randomCity():
    Cities=['Kabul','Herat','Mazar-i-Sharif','Kandahar','Jalalabad','Lashkargah','Kunduz','Taloqan','Puli Khumri',\
     'Sheberghan','Zaranj','Maymana','Ghazni','Khost','Charikar','Fayzabad','Tarinkot','Gardez']
    return(random.choice(Cities))

def generateNode():
    return({'id':randomId(),'loc':randomCity(),'neigh':[]})


def generateNeigh(Vid,mu=10):
    nn=min([poisson.rvs(mu),len(Vid)])           # Poisson dist
    #nn=random.randint(0,min([10,len(Vid)]))     # Uniform
    neigh=set([random.choice(Vid) for ii in range(nn)])
    return(neigh)

def generateNetwork(nsize,mu=10):
    V=[generateNode() for ii in range(nsize)]
    V.append({'id':'0x00000000000000000000','loc':'Kandahar','neigh':[]})
    random.shuffle(V)
    Vp={nn['id']:nn for nn in V}
    Vid=[nn['id'] for nn in V]
    for nn in V:
        rneigh=generateNeigh(Vid,mu)-set([nn['id']])
        nn['neigh']=list(rneigh|set(nn['neigh']))
        for ii in rneigh:
            if not nn['id'] in Vp[ii]['neigh']:
                Vp[ii]['neigh'].append(nn['id'])
    return(Vp)
