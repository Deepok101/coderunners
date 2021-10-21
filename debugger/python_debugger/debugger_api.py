import flask
import debugger
import multiprocessing
import template

app = flask.Flask(__name__)
debugq = multiprocessing.Queue()
outputq = multiprocessing.Queue()
breakpointQ = multiprocessing.Queue()
debugprocess = None
DEBUG_TERMINATED_RES = {"running": False, "localVars": None, "lineNo": -1, "funcName": "", "breakpoints": []}

@app.route('/debug/setup', methods=['POST'])
def setup():
    global debugprocess
    request_json = flask.request.get_json()
    if request_json == None:
        return ("Please set the content-type to application/json and send data in json format", 400)
    
    arrayOfBreakpoints = request_json.get('breakpoints')

    if type(arrayOfBreakpoints) != type([]):
        return ("", 400)
    for i in arrayOfBreakpoints:
        if type(i) != int:
            return ("", 400)

    debuggerObj = debugger.Debugger(arrayOfBreakpoints)
    debugprocess = multiprocessing.Process(target=debuggerObj.debug, args=(breakpointQ, outputq, debugq, template.coderunners_exec))
    debugprocess.start()
    output = outputq.get()
    return (flask.jsonify(output), 200)

@app.route('/debug/set_breakpoint/<lineNo>', methods=['GET'])
def setBreakpoint(lineNo):
    if type(int(lineNo)) != int:
        return ("Line number must be an integer", 400)
    breakpointQ.put(lineNo)
    return ("", 200)

@app.route('/debug/stepin', methods=['GET'])
def stepin():
    if not debugprocess.is_alive():
        return flask.jsonify(DEBUG_TERMINATED_RES)

    debugq.put("step_in")
    output = outputq.get()
    if output["running"] == False:
        return flask.jsonify(DEBUG_TERMINATED_RES)
    return flask.jsonify(output)

@app.route('/debug/stepout', methods=['GET'])
def stepout():
    if not debugprocess.is_alive():
        return flask.jsonify(DEBUG_TERMINATED_RES)

    debugq.put("step_out")
    output = outputq.get()
    if output["running"] == False:
        return flask.jsonify(DEBUG_TERMINATED_RES)
    return flask.jsonify(output)

@app.route('/debug/stepover', methods=['GET'])
def stepover():
    if not debugprocess.is_alive():
        return flask.jsonify(DEBUG_TERMINATED_RES)

    debugq.put("step_over")
    output = outputq.get()
    if output["running"] == False:
        return flask.jsonify(DEBUG_TERMINATED_RES)
    return flask.jsonify(output)

@app.route('/run', methods=['GET'])
def run():
    debugq.put("run")
    output = outputq.get()
    return flask.jsonify(output)

if __name__ == '__main__':
    app.run(host='localhost', port='8000')