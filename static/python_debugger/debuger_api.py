import flask
import debuger
import multiprocessing

app = flask.Flask(__name__)
d = debuger.Debugger()
debugq = multiprocessing.Queue()
outputq = multiprocessing.Queue()
breakpointQ = multiprocessing.Queue()

# TODO: error handling for each route

@app.route('/debug/set_breakpoint/<lineNo>', methods=['GET'])
def setBreakpoint(lineNo):
    if type(lineNo) != int:
        return flask.Response("Line number must be an integer", status=400)
    breakpointQ.put(lineNo)
    return flask.Response(status=200) 

@app.route('/debug/stepin', methods=['GET'])
def stepin():
    debugq.put("step_in")
    output = outputq.get()
    return flask.jsonify(output)

@app.route('/debug/stepout', methods=['GET'])
def stepout():
    debugq.put("step_out")
    output = outputq.get()
    return flask.jsonify(output)

@app.route('/debug/stepover', methods=['GET'])
def stepover():
    debugq.put("step_over")
    output = outputq.get()
    return flask.jsonify(output)

@app.route('run', methods=['GET'])
def run():
    debugq.put("run")
    output = outputq.get()
    return flask.jsonify(output)

if __name__ == '__main__':
    debugprocess = multiprocessing.Process(target=d.debug, args=(breakpointQ, outputq, debugq, debuger.coderunners_main))
    debugprocess.start()
    app.run(host='localhost', port='8000')