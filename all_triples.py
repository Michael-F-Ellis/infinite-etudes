#!/usr/bin/env python
# pylint: disable=R,C
"""
etudes for ear training based on the pentatonic scale plus the 4th degree of the diatonic scale.
"""
import re
from itertools import permutations
from random import shuffle


# LUT for diatonic interval to midi interval
d2m = {1:0, 2:2, 3:4, 4:5, 5:7, 6:9, 7:11}



class Triple:
    def __init__(self, pitchnums):
        self.nums = pitchnums
        self.excursion = excursion(pitchnums)
        self.offset = 0



def shuffledTriples():
    """ Returns a shuffled list of tuples of strings that represent all pentatonic permutations. Each tuple
    is of the form (3,4,1) with elements drawn from the set (1,2,3,5,6) representing
    the scale degrees of the the major pentatonic scale. """
    triples = []
    for p in permutations((1,2,3,4,5,6,7), 3):
        triples.append(Triple(p))
    shuffle(triples)
    return triples

def intervalAsc(p1,p2):
    """ Return the ascending interval between two pitches numbers. """
    ## Convert from base 7 index 1 to base 7 index 0
    a = p1 - 1
    b = p2 - 1
    ## subtract, reduce modulo 7, convert back to 1 indexed.
    return 1 + (b - a) % 7

def intervalReduced(p1, p2):
    """ Return interval constrained a la Tbon, i.e. to nearest 4th """
    assert 1 <= p1 <= 7 and 1 <= p2 <= 7
    i = intervalAsc(p1, p2)
    return i if i <= 4 else -(9 - i)

def midiFirst(keynum, harmonic=False):
    """
    Return midi number of first pitch. Tbon assumes a prior pitch of middle C (60) or
    middle C# (61) depending on the key signature.
    """
    cdegreemidi = { # LUT by keynum
            0:(1,60), 1:(7,60), 2:(7,61), 3:(6,60),
            4:(6,61), 5:(5,60), 6:(4,59), 7:(4,60),
            8:(3,60), 9:(3,61), 10:(2,60), 11:(2,61)
            }
    d0,m0 = cdegreemidi[keynum]
    if keynum == 5 and harmonic:
        m0 += 1
    return m0, d0


def midiNext(m0, d0, d1, harmonic=False):
    """
    Return midi number of diatonic pitch number d1 given preceding midi and diatonic
    pitch numbers m0 and d0.
    """
    if d0 == d1:
        return m0 ## same pitch
    ir = intervalReduced(d0, d1)
    extra = 0
    if harmonic:
        if d0 == 5:
            extra = -1
        elif d1 == 5:
            extra = 1
    if ir >= 1 and d1 > d0:
        ## positive ir, d1 higher
        return m0 + d2m[d1] - d2m[d0] + extra
    elif ir >= 1 and d1 < d0:
        ## positive ir, d1 lower
        return m0 + 12 + d2m[d1] - d2m[d0] + extra
    elif ir < 1 and d1 > d0:
        ## negative ir, d1 higher
        return m0 - 12 + d2m[d1] - d2m[d0] + extra
    elif ir < 1 and d1 < d0:
        ## negative ir, d1 lower
        return m0 + d2m[d1] - d2m[d0] + extra

def computeOctaveOffset(mlo, m, mhi):
    """
    Return an integer number of octaves needed to move midi pitch m within
    the inclusive interval between mlo and mhi.
    """
    assert 0 <= mlo < mhi < 128  ## require valid midi range
    offset12 = 0
    while not(mlo <= offset12 + m <= mhi):
        if offset12 + m < mlo:
            offset12 += 12
        else:
            offset12 -= 12

    # done. convert to octaves and return
    return offset12 // 12

