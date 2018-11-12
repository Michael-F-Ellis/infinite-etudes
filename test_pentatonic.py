#pylint: disable=R,C,W
from pentatonic_triples import *

def test_shuffledTriples():
    assert len(shuffledTriples()) == 60

def test_excursion():
    assert excursion((1,3,5)) == 5
    assert excursion((1,3,5,7)) == 7
    assert excursion((1,5,3)) == -6
    assert excursion((1,4,1)) == 0
    assert excursion((1,5,1)) == 0
    assert excursion((1,5,3,1)) == -8
    assert excursion((1,3,5,1)) == 8
    assert excursion((1,3,5,1,3)) == 10

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
               Triple((1,4,7)),
               Triple((1,4,7)),
               ]
    constrain(triples)
    assert triples[0].offset == 0
    assert triples[1].offset == 0
    assert triples[2].offset == -1
    assert triples[3].offset == -1

def test_octaveMarks():
    assert octaveMarks(0) == ""
    assert octaveMarks(-1) == "/"
    assert octaveMarks(1) == "^"

def test_measure():
    assert measure(Triple((1,4,7))) == "1 4 7 z | "
    assert measure(Triple((1,4,7)), -1) == "/1 4 7 z | "

def test_line():
    assert line(Triple((1,4,7))) ==  "1 4 7 z | /1 4 7 z | /1 4 7 z | /1 4 7 z | "
