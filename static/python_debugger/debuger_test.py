import debuger
import multiprocessing


def hello():
    deepak()
    deepak()
    deepak()
    return 4

def sample(a, b):
    hello()
    x = a + b
    y = x * 2
    hello()
    print('Sample: ' + str(y))

def deepak():
    return 5


d = None
debugq = None
outputq = None

def before_each(steps):
    global d, debugq, outputq
    d = debuger.Debugger()
    debugq = multiprocessing.Queue()
    outputq = multiprocessing.Queue()
    for step in steps: 
        debugq.put(step)

    debugprocess = multiprocessing.Process(target=d.debug, args=(outputq, debugq, sample, (2, 3)))
    debugprocess.start()
    while not debugq.empty():
        pass
    debugprocess.terminate()

def test_stepin():
    steps = ['step_in', 'step_in']
    before_each(steps)

    outputq.get()
    outputq.get()
    output = outputq.get()
    
    funcName = output['funcName']
    lineNo = output['lineNo']
    expectedFuncName = "hello"
    expectedLineNo = 5


    assert funcName == expectedFuncName
    assert lineNo == expectedLineNo
    
def test_stepover():
    steps = ['step_in', 'step_over']
    before_each(steps)

    outputq.get()
    outputq.get()
    output = outputq.get()

    funcName = output['funcName']
    lineNo = output['lineNo']
    expectedFuncName = "sample"
    expectedLineNo = 13

    assert funcName == expectedFuncName
    assert lineNo == expectedLineNo

def test_stepout():
    steps = ['step_in', 'step_in', 'step_out']
    before_each(steps)

    outputq.get()
    outputq.get()
    outputq.get()
    output = outputq.get()

    funcName = output['funcName']
    lineNo = output['lineNo']
    expectedFuncName = "sample"
    expectedLineNo = 13

    assert funcName == expectedFuncName
    assert lineNo == expectedLineNo


def integration_test1(): 
    steps = ['step_in', 'step_in', 'step_in', 'step_over', 'step_over']
    before_each(steps)

    for _ in range(len(steps)):
        outputq.get()
    output = outputq.get()

    funcName = output['funcName']
    lineNo = output['lineNo']
    expectedFuncName = "hello"
    expectedLineNo = 8

    assert funcName == expectedFuncName
    assert lineNo == expectedLineNo
    
def integration_test(steps, expectedFuncName, expectedLineNo):
    before_each(steps)

    for _ in range(len(steps)):
        outputq.get()
    output = outputq.get()

    funcName = output['funcName']
    lineNo = output['lineNo']

    assert funcName == expectedFuncName
    assert lineNo == expectedLineNo

if __name__ == '__main__':
    test_stepin()
    test_stepover()
    test_stepout()

    integration_test(
        steps=['step_in', 'step_in', 'step_in', 'step_over', 'step_in'], 
        expectedFuncName='deepak', 
        expectedLineNo=18
        )
    integration_test(
        steps=['step_in', 'step_in', 'step_in', 'step_over', 'step_in', 'step_in', 'step_out'], 
        expectedFuncName='hello', 
        expectedLineNo=8
        )
    # Jump into hello's definition and step over to line 13
    integration_test(
        steps=['step_in', 'step_in', 'step_over'], 
        expectedFuncName='sample', 
        expectedLineNo=13
        )
    # Jump into hello() and step over through the whole function
    integration_test(
        steps=['step_in', 'step_in', 'step_in', 'step_over', 'step_over', 'step_over', 'step_over', 'step_over'], 
        expectedFuncName='sample', 
        expectedLineNo=13
        )
    # Step over the whole sample function
    integration_test(
        steps=['step_in', 'step_over', 'step_over', 'step_over', 'step_over'], 
        expectedFuncName='sample', 
        expectedLineNo=16
        )
    # Step over the whole sample function
    integration_test(
        steps=['step_in', 'step_over', 'step_over', 'step_over', 'step_in', 'step_in', 'step_in', 'step_in', 'step_out', 'step_out'], 
        expectedFuncName='sample', 
        expectedLineNo=16
        )



