from sys import settrace
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

class Debugger:
    # parentFunc
    # currFunc
    def __init__(self):  
        self.cmd = None
        self.currFrame = None

    def getFrame(self):
        return self.currFrame

    def trace_calls(self, frame, event, arg):
        if event != "call":
            return 
        codeObj = frame.f_code
        funcName = codeObj.co_name
        lineNo = frame.f_lineno

        if self.cmd == "step_out" or self.cmd == "step_over":
            return None

        self.outputq.put({"funcName": funcName, "lineNo": lineNo})
        print(f'Call to {funcName} on line {lineNo} of {codeObj.co_filename}')

        cmd = self.debugq.get()
        self.cmd = cmd
        if cmd == 'step_in':
            return self.trace_lines
        elif cmd == 'step_over':
            return None



    def trace_lines(self, frame, event ,arg):
        if event != 'line' and event != 'return':
            return

        codeObj = frame.f_code
        funcName = codeObj.co_name
        lineNo = frame.f_lineno

        self.outputq.put({"funcName": funcName, "lineNo": lineNo})
        print(f'Line to {funcName} on line {lineNo} of {codeObj.co_filename}')

        cmd = self.debugq.get()
        self.cmd = cmd
        if cmd == "step_in" or cmd == "step_over":
            return self.trace_lines
        elif cmd == "stop":
            exit()
        elif cmd == 'step_out':
            return self.trace_calls

    def debug(self, outputq, debugq, fn, args):
        self.debugq = debugq
        self.outputq = outputq
        settrace(self.trace_calls)
        fn(*args)
        
if __name__ == '__main__':
    d = Debugger()
    debugq = multiprocessing.Queue()
    outputq = multiprocessing.Queue()

    debugprocess = multiprocessing.Process(target=d.debug, args=(outputq, debugq, sample, (2, 3)))
    debugprocess.start()
