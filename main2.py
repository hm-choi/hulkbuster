import numpy as np
import pandas as pd 


b = np.random.normal(0,1,(120,1014))
df = pd.DataFrame(b)
df.to_csv("wow.csv", index=False, header=False)

b_ = np.random.normal(0,1,(120,1200))
for i in range(120):
    b_[i] = [0] * (120-i) + b[i].tolist() + [0] * (66 + i)


b_trans = np.transpose(b_)

print(b_[0][120], b[0][0])

print(b_trans.shape)

result = []
for i in range(120):
    tmp = []
    for j in range(1, 10):
        tmp = tmp + b_trans[j*120+i].tolist()
    tmp = tmp + [0] * (8072-len(tmp)) + b_trans[i].tolist()
    result.append(tmp)

df = pd.DataFrame(np.array(result))
df.to_csv("wow2.csv", index=False, header=False)
