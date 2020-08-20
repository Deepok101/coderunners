import debugger
import multiprocessing
import time

def hello():
    deepak()
    deepak()
    deepak()
    return 4

def sample(a=2, b=3):
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
breakpointQ = None

def before_each(steps, breakpoints=[11]):
    global d, debugq, outputq, breakpointQ
    d = debugger.Debugger(breakpoints)
    debugq = multiprocessing.Queue()
    outputq = multiprocessing.Queue()
    breakpointQ = multiprocessing.Queue()
    for step in steps: 
        debugq.put(step)

    debugprocess = multiprocessing.Process(target=d.debug, args=(breakpointQ, outputq, debugq, sample))
    debugprocess.start()
    print("Running test")
    while (not debugq.empty()) and outputq.qsize() != (len(steps) + 1):
        pass
    time.sleep(0.5)
    debugprocess.terminate()
    print('')

def after_each():
    global d, debugq, outputq, breakpointQ
    debugq.close()
    outputq.close()
    breakpointQ.close()

def test_stepin():
    steps = ['step_in', 'step_in']
    before_each(steps)

    output = outputq.get()
    output = outputq.get()
    output = outputq.get()
    
    funcName = output['funcName']
    lineNo = output['lineNo']
    expectedFuncName = "hello"
    expectedLineNo = 5

    after_each()
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

    after_each()
    assert funcName == expectedFuncName
    assert lineNo == expectedLineNo

def test_stepout():
    steps = ['step_in', 'step_in', 'step_out']
    before_each(steps)

    outputq.get()
    outputq.get()
    outputq.get()
    output = outputq.get()
    print("a")

    funcName = output['funcName']
    lineNo = output['lineNo']
    expectedFuncName = "sample"
    expectedLineNo = 13

    after_each()
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

    after_each()
    assert funcName == expectedFuncName
    assert lineNo == expectedLineNo
    
def integration_test(steps, expectedFuncName, expectedLineNo):
    before_each(steps)

    for _ in range(len(steps)):
        outputq.get()
    output = outputq.get()

    funcName = output['funcName']
    lineNo = output['lineNo']

    after_each()
    assert funcName == expectedFuncName
    assert lineNo == expectedLineNo

def breakpoint_test(steps, breakpoints, expectedFuncName, expectedLineNo):
    before_each(steps, breakpoints)
    for _ in range(len(steps)):
        outputq.get()
    output = outputq.get()

    funcName = output['funcName']
    lineNo = output['lineNo']

    after_each()
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

    breakpoint_test(
        steps=['run'],
        breakpoints=[11,14],
        expectedFuncName="sample",
        expectedLineNo=14
        )



