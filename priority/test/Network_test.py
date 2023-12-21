import RandomNetwork

nt=RandomNetwork.Network(100,0.1)
#  in case of 'import RandomNetwork' use nt=RandomNetwork.Network(100)

nt.printNodeIds(24,34)


nt.printNode('0x00000000000000000000')
bdoc=nt.query('0x00000000000000000000')
print(bdoc['neigh'])