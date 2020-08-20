def coderunners_main():
    x = 2+2
    y=3
    for i in range(3):
        x=x+2
    y=y*2
    z=x+1

from sys import settrace
import multiprocessing

class Debugger:
    # TODO: Improve the breakpoint mechanisms
    # parentFunc
    # currFunc
    def __init__(self, arrayOfBreakPoints):  
        self.cmd = 'run'
        self.currFrame = None
        self.breakPoints = arrayOfBreakPoints

    def trace_calls(self, frame, event, arg):
        if event != "call":
            return 
        codeObj = frame.f_code
        funcName = codeObj.co_name
        lineNo = frame.f_lineno
        localVars = frame.f_locals

        if self.cmd == 'run' and not self.stop_here(frame):
            return self.trace_lines

        if self.cmd == "step_out" or self.cmd == "step_over":
            return None

        self.outputq.put(
                {
                    "funcName": funcName, 
                    "lineNo": lineNo,
                    "localVars": localVars
                }
            )

        print(f'Call to {funcName} on line {lineNo} of {codeObj.co_filename}; {localVars}')

        cmd = self.debugq.get()
        self.cmd = cmd
        self.checkBreakpointQueue()
        if cmd == 'step_in':
            return self.trace_lines
        elif cmd == 'step_over':
            return None
        elif cmd == 'run':
            return self.trace_lines



    def trace_lines(self, frame, event ,arg):
        if event != 'line' and event != 'return':
            return

        if self.cmd == 'run' and not self.stop_here(frame):
            return self.trace_lines

        codeObj = frame.f_code
        funcName = codeObj.co_name
        lineNo = frame.f_lineno
        localVars = frame.f_locals
        self.outputq.put(
                {
                    "funcName": funcName, 
                    "lineNo": lineNo,
                    "localVars": localVars
                }
            )

        print(f'Line to {funcName} on line {lineNo} of {codeObj.co_filename}; {localVars}')

        cmd = self.debugq.get()
        self.cmd = cmd
        self.checkBreakpointQueue()
        if cmd == "step_in" or cmd == "step_over" or cmd == "run":
            return self.trace_lines
        elif cmd == "stop":
            exit()
        elif cmd == 'step_out':
            return self.trace_calls

    def stop_here(self, frame):
        # TODO: Instead of looking at line number, should look if the current stack is above the breakpoint stack location or if both are in the same stack, then we can compare line numbers
        # For now this is ok but less efficient because we need to call a traceback function for every line.
        lineNo = frame.f_lineno
        return lineNo in self.breakPoints

    def checkBreakpointQueue(self):
        if not self.breakpointQ.empty():
            newBreakpoint = self.breakpointQ.get()
            self.breakPoints.append(newBreakpoint)

    def debug(self, breakpointQ, outputq, debugq, fn):
        self.debugq = debugq
        self.outputq = outputq
        self.breakpointQ = breakpointQ
        settrace(self.trace_calls)
        fn()

if __name__ == '__main__':
    d = Debugger([1])
    debugq = multiprocessing.Queue()
    outputq = multiprocessing.Queue()
    breakpointQ = multiprocessing.Queue()
    debugprocess = multiprocessing.Process(target=d.debug, args=(breakpointQ, outputq, debugq, coderunners_main))
    debugprocess.start()
