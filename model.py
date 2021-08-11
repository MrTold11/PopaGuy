# -*- coding: utf-8 -*-
import numpy as np
import nltk
from nltk.util import ngrams
from nltk.probability import ConditionalFreqDist
from nltk.tokenize import word_tokenize
with open("E:/d/a.txt", encoding='utf-8', ) as f:
    text = f.read()

m = 1
#d1 = dict(nltk.FreqDist(text.split()))
#d2 = dict(nltk.FreqDist(nltk.ngrams(text.split(), 2)))
from tqdm.auto import tqdm
text = text.split('\n')
n = int(text[0]) #90134 , 36
keys = []
values = []
d = [{}]
#print(int(a), len(str(a)))
#print(text[36].split()[11])
for i in range(1, len(text)):
    if i % 2 == 1:
        a = text[i].replace('dict_keys[', '').split(" ")
        ng_sz = (i + 1)//2
        keys = []
        for i in range(0, len(a), ng_sz):
            keys.append(tuple(a[i:(i + ng_sz)]))
        #print(a[:10])
    else:
        values = list(map(int, map(str, text[i].split(" ")[:-1])))
        d.append(dict(zip(keys, values)))
print("done")


def get_range(text_sample):
    a = []
    nums = []
    for sen in text_sample:
        num = predict(sen)
        nums.append(num)
        a.append((num, sen))
    nums = sorted(nums)
    

def get_c(s):
    if tuple(s) in d[len(tuple(s))]:
        return d[len(tuple(s))][tuple(s)]
    return 0


def predict(s):
    s = s.split()
    if len(s) < 3:
        gr_sz = len(s)
    else:
        gr_sz = 3
    gr_sz = 1
    trgs = nltk.ngrams(s, gr_sz)
    a = list(trgs)
    prob = 0
    prob += np.log(get_c(a[0])/n)
    last = get_c(a[0])
    currentsentence = []
    for i in range(len(a[0])):
        currentsentence.append(a[0][i])
    for i in range(1, len(a)):

        if len(currentsentence) < m:
            currentsentence.append(a[i][len(a[i]) - 1])
        else:
            currentsentence = currentsentence[1:]
            currentsentence.append(a[i][len(a[i]) - 1])
        cur = get_c(currentsentence)
        if last == 0:
            last += 1
        prob += np.log(cur/last)
        last = cur

    return prob/len(s)


#print(d[1][("привет",)])
#while True: print(predict(input()))