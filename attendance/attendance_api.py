#!/usr/local/bin/python
# pylint: disable=no-member, bare-except, line-too-long
"""
Attendance application is part of OT-Microservices.
This application stores and pull attendance data from MySQL
"""

import json
import os
import logging
import yaml
import mysql.connector
from flask import Flask, jsonify, request
from elasticapm.contrib.flask import ElasticAPM
from elasticapm.handlers.logging import LoggingHandler

ELASTIC_APM_SERVER_URL = os.getenv(
    'ELASTIC_APM_SERVER_URL', 'http://localhost:8200'
)
ELASTIC_APM_SERVICE_NAME = os.getenv('ELASTIC_APM_SERVICE_NAME', 'attendance')
CONFIG_FILE = os.getenv('CONFIG_FILE', '/app/config/config.yaml')

app = Flask(__name__)


def read_config():
    """Method for reading the configuration file for attendance"""
    with open(CONFIG_FILE, 'r', encoding="utf-8") as config_file:
        yaml_values = yaml.load(config_file, Loader=yaml.FullLoader)
    return yaml_values


config_properties = read_config()

apm = ElasticAPM(app,
                 server_url=ELASTIC_APM_SERVER_URL,
                 service_name=ELASTIC_APM_SERVICE_NAME,
                 logging=True)

app.config['MYSQL_HOST'] = config_properties['mysql']['host']
app.config['MYSQL_USERNAME'] = config_properties['mysql']['username']
app.config['MYSQL_PASSWORD'] = config_properties['mysql']['password']
app.config['MYSQL_DATABASE'] = config_properties['mysql']['db_name']


@app.route("/attendance/healthz", methods=['GET'])
def check_health():
    """Method for checking the health of MySQL"""
    try:
        connection = create_mysql_client()
        connection.ping()
        return jsonify(mysql="up", description="MySQL is healthy"), 200
    except Exception:
        app.logger.error("Unable to make mysql connection", exc_info=True)
        return jsonify(mysql="down", description="MySQL is not healthy"), 400


@app.route("/attendance/create", methods=['POST'])
def push_attendance_data():
    """For pushing attendance data inside MySQL Database"""
    record = json.loads(request.data)
    try:
        connection = create_mysql_client()
        cursor = connection.cursor()
        sql = """CREATE TABLE IF NOT EXISTS Employee (
                 id int(6) NOT NULL, status varchar(50) NOT NULL,
                 date varchar(50), PRIMARY KEY (id) )"""
        cursor.execute(sql)
    except Exception:
        app.logger.error("Unable to create table in MySQL", exc_info=True)
        return "Error", 400

    try:
        cursor = connection.cursor()
        sql = """INSERT INTO `Employee` (`id`, `status`, `date`) 
         VALUES (%s, %s, %s)"""
        cursor.execute(sql, (record['id'], record['status'], record['date']))
        app.logger.info("Successfully pushed data in MySQL")
        return jsonify(message="Successfully uploaded the attendance data"), 200
    except Exception:
        app.logger.error("Unable to push data in MySQL", exc_info=True)
        return jsonify(message="Error in pushing attendance data"), 400


@app.route("/attendance/search", methods=['GET'])
def fetch_attendance_data():
    """For pulling attendance data from MySQL Database"""
    try:
        connection = create_mysql_client()
        cursor = connection.cursor()
        sql = "SELECT * FROM Employee ORDER BY id DESC"
        cursor.execute(sql)
        results = cursor.fetchall()
        complete_data = []
        for result in results:
            data = {"id": result[0], "status": result[1], "date": result[2]}
            complete_data.append(data)
        app.logger.info("Successfully pulled data from MySQL")
        return jsonify(complete_data), 200
    except Exception:
        app.logger.error("Unable to pull data from MySQL", exc_info=True)
        return jsonify(message="Error while pulling data for attendance"), 200


def create_mysql_client():
    """For creating the client connection with MySQL"""
    connection = mysql.connector.connect(
        host=app.config['MYSQL_HOST'],
        user=app.config['MYSQL_USERNAME'],
        passwd=app.config['MYSQL_PASSWORD'],
        database=app.config['MYSQL_DATABASE']
    )
    return connection


if __name__ == "__main__":
    handler = LoggingHandler(client=apm.client)
    console_handler = logging.StreamHandler()
    app.logger.addHandler(handler)
    app.logger.addHandler(console_handler)
    app.run()
