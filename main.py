import numpy as np
import pandas as pd 


b = np.random.normal(0,1,(120,1014))
df = pd.DataFrame(b)
df.to_csv("wow3.csv", index=False, header=False)

b_ = np.random.normal(0,1,(120,1200))
for i in range(120):
    b_[i] = [0] * (120-i) + b[i].tolist() + [0] * (66 + i)


b_trans = np.transpose(b_)

print(b_[0][120], b[0][0])

print('b_trans.shape', b_trans.shape)

result = []
for i in range(120):
    tmp = []
    for j in range(10):
        tmp = tmp + b_trans[j*120+i].tolist()
    result.append(tmp)

print('np.array(result).shape', np.array(result).shape)
 

result2 = []
for i in range(20):
    result2.append([]) 
 
for i in range(120):
    result2[i%20] = result2[i%20] + result[i]

print('np.array(result2).shape', np.array(result2).shape)
df = pd.DataFrame(np.array(result2))
df.to_csv("wow4.csv", index=False, header=False)
