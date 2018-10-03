from flask import Flask
from opencensus.trace.exporters.ocagent.trace_exporter import TraceExporter
from opencensus.trace.exporters.transports.background_thread import BackgroundThreadTransport
from opencensus.trace.ext.flask.flask_middleware import FlaskMiddleware
from opencensus.trace.propagation.trace_context_http_header_format import TraceContextPropagator
import os

app = Flask(__name__)
middleware = FlaskMiddleware(app,
	exporter = TraceExporter(
		service_name = 'python-service',
		endpoint = os.getenv('OCAGENT_TRACE_EXPORTER_ENDPOINT'),
		transport = BackgroundThreadTransport,
	),
	propagator = TraceContextPropagator())

@app.route("/")
def hello():
	return "Hello World!"

if __name__ == '__main__':
	import logging
	logger = logging.getLogger('werkzeug')
	logger.setLevel(logging.ERROR)
	app.run(host = 'localhost', port = 8080, threaded = True)
