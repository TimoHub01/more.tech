from http.server import BaseHTTPRequestHandler, HTTPServer
import logging
from ChannelMessages import main as main_e
import asyncio
from news_model import NewsModel
import numpy as np
from csv_to_sql import get_data_for_model
from csv_to_sql import to_db

class S(BaseHTTPRequestHandler):
    def _set_response(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()

    def do_GET(self):
        logging.info("GET request,\nPath: %s\nHeaders:\n%s\n", str(self.path), str(self.headers))
        self._set_response()
        self.wfile.write("GET request for {}".format(self.path).encode('utf-8'))
        if self.path.find("true") > 0:
            asyncio.run(main_e(phone="89859863313", url="https://t.me/rotten_rec"))
            scoring()


    def do_POST(self):
        content_length = int(self.headers['Content-Length'])  # <--- Gets the size of data
        post_data = self.rfile.read(content_length)  # <--- Gets the data itself
        logging.info("POST request,\nPath: %s\nHeaders:\n%s\n\nBody:\n%s\n",
                     str(self.path), str(self.headers), post_data.decode('utf-8'))

        self._set_response()
        self.wfile.write("POST request for {}".format(self.path).encode('utf-8'))


def run(server_class=HTTPServer, handler_class=S, port=8080):
    logging.basicConfig(level=logging.INFO)
    server_address = ('', port)
    httpd = server_class(server_address, handler_class)
    logging.info('Starting httpd...\n')
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        pass
    httpd.server_close()
    logging.info('Stopping httpd...\n')


def scoring():
    dict_keyword_weights = np.load('dict_keyword_weights.npy', allow_pickle='TRUE').item()
    data_from_db = get_data_for_model()
    print("\nDATA FROM DB\n")
    print(data_from_db)
    print("\nDATA FROM DB\n")
    # Получение данные с базы
    model = NewsModel(dict_keyword_weights)
    data_from_model = model.fit_predict(data_from_db.iloc[:500, :])
    print("TRIG")
    to_db(data_from_model)
    print("TRIG")


if __name__ == '__main__':
    from sys import argv

    if len(argv) == 2:
        run(port=int(argv[1]))
    else:
        run()
