# coding: utf-8
# pylint: disable=R,C
"""
This script demonstrates that almost all possible 3-note combinations are found
in the diatonic and harmonic minor scales.
"""
from itertools import combinations
chrom = [c for c in combinations(range(12), 3)]
print("There {} unique 3-note combinations in the chromatic scale.".format(len(chrom)))

## Make sequences containing the diatonic and harmonic scale patterns
## as elements of a 12-tone chromatic scale.
dia = [0, 2, 4, 5, 7, 9, 11]
harm = [0, 2, 4, 5, 8, 9, 11]

## We need to be able to transpose sequences to produce scales in
## all keys.
def transpose(seq, n):
    return sorted([(i+n)%12 for i in seq])

## Now let's find out how many of the chromatic combos are
## exist in the set of all 12 diatonic keys.
found = []
for i in range(12):
    combos = [c for c in combinations(transpose(dia,i),3)]
    for c in combos:
        if c in chrom:
            found.append(c)

print("There {} unique 3-note combinations in the diatonic scales.".format(len(set(found))))

## Next, we include all the combos from harmonic minor scales.
for i in range(12):
    combos = [c for c in combinations(transpose(harm,i),3)]
    for c in combos:
        if c in chrom:
            found.append(c)

nfound = len(set(found))
print("There {} unique 3-note combinations in the diatonic + harmonic minor scales.".format(nfound))

## Finally, print the chromatic combinations not found in diatonic + harmonic minor
print("There are {} chromatic combinations not found in the combined scales.".format(len(chrom) - nfound))
for c in chrom:
    if c not in found:
        print(c)
