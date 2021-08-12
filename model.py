# -*- coding: utf-8 -*-
import numpy as np
import nltk
from nltk.util import ngrams
from nltk.probability import ConditionalFreqDist
from nltk.tokenize import word_tokenize
with open("a.txt", encoding='utf-8', ) as f:
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



def get_range(text_sample):
    a = []
    nums = []
    for sen in text_sample:
        num = predict(sen)
        nums.append(num)
        a.append((num, sen))
    nums = sorted(nums)
    

def get_c(s):
    if len(s) == 0: return n
    for i in range(len(s), 0, -1):
        if tuple(s[len(s) - i:]) in d[i]:
            return d[len(tuple(s[len(s) - i:]))][tuple(s[len(s) - i:])] + 1e-5
    return 1e-5


def predict(s):
    s = s.split()
    m = 3
    cur_sen =[]
    prob = 0
    for i in s:
        if len(cur_sen) == m:
            cur_sen = cur_sen[1:]
        cur_sen.append(i)
        cur = get_c(cur_sen)
        last = float(get_c(cur_sen[:-1])) + 1e-3
        if last < 1: last += 1
        prob += np.log(cur/last)

    return prob