def computeTripleOffset(m0, d0, t, mlo, mhi, harmonic=False):
    """
    m0 : prior midi pitch
    d0 : prior diatonic pitch
    t : Triple instance
    mlo: midi low limit
    mhi: midi high limit
    Sets the offset attribute of t, if possible
    Returns (True, last midi) if an offset value is found that can successfully
    place all pitches within the range from mlo to mhi inclusive.
    Returns (False, None) if not (and leave t.offset unaltered).
    """
    print()
    trialoffset=0
    tries = 0
    while True:
        midinums = []
        m = m0
        d = d0
        for n in t.nums:
            m = midiNext(m, d, n, harmonic)
            d = n
            midinums.append(m)

        offsets = []
        for m in midinums:
            offsets.append(computeOctaveOffset(mlo, m + (12*trialoffset), mhi))

        success = all([offset == offsets[0] for offset in offsets])
        if success:
            t.offset = offsets[0] if tries == 0 else trialoffset
            return success, midinums[2] + 12*t.offset, t.nums[2]

        elif tries < 3:
            print(tries, offsets)
            trialoffset = offsets[tries]
            tries += 1
            continue
        else:
            return success, None, None ## failed

def constrain2(triples, mlo, mhi, keynum, harmonic=False):
    """
    Compute offsets for each triple such that all pitches are within
    the interval mlo:mhi inclusive.
    """
    m, d = midiFirst(keynum, harmonic)
    for t in triples:
        success, m, d = computeTripleOffset(m, d, t, mlo, mhi, harmonic)
        assert success


def excursionMinMax(pitches):
    """ Return a tuple of min and max excursions of a sequence of pitches. """
    x = 0
    xmin = 0
    xmax = 0
    #print("p ir x xmin xmax")
    for i, p in enumerate(pitches[0:-1]):
        ir = intervalReduced(p, pitches[i+1])
        x = x + (ir-1) if ir > 0 else x + (ir + 1)
        xmax = max(x, xmax)
        xmin = min(x, xmin)
        #print(p, ir, x, xmin, xmax)
    xmax = xmax + 1 if xmax >= 0 else xmax - 1
    xmin = xmin + 1 if xmin >= 0 else xmin - 1
    return (xmin, xmax)

def excursion(pitches):
    """ Return the intervallic excursion of a sequence of pitches """
    x = 0
    for i, p in enumerate(pitches[0:-1]):
        ir = intervalReduced(p, pitches[i+1])
        x = x + (ir-1) if ir > 0 else x + (ir + 1)
    return x + 1 if x >= 0 else x - 1

def sequenceFromTriples(triples):
    """ Catenate the pitch numbers from a list of triples """
    seq = []
    for t in triples:
        for n in t.nums:
            seq.append(n)
    return seq

def sumOfOffsets(triples):
    """ Algebraic sum of octave offsets from Triples in list. """
    osum = 0
    for t in triples:
        osum += t.offset
    return osum

def constrain(triples, lo8=-2, hi8=2):
    """ Ensures total excursion does not go outside the octave limits defined by lo8 and hi8 """
    lo = lo8 * 8
    hi = hi8 * 8
    for i, _ in enumerate(triples):
        sublist = triples[0:i+1]
        seq = sequenceFromTriples(sublist)
        xmin, xmax = excursionMinMax(seq)
        so = sumOfOffsets(sublist)
        xlo = xmin + so * 8
        xhi = xmax + so * 8
        #print("{}: {}, {}, {}, {}, {}".format(i,xmin, xmax, so,xlo,xhi))
        if xlo < lo:
            triples[i].offset += 1 # raise last triple by one octave
        elif xhi > hi:
            triples[i].offset -= 1 # lower last triple by one octave


def repetitionOffset(xcursion):
    """ Return an integer octave offset such that a sequence of pitches can be
    repeated without changing the total excursion."""
    if xcursion == 0:
        return 0
    elif xcursion > 0:
        n, r = divmod(xcursion, 8)
        return -n if r + n <= 4 else n - 1
    else:
        n, r = divmod(xcursion, -8)
        return n if r - n >= -4 else n + 1


def octaveMarks(offset):
    """ Return tbon octave marks corresponding to offset. """
    if offset <= 0:
        marks = "/" * abs(offset)
    else:
        marks = "^" * offset
    return marks

