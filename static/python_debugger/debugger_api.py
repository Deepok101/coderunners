import flask
import debugger
import multiprocessing

app = flask.Flask(__name__)
debugq = multiprocessing.Queue()
outputq = multiprocessing.Queue()
breakpointQ = multiprocessing.Queue()

# TODO: error handling for each route
# TODO: handling end of function
@app.route('/debug/setup', methods=['POST'])
def setup():
    global d
    request_json = flask.request.get_json()
    print(request_json)
    arrayOfBreakpoints = request_json.get('breakpoints')

    if type(arrayOfBreakpoints) != type([]):
        return ("", 400)
    for i in arrayOfBreakpoints:
        if type(i) != int:
            return ("", 400)

    d = debugger.Debugger(arrayOfBreakpoints)
    debugprocess = multiprocessing.Process(target=d.debug, args=(breakpointQ, outputq, debugq, debugger.coderunners_main))
    debugprocess.start()
    output = outputq.get()
    return (flask.jsonify(output), 200)


@app.route('/debug/set_breakpoint/<lineNo>', methods=['GET'])
def setBreakpoint(lineNo):
    if type(lineNo) != int:
        return ("Line number must be an integer", 400)
    breakpointQ.put(lineNo)
    return ("", 200)

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

@app.route('/run', methods=['GET'])
def run():
    debugq.put("run")
    output = outputq.get()
    return flask.jsonify(output)

if __name__ == '__main__':
    app.run(host='localhost', port='8000')