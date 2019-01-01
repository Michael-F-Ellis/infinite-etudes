#pylint: disable=R,C,W
from all_triples import *

def test_shuffledTriples():
    assert len(shuffledTriples()) == 210

def test_midiNext():
    assert midiNext(60, 1, 1) == 60
    assert midiNext(64, 2, 4) == 67 # Key of D, e to g
    assert midiNext(64, 2, 7) == 61 # Key of D, e to c#
    assert midiNext(70, 5, 4, harmonic=True) == 67
    assert midiNext(70, 5, 6, harmonic=True) == 71
    assert midiNext(67, 4, 5, harmonic=True) == 70
    assert midiNext(71, 6, 5, harmonic=True) == 70
    assert midiNext(70, 5, 1, harmonic=True) == 74

def test_midiFirst():
    assert midiFirst(0) == (60,1)
    assert midiFirst(2) == (61,7)
    assert midiFirst(7) == (60,4)
    assert midiFirst(5) == (60,5)
    assert midiFirst(5,harmonic=True) == (61,5)

def test_computeOctaveOffset():
    assert computeOctaveOffset(60, 60, 72) == 0
    assert computeOctaveOffset(60, 61, 72) == 0
    assert computeOctaveOffset(60, 72, 72) == 0
    assert computeOctaveOffset(60, 59, 72) == 1
    assert computeOctaveOffset(60, 73, 72) == -1
    assert computeOctaveOffset(60, 47, 72) == 2
    assert computeOctaveOffset(60, 85, 72) == -2

def test_computeTripleOffset():
    m0 = 60
    d0=1
    mlo=48
    mhi=72
    t = Triple((1,2,3))
    assert computeTripleOffset(m0, d0, t, mlo, mhi, False) == (True, 64, 3)
    assert t.offset == 0
    assert computeTripleOffset(36, d0, t, mlo, mhi, False) == (True, 52, 3)
    assert t.offset == 1
    assert computeTripleOffset(72, d0, t, mlo, mhi, False) == (True, 64, 3)
    assert t.offset == -1
    t = Triple((6,2,5))
    assert computeTripleOffset(m0, d0, t, mlo, mhi, False) == (True, 67, 5)
    assert t.offset == 0
    assert computeTripleOffset(36, d0, t, mlo, mhi, False) == (True, 67, 5)
    assert t.offset == 2
    assert computeTripleOffset(72, d0, t, mlo, mhi, False) == (True, 67, 5)
    assert t.offset == -1

def test_excursion():
    assert excursion((1,3,5)) == 5
    assert excursion((1,3,5,7)) == 7
    assert excursion((1,5,3)) == -6
    assert excursion((1,4,1)) == 1
    assert excursion((1,5,1)) == 1
    assert excursion((1,5,3,1)) == -8
    assert excursion((1,3,5,1)) == 8
    assert excursion((1,3,5,1,3)) == 10
    assert excursion((1,6,2)) == 2

def test_excursionMinMax():
    assert excursionMinMax((1,3,5)) == (1,5)
    assert excursionMinMax((1,5,3)) == (-6,1)
    assert excursionMinMax((1,6,2)) == (-3,2)

def test_repetitionOffset():
    assert repetitionOffset(4) == 0
    assert repetitionOffset(-4) == 0
    assert repetitionOffset(5) == -1
    assert repetitionOffset(-5) == 1
    assert repetitionOffset(11) == -1

def test_sequenceFromTriples():
    triples = [Triple((1,5,7)),
               Triple((2,3,4)),
               ]
    assert sequenceFromTriples(triples) == [1,5,7,2,3,4]

def test_constrain():
    triples = [Triple((1,4,7)),
               Triple((1,4,7)),
               Triple((5, 3, 1)),
               Triple((7, 4, 1)),
               Triple((7, 4, 1)),
               ]
    constrain(triples, -1, 1)
    assert triples[0].offset == 0
    assert triples[1].offset == -1
    assert triples[2].offset == 0
    assert triples[3].offset == 0
    assert triples[4].offset == 1
    triples = [
               Triple((1,5,2)),
               Triple((2,1,5)),
               Triple((1,2,3)),
               Triple((1,3,5)),
               ]
    constrain(triples, -1, 1)
    assert triples[0].offset == 0
    assert triples[1].offset == 1
    assert triples[2].offset == 0
    assert triples[3].offset == 0

def test_constrain2():
    triples = [Triple((1,4,7)),
               Triple((1,4,7)),
               Triple((5, 3, 1)),
               Triple((7, 4, 1)),
               Triple((7, 4, 1)),
               ]
    constrain2(triples, 48, 72, 0)
    assert triples[0].offset == 0
    assert triples[1].offset == -1
    assert triples[2].offset == 0
    assert triples[3].offset == 0
    assert triples[4].offset == 1
    triples = [
               Triple((1,5,2)),
               Triple((2,1,5)),
               Triple((1,2,3)),
               Triple((1,3,5)),
               ]
    constrain2(triples, 48, 72, 0)
    assert triples[0].offset == 0
    assert triples[1].offset == 1
    assert triples[2].offset == 0
    assert triples[3].offset == 0



def test_octaveMarks():
    assert octaveMarks(0) == ""
    assert octaveMarks(-1) == "/"
    assert octaveMarks(1) == "^"

def test_measure():
    assert measure(Triple((1,4,7))) == "1 4 7 z | "
    assert measure(Triple((1,4,7)), -1) == "/1 4 7 z | "
    assert measure(Triple((1,6,5))) == "1 6 5 z | "

def test_line():
    assert line(Triple((1,4,7))) ==  "1 4 7 z | /1 4 7 z | /1 4 7 z | /1 4 7 z | "

def test_getKeyNum():
    directives = 'P=1 K=E@ T=120'
    assert getKeyNum(directives) == 3
