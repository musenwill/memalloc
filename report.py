import os
import csv
import sys
import matplotlib.pyplot as plt



filepath = 'mem.csv'
if len(sys.argv) >= 2 and sys.argv[1] != '':
    filepath = sys.argv[1]
dir = os.path.dirname(filepath)
basename = os.path.basename(filepath)


amount, syss, heapSys = [], [], []

alloc, totalAlloc = [], []

heapAlloc, heapIdle, heapInuse, heapReleased = [], [], [], []

stackInuse, stackSys = [], []

mspanInuse, mspanSys, mcacheInuse, mcacheSys, ohterSys, gcSys = [], [], [], [], [], []

count, mallocs, frees, heapObjects, numGC = [], [], [], [], []


with open(filepath, mode ='r') as file:
    csvFile = csv.DictReader(file)
    for line in csvFile:
        amount.append(float(line['amount']))
        syss.append(float(line['sys']))
        heapSys.append(float(line['heapSys']))

        alloc.append(float(line['alloc']))
        totalAlloc.append(float(line['totalAlloc']))

        heapAlloc.append(float(line['heapAlloc']))
        heapIdle.append(float(line['heapIdle']))
        heapInuse.append(float(line['heapInuse']))
        heapReleased.append(float(line['heapReleased']))

        stackInuse.append(float(line['stackInuse']))
        stackSys.append(float(line['stackSys']))

        mspanInuse.append(float(line['mspanInuse']))
        mspanSys.append(float(line['mspanSys']))
        mcacheInuse.append(float( line['mcacheInuse']))
        mcacheSys.append(float( line['mcacheSys']))
        ohterSys.append(float( line['ohterSys']))
        gcSys.append(float( line['gcSys']))

        count.append(float( line['count']))
        mallocs.append(float( line['mallocs']))
        frees.append(float( line['frees']))
        heapObjects.append(float( line['heapObjects']))
        numGC.append(float( line['numGC']))


plt.figure(figsize=(40,35), dpi=256)

plt.subplot(7,1,1)
plt.plot(amount, label='amount')
plt.plot(syss, label='sys')
plt.plot(heapSys, label='heapSys')
plt.title(label="global")
plt.legend()

plt.subplot(7,1,2)
plt.plot(heapAlloc, label='heapAlloc')
plt.plot(heapIdle, label='heapIdle')
plt.plot(heapInuse, label='heapInuse')
plt.plot(heapReleased, label='heapReleased')
plt.title(label="heap")
plt.legend()

plt.subplot(7,1,3)
plt.plot(mspanInuse, label='mspanInuse')
plt.plot(mspanSys, label='mspanSys')
plt.plot(mcacheInuse, label='mcacheInuse')
plt.plot(mcacheSys, label='mcacheSys')
plt.plot(ohterSys, label='ohterSys')
# plt.plot(gcSys, label='gcSys')
plt.title(label="span")
plt.legend()

plt.subplot(7,1,4)
plt.plot(alloc, label='alloc')
plt.plot(totalAlloc, label='totalAlloc')
plt.title(label="cache/alloc")
plt.legend()

plt.subplot(7,1,5)
plt.plot(stackInuse, label='stackInuse')
plt.plot(stackSys, label='stackSys')
plt.plot(numGC, label='numGC')
plt.title(label="stack")
plt.legend()

plt.subplot(7,1,6)
plt.plot(mallocs, label='mallocs')
plt.plot(frees, label='frees')
plt.plot(heapObjects, label='heapObjects')
plt.title(label="mallocs")
plt.legend()

plt.subplot(7,1,7)
plt.plot(count, label='count')
plt.title(label="count")
plt.legend()

filename = os.path.join(dir, basename+".png")
plt.savefig(filename)