def measure(triple, offset=None):
    """ Returns a measure of Tbon notation of the form "3 4 1 z | " """
    if offset is None:
        totaloffset = triple.offset
        print("{} {} {}".format(triple.nums, triple.excursion, triple.offset))
    else:
        totaloffset = offset
    m = [str(n) for n in triple.nums]
    m.append("z | ")
    return octaveMarks(totaloffset) + " ".join(m)


def line(triple):
    """ Return 4 measures of the triple with last 3 offset as needed for no excursion. """
    ofs = repetitionOffset(triple.excursion)
    return measure(triple, offset=None) + measure(triple, offset=ofs) * 3


def bins(triples):
    """ Return tuple with 5 lists of triples:
          0. Pentatonic only (1 2 3 5 6)
          1. Degree 4 but not degree 7
          2. Degree 7 but not 4
          3. Contains 4 and 7
          4. Pentatonic with Degree 5 (to be sharped for harmonic minor)
          5. Degree 5 with at least one of (4,7)
    """
    result = ([], [], [], [], [], [])
    for t in triples:
        tset = set(t.nums)
        if 5 in t.nums:
            if set((4,7)).intersection(tset) == set():
                result[4].append(t)
            else:
                result[5].append(t)
        if set((4,7)).issubset(tset):
            result[3].append(t)
            continue
        if 4 in t.nums:
            result[1].append(t)
            continue
        if 7 in t.nums:
            result[2].append(t)
            continue

        result[0].append(t) # pentatonic

    return result

def getKeyNum(directives):
    """
    Find the K= directive, extract the keyname, return the number from the LUT.
    """
    lut = {
        'C':0,
        'D@':1,
        'D':2,
        'E@':3,
        'E':4,
        'F':5,
        'G@':6,
        'G':7,
        'A@':8,
        'A':9,
        'B@':10,
        'B':11,
        }
    pat = r"[^K]*K=([A-G]@?).*"
    m = re.match(pat, directives)
    if m is None:
        ## not found, assume key of C.
        return 0
    kstring = m.groups()[0]
    if kstring not in lut:
        return 0
    return lut[kstring]



def etude(mlo, mhi, triples, directives="K=E@ T=120", countin="z - - - |", hminor=False):
    """ Return lines of Tbon notation, one line for each possible triple."""
    e = [directives, countin]
    kn = getKeyNum(directives)
    constrain2(triples, mlo, mhi, kn, hminor)
    for t in triples:
        l = line(t)
        if hminor:
            l = l.replace("5", "#5")
        e.append(l)
    return "\n\n".join(e)


def mkEtudes(mlo, mhi, directives="K=E@ T=120", countin="z - - - |"):
    """ Return 4 etudes, 1 for each bin """
    triples = shuffledTriples()
    tbins = bins(triples)
    etudes = [etude(mlo, mhi, tbin, directives, countin) for tbin in tbins[0:4]]
    etudes.append(etude(mlo, mhi, tbins[4], directives, countin, hminor=True))
    etudes.append(etude(mlo, mhi, tbins[5], directives, countin, hminor=True))
    return etudes


if __name__ == "__main__":
    import argparse
    parser = argparse.ArgumentParser(description="Etude Generator")
    parser.add_argument('-d', "--directives", type=str, default="K=E@ T=120",
                        help='Tbon directives. Default="K=E@ T=120"')
    parser.add_argument('-u', '--midihigh', type=int, default=84,
            help='High limit for pitches as a midi number. Default=84')
    parser.add_argument('-l', '--midilow', type=int, default=36,
            help='Low limit for pitches as a midi number. Default=36')

    args = parser.parse_args()
    outfiles = ("pentatonic.tbn", "plus4.tbn", "plus7.tbn", "both47.tbn", "harmonic5.tbn", "harmonic47.tbn")
    for outname, etd in zip(outfiles, mkEtudes(args.midilow, args.midihigh,
                                               directives=args.directives)):
        with open(outname, 'w') as f:
            print(etd, file=f)
